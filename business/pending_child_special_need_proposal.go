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
	"time"

	"github.com/block-vision/sui-go-sdk/constant"
	"github.com/block-vision/sui-go-sdk/sui"
)

type pendingChildSpecialNeedProposalService struct {
	pendingChildSpecialNeedProposalRepo i_repository.IPendingChildSpecialNeedProposalRepository
	redisCache                          cache.IRedisCache
	clients                             map[string]sui.ISuiAPI
	errLogger                           *log.Logger
}

func initializePendingChildSpecialNeedProposalService(
	pendingChildSpecialNeedProposalRepo i_repository.IPendingChildSpecialNeedProposalRepository,
	clients map[string]sui.ISuiAPI,
	errLogger *log.Logger,
) business.IPendingChildSpecialNeedProposalService {
	return &pendingChildSpecialNeedProposalService{
		pendingChildSpecialNeedProposalRepo: pendingChildSpecialNeedProposalRepo,
		redisCache:                          cache.InitializeRedisCache(),
		clients:                             clients,
		errLogger:                           errLogger,
	}
}

func GeneratePendingChildSpecialNeedProposalService() (business.IPendingChildSpecialNeedProposalService, error) {
	var errLogger = util.GetLogConfig(shared.ERROR_LEVEL)

	cnn, err := db.ConnectDB(errLogger, db.InitializePostgreSQL())
	if err != nil {
		return nil, err
	}

	return initializePendingChildSpecialNeedProposalService(
		repository.InitializePendingChildSpecialNeedProposalRepo(cnn, errLogger),
		_networkAliases,
		errLogger,
	), nil
}

