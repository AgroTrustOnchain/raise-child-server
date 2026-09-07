package business

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"raise-child/constants/env"
	"raise-child/constants/env/payment"
	payment_env "raise-child/constants/env/payment"
	"raise-child/constants/noti"
	"raise-child/constants/shared"
	"raise-child/interfaces/business"
	i_repository "raise-child/interfaces/repository"
	"raise-child/model/dtos/request"
	"raise-child/model/dtos/response"
	"raise-child/model/entities"
	"raise-child/repository"
	"raise-child/util"
	"raise-child/util/cache"
	"raise-child/util/db"
	on_chain "raise-child/util/on_chain"
	"slices"
	"time"

	"github.com/block-vision/sui-go-sdk/constant"
	"github.com/block-vision/sui-go-sdk/models"
	"github.com/block-vision/sui-go-sdk/sui"
	"github.com/block-vision/sui-go-sdk/utils"
	"github.com/payOSHQ/payos-lib-golang"
)

const (
	payment_pending_status string = "Pending"
	payment_success_status string = "Success"
	payment_cancel_status  string = "Canceled"
)

type paymentService struct {
	mealSupportDurationRepo i_repository.IMealSupportDurationRepository
	leaderNotiRepo          i_repository.ILeaderNotiRepository
	childTaskDetailRepo     i_repository.IChildTaskDetailRepository
	taskRepo                i_repository.ITaskRepository
	paymentRepo             i_repository.IPaymentRepository
	bankProfileRepo         i_repository.IBankProfileRepository
	profileRepo             i_repository.IProfileRepository
	donationRepo            i_repository.IOffChainDonationRepository
	withdrawRepo            i_repository.IOffChainWithdrawProposalRepository
	redisCache              cache.IRedisCache
	clients                 map[string]sui.ISuiAPI
	errLogger               *log.Logger
}

func initializePaymentService(
	mealSupportDurationRepo i_repository.IMealSupportDurationRepository,
	leaderNotiRepo i_repository.ILeaderNotiRepository,
	childTaskDetailRepo i_repository.IChildTaskDetailRepository,
	taskRepo i_repository.ITaskRepository,
	paymentRepo i_repository.IPaymentRepository,
	bankProfileRepo i_repository.IBankProfileRepository,
	profileRepo i_repository.IProfileRepository,
	donationRepo i_repository.IOffChainDonationRepository,
	withdrawRepo i_repository.IOffChainWithdrawProposalRepository,
	redisCache cache.IRedisCache,
	clients map[string]sui.ISuiAPI,
	errLogger *log.Logger,
) business.IPaymentService {
	return &paymentService{
		mealSupportDurationRepo: mealSupportDurationRepo,
		leaderNotiRepo:          leaderNotiRepo,
		childTaskDetailRepo:     childTaskDetailRepo,
		taskRepo:                taskRepo,
		paymentRepo:             paymentRepo,
		bankProfileRepo:         bankProfileRepo,
		profileRepo:             profileRepo,
		donationRepo:            donationRepo,
		withdrawRepo:            withdrawRepo,
		redisCache:              redisCache,
		clients:                 clients,
		errLogger:               errLogger,
	}
}

func GeneratePaymentService() (business.IPaymentService, error) {
	var errLogger = util.GetLogConfig(shared.ERROR_LEVEL)

	cnn, err := db.ConnectDB(errLogger, db.InitializePostgreSQL())
	if err != nil {
		return nil, err
	}

	return initializePaymentService(
		repository.InitializeMealSupportDurationRepository(cnn, errLogger),
		repository.InitializeLeaderNotiRepository(cnn, errLogger),
		repository.InitializeChildTaskDetailRepository(cnn, errLogger),
		repository.InitializeTaskRepository(cnn, errLogger),
		repository.InitializePaymentRepository(cnn, errLogger),
		repository.InitializeBankProfileRepository(cnn, errLogger),
		repository.InitializeProfileRepository(cnn, errLogger),
		repository.InitializeOffChainDonationRepository(cnn, errLogger),
		repository.InitializeOffChainWithdrawProposalRepository(cnn, errLogger),
		cache.InitializeRedisCache(),
		_networkAliases,
		errLogger,
	), nil
}

// ApprovePayment implements business.IPaymentService.
func (p *paymentService) ApprovePayment(id string, ctx context.Context) error {
	payment, err := p.paymentRepo.GetPaymentById(id, ctx)
	if err != nil {
		return err
	}

	var genericErr error = errors.New(noti.GENERIC_ERROR_WARN_MSG)
	if payment == nil || payment.Method != shared.MANUAL_BANK_METHOD || payment.IsDonateTx {
		return genericErr
	}

	var curTime time.Time = time.Now()
	if curTime.Before(payment.ExpiredAt) {
		return errors.New(noti.STILL_PENDING_PAYMENT_MESSAGE)
	}

	if !payment.IsTransferred {
		return errors.New(noti.PAYMENT_NOT_TRANSFERRED_MESSAGE)
	}

	var sender string = ctx.Value("address").(string)
	var profileId string = ctx.Value("sub").(string)
	var genericRightErr error = errors.New(noti.GENERIC_RIGHT_ACCESS_WARN_MSG)
	if payment.Actor == sender || payment.ProfileID == profileId {
		return genericRightErr
	}

	var client = p.clients[constant.SuiTestnet]
	manage, err := on_chain.GetOnChainObject[entities.Manage](on_chain.GetOnChainObjectRequest{
		Client:    client,
		ObjectId:  os.Getenv(env.MANAGE_OBJECT_ID),
		ErrLogger: p.errLogger,
	}, ctx)
	if err != nil {
		return err
	}

	var internalErr error = errors.New(noti.INTERNALL_ERR_MSG)
	if manage == nil {
		return internalErr
	}

	if !slices.Contains(manage.AdminIds, sender) {
		return genericRightErr
	}

	var module string
	var function string
	var args []interface{}
	proposal, err := on_chain.GetOnChainObject[entities.WithdrawProposal](on_chain.GetOnChainObjectRequest{
		Client:    client,
		ObjectId:  *payment.ProposalID,
		ErrLogger: p.errLogger,
	}, ctx)
	if err != nil {
		return err
	}

	if proposal.Purpose == string(entities.POOL_WITHDRAW_PROPOSAL_PURPOSE) {
		profile, err := p.profileRepo.GetProfileByWalletAddress(proposal.Creator, ctx)
		if err != nil {
			return err
		}

		if err := p.taskRepo.CreateTask(entities.Task{
			ID:                util.GenerateId(),
			CreatedBy:         "System",
			AssignedProfileID: &profile.ID,
			AssignedStaff:     &proposal.Creator,
			Region:            proposal.PoolName,
			StartPeriod:       util.ToStartOfDate(curTime),
			EndPeriod:         util.ToEndOfDate(curTime),
			CreatedAt:         curTime,
			UpdatedAt:         curTime,
		}, ctx); err != nil {
			return err
		}

		var poolModule = on_chain.InitializeModulePool()
		module = poolModule.GetModule()
		function = poolModule.GetFunctionWithdrawFromPool()

		var localPoolId string
		if proposal.TargetID == os.Getenv(env.POOL_ID) {
			localPoolId = os.Getenv(env.SHARED_LOCAL_POOL_ID)
		} else {
			localPoolId = proposal.TargetID
		}

		args = poolModule.ToWithdrawFromPoolArguments(on_chain.WithdrawFromPoolArguments{
			LocalPoolId:        localPoolId,
			WithdrawProposalId: *payment.ProposalID,
			Sender:             payment.Actor,
		})

	} else if proposal.Purpose == string(entities.CAMPAIGN_PURPOSE) {
		var campaignModule = on_chain.InitializeModuleCampaign()
		module = campaignModule.GetModule()
		function = campaignModule.GetFunctionWithdrawFromCampaign()
		args = campaignModule.ToWithdrawFromCampaignArguments(on_chain.WithdrawFromCampaignArguments{
			CampaignID: proposal.TargetID,
			ProposalID: *payment.ProposalID,
			Sender:     payment.Actor,
		})
	} else {
		var childModule = on_chain.InitializeModuleChild()
		module = childModule.GetModule()
		switch proposal.Purpose {
		case string(entities.BOOKS_NEED_WITHDRAW_PROPOSAL_PURPOSE):
			function = childModule.GetFunctionWithdrawFromBooksNeedProposal()
		case string(entities.HEALTH_INSURANCE_NEED_WITHDRAW_PROPOSAL_PURPOSE):
			function = childModule.GetFunctionWithdrawFromHealthInsuranceNeedProposal()
		case string(entities.MEAL_NEED_WITHDRAW_PROPOSAL_PURPOSE):
			function = childModule.GetFunctionWithdrawFromMealNeedProposal()
		case string(entities.SPECIAL_NEED_CAMPAIGN_WITHDRAW_PROPOSAL_PURPOSE):
			function = childModule.GetFunctionWithdrawFromSpecialNeedCampaign()
		}

		args = childModule.ToWithdrawFromNeedArgumentsV2(on_chain.WithdrawFromNeedArgumentsV2{
			LocalPool:     proposal.PoolID,
			TargetID:      proposal.TargetID,
			ProposalID:    *payment.ProposalID,
			TransferredAt: util.ToMilliseconds(*payment.TransferredAt),
			Creator:       payment.Actor,
			Sender:        sender,
		})
	}

	payment.ReviewedBy = &sender
	payment.ReviewStatus = request_approved_status
	payment.UpdatedAt = curTime
	if err := p.paymentRepo.UpdatePayment(*payment, ctx); err != nil {
		return err
	}

	var req = on_chain.ExecuteTransactionRequestV2{
		Client:    client,
		Module:    module,
		Function:  function,
		Arguments: args,
		ErrLogger: p.errLogger,
	}

	for i := 1; i <= 3; i++ {
		if _, err := on_chain.ExecuteTransactionV2(req, ctx); err == nil {
			return nil
		}
	}

	return internalErr
}

