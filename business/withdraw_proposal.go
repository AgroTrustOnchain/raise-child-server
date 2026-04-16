package business

import (
	"context"
	"errors"
	"fmt"
	"log"
	"math"
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
	"raise-child/util/cache"
	"raise-child/util/db"
	on_chain "raise-child/util/on_chain"
	"slices"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/block-vision/sui-go-sdk/constant"
	"github.com/block-vision/sui-go-sdk/sui"
	"github.com/payOSHQ/payos-lib-golang"
)

type withdrawProposalService struct {
	paymentRepo     i_repository.IPaymentRepository
	bankProfileRepo i_repository.IBankProfileRepository
	withdrawRepo    i_repository.IOffChainWithdrawProposalRepository
	redisCache      cache.IRedisCache
	clients         map[string]sui.ISuiAPI
	errLogger       *log.Logger
}

func initializeWithdrawProposalService(
	paymentRepo i_repository.IPaymentRepository,
	bankProfileRepo i_repository.IBankProfileRepository,
	withdrawRepo i_repository.IOffChainWithdrawProposalRepository,
	clients map[string]sui.ISuiAPI,
	errLogger *log.Logger,
) business.IWithdrawProposalService {
	return &withdrawProposalService{
		paymentRepo:     paymentRepo,
		bankProfileRepo: bankProfileRepo,
		withdrawRepo:    withdrawRepo,
		redisCache:      cache.InitializeRedisCache(),
		clients:         clients,
		errLogger:       errLogger,
	}
}

func GenerateWithdrawProposalService() (business.IWithdrawProposalService, error) {
	var errLogger = util.GetLogConfig(shared.ERROR_LEVEL)

	cnn, err := db.ConnectDB(errLogger, db.InitializePostgreSQL())
	if err != nil {
		return nil, err
	}

	return initializeWithdrawProposalService(
		repository.InitializePaymentRepository(cnn, errLogger),
		repository.InitializeBankProfileRepository(cnn, errLogger),
		repository.InitializeOffChainWithdrawProposalRepository(cnn, errLogger),
		_networkAliases,
		errLogger,
	), nil
}

const (
	withdraw_proposal_records_limit    int     = 10
	min_withdraw_proposal_amount_value int64   = 10_000
	min_withdraw_proposal_approve_rate float32 = 0.8
)

