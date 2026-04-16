package business

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"raise-child/constants/env"
	"raise-child/constants/env/payment"
	"raise-child/constants/noti"
	"raise-child/constants/shared"
	"raise-child/interfaces/business"
	i_repository "raise-child/interfaces/repository"
	"raise-child/model/dtos/request"
	"raise-child/model/dtos/response"
	"raise-child/model/entities"
	"raise-child/repository"
	"raise-child/util"
	"raise-child/util/db"
	on_chain "raise-child/util/on_chain"
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
	profileRepo             i_repository.IProfileRepository
	donationRepo            i_repository.IOffChainDonationRepository
	withdrawRepo            i_repository.IOffChainWithdrawProposalRepository
	clients                 map[string]sui.ISuiAPI
	errLogger               *log.Logger
}

func initializePaymentService(
	mealSupportDurationRepo i_repository.IMealSupportDurationRepository,
	leaderNotiRepo i_repository.ILeaderNotiRepository,
	childTaskDetailRepo i_repository.IChildTaskDetailRepository,
	taskRepo i_repository.ITaskRepository,
	paymentRepo i_repository.IPaymentRepository,
	profileRepo i_repository.IProfileRepository,
	donationRepo i_repository.IOffChainDonationRepository,
	withdrawRepo i_repository.IOffChainWithdrawProposalRepository,
	clients map[string]sui.ISuiAPI,
	errLogger *log.Logger,
) business.IPaymentService {
	return &paymentService{
		mealSupportDurationRepo: mealSupportDurationRepo,
		leaderNotiRepo:          leaderNotiRepo,
		childTaskDetailRepo:     childTaskDetailRepo,
		taskRepo:                taskRepo,
		paymentRepo:             paymentRepo,
		profileRepo:             profileRepo,
		donationRepo:            donationRepo,
		withdrawRepo:            withdrawRepo,
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
		repository.InitializeProfileRepository(cnn, errLogger),
		repository.InitializeOffChainDonationRepository(cnn, errLogger),
		repository.InitializeOffChainWithdrawProposalRepository(cnn, errLogger),
		_networkAliases,
		errLogger,
	), nil
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

// // CallbackTx implements business.IPaymentService.
// func (p *paymentService) CallbackTx(id string, ctx context.Context) (response.BuildTransactionResponse, error) {
// 	var genericErr error = errors.New(noti.GENERIC_ERROR_WARN_MSG)

// 	payment, err := p.paymentRepo.GetPaymentById(id, ctx)
// 	if err != nil {
// 		return response.BuildTransactionResponse{}, err
// 	}

// 	if payment == nil {
// 		return response.BuildTransactionResponse{}, genericErr
// 	}

// 	if !payment.IsDonateTx {
// 		return response.BuildTransactionResponse{}, genericErr
// 	}

// 	// Expired
// 	if payment.ExpiredAt.Before(time.Now()) {
// 		return response.BuildTransactionResponse{}, errors.New(noti.PAYMENT_EXPIRED_MESSAGE)
// 	}

// 	data, err := payos.GetPaymentLinkInformation(payment.TransactionId)
// 	if err != nil {
// 		p.errLogger.Println("Error while get payos payment link information: " + err.Error())
// 		return response.BuildTransactionResponse{}, errors.New(noti.INTERNALL_ERR_MSG)
// 	}

// 	switch data.Status {
// 	case shared.PAYOS_PAID_STATUS:
// 		payment.Status = shared.PAYMENT_SUCCESS_STATUS
// 	case shared.PAYOS_CANCELLED_STATUS:
// 		payment.Status = shared.PAYMENT_CANCELED_STATUS
// 		if data.CancellationReason != nil {
// 			payment.CancelReason = *data.CancellationReason
// 		}
// 	default:
// 		return response.BuildTransactionResponse{}, genericErr
// 	}

// 	payment.UpdatedAt = time.Now()
// 	if err := p.paymentRepo.UpdatePayment(*payment, ctx); err != nil {
// 		return response.BuildTransactionResponse{}, err
// 	}

// 	if data.Status == shared.PAYOS_CANCELLED_STATUS {
// 		return response.BuildTransactionResponse{}, errors.New(noti.PAYMENT_CANCEL_MESSAGE)
// 	}

// 	profile, err := p.profileRepo.GetProfile(payment.Sub, ctx)
// 	if err != nil {
// 		return response.BuildTransactionResponse{}, err
// 	}

// 	if profile == nil {
// 		return response.BuildTransactionResponse{}, genericErr
// 	}

// 	var module = on_chain.InitializeModulePool()
// 	var function string
// 	var args []interface{}
// 	if payment.Target != os.Getenv(env.POOL_ID) { // Local Pool
// 		function = module.GetFunctionDonateToLocalPool()
// 		args = module.ToDonateToLocalPoolArguments(on_chain.DonateToLocalPoolArguments{
// 			LocalPoolId: payment.Target,
// 			DonateToPoolArguments: on_chain.DonateToPoolArguments{
// 				Amount:      payment.Amount,
// 				FirstName:   profile.FirstName,
// 				LastName:    profile.LastName,
// 				Gender:      profile.Gender,
// 				PhoneNumber: profile.PhoneNumber,
// 				Email:       profile.Email,
// 				Message:     payment.Message,
// 			},
// 		})
// 	} else { // Main pool
// 		function = module.GetFunctionDonateToPool()
// 		args = module.ToDonateToPoolArguments(on_chain.DonateToPoolArguments{
// 			Amount:      payment.Amount,
// 			FirstName:   profile.FirstName,
// 			LastName:    profile.LastName,
// 			Gender:      profile.Gender,
// 			PhoneNumber: profile.PhoneNumber,
// 			Email:       profile.Email,
// 			Message:     payment.Message,
// 		})
// 	}

// 	txBytes, err := on_chain.BuildTransaction(on_chain.BuildTransactionRequest{
// 		Client:    p.clients[constant.SuiTestnet],
// 		Sender:    payment.Actor,
// 		Module:    module.GetModule(),
// 		Function:  function,
// 		ErrLogger: p.errLogger,
// 		Arguments: args,
// 	}, ctx)

// 	return response.BuildTransactionResponse{
// 		TxBytes: txBytes,
// 	}, err
// }

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
			ExpiredAt:     time.Unix(int64(*data.ExpiredAt), 0),
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

		if nfts != nil && len(nfts) > 0 {
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
				// switch detail.Purpose {
				// case string(entities.BOOKS_NEED_PURPOSE):
				// 	need, err := on_chain.GetOnChainObject[entities.BooksNeed](on_chain.GetOnChainObjectRequest{
				// 		Client:    client,
				// 		ObjectId:  detail.Target,
				// 		ErrLogger: p.errLogger,
				// 	}, ctx)
				// 	if err != nil {
				// 		return util.GeneratePaymentRedirectUrl(util.GenerateRedirectParamRequest{
				// 			OrderCode: payment.TransactionId,
				// 			Status:    data.Status,
				// 			Message:   err.Error(),
				// 		}), nil
				// 	}

				// 	child, err := on_chain.GetOnChainObject[entities.Child](on_chain.GetOnChainObjectRequest{
				// 		Client:    client,
				// 		ObjectId:  need.ChildID,
				// 		ErrLogger: p.errLogger,
				// 	}, ctx)

				// 	pool, err := on_chain.GetOnChainObject[entities.MainPool](on_chain.GetOnChainObjectRequest{
				// 		Client:    client,
				// 		ObjectId:  os.Getenv(env.POOL_ID),
				// 		ErrLogger: p.errLogger,
				// 	}, ctx)

				// 	withdrawDates, err := on_chain.GetOnChainObject[entities.BooksNeedWithdrawDates](on_chain.GetOnChainObjectRequest{
				// 		Client:    client,
				// 		ObjectId:  os.Getenv(env.BOOKS_NEED_WITHDRAW_DATES_ID),
				// 		ErrLogger: p.errLogger,
				// 	}, ctx)

				// 	var withdrawDate string
				// 	if need.Semster == "1" {
				// 		withdrawDate = withdrawDates.FirstSemesterDate
				// 	} else {
				// 		withdrawDate = withdrawDates.SecondSemesterDate
				// 	}

				// 	localPools, err := on_chain.GetOnChainObjects[entities.LocalPool](on_chain.GetOnChainObjectsRequest{
				// 		Client:    client,
				// 		ObjectIds: pool.LocalPools,
				// 		ErrLogger: p.errLogger,
				// 	}, ctx)

				// 	var leaders []string
				// 	var localPoolId string
				// 	for _, localPool := range localPools {
				// 		if localPool.Region == child.Region {
				// 			leaders = localPool.Mods
				// 			localPoolId = localPool.ID.ID
				// 			break
				// 		}
				// 	}

				// 	var rawExpectedDate string = withdrawDate + "/" + need.Year
				// 	var content string = fmt.Sprintf("Withdraw books need semester %s - %s for child %s %s - %s", need.Semster, need.Year, child.LastName, child.FirstName, util.FormatAddress(child.ID.ID))
				// 	if leaderNoti == nil {
				// 		var curTime time.Time = time.Now()
				// 		p.leaderNotiRepo.CreateNoti(entities.LeaderNoti{
				// 			ID:                      util.GenerateId(),
				// 			NeedID:                  detail.Target,
				// 			NeedType:                string(entities.BOOKS_NEED_PURPOSE),
				// 			ChildID:                 need.ChildID,
				// 			Region:                  child.Region,
				// 			AssignedLeaders:         leaders,
				// 			ExpectedWithdrawPeriods: []string{rawExpectedDate},
				// 			Contents:                []string{content},
				// 			GeneralContent:          fmt.Sprintf("Withdraw books need semester %s for child %s %s - %s", need.Semster, child.LastName, child.FirstName, util.FormatAddress(child.ID.ID)),
				// 			CreatedAt:               curTime,
				// 			UpdatedAt:               curTime,
				// 		}, ctx)
				// 	} else {
				// 		leaderNoti.ExpectedWithdrawPeriods = append(leaderNoti.ExpectedWithdrawPeriods, rawExpectedDate)
				// 		leaderNoti.Contents = append(leaderNoti.Contents, content)
				// 		p.leaderNotiRepo.UpdateNoti(*leaderNoti, ctx)
				// 	}

				// 	function = childModule.GetFunctionSupportChildBooksNeed()
				// 	args = childModule.ToSupportChildBooksNeedArguments(on_chain.SupportChildBooksNeedArguments{
				// 		NeedID:      detail.Target,
				// 		LocalPool:   localPoolId,
				// 		ChildID:     need.ChildID,
				// 		DonorNft:    nftId,
				// 		Amount:      payment.Amount,
				// 		FirstName:   profile.FirstName,
				// 		LastName:    profile.LastName,
				// 		Gender:      profile.Gender,
				// 		PhoneNumber: profile.PhoneNumber,
				// 		Email:       profile.Email,
				// 		Message:     payment.Message,
				// 	})
				// case string(entities.HEALTH_INSURANCE_NEED_PURPOSE):
				// 	need, err := on_chain.GetOnChainObject[entities.HealthInsuranceNeed](on_chain.GetOnChainObjectRequest{
				// 		Client:    client,
				// 		ObjectId:  id,
				// 		ErrLogger: p.errLogger,
				// 	}, ctx)
				// 	if err != nil {
				// 		return util.GeneratePaymentRedirectUrl(util.GenerateRedirectParamRequest{
				// 			OrderCode: payment.TransactionId,
				// 			Status:    data.Status,
				// 			Message:   err.Error(),
				// 		}), nil
				// 	}

				// 	child, err := on_chain.GetOnChainObject[entities.Child](on_chain.GetOnChainObjectRequest{
				// 		Client:    client,
				// 		ObjectId:  need.ChildID,
				// 		ErrLogger: p.errLogger,
				// 	}, ctx)

				// 	pool, err := on_chain.GetOnChainObject[entities.MainPool](on_chain.GetOnChainObjectRequest{
				// 		Client:    client,
				// 		ObjectId:  os.Getenv(env.POOL_ID),
				// 		ErrLogger: p.errLogger,
				// 	}, ctx)

				// 	withdrawDate, err := on_chain.GetOnChainObject[entities.HealthInsuranceNeedWithdrawDate](on_chain.GetOnChainObjectRequest{
				// 		Client:    client,
				// 		ObjectId:  os.Getenv(""),
				// 		ErrLogger: p.errLogger,
				// 	}, ctx)

				// 	localPools, err := on_chain.GetOnChainObjects[entities.LocalPool](on_chain.GetOnChainObjectsRequest{
				// 		Client:    client,
				// 		ObjectIds: pool.LocalPools,
				// 		ErrLogger: p.errLogger,
				// 	}, ctx)

				// 	var leaders []string
				// 	//var localPoolId string
				// 	for _, localPool := range localPools {
				// 		if localPool.Region == child.Region {
				// 			leaders = localPool.Mods
				// 			//localPoolId = localPool.ID.ID
				// 			break
				// 		}
				// 	}

				// 	var rawExpectedDate string = withdrawDate.ExpectedDate + "/" + need.Year
				// 	var content string = fmt.Sprintf("Withdraw health insurance need %s for child %s %s - %s", need.Year, child.LastName, child.FirstName, util.FormatAddress(child.ID.ID))
				// 	if leaderNoti == nil {
				// 		var curTime time.Time = time.Now()
				// 		p.leaderNotiRepo.CreateNoti(entities.LeaderNoti{
				// 			ID:                      util.GenerateId(),
				// 			NeedID:                  detail.Target,
				// 			NeedType:                string(entities.HEALTH_INSURANCE_NEED_PURPOSE),
				// 			ChildID:                 need.ChildID,
				// 			Region:                  child.Region,
				// 			AssignedLeaders:         leaders,
				// 			ExpectedWithdrawPeriods: []string{rawExpectedDate},
				// 			Contents:                []string{content},
				// 			GeneralContent:          fmt.Sprintf("Withdraw health insurance need for child %s %s - %s", child.LastName, child.FirstName, util.FormatAddress(child.ID.ID)),
				// 			CreatedAt:               curTime,
				// 			UpdatedAt:               curTime,
				// 		}, ctx)
				// 	} else {
				// 		leaderNoti.ExpectedWithdrawPeriods = append(leaderNoti.ExpectedWithdrawPeriods, rawExpectedDate)
				// 		leaderNoti.Contents = append(leaderNoti.Contents, content)
				// 		p.leaderNotiRepo.UpdateNoti(*leaderNoti, ctx)
				// 	}

				// 	function = ""
				// case string(entities.MEAL_NEED_PURPOSE):
				// 	need, err := on_chain.GetOnChainObject[entities.MealNeed](on_chain.GetOnChainObjectRequest{
				// 		Client:    client,
				// 		ObjectId:  detail.Target,
				// 		ErrLogger: p.errLogger,
				// 	}, ctx)
				// 	if err != nil {
				// 		return util.GeneratePaymentRedirectUrl(util.GenerateRedirectParamRequest{
				// 			OrderCode: payment.TransactionId,
				// 			Status:    data.Status,
				// 			Message:   err.Error(),
				// 		}), nil
				// 	}

				// 	child, err := on_chain.GetOnChainObject[entities.Child](on_chain.GetOnChainObjectRequest{
				// 		Client:    client,
				// 		ObjectId:  need.ChildID,
				// 		ErrLogger: p.errLogger,
				// 	}, ctx)

				// 	pool, err := on_chain.GetOnChainObject[entities.MainPool](on_chain.GetOnChainObjectRequest{
				// 		Client:    client,
				// 		ObjectId:  os.Getenv(env.POOL_ID),
				// 		ErrLogger: p.errLogger,
				// 	}, ctx)

				// 	duration, err := p.mealSupportDurationRepo.GetMealSupportDuration(detail.MealDurationID, ctx)

				// 	var expectedWithdrawDates, contents []string
				// 	var startPeriod time.Time = util.ToStartOfDate(util.RawDateToTime(duration.StartPeriod))
				// 	var endPeriod time.Time = util.ToEndOfDate(util.RawDateToTime(duration.EndPeriod))
				// 	var endMonth int = int(endPeriod.Month())
				// 	if endMonth == 1 {
				// 		endMonth = 13
				// 	}

				// 	var initExpectedDate time.Time = startPeriod.AddDate(0, 0, -3)
				// 	var childFormattedAddr string = util.FormatAddress(child.ID.ID)
				// 	var contentFormat string = "Withdraw meal need %s - %s for child " + child.LastName + " " + child.FirstName + " - " + childFormattedAddr
				// 	for i := 0; i < endMonth-int(startPeriod.Month()); i++ {
				// 		var date time.Time = initExpectedDate.AddDate(0, i, 0)
				// 		expectedWithdrawDates = append(expectedWithdrawDates, util.TimeToRawDate(date))

				// 		var startProvidePeriod time.Time = startPeriod.AddDate(0, i, 0)
				// 		var endProvidePeriod time.Time = startProvidePeriod.AddDate(0, 1, 0)
				// 		contents = append(contents, fmt.Sprintf(contentFormat, util.TimeToRawDate(startProvidePeriod), util.TimeToRawDate(endProvidePeriod)))
				// 	}

				// 	localPools, err := on_chain.GetOnChainObjects[entities.LocalPool](on_chain.GetOnChainObjectsRequest{
				// 		Client:    client,
				// 		ObjectIds: pool.LocalPools,
				// 		ErrLogger: p.errLogger,
				// 	}, ctx)

				// 	var leaders []string
				// 	var localPoolId string
				// 	for _, localPool := range localPools {
				// 		if localPool.Region == child.Region {
				// 			leaders = localPool.Mods
				// 			localPoolId = localPool.ID.ID
				// 			break
				// 		}
				// 	}

				// 	if leaderNoti == nil {
				// 		var curTime time.Time = time.Now()
				// 		p.leaderNotiRepo.CreateNoti(entities.LeaderNoti{
				// 			ID:                      util.GenerateId(),
				// 			NeedID:                  detail.Target,
				// 			NeedType:                string(entities.MEAL_NEED_PURPOSE),
				// 			ChildID:                 need.ChildID,
				// 			Region:                  child.Region,
				// 			AssignedLeaders:         leaders,
				// 			ExpectedWithdrawPeriods: expectedWithdrawDates,
				// 			Contents:                contents,
				// 			GeneralContent:          fmt.Sprintf("Withdraw meal need for child %s %s - %s", child.LastName, child.FirstName, childFormattedAddr),
				// 			CreatedAt:               curTime,
				// 			UpdatedAt:               curTime,
				// 		}, ctx)
				// 	} else {
				// 		leaderNoti.ExpectedWithdrawPeriods = append(leaderNoti.ExpectedWithdrawPeriods, expectedWithdrawDates...)
				// 		leaderNoti.Contents = append(leaderNoti.Contents, contents...)
				// 		p.leaderNotiRepo.UpdateNoti(*leaderNoti, ctx)
				// 	}

				// 	manageObj, err := on_chain.GetOnChainObject[entities.Manage](on_chain.GetOnChainObjectRequest{
				// 		Client:    client,
				// 		ObjectId:  os.Getenv(env.PACKAGE_ID),
				// 		ErrLogger: p.errLogger,
				// 	}, ctx)

				// 	volunteers, err := on_chain.GetOnChainObjects[entities.StaffNft](on_chain.GetOnChainObjectsRequest{
				// 		Client:    client,
				// 		ObjectIds: manageObj.VolunteerNfts,
				// 		ErrLogger: p.errLogger,
				// 	}, ctx)

				// 	var volunteerAddresses []string
				// 	for i, volunteer := range volunteers {
				// 		if volunteer.Region == child.Region {
				// 			volunteerAddresses = append(volunteerAddresses, manageObj.VolunteerIds[i])
				// 		}
				// 	}

				// 	p.taskRepo.CreateNoti(entities.VolunteerNoti{
				// 		ID:                 util.GenerateId(),
				// 		ChildID:            need.ChildID,
				// 		Region:             child.Region,
				// 		AssginedVolunteers: volunteerAddresses,
				// 		Content:            fmt.Sprintf("Provide meal for child %s %s - %s from %s to %s", child.LastName, child.FirstName, childFormattedAddr, duration.StartPeriod, duration.EndPeriod),
				// 		StartPeriod:        startPeriod,
				// 		EndPeriod:          endPeriod,
				// 	}, ctx)

				// 	function = childModule.GetFunctionSupportChildMealNeed()
				// 	args = childModule.ToSupportChildMealNeedArguments(on_chain.SupportChildMealNeedArguments{
				// 		StartPeriod: duration.StartPeriod,
				// 		EndPeriod:   duration.EndPeriod,
				// 		SupportChildBooksNeedArguments: on_chain.SupportChildBooksNeedArguments{
				// 			NeedID:      detail.Target,
				// 			LocalPool:   localPoolId,
				// 			ChildID:     need.ChildID,
				// 			DonorNft:    nftId,
				// 			Amount:      payment.Amount,
				// 			FirstName:   profile.FirstName,
				// 			LastName:    profile.LastName,
				// 			Gender:      profile.Gender,
				// 			PhoneNumber: profile.PhoneNumber,
				// 			Email:       profile.Email,
				// 			Message:     payment.Message,
				// 		},
				// 	})
				// case string(entities.SPECIAL_NEED_PURPOSE):
				// 	campaign, err := on_chain.GetOnChainObject[entities.SpecialNeedCampaign](on_chain.GetOnChainObjectRequest{
				// 		Client:    client,
				// 		ObjectId:  detail.Target,
				// 		ErrLogger: p.errLogger,
				// 	}, ctx)
				// 	if err != nil {
				// 		return util.GeneratePaymentRedirectUrl(util.GenerateRedirectParamRequest{
				// 			OrderCode: payment.TransactionId,
				// 			Status:    data.Status,
				// 			Message:   err.Error(),
				// 		}), nil
				// 	}

				// 	function = childModule.GetFunctionSupportChildSpecialNeedCampaign()
				// 	args = childModule.ToSupportChildSpeicalNeedArguments(on_chain.SupportChildSpeicalNeedArguments{
				// 		CampaignID:  detail.Target,
				// 		LocalPool:   detail.LocalPoolID,
				// 		ChildID:     campaign.ChildID,
				// 		DonorNft:    nftId,
				// 		Amount:      payment.Amount,
				// 		FirstName:   profile.FirstName,
				// 		LastName:    profile.LastName,
				// 		Gender:      profile.Gender,
				// 		PhoneNumber: profile.PhoneNumber,
				// 		Email:       profile.Email,
				// 		Message:     payment.Message,
				// 	})
				// }

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
					contentFormats = append(contentFormats, "Withdraw books need semester "+need.Semster+" - "+need.Year+" for child %s %s - %s")
					generalContentFormat = "Withdraw books need semester " + need.Semster + " for child %s %s - %s"
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
						ObjectId:  os.Getenv(env.HEALTH_INSURANCE_NEED_wITHDRAW_DATE_ID),
						ErrLogger: p.errLogger,
					}, ctx)

					expectedWithdrawPeriods = append(expectedWithdrawPeriods, withdrawDate.ExpectedDate+"/"+need.Year)
					contentFormats = append(contentFormats, "Withdraw health insurance need "+need.Year+" for child %s %s - %s")
					generalContentFormat = "Withdraw health insurance need for child %s %s - %s"
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

					taskContentFormat = "Provide meal for child %s %s - %s from " + duration.StartPeriod + " to " + duration.EndPeriod
					mealSupportStartPeriod = duration.StartPeriod
					mealSupportEndPeriod = duration.EndPeriod
					generalContentFormat = "Withdraw meal need for child %s %s - %s"
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
	} else { // Never has case withdraw with this call back endpoint
		// detail, err := p.withdrawRepo.GetOffChainWithdrawProposal(*payment.ProposalID, ctx)
		// if err != nil {
		// 	return util.GeneratePaymentRedirectUrl(util.GenerateRedirectParamRequest{
		// 		OrderCode: payment.TransactionId,
		// 		Status:    data.Status,
		// 		Message:   err.Error(),
		// 	}), nil
		// }

		// if detail.Purpose == string(entities.WITHDRAW_PURPOSE) {
		// 	var poolModule = on_chain.InitializeModulePool()
		// 	var localPoolId string
		// 	if detail.Target != os.Getenv(env.POOL_ID) {
		// 		localPoolId = detail.Target
		// 	} else {
		// 		localPoolId = os.Getenv(env.SHARED_LOCAL_POOL_ID)
		// 	}

		// 	module = poolModule.GetModule()
		// 	function = poolModule.GetFunctionWithdrawFromPool()
		// 	args = poolModule.ToWithdrawFromPoolArguments(on_chain.WithdrawFromPoolArguments{
		// 		LocalPoolId:        localPoolId,
		// 		WithdrawProposalId: detail.ProposalID,
		// 	})
		// } else {
		// 	proposal, err := on_chain.GetOnChainObject[entities.WithdrawProposal](on_chain.GetOnChainObjectRequest{
		// 		Client:    client,
		// 		ObjectId:  detail.ProposalID,
		// 		ErrLogger: p.errLogger,
		// 	}, ctx)
		// 	if err != nil {
		// 		return util.GeneratePaymentRedirectUrl(util.GenerateRedirectParamRequest{
		// 			OrderCode: payment.TransactionId,
		// 			Status:    data.Status,
		// 			Message:   err.Error(),
		// 		}), nil
		// 	}

		// 	pool, err := on_chain.GetOnChainObject[entities.MainPool](on_chain.GetOnChainObjectRequest{
		// 		Client:    client,
		// 		ObjectId:  os.Getenv(env.POOL_ID),
		// 		ErrLogger: p.errLogger,
		// 	}, ctx)
		// 	if err != nil {
		// 		return util.GeneratePaymentRedirectUrl(util.GenerateRedirectParamRequest{
		// 			OrderCode: payment.TransactionId,
		// 			Status:    data.Status,
		// 			Message:   err.Error(),
		// 		}), nil
		// 	}

		// 	localPools, err := on_chain.GetOnChainObjects[entities.LocalPool](on_chain.GetOnChainObjectsRequest{
		// 		Client:    client,
		// 		ObjectIds: pool.LocalPools,
		// 		ErrLogger: p.errLogger,
		// 	}, ctx)
		// 	if err != nil {
		// 		return util.GeneratePaymentRedirectUrl(util.GenerateRedirectParamRequest{
		// 			OrderCode: payment.TransactionId,
		// 			Status:    data.Status,
		// 			Message:   err.Error(),
		// 		}), nil
		// 	}

		// 	var localPoolId string
		// 	for _, localPool := range localPools {
		// 		if localPool.Region == proposal.PoolName {
		// 			localPoolId = localPool.ID.ID
		// 			break
		// 		}
		// 	}

		// 	var childModule = on_chain.InitializeModuleChild()
		// 	switch detail.Purpose {
		// 	case string(entities.BOOKS_NEED_PURPOSE):
		// 		function = childModule.GetFunctionWithdrawFromBooksNeedProposal()
		// 	case string(entities.MEAL_NEED_PURPOSE):
		// 		function = childModule.GetFunctionWithdrawFromMealNeedProposal()
		// 	case string(entities.SPECIAL_NEED_PURPOSE):
		// 		function = childModule.GetFunctionWithdrawFromSpecialNeedCampaign()
		// 	}

		// 	args = childModule.ToWithdrawFromNeedArguments(on_chain.WithdrawFromNeedArguments{
		// 		LocalPool:  localPoolId,
		// 		TargetID:   detail.Target,
		// 		ProposalID: detail.ProposalID,
		// 	})
		// }
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
			// proposal, err := on_chain.GetOnChainObject[entities.WithdrawProposal](on_chain.GetOnChainObjectRequest{
			// 	Client:    client,
			// 	ObjectId:  detail.ProposalID,
			// 	ErrLogger: p.errLogger,
			// }, ctx)
			// if err != nil {
			// 	return util.GeneratePaymentRedirectUrl(util.GenerateRedirectParamRequest{
			// 		OrderCode: payment.TransactionId,
			// 		Status:    data.Status,
			// 		Message:   err.Error(),
			// 	}), nil
			// }

			// pool, err := on_chain.GetOnChainObject[entities.MainPool](on_chain.GetOnChainObjectRequest{
			// 	Client:    client,
			// 	ObjectId:  os.Getenv(env.POOL_ID),
			// 	ErrLogger: p.errLogger,
			// }, ctx)
			// if err != nil {
			// 	return util.GeneratePaymentRedirectUrl(util.GenerateRedirectParamRequest{
			// 		OrderCode: payment.TransactionId,
			// 		Status:    data.Status,
			// 		Message:   err.Error(),
			// 	}), nil
			// }

			// localPools, err := on_chain.GetOnChainObjects[entities.LocalPool](on_chain.GetOnChainObjectsRequest{
			// 	Client:    client,
			// 	ObjectIds: pool.LocalPools,
			// 	ErrLogger: p.errLogger,
			// }, ctx)
			// if err != nil {
			// 	return util.GeneratePaymentRedirectUrl(util.GenerateRedirectParamRequest{
			// 		OrderCode: payment.TransactionId,
			// 		Status:    data.Status,
			// 		Message:   err.Error(),
			// 	}), nil
			// }

			// var localPoolId string
			// for _, localPool := range localPools {
			// 	if localPool.Region == proposal.PoolName {
			// 		localPoolId = localPool.ID.ID
			// 		break
			// 	}
			// }

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