// GetPayment implements business.IPaymentService.
func (p *paymentService) GetPayment(id string, ctx context.Context) (*entities.Payment, error) {
	return p.paymentRepo.GetPaymentById(id, ctx)
}

// GetPayments implements business.IPaymentService.
func (p *paymentService) GetPayments(req request.GetPaymentsRequest, ctx context.Context) (response.PaginationDataResponse, error) {
	var genericErr error = errors.New(noti.GENERIC_ERROR_WARN_MSG)
	if req.Actor != "" {
		if !util.IsValidSuiAddressStrict(req.Actor) {
			return response.PaginationDataResponse{}, genericErr
		}
	}

	if req.MaxAmount != nil {
		if *req.MaxAmount < 0 {
			return response.PaginationDataResponse{}, nil
		}
	}

	if req.MinAmount != nil {
		if req.MaxAmount != nil {
			if *req.MinAmount > *req.MaxAmount {
				return response.PaginationDataResponse{}, nil
			}
		} else {
			if *req.MinAmount < 0 {
				return response.PaginationDataResponse{}, nil
			}
		}
	}

	req.SortOrder = util.StandardizeSortOrder(req.SortOrder)
	req.Keyword = util.StandardizeString(req.Keyword)
	req.FilterProp = util.StandardizeSortCriteria(req.FilterProp)
	if req.Page < 1 {
		req.Page = 1
	}

	if req.PageSize < 1 {
		req.PageSize = default_page_size
	}

	var res response.PaginationDataResponse
	// var redisKey string = p.getGetPaymentsRedisKey(req)
	// if p.redisCache.Get(redisKey, &res, ctx) {
	// 	return res, nil
	// }

	data, pages, err := p.paymentRepo.GetPayments(req, ctx)
	if err != nil {
		return response.PaginationDataResponse{}, err
	}

	var amount int
	if len(data) == 0 {
		amount = 0
	} else {
		amount = len(data)
	}

	res = response.PaginationDataResponse{
		Data:       data,
		Amount:     amount,
		Page:       req.Page,
		TotalPages: pages,
	}

	// p.redisCache.Set(redisKey, res, time.Minute*2, ctx)

	return res, nil
}

// RefusePayment implements business.IPaymentService.
func (p *paymentService) RefusePayment(id string, ctx context.Context) error {
	payment, err := p.paymentRepo.GetPaymentById(id, ctx)
	if err != nil {
		return err
	}

	var genericErr error = errors.New(noti.GENERIC_ERROR_WARN_MSG)
	if payment == nil || payment.Method != shared.MANUAL_BANK_METHOD {
		return genericErr
	}

	var curTime time.Time = time.Now()
	if curTime.Before(payment.ExpiredAt) {
		return errors.New(noti.STILL_PENDING_PAYMENT_MESSAGE)
	}

	var sender string = ctx.Value("address").(string)
	var profileId string = ctx.Value("sub").(string)
	if payment.Actor == sender || payment.ProfileID == profileId {
		return errors.New(noti.GENERIC_RIGHT_ACCESS_WARN_MSG)
	}

	var client = p.clients[constant.SuiTestnet]
	manage, err := on_chain.GetOnChainObject[entities.Manage](on_chain.GetOnChainObjectRequest{
		Client:    client,
		ObjectId:  os.Getenv(env.MANAGE_OBJECT_ID),
		ErrLogger: p.errLogger,
	}, ctx)

	if slices.Contains(manage.AdminIds, sender) {
		return errors.New(noti.GENERIC_RIGHT_ACCESS_WARN_MSG)
	}

	payment.UpdatedAt = curTime
	payment.ReviewStatus = request_refused_status
	payment.ReviewedBy = &sender

	return p.paymentRepo.UpdatePayment(*payment, ctx)
}