// CreateWithdrawProposal implements business.IWithdrawProposalService.
func (w *withdrawProposalService) CreateWithdrawProposal(req request.CreateWithdrawProposalRequest, ctx context.Context) (response.BuildTransactionResponse, error) {
	// if !utils.IsValidSuiAddress(models.SuiAddress(req.PoolID)) || !utils.IsValidSuiAddress(models.SuiAddress(sender)) {
	// 	return response.BuildTransactionResponse{}, genericErr
	// }

	var client = w.clients[constant.SuiTestnet]
	manageObj, err := on_chain.GetOnChainObject[entities.Manage](on_chain.GetOnChainObjectRequest{
		Client:    client,
		ObjectId:  os.Getenv(env.MANAGE_OBJECT_ID),
		ErrLogger: w.errLogger,
	}, ctx)
	if err != nil {
		return response.BuildTransactionResponse{}, err
	}

	var sender string = ctx.Value("address").(string)
	var isAdmin bool = slices.Contains(manageObj.AdminIds, sender)
	var isLeader bool = slices.Contains(manageObj.LocalLeaderIds, sender)
	var genericRightErr error = errors.New(noti.GENERIC_RIGHT_ACCESS_WARN_MSG)
	if !isAdmin && !isLeader {
		return response.BuildTransactionResponse{}, genericRightErr
	}

	var mainPoolId string = os.Getenv(env.POOL_ID)
	var reqPoolId string = strings.TrimSpace(req.PoolID)
	var isMainPoolRequested bool = mainPoolId == reqPoolId
	if isLeader {
		if isMainPoolRequested {
			return response.BuildTransactionResponse{}, genericRightErr
		}
	}

	var poolNotEnoughBalenceErr error = errors.New(noti.POOL_CURRENTLY_NOT_ENOUGH_BALENCE)
	var internalErr error = errors.New(noti.INTERNALL_ERR_MSG)
	var localPoolId string
	if !isMainPoolRequested {
		localPool, _ := on_chain.GetOnChainObject[entities.LocalPool](on_chain.GetOnChainObjectRequest{
			Client:    client,
			ObjectId:  reqPoolId,
			ErrLogger: w.errLogger,
		}, ctx)
		if localPool == nil {
			return response.BuildTransactionResponse{}, internalErr
		}

		if isLeader {
			if !slices.Contains(localPool.Mods, sender) {
				return response.BuildTransactionResponse{}, genericRightErr
			}
		}

		bankProfile, err := w.bankProfileRepo.GetBankProfileByOwner(sender, ctx)
		if err != nil {
			return response.BuildTransactionResponse{}, err
		}

		if bankProfile == nil {
			return response.BuildTransactionResponse{}, errors.New(noti.LEADER_NOT_UPLOAD_BANK_PROFILE_MESSAGE)
		}

		totalAmount, _ := strconv.ParseInt(localPool.TotalAmount, 10, 64)
		if totalAmount < req.WithdrawAmount {
			return response.BuildTransactionResponse{}, poolNotEnoughBalenceErr
		}

		localPoolId = reqPoolId
	} else {
		mainPool, _ := on_chain.GetOnChainObject[entities.MainPool](on_chain.GetOnChainObjectRequest{
			Client:    client,
			ObjectId:  mainPoolId,
			ErrLogger: w.errLogger,
		}, ctx)
		if mainPool == nil {
			return response.BuildTransactionResponse{}, internalErr
		}

		localPools, _ := on_chain.GetOnChainObjects[entities.LocalPool](on_chain.GetOnChainObjectsRequest{
			Client:    client,
			ObjectIds: mainPool.LocalPools,
			ErrLogger: w.errLogger,
		}, ctx)
		if localPools == nil || len(localPools) == 0 {
			return response.BuildTransactionResponse{}, internalErr
		}

		var localPoolsTotalAmount int64
		for _, localPool := range localPools {
			totalAmount, _ := strconv.ParseInt(localPool.TotalAmount, 10, 64)
			localPoolsTotalAmount += totalAmount
		}

		totalAmount, _ := strconv.ParseInt(mainPool.TotalAmount, 10, 64)
		var mainPoolAmount int64 = totalAmount - localPoolsTotalAmount
		if mainPoolAmount < req.WithdrawAmount {
			return response.BuildTransactionResponse{}, poolNotEnoughBalenceErr
		}

		localPoolId = mainPool.LocalPools[0]
	}

	var description string = strings.TrimSpace(req.Description)
	if description == "" {
		description = "Withdraw"
	}

	var module = on_chain.InitializeModulePool()
	txBytes, err := on_chain.BuildTransaction(on_chain.BuildTransactionRequest{
		Client:    client,
		Sender:    sender,
		Module:    module.GetModule(),
		Function:  module.GetFunctionCreateWithdrawProposal(),
		ErrLogger: w.errLogger,
		Arguments: module.ToCreateWithdrawProposalArguments(on_chain.CreateWithdrawProposalArguments{
			LocalPoolId:     localPoolId,
			WithdrawAmount:  req.WithdrawAmount,
			Description:     description,
			IsFromLocalPool: reqPoolId != mainPoolId,
			ClosedAt:        util.ToMilliseconds(util.GetRequestDuration()),
		}),
	}, ctx)
	if err != nil {
		return response.BuildTransactionResponse{}, err
	}

	var proposalId string = util.GenerateId()
	return response.BuildTransactionResponse{
			TxBytes:    txBytes,
			ProposalId: proposalId,
		}, w.withdrawRepo.CreateOffChainWithdrawProposal(entities.OffChainWithdrawProposal{
			ID:        proposalId,
			Purpose:   string(entities.WITHDRAW_PURPOSE),
			Target:    req.PoolID,
			CreatedAt: time.Now(),
		}, ctx)
}

