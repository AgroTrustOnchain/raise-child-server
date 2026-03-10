package business

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"raise-child/constants/env"
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
	"strings"
	"time"

	"github.com/block-vision/sui-go-sdk/constant"
	"github.com/block-vision/sui-go-sdk/sui"
)

type pendingWithdrawProposalService struct {
	pendingWithdrawProposalRepo i_repository.IPendingWithdrawProposalRepository
	offWithdrawProposalRepo     i_repository.IOffChainWithdrawProposalRepository
	bankProfileRepo             i_repository.IBankProfileRepository
	redisCache                  cache.IRedisCache
	clients                     map[string]sui.ISuiAPI
	errLogger                   *log.Logger
}

func initializePendingWithdrawProposalService(
	pendingWithdrawProposalRepo i_repository.IPendingWithdrawProposalRepository,
	bankProfileRepo i_repository.IBankProfileRepository,
	clients map[string]sui.ISuiAPI,
	errLogger *log.Logger,
) business.IPendingWithdrawProposalService {
	return &pendingWithdrawProposalService{
		pendingWithdrawProposalRepo: pendingWithdrawProposalRepo,
		bankProfileRepo:             bankProfileRepo,
		redisCache:                  cache.InitializeRedisCache(),
		clients:                     clients,
		errLogger:                   errLogger,
	}
}

func GeneratePendingWithdrawProposalService() (business.IPendingWithdrawProposalService, error) {
	var errLogger = util.GetLogConfig(shared.ERROR_LEVEL)

	cnn, err := db.ConnectDB(errLogger, db.InitializePostgreSQL())
	if err != nil {
		return nil, err
	}

	return initializePendingWithdrawProposalService(repository.InitializePendingWithdrawProposalRepo(cnn, errLogger), repository.InitializeBankProfileRepository(cnn, errLogger), _networkAliases, errLogger), nil
}

// ApprovePendingWithdrawProposal implements business.IPendingWithdrawProposalService.
func (p *pendingWithdrawProposalService) ApprovePendingWithdrawProposal(id string, ctx context.Context) (response.BuildTransactionResponse, error) {
	proposal, err := p.pendingWithdrawProposalRepo.GetPendingWithdrawProposal(id, ctx)
	if err != nil {
		return response.BuildTransactionResponse{}, err
	}

	if proposal == nil {
		return response.BuildTransactionResponse{}, errors.New(noti.GENERIC_ERROR_WARN_MSG)
	}

	if proposal.ReviewedBy != nil || proposal.Status != request_pending_status {
		return response.BuildTransactionResponse{}, errors.New(noti.REQUEST_REVIEWED_MESSAGE)
	}

	var client = p.clients[constant.SuiTestnet]
	var reviewer string = ctx.Value("address").(string)
	manageObj, err := on_chain.GetOnChainObject[entities.Manage](on_chain.GetOnChainObjectRequest{
		Client:    client,
		ObjectId:  os.Getenv(env.PACKAGE_ID),
		ErrLogger: p.errLogger,
	}, ctx)
	if err != nil {
		return response.BuildTransactionResponse{}, err
	}

	if !slices.Contains(manageObj.AdminIds, reviewer) {
		return response.BuildTransactionResponse{}, errors.New(noti.GENERIC_RIGHT_ACCESS_WARN_MSG)
	}

	proposal.ReviewedBy = &reviewer
	proposal.Status = request_approved_status
	if err := p.pendingWithdrawProposalRepo.UpdatePendingWithdrawProposal(*proposal, ctx); err != nil {
		return response.BuildTransactionResponse{}, err
	}

	var proposalid string = util.GenerateId()

	var isFromLocalPool bool = proposal.PoolID != os.Getenv(env.POOL_ID)
	var localPoolId string
	if isFromLocalPool {
		localPoolId = proposal.PoolID
	} else {
		localPoolId = os.Getenv(env.SHARED_LOCAL_POOL_ID)
	}

	var module = on_chain.InitializeModulePool()
	txBytes, err := on_chain.BuildTransaction(on_chain.BuildTransactionRequest{
		Client:    p.clients[constant.SuiTestnet],
		Sender:    reviewer,
		Module:    module.GetModule(),
		Function:  module.GetFunctionCreateWithdrawProposalV2(),
		ErrLogger: p.errLogger,
		Arguments: module.ToCreateWithdrawProposalV2Arguments(on_chain.CreateWithdrawProposalV2Arguments{
			LocalPoolId:     localPoolId,
			WithdrawAmount:  proposal.WithdrawAmount,
			Description:     proposal.Description,
			IsFromLocalPool: isFromLocalPool,
			ClosedAt:        util.ToMilliseconds(util.GetRequestDuration()),
			Creator:         proposal.Creator,
		}),
	}, ctx)
	if err != nil {
		return response.BuildTransactionResponse{}, err
	}

	return response.BuildTransactionResponse{
			TxBytes:    txBytes,
			ProposalId: proposalid,
		}, p.offWithdrawProposalRepo.CreateOffChainWithdrawProposal(entities.OffChainWithdrawProposal{
			ID:        proposalid,
			Purpose:   proposal.Purpose,
			Target:    proposal.PoolID,
			CreatedAt: time.Now(),
		}, ctx)
}