// CallbackV2 implements business.IPaymentService.
func (p *paymentService) CallbackV2(id string, ctx context.Context) error {
	var payment *entities.Payment
	for i := 1; i <= 5; i++ {
		res, _ := p.paymentRepo.GetPaymentById(id, ctx)
		if res != nil {
			payment = res
			break
		} else {
			var duration time.Duration = time.Duration(i) * 200 * time.Millisecond
			time.Sleep(duration)
		}
	}

	var internalErr error = errors.New(noti.INTERNALL_ERR_MSG)
	if payment == nil {
		return internalErr
	}

	var genericErr error = errors.New(noti.GENERIC_ERROR_WARN_MSG)
	if payment.Method == shared.MANUAL_BANK_METHOD || payment.Status != payment_pending_status || payment.IsTransferred || payment.TransferredAt != nil {
		p.errLogger.Println("Fail check 1")
		return genericErr
	}

	// Expired
	var curTime time.Time = time.Now()
	if payment.ExpiredAt.Before(curTime) {
		return nil
	}

	var data *payos.PaymentLinkDataType
	if payment.IsDonateTx {
		var errRes error
		for i := 1; i <= 5; i++ {
			res, err := payos.GetPaymentLinkInformation(payment.TransactionId)
			if res != nil {
				data = res
				break
			} else {
				errRes = err
				var duration time.Duration = time.Duration(i) * 200 * time.Millisecond
				time.Sleep(duration)
			}
		}

		if errRes != nil {
			p.errLogger.Println("Error while get payment info:", errRes.Error())
			return internalErr
		}
	} else {
		var proposal *entities.WithdrawProposal
		for i := 1; i <= 5; i++ {
			res, _ := on_chain.GetOnChainObject[entities.WithdrawProposal](on_chain.GetOnChainObjectRequest{
				Client:    p.clients[constant.SuiTestnet],
				ObjectId:  *payment.ProposalID,
				ErrLogger: p.errLogger,
			}, ctx)
			if res != nil {
				proposal = res
				break
			} else {
				var duration time.Duration = time.Duration(i) * 200 * time.Millisecond
				time.Sleep(duration)
			}
		}

		if proposal == nil {
			return internalErr
		}

		var bankProfile *entities.BankProfile
		for i := 1; i <= 5; i++ {
			res, _ := p.bankProfileRepo.GetBankProfileByOwner(proposal.Creator, ctx)
			if res != nil {
				bankProfile = res
				break
			} else {
				var duration time.Duration = time.Duration(i) * 200 * time.Millisecond
				time.Sleep(duration)
			}
		}

		if bankProfile == nil {
			return internalErr
		}

		var errSetPayos error
		for i := 1; i <= 5; i++ {
			if err := payos.Key(util.Decrypt(bankProfile.PayosClientID), util.Decrypt(bankProfile.PayosApiKey), util.Decrypt(bankProfile.PayosCheckSumKey)); err != nil {
				errSetPayos = err
				var duration time.Duration = time.Duration(i) * 200 * time.Millisecond
				time.Sleep(duration)
			} else {
				break
			}
		}

		if errSetPayos != nil {
			return internalErr
		}

		var errRes error
		for i := 1; i <= 5; i++ {
			res, err := payos.GetPaymentLinkInformation(payment.TransactionId)
			if res != nil {
				data = res
				break
			} else {
				errRes = err
				var duration time.Duration = time.Duration(i) * 200 * time.Millisecond
				time.Sleep(duration)
			}
		}

		if errRes != nil {
			p.errLogger.Println("Error while get payment info:", errRes.Error())
			return internalErr
		}

		for i := 1; i <= 5; i++ {
			if err := payos.Key(os.Getenv(payment_env.PAYOS_CLIENT_ID), os.Getenv(payment_env.PAYOS_API_KEY), os.Getenv(payment_env.PAYOS_CHECKSUM_KEY)); err != nil {
				errSetPayos = err
				var duration time.Duration = time.Duration(i) * 200 * time.Millisecond
				time.Sleep(duration)
			} else {
				break
			}
		}

		if errSetPayos != nil {
			return internalErr
		}
	}

	switch data.Status {
	case shared.PAYOS_PAID_STATUS:
		payment.ReviewStatus = request_approved_status
		payment.Status = shared.PAYMENT_SUCCESS_STATUS
		payment.IsTransferred = true
		payment.TransferredAt = &curTime
	case shared.PAYOS_CANCELLED_STATUS:
		payment.Status = shared.PAYMENT_CANCELED_STATUS
		payment.ReviewStatus = request_refused_status
		if data.CancellationReason != nil {
			payment.CancelReason = data.CancellationReason
		}
	default:
		p.errLogger.Println("Fail check 2")
		return genericErr
	}

	var system string = "System"
	payment.ReviewedBy = &system
	payment.UpdatedAt = curTime

	if err := p.paymentRepo.UpdatePayment(*payment, ctx); err != nil {
		return err
	}

	if payment.ReviewStatus == request_refused_status {
		return nil
	}

	var profile *entities.Profile
	for i := 1; i <= 5; i++ {
		res, _ := p.profileRepo.GetProfile(payment.ProfileID, ctx)
		if res != nil {
			profile = res
			break
		}

		var duration time.Duration = time.Duration(i) * 200 * time.Millisecond
		time.Sleep(duration)
	}

	if profile == nil {
		return internalErr
	}

	var client = p.clients[constant.SuiTestnet]
	var module string
	var function string
	var args []interface{}
	if payment.IsDonateTx {
		var detail *entities.OffChainDonation
		for i := 1; i <= 5; i++ {
			res, _ := p.donationRepo.GetDonation(*payment.DonationID, ctx)
			if res != nil {
				detail = res
				break
			}

			var duration time.Duration = time.Duration(i) * 200 * time.Millisecond
			time.Sleep(duration)
		}

		if detail == nil {
			return internalErr
		}

		manage, err := on_chain.GetOnChainObject[entities.Manage](on_chain.GetOnChainObjectRequest{
			Client:    client,
			ObjectId:  os.Getenv(env.MANAGE_OBJECT_ID),
			ErrLogger: p.errLogger,
		}, ctx)
		if err != nil {
			return err
		}

		if manage == nil {
			return errors.New(noti.INTERNALL_ERR_MSG)
		}

		var nftId string
		for i, donor := range manage.DonorIds {
			if donor == payment.Actor {
				nftId = manage.DonorNfts[i]
				break
			}
		}

		if nftId == "" {
			nftId = os.Getenv(env.PUB_DONOR_NFT_ID)
		}

		if detail.Purpose == string(entities.DONATE_PURPOSE) {
			var poolModule = on_chain.InitializeModulePool()
			module = poolModule.GetModule()
			if detail.Target != os.Getenv(env.POOL_ID) { // Donate to local pool
				function = poolModule.GetFunctionDonateToLocalPool()
				args = poolModule.ToDonateToLocalPoolArguments(on_chain.DonateToLocalPoolArguments{
					LocalPoolId: detail.Target,
					DonateToPoolArguments: on_chain.DonateToPoolArguments{
						DonorID:     nftId,
						Amount:      payment.Amount,
						FirstName:   *profile.FirstName,
						LastName:    *profile.LastName,
						Gender:      *profile.Gender,
						PhoneNumber: *profile.PhoneNumber,
						Email:       *profile.Email,
						Message:     payment.Message,
						Sender:      payment.Actor,
					},
				})
			} else {
				function = poolModule.GetFunctionDonateToPool()
				args = poolModule.ToDonateToPoolArguments(on_chain.DonateToPoolArguments{
					DonorID:     nftId,
					Amount:      payment.Amount,
					FirstName:   *profile.FirstName,
					LastName:    *profile.LastName,
					Gender:      *profile.Gender,
					PhoneNumber: *profile.PhoneNumber,
					Email:       *profile.Email,
					Message:     payment.Message,
					Sender:      payment.Actor,
				})
			}
		} else if detail.Purpose == string(entities.CAMPAIGN_PURPOSE) {
			var localPoolId string
			for i := 1; i <= 3; i++ {
				if campaign, err := on_chain.GetOnChainObject[entities.OnChainCampaign](on_chain.GetOnChainObjectRequest{
					Client:    client,
					ObjectId:  detail.Target,
					ErrLogger: p.errLogger,
				}, ctx); err == nil {
					if campaign.PoolID == os.Getenv(env.POOL_ID) {
						localPoolId = os.Getenv(env.SHARED_LOCAL_POOL_ID)
					} else {
						localPoolId = campaign.PoolID
					}
					break
				}
			}

			var campaignModule = on_chain.InitializeModuleCampaign()
			module = campaignModule.GetModule()
			function = campaignModule.GetFunctionSupportCampaign()

			args = campaignModule.ToSupportCampaignArguments(on_chain.SupportCampaignArguments{
				LocalPoolID: localPoolId,
				CampaignID:  detail.Target,
				DonorNFT:    nftId,
				Amount:      payment.Amount,
				FirstName:   *profile.FirstName,
				LastName:    *profile.LastName,
				Gender:      *profile.Gender,
				PhoneNumber: *profile.PhoneNumber,
				Email:       *profile.Email,
				Message:     payment.Message,
				Sender:      payment.Actor,
			})
		} else {
			var childModule = on_chain.InitializeModuleChild()
			module = childModule.GetModule()

			var childId, generalContentFormat, taskContentFormat, mealSupportStartPeriod, mealSupportEndPeriod string
			var expectedWithdrawPeriods, contentFormats []string
			var task entities.Task

			switch detail.Purpose {
			case string(entities.BOOKS_NEED_PURPOSE):
				need, err := on_chain.GetOnChainObject[entities.BooksNeed](on_chain.GetOnChainObjectRequest{
					Client:    client,
					ObjectId:  detail.Target,
					ErrLogger: p.errLogger,
				}, ctx)
				if err != nil {
					return nil
				}

				withdrawDates, err := on_chain.GetOnChainObject[entities.BooksNeedWithdrawDates](on_chain.GetOnChainObjectRequest{
					Client:    client,
					ObjectId:  os.Getenv(env.BOOKS_NEED_WITHDRAW_DATES_ID),
					ErrLogger: p.errLogger,
				}, ctx)

				var withdrawDate string
				if need.Semster == "1" {
					withdrawDate = withdrawDates.FirstSemesterDate
				} else {
					withdrawDate = withdrawDates.SecondSemesterDate
				}

				withdrawDate += fmt.Sprintf("/%d", curTime.Year())
				expectedWithdrawPeriods = append(expectedWithdrawPeriods, withdrawDate)
				contentFormats = append(contentFormats, "Rút tiền mua sách cho học kì - "+need.Semster+" - "+need.Year+" for child %s %s - %s")
				generalContentFormat = "Rút tiền mua sách cho học kì - " + need.Semster + " for child %s %s - %s"
				childId = need.ChildID
				function = childModule.GetFunctionSupportChildBooksNeed()
			case string(entities.HEALTH_INSURANCE_NEED_PURPOSE):
				need, err := on_chain.GetOnChainObject[entities.HealthInsuranceNeed](on_chain.GetOnChainObjectRequest{
					Client:    client,
					ObjectId:  detail.Target,
					ErrLogger: p.errLogger,
				}, ctx)
				if err != nil {
					return nil
				}

				withdrawDate, err := on_chain.GetOnChainObject[entities.HealthInsuranceNeedWithdrawDate](on_chain.GetOnChainObjectRequest{
					Client:    client,
					ObjectId:  os.Getenv(env.HEALTH_INSURANCE_NEED_WITHDRAW_DATE_ID),
					ErrLogger: p.errLogger,
				}, ctx)

				expectedWithdrawPeriods = append(expectedWithdrawPeriods, withdrawDate.ExpectedDate+"/"+need.Year)
				contentFormats = append(contentFormats, "Rút tiền mua bảo hiểm sức khoẻ - "+need.Year+" for child %s %s - %s")
				generalContentFormat = "Rút tiền mua bảo hiểm sức khoẻ cho trẻ - %s %s - %s"
				childId = need.ChildID
				function = childModule.GetFunctionSupportChildHealthInsuranceNeed()
			case string(entities.MEAL_NEED_PURPOSE):
				need, err := on_chain.GetOnChainObject[entities.MealNeed](on_chain.GetOnChainObjectRequest{
					Client:    client,
					ObjectId:  detail.Target,
					ErrLogger: p.errLogger,
				}, ctx)
				if err != nil {
					return nil
				}

				var duration *entities.OffChainMealSupportDuration
				for i := 1; i <= 5; i++ {
					res, _ := p.mealSupportDurationRepo.GetMealSupportDuration(*detail.MealDurationID, ctx)
					if res != nil {
						duration = res
						break
					}

					var duration time.Duration = time.Duration(i) * 200 * time.Millisecond
					time.Sleep(duration)
				}

				if duration == nil {
					return internalErr
				}

				var startPeriod time.Time = util.ToStartOfDate(util.RawDateToTime(duration.StartPeriod))
				var endPeriod time.Time = util.ToEndOfDate(util.RawDateToTime(duration.EndPeriod))
				var endMonth int = int(endPeriod.Month())
				if endMonth == 1 {
					endMonth = 13
				}

				var initExpectedDate time.Time = startPeriod.AddDate(0, 0, -3)
				var contentFormat string = "Withdraw meal need %s - %s for child %s %s - %s"
				for i := 0; i < endMonth-int(startPeriod.Month()); i++ {
					var date time.Time = initExpectedDate.AddDate(0, i, 0)
					expectedWithdrawPeriods = append(expectedWithdrawPeriods, util.TimeToRawDate(date))

					var startProvidePeriod time.Time = startPeriod.AddDate(0, i, 0)
					var endProvidePeriod time.Time = startProvidePeriod.AddDate(0, 1, 0)

					// var appendChildInfoContentFormat string = fmt.Sprintf(contentFormat, util.TimeToRawDate(startProvidePeriod), util.TimeToRawDate(endProvidePeriod))
					// appendChildInfoContentFormat += "%s %s - %s"
					var format string = fmt.Sprintf(contentFormat, util.TimeToRawDate(startProvidePeriod), util.TimeToRawDate(endProvidePeriod), "%s", "%s", "%s")
					contentFormats = append(contentFormats, format)
				}

				task = entities.Task{
					ID:          util.GenerateId(),
					CreatedBy:   "System",
					StartPeriod: startPeriod,
					EndPeriod:   endPeriod,
				}

				taskContentFormat = "Cho trẻ ăn -  %s %s - %s từ " + duration.StartPeriod + " đến " + duration.EndPeriod
				mealSupportStartPeriod = duration.StartPeriod
				mealSupportEndPeriod = duration.EndPeriod
				generalContentFormat = "Rút tiền ăn cho trẻ - %s %s - %s"
				childId = need.ChildID
				function = childModule.GetFunctionSupportChildMealNeed()
			case string(entities.SPECIAL_NEED_PURPOSE):
				campaign, err := on_chain.GetOnChainObject[entities.SpecialNeedCampaign](on_chain.GetOnChainObjectRequest{
					Client:    client,
					ObjectId:  detail.Target,
					ErrLogger: p.errLogger,
				}, ctx)
				if err != nil {
					return nil
				}

				childId = campaign.ChildID
				function = childModule.GetFunctionSupportChildSpecialNeedCampaign()
			}

			child, err := on_chain.GetOnChainObject[entities.Child](on_chain.GetOnChainObjectRequest{
				Client:    client,
				ObjectId:  childId,
				ErrLogger: p.errLogger,
			}, ctx)
			if err != nil {
				return nil
			}

			pool, err := on_chain.GetOnChainObject[entities.MainPool](on_chain.GetOnChainObjectRequest{
				Client:    client,
				ObjectId:  os.Getenv(env.POOL_ID),
				ErrLogger: p.errLogger,
			}, ctx)
			if err != nil {
				return nil
			}

			localPools, err := on_chain.GetOnChainObjects[entities.LocalPool](on_chain.GetOnChainObjectsRequest{
				Client:    client,
				ObjectIds: pool.LocalPools,
				ErrLogger: p.errLogger,
			}, ctx)
			if err != nil {
				return nil
			}

			var leaders []string
			var localPoolId string
			for _, localPool := range localPools {
				if localPool.Region == child.Region {
					leaders = localPool.Mods
					localPoolId = localPool.ID.ID
					break
				}
			}

			if detail.Purpose != string(entities.SPECIAL_NEED_PURPOSE) {
				leaderNoti, err := p.leaderNotiRepo.GetNotiByNeed(detail.Target, ctx)
				if err != nil {
					return nil
				}

				var childFormattedAddr string = util.FormatAddress(childId)
				var contents []string
				for _, format := range contentFormats {
					contents = append(contents, fmt.Sprintf(format, child.LastName, child.FirstName, childFormattedAddr))
				}

				var curTime time.Time = time.Now()
				if leaderNoti == nil {
					p.leaderNotiRepo.CreateNoti(entities.LeaderNoti{
						ID:                      util.GenerateId(),
						NeedID:                  detail.Target,
						NeedType:                detail.Purpose,
						ChildID:                 childId,
						Region:                  child.Region,
						AssignedLeaders:         leaders,
						ExpectedWithdrawPeriods: expectedWithdrawPeriods,
						Contents:                contents,
						GeneralContent:          fmt.Sprintf(generalContentFormat, child.LastName, child.FirstName, childFormattedAddr),
						CreatedAt:               curTime,
						UpdatedAt:               curTime,
					}, ctx)
				} else {
					leaderNoti.ExpectedWithdrawPeriods = append(leaderNoti.ExpectedWithdrawPeriods, expectedWithdrawPeriods...)
					leaderNoti.Contents = append(leaderNoti.Contents, contents...)
					p.leaderNotiRepo.UpdateNoti(*leaderNoti, ctx)
				}

				switch detail.Purpose {
				case string(entities.BOOKS_NEED_PURPOSE):
					args = childModule.ToSupportChildBooksNeedArguments(on_chain.SupportChildBooksNeedArguments{
						NeedID:      detail.Target,
						LocalPool:   localPoolId,
						ChildID:     childId,
						DonorNft:    nftId,
						Amount:      payment.Amount,
						FirstName:   *profile.FirstName,
						LastName:    *profile.LastName,
						Gender:      *profile.Gender,
						PhoneNumber: *profile.PhoneNumber,
						Email:       *profile.Email,
						Message:     payment.Message,
						Sender:      payment.Actor,
					})
				case string(entities.MEAL_NEED_PURPOSE):
					var detailId string = util.GenerateId()
					p.childTaskDetailRepo.CreateChildTaskDetail(entities.ChildTaskDetail{
						ID:      detailId,
						ChildID: childId,
						Target:  detail.Target,
					}, ctx)

					task.IsChildTask = true
					task.ChildTaskDetailID = &detailId
					task.Region = child.Region
					task.Description = fmt.Sprintf(taskContentFormat, child.LastName, child.FirstName, childFormattedAddr)
					task.CreatedAt = curTime
					task.UpdatedAt = curTime
					p.taskRepo.CreateTask(task, ctx)

					args = childModule.ToSupportChildMealNeedArguments(on_chain.SupportChildMealNeedArguments{
						StartPeriod: mealSupportStartPeriod,
						EndPeriod:   mealSupportEndPeriod,
						SupportChildBooksNeedArguments: on_chain.SupportChildBooksNeedArguments{
							NeedID:      detail.Target,
							LocalPool:   localPoolId,
							ChildID:     childId,
							DonorNft:    nftId,
							Amount:      payment.Amount,
							FirstName:   *profile.FirstName,
							LastName:    *profile.LastName,
							Gender:      *profile.Gender,
							PhoneNumber: *profile.PhoneNumber,
							Email:       *profile.Email,
							Message:     payment.Message,
							Sender:      payment.Actor,
						},
					})
				case string(entities.HEALTH_INSURANCE_NEED_PURPOSE):
					args = childModule.ToSupportChildHealthInsuranceNeedArguments(on_chain.SupportChildHealthInsuranceNeedArguments{
						NeedID:      detail.Target,
						LocalPool:   localPoolId,
						ChildID:     childId,
						DonorNft:    nftId,
						Amount:      payment.Amount,
						FirstName:   *profile.FirstName,
						LastName:    *profile.LastName,
						Gender:      *profile.Gender,
						PhoneNumber: *profile.PhoneNumber,
						Email:       *profile.Email,
						Message:     payment.Message,
						Sender:      payment.Actor,
					})
				}
			} else {
				args = childModule.ToSupportChildSpeicalNeedArguments(on_chain.SupportChildSpeicalNeedArguments{
					CampaignID:  detail.Target,
					LocalPool:   localPoolId,
					ChildID:     childId,
					DonorNft:    nftId,
					Amount:      payment.Amount,
					FirstName:   *profile.FirstName,
					LastName:    *profile.LastName,
					Gender:      *profile.Gender,
					PhoneNumber: *profile.PhoneNumber,
					Email:       *profile.Email,
					Message:     payment.Message,
					Sender:      payment.Actor,
				})
			}
		}
	} else {
		var proposal *entities.WithdrawProposal
		for i := 1; i <= 3; i++ {
			if res, _ := on_chain.GetOnChainObject[entities.WithdrawProposal](on_chain.GetOnChainObjectRequest{
				Client:    client,
				ObjectId:  *payment.ProposalID,
				ErrLogger: p.errLogger,
			}, ctx); res != nil {
				proposal = res
				break
			}
		}

		if proposal.Purpose == string(entities.POOL_WITHDRAW_PROPOSAL_PURPOSE) {
			var poolModule = on_chain.InitializeModulePool()
			module = poolModule.GetModule()
			function = poolModule.GetFunctionWithdrawFromPool()

			var localPoolId string
			if proposal.TargetID == os.Getenv(env.POOL_ID) {
				localPoolId = os.Getenv(env.SHARED_LOCAL_POOL_ID)
			} else {
				localPoolId = proposal.TargetID
			}

			args = poolModule.ToWithdrawFromPoolArguments(on_chain.WithdrawFromPoolArguments{
				LocalPoolId:        localPoolId,
				WithdrawProposalId: *payment.ProposalID,
				Sender:             payment.Actor,
			})

			var profile *entities.Profile
			for i := 1; i <= 5; i++ {
				res, _ := p.profileRepo.GetProfileByWalletAddress(proposal.Creator, ctx)
				if res != nil {
					profile = res
					break
				}

				var duration time.Duration = time.Duration(i) * 200 * time.Millisecond
				time.Sleep(duration)
			}

			if profile == nil {
				return internalErr
			}

			if err := p.taskRepo.CreateTask(entities.Task{
				ID:                util.GenerateId(),
				CreatedBy:         "System",
				AssignedProfileID: &profile.ID,
				AssignedStaff:     &proposal.Creator,
				Region:            proposal.PoolName,
				StartPeriod:       util.ToStartOfDate(curTime),
				EndPeriod:         util.ToEndOfDate(curTime),
				CreatedAt:         curTime,
				UpdatedAt:         curTime,
			}, ctx); err != nil {
				return err
			}
		} else if proposal.Purpose == string(entities.POOL_CAMPAIGN_WITHDRAW_PROPOSAL_PURPOSE) {
			var campaignModule = on_chain.InitializeModuleCampaign()
			module = campaignModule.GetModule()
			function = campaignModule.GetFunctionWithdrawFromCampaign()
			args = campaignModule.ToWithdrawFromCampaignArguments(on_chain.WithdrawFromCampaignArguments{
				CampaignID: proposal.TargetID,
				ProposalID: *payment.ProposalID,
				Sender:     payment.Actor,
			})
		} else {
			var childModule = on_chain.InitializeModuleChild()
			module = childModule.GetModule()
			switch proposal.Purpose {
			case string(entities.BOOKS_NEED_WITHDRAW_PROPOSAL_PURPOSE):
				function = childModule.GetFunctionWithdrawFromBooksNeedProposal()
			case string(entities.HEALTH_INSURANCE_NEED_WITHDRAW_PROPOSAL_PURPOSE):
				function = childModule.GetFunctionWithdrawFromHealthInsuranceNeedProposal()
			case string(entities.MEAL_NEED_WITHDRAW_PROPOSAL_PURPOSE):
				function = childModule.GetFunctionWithdrawFromMealNeedProposal()
			case string(entities.SPECIAL_NEED_CAMPAIGN_WITHDRAW_PROPOSAL_PURPOSE):
				function = childModule.GetFunctionWithdrawFromSpecialNeedCampaign()
			}

			args = childModule.ToWithdrawFromNeedArguments(on_chain.WithdrawFromNeedArguments{
				LocalPool:  proposal.PoolID,
				TargetID:   proposal.TargetID,
				ProposalID: *payment.ProposalID,
				Sender:     payment.Actor,
			})
		}
	}

	for i := 1; i <= 3; i++ {
		if _, err := on_chain.ExecuteTransactionV2(on_chain.ExecuteTransactionRequestV2{
			Client:    client,
			Module:    module,
			Function:  function,
			Arguments: args,
			ErrLogger: p.errLogger,
		}, ctx); err == nil {
			return nil
		}
	}

	return errors.New(noti.INTERNALL_ERR_MSG)
}