// ConfirmWithdrawProposal implements business.IWithdrawProposalService.
func (w *withdrawProposalService) ConfirmWithdrawProposal(id string, ctx context.Context) (map[string]interface{}, error) {
	var genericErr error = errors.New(noti.GENERIC_ERROR_WARN_MSG)

	var sender string = ctx.Value("address").(string)
	if !util.IsValidSuiAddressStrict(id) {
		return nil, genericErr
	}

	var client = w.clients[constant.SuiTestnet]
	var manageModule = on_chain.InitializeModuleManage()
	if nfts, err := on_chain.GetOnChainOwnedObjects[entities.AdminNft](on_chain.GetOnChainOwnedObjectsRequest{
		Client:       client,
		OwnerAddress: sender,
		StructType:   fmt.Sprintf("%s::%s::%s", os.Getenv(env.PACKAGE_ID), manageModule.GetModule(), manageModule.GetAdminNftStruct()),
		ErrLogger:    w.errLogger,
	}, ctx); err != nil {
		return nil, err
	} else {
		if nfts == nil || len(nfts) == 0 {
			return nil, errors.New(noti.GENERIC_RIGHT_ACCESS_WARN_MSG)
		}
	}

	// var internalErr error = errors.New(noti.INTERNALL_ERR_MSG)
	// manageObj, _ := on_chain.GetOnChainObject[entities.Manage](on_chain.GetOnChainObjectRequest{
	// 	Client:    client,
	// 	ObjectId:  os.Getenv(env.MANAGE_OBJECT_ID),
	// 	ErrLogger: w.errLogger,
	// }, ctx)
	// if manageObj == nil {
	// 	return nil, internalErr
	// }

	// if !slices.Contains(manageObj.AdminIds, sender) {
	// 	return nil, errors.New(noti.GENERIC_RIGHT_ACCESS_WARN_MSG)
	// }

	proposal, _ := on_chain.GetOnChainObject[entities.WithdrawProposal](on_chain.GetOnChainObjectRequest{
		Client:    client,
		ObjectId:  id,
		ErrLogger: w.errLogger,
	}, ctx)
	if proposal == nil {
		return nil, genericErr
	}

	closedAt, _ := strconv.ParseInt(proposal.ClosedAt, 10, 64)
	if time.Now().Before(util.MilliSecToTime(closedAt)) {
		return nil, errors.New(noti.STILL_PENDING_REQUEST_MESSAGE)
	}

	if proposal.IsExecuted {
		return nil, errors.New(noti.WITHDRAW_PROPOSAL_EXECUTED_MESSSAGE)
	}

	if !proposal.IsFromLocalPool {
		return nil, genericErr
	}

	approveWeight, _ := strconv.ParseInt(proposal.ApproveWeight, 10, 64)
	refuseWeight, _ := strconv.ParseInt(proposal.RefuseWeight, 10, 64)
	if approveWeight/(approveWeight+refuseWeight) < min_withdraw_proposal_amount_value {
		return nil, errors.New(noti.WITHDRAW_PROPOSAL_FAIL_CONDITION_MESSAGE)
	}

	offChainProposal, err := w.withdrawRepo.GetOffChainWithdrawProposalByProposal(id, ctx)
	if err != nil {
		return nil, err
	}

	isProcessed, err := w.paymentRepo.IsWithdrawalPaymentInProcess(offChainProposal.ID, ctx)
	if err != nil {
		return nil, err
	}

	// Withdraw Proposal is being executed or been success
	if isProcessed {
		return nil, errors.New(noti.WITHDRAW_PROPOSAL_IN_PROCESS_MESSAGE)
	}

	// pool, err := on_chain.GetOnChainObject[entities.MainPool](on_chain.GetOnChainObjectRequest{
	// 	Client:    client,
	// 	ObjectId:  os.Getenv(env.POOL_ID),
	// 	ErrLogger: w.errLogger,
	// }, ctx)
	// if err != nil {
	// 	return nil, err
	// }

	// localPools, err := on_chain.GetOnChainObjects[entities.LocalPool](on_chain.GetOnChainObjectsRequest{
	// 	Client:    client,
	// 	ObjectIds: pool.LocalPools,
	// 	ErrLogger: w.errLogger,
	// }, ctx)
	// if err != nil {
	// 	return nil, err
	// }

	// var mods []string
	// for _, localPool := range localPools {
	// 	if localPool.Region == proposal.PoolName {
	// 		mods = localPool.Mods
	// 		break
	// 	}
	// }
	// var profileOwner string
	// for i := 0; i < len(manageObj.LocalRegions); i++ {
	// 	if proposal.PoolName == manageObj.LocalRegions[i] {
	// 		profileOwner = manageObj.LocalLeaderIds[i]
	// 		break
	// 	}
	// }

	// bankProfile, err := w.bankProfileRepo.GetBankProfileByOwner(profileOwner, ctx)
	// if err != nil {
	// 	return nil, err
	// }

	bankProfile, err := w.bankProfileRepo.GetBankProfileByOwner(proposal.Creator, ctx)
	if err != nil {
		return nil, err
	}

	var paymentId string = util.GenerateId()
	var orderCode int = util.GenerateNumber()
	//var callbackUrl string = fmt.Sprintf("%s/%s/%d", os.Getenv(payment.PAYMENT_CALLBACK_URL), paymentId, orderCode)
	var callbackUrl string = os.Getenv(payment.PAYMENT_CALLBACK_URL) + paymentId
	var curTime time.Time = time.Now()
	var expiredAt time.Time
	var paymentMethod string
	withdrawAmount, _ := strconv.ParseInt(proposal.WithdrawAmount, 10, 64)
	var res map[string]interface{} = make(map[string]interface{})
	var isPayosAvailable bool = bankProfile.PayosApiKey != "" && bankProfile.PayosCheckSumKey != "" && bankProfile.PayosClientID != ""

	var paymentDescription string
	switch offChainProposal.Purpose {
	case string(entities.BOOKS_NEED_PURPOSE):
		paymentDescription = entities.BOOKS_NEED_PAYMENT_DESCRIPTION.GenerateWithdrawPaymentDescription()
	case string(entities.MEAL_NEED_PURPOSE):
		paymentDescription = entities.MEAL_NEED_PAYMENT_DESCRIPTION.GenerateWithdrawPaymentDescription()
	case string(entities.HEALTH_INSURANCE_NEED_PURPOSE):
		paymentDescription = entities.HEALTH_INSRUANCE_PAYMENT_DESCRIPTION.GenerateWithdrawPaymentDescription()
	case string(entities.SPECIAL_NEED_PURPOSE):
		paymentDescription = entities.SPECIAL_NEED_CAMPAIGN_PAYMENT_DESCRIPTION.GenerateWithdrawPaymentDescription()
	case string(entities.CAMPAIGN_PURPOSE):
		paymentDescription = entities.POOL_CAMPAIGN_PAYMENT_DESCRIPTION.GenerateWithdrawPaymentDescription()
	}

	if isPayosAvailable {
		if err := payos.Key(util.Decrypt(bankProfile.PayosClientID), util.Decrypt(bankProfile.PayosApiKey), util.Decrypt(bankProfile.PayosCheckSumKey)); err != nil {
			w.errLogger.Println(fmt.Sprintf(noti.PAYMENT_INIT_ENV_ERR_MSG, "payos") + err.Error())
			isPayosAvailable = false
		} else {
			// Set payos back to app default
			defer payos.Key(os.Getenv(payment.PAYOS_CLIENT_ID), os.Getenv(payment.PAYOS_API_KEY), os.Getenv(payment.PAYOS_CHECKSUM_KEY))

			data, err := payos.CreatePaymentLink(payos.CheckoutRequestType{
				OrderCode:   int64(orderCode),
				Amount:      int(withdrawAmount),
				Description: paymentDescription,
				ReturnUrl:   callbackUrl,
				CancelUrl:   callbackUrl,
			})

			if err != nil {
				w.errLogger.Println("Err: ", err.Error())
				return nil, errors.New(noti.INTERNALL_ERR_MSG)
			}

			expiredAt = time.Unix(int64(*data.ExpiredAt), 0)
			paymentMethod = shared.PAYMENT_PAYOS_METHOD
			res["checkout_url"] = data.CheckoutUrl
		}
	}

	if !isPayosAvailable {
		expiredAt = util.GetBankTransactionDuration()
		paymentMethod = shared.MANUAL_BANK_METHOD
		res["owner"] = bankProfile.OwnerName
		res["bank_org"] = bankProfile.BankOrg
		res["bank_code"] = bankProfile.BankCode
		res["amount"] = proposal.WithdrawAmount
		res["payment_id"] = fmt.Sprint(paymentId)
		res["description"] = proposal.Description
	}

	// detail, err := w.withdrawRepo.GetOffChainWithdrawProposalByProposal(id, ctx)
	// if err != nil {
	// 	return nil, err
	// }

	return res, w.paymentRepo.CreatePayment(entities.Payment{
		ID:            paymentId,
		Actor:         sender,
		ProfileID:     ctx.Value("sub").(string),
		ProposalID:    &offChainProposal.ID,
		IsDonateTx:    false,
		TransactionId: fmt.Sprint(orderCode),
		Amount:        withdrawAmount,
		Currency:      shared.VIETNAMDONG_CURRENCY,
		Status:        payment_pending_status,
		Method:        paymentMethod,
		Message:       proposal.Description,
		ExpiredAt:     expiredAt,
		CreatedAt:     curTime,
		UpdatedAt:     curTime,
	}, ctx)
}