// CreatePendingWithdrawProposal implements business.IPendingWithdrawProposalService.
func (p *pendingWithdrawProposalService) CreatePendingWithdrawProposal(req request.CreatePendingWithdrawProposalRequest, ctx context.Context) (*entities.PendingWithdrawProposal, error) {
	var client = p.clients[constant.SuiTestnet]
	manageObj, err := on_chain.GetOnChainObject[entities.Manage](on_chain.GetOnChainObjectRequest{
		Client:    client,
		ObjectId:  os.Getenv(env.MANAGE_OBJECT_ID),
		ErrLogger: p.errLogger,
	}, ctx)
	if err != nil {
		return nil, err
	}

	var sender string = ctx.Value("address").(string)
	var isLeader bool = slices.Contains(manageObj.LocalLeaderIds, sender)
	var genericRightErr error = errors.New(noti.GENERIC_RIGHT_ACCESS_WARN_MSG)
	if !slices.Contains(manageObj.AdminIds, sender) && !isLeader {
		return nil, genericRightErr
	}

	var mainPoolId string = os.Getenv(env.POOL_ID)
	var reqPoolId string = strings.TrimSpace(req.PoolID)
	var isMainPoolRequested bool = mainPoolId == reqPoolId
	var poolName string
	if isMainPoolRequested {
		if isLeader {
			return nil, genericRightErr
		}

		poolName = "Main Pool"
	} else {
		localPool, err := on_chain.GetOnChainObject[entities.LocalPool](on_chain.GetOnChainObjectRequest{
			Client:    client,
			ObjectId:  reqPoolId,
			ErrLogger: p.errLogger,
		}, ctx)
		if err != nil {
			return nil, err
		}

		if localPool == nil {
			return nil, errors.New(noti.GENERIC_ERROR_WARN_MSG)
		}

		if isLeader {
			if !slices.Contains(localPool.Mods, sender) {
				return nil, genericRightErr
			}

			bankProfile, err := p.bankProfileRepo.GetBankProfileByOwner(sender, ctx)
			if err != nil {
				return nil, err
			}

			if bankProfile == nil {
				return nil, errors.New(noti.LEADER_NOT_UPLOAD_BANK_PROFILE_MESSAGE)
			}
		}

		poolName = localPool.Region
	}

	// todo: AI validation
	var curTime time.Time = time.Now()
	var proposal = entities.PendingWithdrawProposal{
		ID:             util.GenerateId(),
		ProfileID:      ctx.Value("sub").(string),
		Creator:        sender,
		PoolID:         reqPoolId,
		PoolName:       poolName,
		Purpose:        string(entities.WITHDRAW_PURPOSE),
		Target:         reqPoolId,
		WithdrawAmount: req.WithdrawAmount,
		ProofBlobID:    req.ProofBlobID,
		Description:    strings.TrimSpace(req.Description),
		Status:         request_pending_status,
		AIEvaluation:   "",
		CreatedAt:      curTime,
		UpdatedAt:      curTime,
	}

	return &proposal, p.pendingWithdrawProposalRepo.CreatePendingWithdrawProposal(proposal, ctx)
}

// GetPendingWithdrawProposal implements business.IPendingWithdrawProposalService.
func (p *pendingWithdrawProposalService) GetPendingWithdrawProposal(id string, ctx context.Context) (*entities.PendingWithdrawProposal, error) {
	res, err := p.pendingWithdrawProposalRepo.GetPendingWithdrawProposal(id, ctx)
	if err != nil {
		return nil, err
	}

	if res == nil {
		return nil, errors.New(noti.GENERIC_ERROR_WARN_MSG)
	}

	var sender string = ctx.Value("address").(string)
	if res.Creator != sender {
		manageObj, err := on_chain.GetOnChainObject[entities.Manage](on_chain.GetOnChainObjectRequest{
			Client:    p.clients[constant.SuiTestnet],
			ObjectId:  os.Getenv(env.PACKAGE_ID),
			ErrLogger: p.errLogger,
		}, ctx)
		if err != nil {
			return nil, err
		}

		if !slices.Contains(manageObj.AdminIds, sender) {
			return nil, errors.New(noti.GENERIC_RIGHT_ACCESS_WARN_MSG)
		}
	}

	return res, nil
}