// CallbackWithAuthV2 implements business.IPaymentService.
func (p *paymentService) CallbackWithAuthV2(id string, req request.PaymentAuthCallbackRequest, ctx context.Context) error {
	payment, err := p.paymentRepo.GetPaymentById(id, ctx)
	if err != nil {
		return err
	}

	var genericErr error = errors.New(noti.GENERIC_ERROR_WARN_MSG)
	if payment == nil || payment.Method != shared.MANUAL_BANK_METHOD || payment.IsDonateTx {
		return genericErr
	}

	var sender string = ctx.Value("address").(string)
	var profileId string = ctx.Value("sub").(string)
	if payment.Actor != sender || payment.ProfileID != profileId {
		return errors.New(noti.GENERIC_RIGHT_ACCESS_WARN_MSG)
	}

	// Expired
	if payment.ExpiredAt.Before(time.Now()) {
		return nil
	}

	var curTime time.Time = time.Now()
	payment.UpdatedAt = curTime
	payment.IsTransferred = true
	payment.TransferredAt = &curTime
	payment.ProofBlobID = &req.ProofBlobID

	return p.paymentRepo.UpdatePayment(*payment, ctx)
}

// DonateV2 implements business.IPaymentService.
func (p *paymentService) DonateV2(req request.DonateRequest, ctx context.Context) (response.PaymentUrlResponse, error) {
	var genericErr error = errors.New(noti.GENERIC_ERROR_WARN_MSG)
	if !util.IsValidSuiAddressStrict(req.PoolId) {
		return response.PaymentUrlResponse{}, genericErr
	}

	profile, err := p.profileRepo.GetProfile(ctx.Value("sub").(string), ctx)
	if err != nil {
		return response.PaymentUrlResponse{}, err
	}

	if profile == nil {
		return response.PaymentUrlResponse{}, genericErr
	}

	if profile.IdentityCode == nil {
		return response.PaymentUrlResponse{}, errors.New(noti.NOT_UPLOADED_PROFILE_MESSAGE)
	}

	var paymentId string = util.GenerateId()
	var orderCode int = util.GenerateNumber()
	var callbackUrl string = os.Getenv(payment.PAYMENT_CALLBACK_URL) + paymentId
	var description string = req.Message
	if description == "" {
		description = fmt.Sprint(orderCode)
	}

	data, err := payos.CreatePaymentLink(payos.CheckoutRequestType{
		OrderCode:   int64(orderCode),
		Amount:      int(req.Amount),
		Description: description,
		ReturnUrl:   callbackUrl,
		CancelUrl:   callbackUrl,
	})

	if err != nil {
		p.errLogger.Println("Err: ", err.Error())
		return response.PaymentUrlResponse{}, errors.New(noti.INTERNALL_ERR_MSG)
	}

	var donationId string = util.GenerateId()
	var curTime time.Time = time.Now()
	if err := p.donationRepo.CreateDonation(entities.OffChainDonation{
		ID:        donationId,
		Purpose:   string(entities.DONATE_PURPOSE),
		Target:    req.PoolId,
		CreatedAt: curTime,
	}, ctx); err != nil {
		return response.PaymentUrlResponse{}, err
	}

	var expiredAt time.Time
	if data.ExpiredAt != nil {
		expiredAt = time.Unix(int64(*data.ExpiredAt), 0)
	} else {
		expiredAt = time.Now().Add(2 * time.Minute) // Default 15p nếu PayOS ko trả về
	}

	return response.PaymentUrlResponse{
			Url:       data.CheckoutUrl,
			PaymentID: paymentId,
		}, p.paymentRepo.CreatePayment(entities.Payment{
			ID:            paymentId,
			Actor:         ctx.Value("address").(string),
			ProfileID:     profile.ID,
			DonationID:    &donationId,
			IsDonateTx:    true,
			TransactionId: fmt.Sprint(orderCode),
			Amount:        req.Amount,
			Currency:      shared.VIETNAMDONG_CURRENCY,
			Status:        payment_pending_status,
			Method:        shared.PAYMENT_PAYOS_METHOD,
			Message:       description,
			ExpiredAt:     expiredAt,
			CreatedAt:     curTime,
			UpdatedAt:     curTime,
		}, ctx)
}