// ConfirmMainPoolWithdrawProposal implements business.IWithdrawProposalService.
func (w *withdrawProposalService) ConfirmMainPoolWithdrawProposal(id string, capturedImgBlobId string, ctx context.Context) (response.BuildTransactionResponse, error) {
	var genericErr error = errors.New(noti.GENERIC_ERROR_WARN_MSG)
	var sender string = ctx.Value("address").(string)
	if !util.IsValidSuiAddressStrict(id) || capturedImgBlobId == "" {
		return response.BuildTransactionResponse{}, genericErr
	}

	var client = w.clients[constant.SuiTestnet]
	var internalErr error = errors.New(noti.INTERNALL_ERR_MSG)
	manageObj, _ := on_chain.GetOnChainObject[entities.Manage](on_chain.GetOnChainObjectRequest{
		Client:    client,
		ObjectId:  os.Getenv(env.MANAGE_OBJECT_ID),
		ErrLogger: w.errLogger,
	}, ctx)
	if manageObj == nil {
		return response.BuildTransactionResponse{}, internalErr
	}

	if !slices.Contains(manageObj.AdminIds, sender) {
		return response.BuildTransactionResponse{}, errors.New(noti.GENERIC_RIGHT_ACCESS_WARN_MSG)
	}

	proposal, _ := on_chain.GetOnChainObject[entities.WithdrawProposal](on_chain.GetOnChainObjectRequest{
		Client:    client,
		ObjectId:  id,
		ErrLogger: w.errLogger,
	}, ctx)
	if proposal == nil {
		return response.BuildTransactionResponse{}, genericErr
	}

	closedAt, _ := strconv.ParseInt(proposal.ClosedAt, 10, 64)
	if time.Now().Before(util.MilliSecToTime(closedAt)) {
		return response.BuildTransactionResponse{}, errors.New(noti.STILL_PENDING_REQUEST_MESSAGE)
	}

	if proposal.IsExecuted {
		return response.BuildTransactionResponse{}, errors.New(noti.WITHDRAW_PROPOSAL_EXECUTED_MESSSAGE)
	}

	if proposal.IsFromLocalPool {
		return response.BuildTransactionResponse{}, genericErr
	}

	mainPool, _ := on_chain.GetOnChainObject[entities.MainPool](on_chain.GetOnChainObjectRequest{
		Client:    client,
		ObjectId:  os.Getenv(env.POOL_ID),
		ErrLogger: w.errLogger,
	}, ctx)
	if mainPool == nil {
		return response.BuildTransactionResponse{}, internalErr
	}

	approveWeight, _ := strconv.ParseInt(proposal.ApproveWeight, 10, 64)
	refuseWeight, _ := strconv.ParseInt(proposal.RefuseWeight, 10, 64)
	if approveWeight/(approveWeight+refuseWeight) < min_withdraw_proposal_amount_value {
		return response.BuildTransactionResponse{}, errors.New(noti.WITHDRAW_PROPOSAL_FAIL_CONDITION_MESSAGE)
	}

	// todo: validate transaction capture image blob id

	var module = on_chain.InitializeModulePool()
	txBytes, err := on_chain.BuildTransaction(on_chain.BuildTransactionRequest{
		Client:    client,
		Sender:    sender,
		Module:    module.GetModule(),
		Function:  module.GetFunctionWithdrawFromPool(),
		ErrLogger: w.errLogger,
		Arguments: module.ToWithdrawFromPoolArguments(on_chain.WithdrawFromPoolArguments{
			LocalPoolId:        mainPool.LocalPools[0],
			WithdrawProposalId: id,
		}),
	}, ctx)

	return response.BuildTransactionResponse{
		TxBytes: txBytes,
	}, err
}