// GetPendingWithdrawProposals implements business.IPendingWithdrawProposalService.
func (p *pendingWithdrawProposalService) GetPendingWithdrawProposals(req request.GetPendingWithdrawProposalsRequest, ctx context.Context) (response.PaginationDataResponse, error) {
	var genericErr error = errors.New(noti.GENERIC_ERROR_WARN_MSG)

	req.Creator = strings.TrimSpace(req.Creator)
	if req.Creator != "" {
		if !util.IsValidSuiAddressStrict(req.Creator) {
			return response.PaginationDataResponse{}, genericErr
		}
	}

	manageObj, err := on_chain.GetOnChainObject[entities.Manage](on_chain.GetOnChainObjectRequest{
		Client:    p.clients[constant.SuiTestnet],
		ObjectId:  os.Getenv(env.PACKAGE_ID),
		ErrLogger: p.errLogger,
	}, ctx)
	if err != nil {
		return response.PaginationDataResponse{}, err
	}

	var sender string = ctx.Value("address").(string)
	if !slices.Contains(manageObj.AdminIds, sender) {
		if sender != req.Creator {
			return response.PaginationDataResponse{}, errors.New(noti.GENERIC_RIGHT_ACCESS_WARN_MSG)
		}
	}

	req.Reviewer = strings.TrimSpace(req.Reviewer)
	if req.Reviewer != "" {
		if !util.IsValidSuiAddressStrict(req.Reviewer) {
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

	req.SortOrder = util.StanderizeSortOrder(req.SortOrder)
	req.Keyword = util.StanderizeString(req.Keyword)
	if req.Page < 1 {
		req.Page = 1
	}

	if req.PageSize < 1 {
		req.PageSize = default_page_size
	}

	var res response.PaginationDataResponse
	var redisKey string = p.getGetPendingWithdrawProposalsRedisKey(req)
	if p.redisCache.Get(redisKey, &res, ctx) {
		return res, nil
	}

	data, pages, err := p.pendingWithdrawProposalRepo.GetPendingWithdrawProposals(req, ctx)
	if err != nil {
		return response.PaginationDataResponse{}, err
	}

	var amount int
	if data == nil || len(data) == 0 {
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

	p.redisCache.Set(redisKey, res, time.Minute, ctx)

	return res, nil
}

// RefusePendingWithdrawProposal implements business.IPendingWithdrawProposalService.
func (p *pendingWithdrawProposalService) RefusePendingWithdrawProposal(id string, ctx context.Context) error {
	proposal, err := p.pendingWithdrawProposalRepo.GetPendingWithdrawProposal(id, ctx)
	if err != nil {
		return err
	}

	if proposal == nil {
		return errors.New(noti.GENERIC_ERROR_WARN_MSG)
	}

	if proposal.ReviewedBy != nil || proposal.Status != request_pending_status {
		return errors.New(noti.REQUEST_REVIEWED_MESSAGE)
	}

	var client = p.clients[constant.SuiTestnet]
	var reviewer string = ctx.Value("address").(string)
	manageObj, err := on_chain.GetOnChainObject[entities.Manage](on_chain.GetOnChainObjectRequest{
		Client:    client,
		ObjectId:  os.Getenv(env.PACKAGE_ID),
		ErrLogger: p.errLogger,
	}, ctx)
	if err != nil {
		return err
	}

	if !slices.Contains(manageObj.AdminIds, reviewer) {
		return errors.New(noti.GENERIC_RIGHT_ACCESS_WARN_MSG)
	}

	proposal.ReviewedBy = &reviewer
	proposal.Status = request_refused_status
	return p.pendingWithdrawProposalRepo.UpdatePendingWithdrawProposal(*proposal, ctx)
}

func (p *pendingWithdrawProposalService) getGetPendingWithdrawProposalsRedisKey(req request.GetPendingWithdrawProposalsRequest) string {
	var keyword string = "empty"
	if req.Keyword != "" {
		keyword = req.Keyword
	}

	var creator string = "empty"
	if req.Creator != "" {
		creator = req.Creator
	}

	var reviewer string = "empty"
	if req.Reviewer != "" {
		reviewer = req.Reviewer
	}

	var minAmount string = "empty"
	if req.MinAmount != nil {
		minAmount = fmt.Sprintf("%d", *req.MinAmount)
	}

	var maxAmount string = "empty"
	if req.MaxAmount != nil {
		maxAmount = fmt.Sprintf("%d", *req.MaxAmount)
	}

	var status string = "empty"
	if req.Status != "" {
		status = req.Status
	}

	var sortCriteria string = "empty"
	if req.SortCriteria != "" {
		sortCriteria = req.SortCriteria
	}

	return fmt.Sprintf("pending_withdraw_proposal:kw:%s:of:%s:reviewed:%s:min:%s:max:%s:status:%s:sc:%s:o:%s:s:%d:p:%d",
		keyword, creator, reviewer, minAmount, maxAmount, status, sortCriteria, req.SortOrder, req.PageSize, req.Page)
}