// ConfirmWithdrawProposal implements business.IPaymentService.
func (p *paymentService) ConfirmWithdrawProposal(id string, ctx context.Context) (map[string]interface{}, error) {
	var genericErr error = errors.New(noti.GENERIC_ERROR_WARN_MSG)

	var sender string = ctx.Value("address").(string)
	if !utils.IsValidSuiAddress(models.SuiAddress(sender)) || !utils.IsValidSuiAddress(models.SuiAddress(id)) {
		return nil, genericErr
	}

	return nil, nil
}

// Donate implements business.IPaymentService.
func (p *paymentService) Donate(req request.DonateRequest, ctx context.Context) (response.UrlAPIResponse, error) {
	var genericErr error = errors.New(noti.GENERIC_ERROR_WARN_MSG)
	if !util.IsValidSuiAddressStrict(req.PoolId) {
		return response.UrlAPIResponse{}, genericErr
	}
	profile, err := p.profileRepo.GetProfile(ctx.Value("sub").(string), ctx)
	if err != nil {
		return response.UrlAPIResponse{}, err
	}

	if profile == nil {
		return response.UrlAPIResponse{}, genericErr
	}

	if profile.IdentityCode == nil {
		return response.UrlAPIResponse{}, errors.New(noti.NOT_UPLOADED_PROFILE_MESSAGE)
	}

	var paymentId string = util.GenerateId()
	var orderCode int = util.GenerateNumber()
	var callbackUrl string = os.Getenv(payment.PAYMENT_CALLBACK_URL) + paymentId
	var description string = req.Message
	if description == "" {
		description = fmt.Sprint(orderCode)
	}

	data, err := payos.CreatePaymentLink(payos.CheckoutRequestType{
		OrderCode:   int64(orderCode),
		Amount:      int(req.Amount),
		Description: description,
		ReturnUrl:   callbackUrl,
		CancelUrl:   callbackUrl,
	})

	if err != nil {
		p.errLogger.Println("Err: ", err.Error())
		return response.UrlAPIResponse{}, errors.New(noti.INTERNALL_ERR_MSG)
	}

	var donationId string = util.GenerateId()
	var curTime time.Time = time.Now()
	if err := p.donationRepo.CreateDonation(entities.OffChainDonation{
		ID:        donationId,
		Purpose:   string(entities.DONATE_PURPOSE),
		Target:    req.PoolId,
		CreatedAt: curTime,
	}, ctx); err != nil {
		return response.UrlAPIResponse{}, err
	}

	var expiredAt time.Time
	if data.ExpiredAt != nil {
		expiredAt = time.Unix(int64(*data.ExpiredAt), 0)
	} else {
		expiredAt = time.Now().Add(1 * time.Minute) // Default 15p nếu PayOS ko trả về
	}

	return response.UrlAPIResponse{
			Url: data.CheckoutUrl,
		}, p.paymentRepo.CreatePayment(entities.Payment{
			ID:            paymentId,
			Actor:         ctx.Value("address").(string),
			ProfileID:     profile.ID,
			DonationID:    &donationId,
			IsDonateTx:    true,
			TransactionId: fmt.Sprint(orderCode),
			Amount:        req.Amount,
			Currency:      shared.VIETNAMDONG_CURRENCY,
			Status:        payment_pending_status,
			Method:        shared.PAYMENT_PAYOS_METHOD,
			Message:       description,
			ExpiredAt:     expiredAt,
			CreatedAt:     curTime,
			UpdatedAt:     curTime,
		}, ctx)
}