// GetWithdrawProposal implements business.IWithdrawProposalService.
func (w *withdrawProposalService) GetWithdrawProposal(id string, ctx context.Context) (response.WithdrawProposalResponse, error) {
	var genericErr error = errors.New(noti.GENERIC_ERROR_WARN_MSG)
	if !!util.IsValidSuiAddressStrict(id) {
		return response.WithdrawProposalResponse{}, genericErr
	}

	res, _ := on_chain.GetOnChainObject[entities.WithdrawProposal](on_chain.GetOnChainObjectRequest{
		Client:    w.clients[constant.SuiTestnet],
		ObjectId:  id,
		ErrLogger: w.errLogger,
	}, ctx)

	return res.ToWithdrawProposalResponse(), genericErr
}

// GetWithdrawProposals implements business.IWithdrawProposalService.
func (w *withdrawProposalService) GetWithdrawProposals(req request.GetWithdrawProposalsRequest, ctx context.Context) (response.PaginationDataResponse, error) {
	var genericErr error = errors.New(noti.GENERIC_ERROR_WARN_MSG)

	var creator string = strings.TrimSpace(req.Creator)
	if creator != "" {
		if !util.IsValidSuiAddressStrict(creator) {
			return response.PaginationDataResponse{}, genericErr
		}
	}

	if req.MaxAmount != nil {
		if *req.MaxAmount < min_withdraw_proposal_amount_value {
			return response.PaginationDataResponse{}, nil
		}

		if req.MinAmount != nil {
			if *req.MaxAmount <= *req.MinAmount {
				return response.PaginationDataResponse{}, nil
			}
		}
	}

	req.SortOrder = util.StandardizeSortOrder(req.SortOrder)
	req.SortCriteria = util.StandardizeSortCriteria(req.SortCriteria)
	req.Keyword = util.StandardizeString(req.Keyword)
	if req.Page < 1 {
		req.Page = 1
	}

	if req.PageSize < 1 {
		req.PageSize = default_page_size
	}

	var res response.PaginationDataResponse
	var redisKey string = w.getGetWithdrawProposalsRedisKey(req)
	if w.redisCache.Get(redisKey, &res, ctx) {
		return res, nil
	}

	var client = w.clients[constant.SuiTestnet]
	pool, _ := on_chain.GetOnChainObject[entities.MainPool](on_chain.GetOnChainObjectRequest{
		Client:    client,
		ObjectId:  os.Getenv(env.POOL_ID),
		ErrLogger: w.errLogger,
	}, ctx)
	if pool == nil {
		return response.PaginationDataResponse{}, nil
	}

	proposals, _ := on_chain.GetOnChainObjects[entities.WithdrawProposal](on_chain.GetOnChainObjectsRequest{
		Client:    client,
		ObjectIds: pool.WithdrawProposals,
		ErrLogger: w.errLogger,
	}, ctx)
	if proposals == nil || len(proposals) == 0 {
		return response.PaginationDataResponse{}, nil
	}

	var filteredProposals []entities.WithdrawProposal
	var curTime time.Time = time.Now()
	for i := len(pool.WithdrawProposals) - 1; i >= 0; i-- {
		var proposal = proposals[i]
		if creator != "" {
			if proposal.Creator != creator { // Not matched
				continue
			}
		}

		if req.Keyword != "" {
			var poolName string = util.StandardizeString(proposal.PoolName)
			var description string = util.StandardizeString(proposal.Description)
			if !strings.Contains(proposal.PoolID, req.Keyword) && !strings.Contains(poolName, req.Keyword) && !strings.Contains(description, req.Keyword) {
				continue
			}
		}

		withdrawAmount, _ := strconv.ParseInt(proposal.WithdrawAmount, 10, 64)
		if req.MinAmount != nil {
			if withdrawAmount < *req.MinAmount {
				continue
			}
		}

		if req.MaxAmount != nil {
			if withdrawAmount > *req.MaxAmount {
				continue
			}
		}

		if req.IsExecuted != nil {
			if proposal.IsExecuted != *req.IsExecuted {
				continue
			}
		}

		if req.IsClosed != nil {
			closedAt, _ := strconv.ParseInt(proposal.ClosedAt, 10, 64)
			var closedPeriod time.Time = util.MilliSecToTime(closedAt)

			if *req.IsClosed {
				if curTime.Before(closedPeriod) {
					continue
				}
			} else {
				if curTime.After(closedPeriod) {
					continue
				}
			}
		}

		filteredProposals = append(filteredProposals, proposal)
	}

	sort.Slice(filteredProposals, func(i, j int) bool {
		if req.SortCriteria == "withdraw_amount" {
			withdrawAmount1, _ := strconv.ParseInt(filteredProposals[i].WithdrawAmount, 10, 64)
			withdrawAmount2, _ := strconv.ParseInt(filteredProposals[j].WithdrawAmount, 10, 64)
			if req.SortOrder == "DESC" {
				return withdrawAmount2 > withdrawAmount1
			}

			return withdrawAmount2 < withdrawAmount1
		} else if req.SortCriteria == "closed_at" {
			closedAt1, _ := strconv.ParseInt(filteredProposals[i].ClosedAt, 10, 64)
			closedAt2, _ := strconv.ParseInt(filteredProposals[j].ClosedAt, 10, 64)
			var closedPeriod1 = util.MilliSecToTime(closedAt1)
			var closedPeriod2 = util.MilliSecToTime(closedAt2)
			if req.SortOrder == "DESC" {
				return closedPeriod2.After(closedPeriod1)
			}

			return closedPeriod2.Before(closedPeriod1)
		}

		if req.SortOrder == "ASC" {
			return false
		}

		return true
	})

	var skippedRecords int = (req.Page - 1) * req.PageSize
	if len(filteredProposals) <= skippedRecords {
		return response.PaginationDataResponse{}, nil
	}

	var data []response.WithdrawProposalResponse
	for i := skippedRecords; i < len(filteredProposals); i++ {
		data = append(data, filteredProposals[i].ToMinimumWithdrawProposalResponse())
		if len(data) == req.PageSize {
			break
		}
	}

	res = response.PaginationDataResponse{
		Data:       data,
		Amount:     len(data),
		Page:       req.Page,
		TotalPages: int(math.Ceil(float64(len(filteredProposals)) / float64(req.PageSize))),
	}

	w.redisCache.Set(redisKey, res, time.Minute*5, ctx)

	return res, nil
}