// ApprovePendingChildSpecialNeedProposal implements business.IPendingChildSpecialNeedProposalService.
func (p *pendingChildSpecialNeedProposalService) ApprovePendingChildSpecialNeedProposal(id string, ctx context.Context) error {
	proposal, err := p.pendingChildSpecialNeedProposalRepo.GetPendingChildSpecialNeedProposal(id, ctx)
	if err != nil {
		return err
	}

	if proposal == nil {
		return errors.New(noti.GENERIC_ERROR_WARN_MSG)
	}

	if proposal.ReviewedBy != nil || proposal.ReviewStatus != request_pending_status {
		return errors.New(noti.REQUEST_REVIEWED_MESSAGE)
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

	var reviewer string = ctx.Value("address").(string)
	if !slices.Contains(manage.AdminIds, reviewer) {
		return errors.New(noti.GENERIC_RIGHT_ACCESS_WARN_MSG)
	}

	var pool *entities.MainPool
	var mainPoolId string = os.Getenv(env.POOL_ID)
	pool, err = on_chain.GetOnChainObject[entities.MainPool](on_chain.GetOnChainObjectRequest{
		Client:    client,
		ObjectId:  mainPoolId,
		ErrLogger: p.errLogger,
	}, ctx)
	if err != nil {
		return err
	}

	if pool == nil {
		return errors.New(noti.INTERNALL_ERR_MSG)
	}

	localPools, err := on_chain.GetOnChainObjects[entities.LocalPool](on_chain.GetOnChainObjectsRequest{
		Client:    client,
		ObjectIds: pool.LocalPools,
		ErrLogger: p.errLogger,
	}, ctx)
	if err != nil {
		return err
	}

	var localPoolId string
	for _, localPool := range localPools {
		if localPool.Region == proposal.Region {
			localPoolId = localPool.ID.ID
			break
		}
	}

	proposal.ReviewedBy = &reviewer
	proposal.ReviewStatus = request_approved_status
	proposal.UpdatedAt = time.Now()
	if err := p.pendingChildSpecialNeedProposalRepo.UpdatePendingChildSpecialNeedProposal(*proposal, ctx); err != nil {
		return err
	}

	var childModule = on_chain.InitializeModuleChild()
	var req = on_chain.ExecuteTransactionRequestV2{
		Client:    client,
		Module:    childModule.GetModule(),
		Function:  childModule.GetFunctionCreateChildSpecialNeedProposalV2(),
		ErrLogger: p.errLogger,
		Arguments: childModule.ToCreateChildSpecialNeedProposalArgumentsV2(on_chain.CreateChildSpecialNeedProposalArgumentsV2{
			ChildID:     proposal.ChildID,
			LocalPool:   localPoolId,
			Target:      proposal.Target,
			Description: proposal.Description,
			ProofBlobID: proposal.ProofBlobID,
			ClosedAt:    util.ToMilliseconds(util.GetRequestDuration()),
			Creator:     proposal.ActorAddress,
			Sender:      reviewer,
		}),
	}

	for i := 1; i <= 3; i++ {
		if _, err := on_chain.ExecuteTransactionV2(req, ctx); err == nil {
			return nil
		}
	}

	return internalErr
}

// GetPendingChildSpecialNeedProposal implements business.IPendingChildSpecialNeedProposalService.
func (p *pendingChildSpecialNeedProposalService) GetPendingChildSpecialNeedProposal(id string, ctx context.Context) (*entities.PendingChildSpecialNeedProposal, error) {
	res, err := p.pendingChildSpecialNeedProposalRepo.GetPendingChildSpecialNeedProposal(id, ctx)
	if err != nil {
		return nil, err
	}

	if res == nil {
		return nil, errors.New(noti.GENERIC_ERROR_WARN_MSG)
	}

	var sender string = ctx.Value("address").(string)
	if res.ActorAddress != sender {
		manageObj, err := on_chain.GetOnChainObject[entities.Manage](on_chain.GetOnChainObjectRequest{
			Client:    p.clients[constant.SuiTestnet],
			ObjectId:  os.Getenv(env.MANAGE_OBJECT_ID),
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

// GetPendingChildSpecialNeedProposals implements business.IPendingChildSpecialNeedProposalService.
func (p *pendingChildSpecialNeedProposalService) GetPendingChildSpecialNeedProposals(req request.GetPendingChildSpecialNeedProposalsRequest, ctx context.Context) (response.PaginationDataResponse, error) {
	var genericErr error = errors.New(noti.GENERIC_ERROR_WARN_MSG)

	if req.Creator != "" {
		if !util.IsValidSuiAddressStrict(req.Creator) {
			return response.PaginationDataResponse{}, genericErr
		}
	}

	if req.Reviewer != "" {
		if !util.IsValidSuiAddressStrict(req.Reviewer) {
			return response.PaginationDataResponse{}, genericErr
		}
	}

	var sender string = ctx.Value("address").(string)
	if sender != req.Creator {
		manageObj, err := on_chain.GetOnChainObject[entities.Manage](on_chain.GetOnChainObjectRequest{
			Client:    p.clients[constant.SuiTestnet],
			ObjectId:  os.Getenv(env.MANAGE_OBJECT_ID),
			ErrLogger: p.errLogger,
		}, ctx)
		if err != nil {
			return response.PaginationDataResponse{}, err
		}

		if !slices.Contains(manageObj.AdminIds, sender) {
			if sender != req.Creator {
				return response.PaginationDataResponse{}, errors.New(noti.GENERIC_RIGHT_ACCESS_WARN_MSG)
			}
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
	req.SortCriteria = util.StandardizeSortCriteria(req.SortCriteria)
	if req.Page < 1 {
		req.Page = 1
	}

	if req.PageSize < 1 {
		req.PageSize = default_page_size
	}

	var res response.PaginationDataResponse
	// var redisKey string = p.getGetPendingChildSpecialNeedProposalsRedisKey(req)
	// if p.redisCache.Get(redisKey, &res, ctx) {
	// 	return res, nil
	// }

	data, pages, err := p.pendingChildSpecialNeedProposalRepo.GetPendingChildSpecialNeedProposals(req, ctx)
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

// RefusePendingChildSpecialNeedProposal implements business.IPendingChildSpecialNeedProposalService.
func (p *pendingChildSpecialNeedProposalService) RefusePendingChildSpecialNeedProposal(id string, ctx context.Context) error {
	proposal, err := p.pendingChildSpecialNeedProposalRepo.GetPendingChildSpecialNeedProposal(id, ctx)
	if err != nil {
		return err
	}

	if proposal == nil {
		return errors.New(noti.GENERIC_ERROR_WARN_MSG)
	}

	if proposal.ReviewedBy != nil || proposal.ReviewStatus != request_pending_status {
		return errors.New(noti.REQUEST_REVIEWED_MESSAGE)
	}

	var client = p.clients[constant.SuiTestnet]
	var reviewer string = ctx.Value("address").(string)
	manageObj, err := on_chain.GetOnChainObject[entities.Manage](on_chain.GetOnChainObjectRequest{
		Client:    client,
		ObjectId:  os.Getenv(env.MANAGE_OBJECT_ID),
		ErrLogger: p.errLogger,
	}, ctx)
	if err != nil {
		return err
	}

	if !slices.Contains(manageObj.AdminIds, reviewer) {
		return errors.New(noti.GENERIC_RIGHT_ACCESS_WARN_MSG)
	}

	proposal.ReviewedBy = &reviewer
	proposal.ReviewStatus = request_refused_status
	return p.pendingChildSpecialNeedProposalRepo.UpdatePendingChildSpecialNeedProposal(*proposal, ctx)
}

func (p *pendingChildSpecialNeedProposalService) getGetPendingChildSpecialNeedProposalsRedisKey(req request.GetPendingChildSpecialNeedProposalsRequest) string {
	var keyword string = "empty"
	if req.Keyword != "" {
		keyword = req.Keyword
	}

	var region string = "empty"
	if req.Region != "" {
		region = req.Region
	}

	var status string = "empty"
	if req.Status != "" {
		status = req.Status
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

	var sortCriteria string = "empty"
	if req.SortCriteria != "" {
		sortCriteria = req.SortCriteria
	}

	return fmt.Sprintf("pending_special_need_proposal:kw:%s:r:%s:status:%s:of:%s:reviewed:%s:min:%s:max:%s:sc:%s:o:%s:s:%d:p:%d",
		keyword, region, status, creator, reviewer, minAmount, maxAmount, sortCriteria, req.SortOrder, req.PageSize, req.Page)
}