// Callback implements business.IPaymentService.
func (p *paymentService) Callback(id string, ctx context.Context) (string, error) {
	payment, err := p.paymentRepo.GetPaymentById(id, ctx)
	if err != nil {
		return "", err
	}

	var genericErr error = errors.New(noti.GENERIC_ERROR_WARN_MSG)
	if payment == nil {
		return "", genericErr
	}

	// Expired
	if payment.ExpiredAt.Before(time.Now()) {
		return util.GeneratePaymentRedirectUrl(util.GenerateRedirectParamRequest{
			OrderCode: payment.TransactionId,
			Status:    shared.PAYAMENT_EXPIRED_STATUS,
			Message:   noti.PAYMENT_EXPIRED_MESSAGE,
		}), nil
	}

	data, err := payos.GetPaymentLinkInformation(payment.TransactionId)
	if err != nil {
		p.errLogger.Println("Error while get payos payment link information: " + err.Error())
		return "", errors.New(noti.INTERNALL_ERR_MSG)
	}

	switch data.Status {
	case shared.PAYOS_PAID_STATUS:
		payment.Status = shared.PAYMENT_SUCCESS_STATUS
	case shared.PAYOS_CANCELLED_STATUS:
		payment.Status = shared.PAYMENT_CANCELED_STATUS
		if data.CancellationReason != nil {
			payment.CancelReason = data.CancellationReason
		}
	default:
		return "", genericErr
	}

	payment.UpdatedAt = time.Now()
	if err := p.paymentRepo.UpdatePayment(*payment, ctx); err != nil {
		return "", err
	}

	if data.Status == shared.PAYMENT_CANCELED_STATUS {
		return util.GeneratePaymentRedirectUrl(util.GenerateRedirectParamRequest{
			OrderCode: payment.TransactionId,
			Status:    data.Status,
			Message:   noti.PAYMENT_CANCEL_MESSAGE,
		}), nil
	}

	profile, err := p.profileRepo.GetProfile(payment.ProfileID, ctx)
	if err != nil {
		return "", err
	}

	if profile == nil {
		return "", genericErr
	}

	var client = p.clients[constant.SuiTestnet]
	var module string
	var function string
	var args []interface{}
	if payment.IsDonateTx {
		detail, err := p.donationRepo.GetDonation(*payment.DonationID, ctx)
		if err != nil {
			return util.GeneratePaymentRedirectUrl(util.GenerateRedirectParamRequest{
				OrderCode: payment.TransactionId,
				Status:    data.Status,
				Message:   err.Error(),
			}), nil
		}

		var donorModule = on_chain.InitializeModuleDonor()
		var nftId string
		nfts, err := on_chain.GetOnChainOwnedObjects[entities.Donor](on_chain.GetOnChainOwnedObjectsRequest{
			Client:       client,
			OwnerAddress: payment.Actor,
			StructType:   fmt.Sprintf("%s::%s::%s", os.Getenv(env.PACKAGE_ID), donorModule.GetModule(), donorModule.GetDonorNftStruct()),
		}, ctx)
		if err != nil {
			return "", err
		}

		if len(nfts) > 0 {
			nftId = nfts[0].ID.ID
		} else {
			nftId = os.Getenv(env.PUB_DONOR_NFT_ID)
		}

		if detail.Purpose != string(entities.DONATE_PURPOSE) {
			if detail.Purpose == string(entities.CAMPAIGN_PURPOSE) {
				var localPoolId string
				for i := 1; i <= 3; i++ {
					if campaign, err := on_chain.GetOnChainObject[entities.OnChainCampaign](on_chain.GetOnChainObjectRequest{
						Client:    client,
						ObjectId:  detail.Target,
						ErrLogger: p.errLogger,
					}, ctx); err == nil {
						if campaign.PoolID == os.Getenv(env.POOL_ID) {
							localPoolId = os.Getenv(env.SHARED_LOCAL_POOL_ID)
						} else {
							localPoolId = campaign.PoolID
						}
						break
					}
				}

				var campaignModule = on_chain.InitializeModuleCampaign()
				module = campaignModule.GetModule()
				function = campaignModule.GetFunctionSupportCampaign()

				args = campaignModule.ToSupportCampaignArguments(on_chain.SupportCampaignArguments{
					LocalPoolID: localPoolId,
					CampaignID:  detail.Target,
					DonorNFT:    nftId,
					Amount:      payment.Amount,
					FirstName:   *profile.FirstName,
					LastName:    *profile.LastName,
					Gender:      *profile.Gender,
					PhoneNumber: *profile.PhoneNumber,
					Email:       *profile.Email,
					Message:     payment.Message,
				})
			} else {
				var childModule = on_chain.InitializeModuleChild()
				module = childModule.GetModule()
				var childId, generalContentFormat, taskContentFormat, mealSupportStartPeriod, mealSupportEndPeriod string
				var expectedWithdrawPeriods, contentFormats []string
				var task entities.Task
				switch detail.Purpose {
				case string(entities.BOOKS_NEED_PURPOSE):
					need, err := on_chain.GetOnChainObject[entities.BooksNeed](on_chain.GetOnChainObjectRequest{
						Client:    client,
						ObjectId:  detail.Target,
						ErrLogger: p.errLogger,
					}, ctx)
					if err != nil {
						return util.GeneratePaymentRedirectUrl(util.GenerateRedirectParamRequest{
							OrderCode: payment.TransactionId,
							Status:    data.Status,
							Message:   err.Error(),
						}), nil
					}

					withdrawDates, err := on_chain.GetOnChainObject[entities.BooksNeedWithdrawDates](on_chain.GetOnChainObjectRequest{
						Client:    client,
						ObjectId:  os.Getenv(env.BOOKS_NEED_WITHDRAW_DATES_ID),
						ErrLogger: p.errLogger,
					}, ctx)

					var withdrawDate string
					if need.Semster == "1" {
						withdrawDate = withdrawDates.FirstSemesterDate
					} else {
						withdrawDate = withdrawDates.SecondSemesterDate
					}

					expectedWithdrawPeriods = append(expectedWithdrawPeriods, withdrawDate)
					contentFormats = append(contentFormats, "Rút tiền mua sách cho học kì - "+need.Semster+" - "+need.Year+" cho trẻ %s %s - %s")
					generalContentFormat = "Rút tiền mua sách cho học kì - " + need.Semster + " cho trẻ %s %s - %s"
					childId = need.ChildID
					function = childModule.GetFunctionSupportChildBooksNeed()
				case string(entities.HEALTH_INSURANCE_NEED_PURPOSE):
					need, err := on_chain.GetOnChainObject[entities.HealthInsuranceNeed](on_chain.GetOnChainObjectRequest{
						Client:    client,
						ObjectId:  id,
						ErrLogger: p.errLogger,
					}, ctx)
					if err != nil {
						return util.GeneratePaymentRedirectUrl(util.GenerateRedirectParamRequest{
							OrderCode: payment.TransactionId,
							Status:    data.Status,
							Message:   err.Error(),
						}), nil
					}

					withdrawDate, err := on_chain.GetOnChainObject[entities.HealthInsuranceNeedWithdrawDate](on_chain.GetOnChainObjectRequest{
						Client:    client,
						ObjectId:  os.Getenv(env.HEALTH_INSURANCE_NEED_WITHDRAW_DATE_ID),
						ErrLogger: p.errLogger,
					}, ctx)

					expectedWithdrawPeriods = append(expectedWithdrawPeriods, withdrawDate.ExpectedDate+"/"+need.Year)
					contentFormats = append(contentFormats, "Rút tiền mua bảo hiểm sức khoẻ - "+need.Year+" cho trẻ %s %s - %s")
					generalContentFormat = "Rút tiền mua bảo hiểm sức khoẻ cho trẻ - %s %s - %s"
					childId = need.ChildID
					function = childModule.GetFunctionSupportChildHealthInsuranceNeed()
				case string(entities.MEAL_NEED_PURPOSE):
					need, err := on_chain.GetOnChainObject[entities.MealNeed](on_chain.GetOnChainObjectRequest{
						Client:    client,
						ObjectId:  detail.Target,
						ErrLogger: p.errLogger,
					}, ctx)
					if err != nil {
						return util.GeneratePaymentRedirectUrl(util.GenerateRedirectParamRequest{
							OrderCode: payment.TransactionId,
							Status:    data.Status,
							Message:   err.Error(),
						}), nil
					}

					duration, err := p.mealSupportDurationRepo.GetMealSupportDuration(*detail.MealDurationID, ctx)

					var startPeriod time.Time = util.ToStartOfDate(util.RawDateToTime(duration.StartPeriod))
					var endPeriod time.Time = util.ToEndOfDate(util.RawDateToTime(duration.EndPeriod))
					var endMonth int = int(endPeriod.Month())
					if endMonth == 1 {
						endMonth = 13
					}

					var initExpectedDate time.Time = startPeriod.AddDate(0, 0, -3)
					var contentFormat string = "Rút tiền ăn %s - %s cho trẻ %s %s - %s"
					for i := 0; i < endMonth-int(startPeriod.Month()); i++ {
						var date time.Time = initExpectedDate.AddDate(0, i, 0)
						expectedWithdrawPeriods = append(expectedWithdrawPeriods, util.TimeToRawDate(date))

						var startProvidePeriod time.Time = startPeriod.AddDate(0, i, 0)
						var endProvidePeriod time.Time = startProvidePeriod.AddDate(0, 1, 0)
						var format string = fmt.Sprintf(contentFormat, util.TimeToRawDate(startProvidePeriod), util.TimeToRawDate(endProvidePeriod), "%s", "%s", "%s")
						contentFormats = append(contentFormats, format)
					}

					task = entities.Task{
						ID:          util.GenerateId(),
						CreatedBy:   "System",
						StartPeriod: startPeriod,
						EndPeriod:   endPeriod,
					}

					taskContentFormat = "Cho trẻ ăn - %s %s - %s từ " + duration.StartPeriod + " đến " + duration.EndPeriod
					mealSupportStartPeriod = duration.StartPeriod
					mealSupportEndPeriod = duration.EndPeriod
					generalContentFormat = "Rút tiền ăn cho trẻ - %s %s - %s"
					childId = need.ChildID
					function = childModule.GetFunctionSupportChildMealNeed()
				case string(entities.SPECIAL_NEED_PURPOSE):
					campaign, err := on_chain.GetOnChainObject[entities.SpecialNeedCampaign](on_chain.GetOnChainObjectRequest{
						Client:    client,
						ObjectId:  detail.Target,
						ErrLogger: p.errLogger,
					}, ctx)
					if err != nil {
						return util.GeneratePaymentRedirectUrl(util.GenerateRedirectParamRequest{
							OrderCode: payment.TransactionId,
							Status:    data.Status,
							Message:   err.Error(),
						}), nil
					}

					childId = campaign.ChildID
					function = childModule.GetFunctionSupportChildSpecialNeedCampaign()
				}

				child, err := on_chain.GetOnChainObject[entities.Child](on_chain.GetOnChainObjectRequest{
					Client:    client,
					ObjectId:  childId,
					ErrLogger: p.errLogger,
				}, ctx)
				if err != nil {
					return util.GeneratePaymentRedirectUrl(util.GenerateRedirectParamRequest{
						OrderCode: payment.TransactionId,
						Status:    data.Status,
						Message:   err.Error(),
					}), nil
				}

				pool, err := on_chain.GetOnChainObject[entities.MainPool](on_chain.GetOnChainObjectRequest{
					Client:    client,
					ObjectId:  os.Getenv(env.POOL_ID),
					ErrLogger: p.errLogger,
				}, ctx)
				if err != nil {
					return util.GeneratePaymentRedirectUrl(util.GenerateRedirectParamRequest{
						OrderCode: payment.TransactionId,
						Status:    data.Status,
						Message:   err.Error(),
					}), nil
				}

				localPools, err := on_chain.GetOnChainObjects[entities.LocalPool](on_chain.GetOnChainObjectsRequest{
					Client:    client,
					ObjectIds: pool.LocalPools,
					ErrLogger: p.errLogger,
				}, ctx)
				if err != nil {
					return util.GeneratePaymentRedirectUrl(util.GenerateRedirectParamRequest{
						OrderCode: payment.TransactionId,
						Status:    data.Status,
						Message:   err.Error(),
					}), nil
				}

				var leaders []string
				var localPoolId string
				for _, localPool := range localPools {
					if localPool.Region == child.Region {
						leaders = localPool.Mods
						localPoolId = localPool.ID.ID
						break
					}
				}

				if detail.Purpose != string(entities.SPECIAL_NEED_PURPOSE) {
					leaderNoti, err := p.leaderNotiRepo.GetNotiByNeed(detail.Target, ctx)
					if err != nil {
						return util.GeneratePaymentRedirectUrl(util.GenerateRedirectParamRequest{
							OrderCode: payment.TransactionId,
							Status:    data.Status,
							Message:   err.Error(),
						}), nil
					}

					var childFormattedAddr string = util.FormatAddress(childId)
					var contents []string
					for _, format := range contentFormats {
						contents = append(contents, fmt.Sprintf(format, child.LastName, child.FirstName, childFormattedAddr))
					}

					var curTime time.Time = time.Now()
					if leaderNoti == nil {
						p.leaderNotiRepo.CreateNoti(entities.LeaderNoti{
							ID:                      util.GenerateId(),
							NeedID:                  detail.Target,
							NeedType:                detail.Purpose,
							ChildID:                 childId,
							Region:                  child.Region,
							AssignedLeaders:         leaders,
							ExpectedWithdrawPeriods: expectedWithdrawPeriods,
							Contents:                contents,
							GeneralContent:          fmt.Sprintf(generalContentFormat, child.LastName, child.FirstName, childFormattedAddr),
							CreatedAt:               curTime,
							UpdatedAt:               curTime,
						}, ctx)
					} else {
						leaderNoti.ExpectedWithdrawPeriods = append(leaderNoti.ExpectedWithdrawPeriods, expectedWithdrawPeriods...)
						leaderNoti.Contents = append(leaderNoti.Contents, contents...)
						p.leaderNotiRepo.UpdateNoti(*leaderNoti, ctx)
					}

					switch detail.Purpose {
					case string(entities.BOOKS_NEED_PURPOSE):
						args = childModule.ToSupportChildBooksNeedArguments(on_chain.SupportChildBooksNeedArguments{
							NeedID:      detail.Target,
							LocalPool:   localPoolId,
							ChildID:     childId,
							DonorNft:    nftId,
							Amount:      payment.Amount,
							FirstName:   *profile.FirstName,
							LastName:    *profile.LastName,
							Gender:      *profile.Gender,
							PhoneNumber: *profile.PhoneNumber,
							Email:       *profile.Email,
							Message:     payment.Message,
						})
					case string(entities.MEAL_NEED_PURPOSE):
						var detailId string = util.GenerateId()
						p.childTaskDetailRepo.CreateChildTaskDetail(entities.ChildTaskDetail{
							ID:      detailId,
							ChildID: childId,
							Purpose: string(entities.MEAL_NEED_PURPOSE),
							Target:  detail.Target,
						}, ctx)

						task.IsChildTask = true
						task.ChildTaskDetailID = &detailId
						task.Region = child.Region
						task.Description = fmt.Sprintf(taskContentFormat, child.LastName, child.FirstName, childFormattedAddr)
						task.CreatedAt = curTime
						task.UpdatedAt = curTime
						p.taskRepo.CreateTask(task, ctx)

						args = childModule.ToSupportChildMealNeedArguments(on_chain.SupportChildMealNeedArguments{
							StartPeriod: mealSupportStartPeriod,
							EndPeriod:   mealSupportEndPeriod,
							SupportChildBooksNeedArguments: on_chain.SupportChildBooksNeedArguments{
								NeedID:      detail.Target,
								LocalPool:   localPoolId,
								ChildID:     childId,
								DonorNft:    nftId,
								Amount:      payment.Amount,
								FirstName:   *profile.FirstName,
								LastName:    *profile.LastName,
								Gender:      *profile.Gender,
								PhoneNumber: *profile.PhoneNumber,
								Email:       *profile.Email,
								Message:     payment.Message,
							},
						})
					case string(entities.HEALTH_INSURANCE_NEED_PURPOSE):
						args = childModule.ToSupportChildHealthInsuranceNeedArguments(on_chain.SupportChildHealthInsuranceNeedArguments{
							NeedID:      detail.Target,
							LocalPool:   localPoolId,
							ChildID:     childId,
							DonorNft:    nftId,
							Amount:      payment.Amount,
							FirstName:   *profile.FirstName,
							LastName:    *profile.LastName,
							Gender:      *profile.Gender,
							PhoneNumber: *profile.PhoneNumber,
							Email:       *profile.Email,
							Message:     payment.Message,
						})
					}
				} else {
					args = childModule.ToSupportChildSpeicalNeedArguments(on_chain.SupportChildSpeicalNeedArguments{
						CampaignID:  detail.Target,
						LocalPool:   localPoolId,
						ChildID:     childId,
						DonorNft:    nftId,
						Amount:      payment.Amount,
						FirstName:   *profile.FirstName,
						LastName:    *profile.LastName,
						Gender:      *profile.Gender,
						PhoneNumber: *profile.PhoneNumber,
						Email:       *profile.Email,
						Message:     payment.Message,
					})
				}
			}
		} else {
			var poolModule = on_chain.InitializeModulePool()
			module = poolModule.GetModule()
			if detail.Target != os.Getenv(env.POOL_ID) { // Donate to local pool
				function = poolModule.GetFunctionDonateToLocalPool()
				args = poolModule.ToDonateToLocalPoolArguments(on_chain.DonateToLocalPoolArguments{
					LocalPoolId: detail.Target,
					DonateToPoolArguments: on_chain.DonateToPoolArguments{
						DonorID:     nftId,
						Amount:      payment.Amount,
						FirstName:   *profile.FirstName,
						LastName:    *profile.LastName,
						Gender:      *profile.Gender,
						PhoneNumber: *profile.PhoneNumber,
						Email:       *profile.Email,
						Message:     payment.Message,
					},
				})
			} else {
				function = poolModule.GetFunctionDonateToPool()
				args = poolModule.ToDonateToPoolArguments(on_chain.DonateToPoolArguments{
					DonorID:     nftId,
					Amount:      payment.Amount,
					FirstName:   *profile.FirstName,
					LastName:    *profile.LastName,
					Gender:      *profile.Gender,
					PhoneNumber: *profile.PhoneNumber,
					Email:       *profile.Email,
					Message:     payment.Message,
				})
			}
		}
	}

	txBytes, err := on_chain.BuildTransaction(on_chain.BuildTransactionRequest{
		Client:    client,
		Sender:    payment.Actor,
		Module:    module,
		Function:  function,
		ErrLogger: p.errLogger,
		Arguments: args,
	}, ctx)

	var req = util.GenerateRedirectParamRequest{
		OrderCode: payment.TransactionId,
		Status:    data.Status,
		TxBytes:   txBytes,
	}

	if err != nil {
		req.Message = err.Error()
	} else {
		req.Message = "Success"
	}

	return util.GeneratePaymentRedirectUrl(req), nil
}