// VoteWithdrawProposal implements business.IWithdrawProposalService.
func (w *withdrawProposalService) VoteWithdrawProposal(id string, req request.VoteRequest, ctx context.Context) (response.BuildTransactionResponse, error) {
	var genericErr error = errors.New(noti.GENERIC_ERROR_WARN_MSG)
	var sender string = ctx.Value("address").(string)
	if !util.IsValidSuiAddressStrict(id) {
		return response.BuildTransactionResponse{}, genericErr
	}

	// todo: add admin nft and get nft of wallet to check
	var client = w.clients[constant.SuiTestnet]
	proposal, _ := on_chain.GetOnChainObject[entities.WithdrawProposal](on_chain.GetOnChainObjectRequest{
		Client:    client,
		ObjectId:  id,
		ErrLogger: w.errLogger,
	}, ctx)
	if proposal == nil {
		return response.BuildTransactionResponse{}, genericErr
	}

	if !proposal.ToWithdrawProposalResponse().ClosedAt.After(time.Now()) {
		return response.BuildTransactionResponse{}, errors.New(noti.WITHDRAW_PROPOSAL_CLOSED_MESSAGE)
	}

	if slices.Contains(proposal.Approvers, sender) || slices.Contains(proposal.Refusers, sender) {
		return response.BuildTransactionResponse{}, errors.New(noti.ALREADY_VOTE_MESSAGE)
	}

	var donorModule = on_chain.InitializeModuleDonor()
	nfts, _ := on_chain.GetOnChainOwnedObjects[entities.Donor](on_chain.GetOnChainOwnedObjectsRequest{
		Client:       client,
		OwnerAddress: sender,
		StructType:   fmt.Sprintf("%s::%s::%s", os.Getenv(env.PACKAGE_ID), donorModule.GetModule(), donorModule.GetDonorNftStruct()),
		ErrLogger:    w.errLogger,
	}, ctx)
	if nfts == nil || len(nfts) == 0 {
		return response.BuildTransactionResponse{}, errors.New(noti.HAVE_TO_DONATE_TO_VOTE)
	}

	var refuseReason string = strings.TrimSpace(req.RefuseReason)
	if refuseReason == "" {
		refuseReason = "Refuse"
	}

	var poolModule = on_chain.InitializeModulePool()
	txBytes, err := on_chain.BuildTransaction(on_chain.BuildTransactionRequest{
		Client:    client,
		Sender:    sender,
		Module:    poolModule.GetModule(),
		Function:  poolModule.GetFunctionVoteWithdrawProposal(),
		ErrLogger: w.errLogger,
		Arguments: poolModule.ToVoteWithdrawProposalArguments(on_chain.VoteWithdrawProposalArguments{
			ProposalId:   id,
			DonorId:      nfts[0].ID.ID,
			IsApprove:    req.IsVoteYes,
			RefuseReason: refuseReason,
		}),
	}, ctx)

	return response.BuildTransactionResponse{
		TxBytes: txBytes,
	}, err
}