// CallbackWithAuth implements business.IPaymentService.
func (p *paymentService) CallbackWithAuth(id string, capturedImgBlobId string, ctx context.Context) (string, error) {
	var genericErr error = errors.New(noti.GENERIC_ERROR_WARN_MSG)

	payment, err := p.paymentRepo.GetPaymentById(id, ctx)
	if err != nil {
		return "", err
	}

	if payment == nil {
		return "", genericErr
	}

	var sender string = ctx.Value("address").(string)
	if payment.Actor != sender {
		return "", errors.New(noti.GENERIC_RIGHT_ACCESS_WARN_MSG)
	}

	// Expired
	if payment.ExpiredAt.Before(time.Now()) {
		return util.GeneratePaymentRedirectUrl(util.GenerateRedirectParamRequest{
			OrderCode: payment.TransactionId,
			Status:    shared.PAYAMENT_EXPIRED_STATUS,
			Message:   noti.PAYMENT_EXPIRED_MESSAGE,
		}), nil
	}

	data, err := payos.GetPaymentLinkInformation(payment.TransactionId)
	if err != nil {
		p.errLogger.Println("Error while get payos payment link information: " + err.Error())
		return "", errors.New(noti.INTERNALL_ERR_MSG)
	}

	switch data.Status {
	case shared.PAYOS_PAID_STATUS:
		payment.Status = shared.PAYMENT_SUCCESS_STATUS
	case shared.PAYOS_CANCELLED_STATUS:
		payment.Status = shared.PAYMENT_CANCELED_STATUS
		if data.CancellationReason != nil {
			payment.CancelReason = data.CancellationReason
		}
	default:
		return "", genericErr
	}

	payment.UpdatedAt = time.Now()
	if err := p.paymentRepo.UpdatePayment(*payment, ctx); err != nil {
		return "", err
	}

	if data.Status == shared.PAYMENT_CANCELED_STATUS {
		return util.GeneratePaymentRedirectUrl(util.GenerateRedirectParamRequest{
			OrderCode: payment.TransactionId,
			Status:    data.Status,
			Message:   noti.PAYMENT_CANCEL_MESSAGE,
		}), nil
	}

	profile, err := p.profileRepo.GetProfile(payment.ProfileID, ctx)
	if err != nil {
		return "", err
	}

	if profile == nil {
		return "", genericErr
	}

	var client = p.clients[constant.SuiTestnet]
	var module string
	var function string
	var args []interface{}
	if payment.IsDonateTx { // Never has this case with this endpoint
	} else {
		detail, err := p.withdrawRepo.GetOffChainWithdrawProposal(*payment.ProposalID, ctx)
		if err != nil {
			return util.GeneratePaymentRedirectUrl(util.GenerateRedirectParamRequest{
				OrderCode: payment.TransactionId,
				Status:    data.Status,
				Message:   err.Error(),
			}), nil
		}

		if detail.Purpose == string(entities.WITHDRAW_PURPOSE) {
			var poolModule = on_chain.InitializeModulePool()
			module = poolModule.GetModule()
			function = poolModule.GetFunctionWithdrawFromPool()
			args = poolModule.ToWithdrawFromPoolArguments(on_chain.WithdrawFromPoolArguments{
				LocalPoolId:        detail.LocalPoolID,
				WithdrawProposalId: *detail.ProposalID,
			})
		} else {
			if detail.Purpose == string(entities.CAMPAIGN_PURPOSE) {
				var campaignModule = on_chain.InitializeModuleCampaign()
				module = campaignModule.GetModule()
				function = campaignModule.GetFunctionWithdrawFromCampaign()
				args = campaignModule.ToWithdrawFromCampaignArguments(on_chain.WithdrawFromCampaignArguments{
					CampaignID: detail.Target,
					ProposalID: *detail.ProposalID,
				})
			} else {
				var childModule = on_chain.InitializeModuleChild()
				module = childModule.GetModule()
				switch detail.Purpose {
				case string(entities.BOOKS_NEED_PURPOSE):
					function = childModule.GetFunctionWithdrawFromBooksNeedProposal()
				case string(entities.HEALTH_INSURANCE_NEED_PURPOSE):
					function = childModule.GetFunctionWithdrawFromHealthInsuranceNeedProposal()
				case string(entities.MEAL_NEED_PURPOSE):
					function = childModule.GetFunctionWithdrawFromMealNeedProposal()
				case string(entities.SPECIAL_NEED_PURPOSE):
					function = childModule.GetFunctionWithdrawFromSpecialNeedCampaign()
				}

				args = childModule.ToWithdrawFromNeedArguments(on_chain.WithdrawFromNeedArguments{
					LocalPool:  detail.LocalPoolID,
					TargetID:   detail.Target,
					ProposalID: *detail.ProposalID,
				})
			}
		}
	}

	txBytes, err := on_chain.BuildTransaction(on_chain.BuildTransactionRequest{
		Client:    client,
		Sender:    payment.Actor,
		Module:    module,
		Function:  function,
		ErrLogger: p.errLogger,
		Arguments: args,
	}, ctx)

	var req = util.GenerateRedirectParamRequest{
		OrderCode: payment.TransactionId,
		Status:    data.Status,
		TxBytes:   txBytes,
	}

	if err != nil {
		req.Message = err.Error()
	} else {
		req.Message = "Success"
	}

	return util.GeneratePaymentRedirectUrl(req), nil
}