func (w *withdrawProposalService) getGetWithdrawProposalsRedisKey(req request.GetWithdrawProposalsRequest) string {
	var keyword string = "empty"
	if req.Keyword != "" {
		keyword = req.Keyword
	}

	var creator string = "empty"
	if req.Creator != "" {
		creator = req.Creator
	}

	var minAmount string = "empty"
	if req.MinAmount != nil {
		minAmount = fmt.Sprintf("%d", *req.MinAmount)
	}

	var maxAmount string = "empty"
	if req.MaxAmount != nil {
		maxAmount = fmt.Sprintf("%d", *req.MaxAmount)
	}

	var isExecuted string = "empty"
	if req.IsExecuted != nil {
		isExecuted = fmt.Sprintf("%v", *req.IsExecuted)
	}

	var isClosed string = "empty"
	if req.IsClosed != nil {
		isClosed = fmt.Sprintf("%v", *req.IsClosed)
	}

	var sortCriteria string = "empty"
	if req.SortCriteria != "" {
		sortCriteria = req.SortCriteria
	}

	return fmt.Sprintf("off_withdraw_proposal:kw:%s:of:%s:min:%s:max:%s:executed:%s:closed:%s:sc:%s:o:%s:s:%d:p:%d",
		keyword, creator, minAmount, maxAmount, isExecuted, isClosed, sortCriteria, req.SortOrder, req.PageSize, req.Page)
}