func (p *paymentService) getGetPaymentsRedisKey(req request.GetPaymentsRequest) string {
	var keyword string = "empty"
	if req.Keyword != "" {
		keyword = req.Keyword
	}

	var status string = "empty"
	if req.Actor != "" {
		status = req.Status
	}

	var method string = "empty"
	if req.Actor != "" {
		method = req.Method
	}

	var actor string = "empty"
	if req.Actor != "" {
		actor = req.Actor
	}

	var minAmount string = "empty"
	if req.MinAmount != nil {
		minAmount = fmt.Sprintf("%d", *req.MinAmount)
	}

	var maxAmount string = "empty"
	if req.MaxAmount != nil {
		maxAmount = fmt.Sprintf("%d", *req.MaxAmount)
	}

	var isDonateTx string = "empty"
	if req.IsDonatePayment != nil {
		isDonateTx = fmt.Sprintf("%v", *req.IsDonatePayment)
	}

	var isExpired string = "empty"
	if req.IsPaymentExpired != nil {
		isExpired = fmt.Sprintf("%v", *req.IsPaymentExpired)
	}

	var sortCriteria string = "empty"
	if req.FilterProp != "" {
		sortCriteria = req.FilterProp
	}

	return fmt.Sprintf("payment:kw:%s:status:%s:method:%s:of:%s:min:%s:max:%s:is_donate_tx:%s:expired:%s:sc:%s:o:%s:s:%d:p:%d",
		keyword, status, method, actor, minAmount, maxAmount, isDonateTx, isExpired, sortCriteria, req.SortOrder, req.PageSize, req.Page)
}
