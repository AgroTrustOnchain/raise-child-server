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
	"raise-child/util/ai"
	"raise-child/util/cache"
	"raise-child/util/db"
	on_chain "raise-child/util/on_chain"
	walrus_pkg "raise-child/util/walrus_pkg"
	"slices"
	"strings"
	"time"

	"github.com/block-vision/sui-go-sdk/constant"
	"github.com/block-vision/sui-go-sdk/sui"
)

type pendingCampaignService struct {
	pendingCampaignRepo i_repository.IPendingCampaignRepository
	bankProfileRepo     i_repository.IBankProfileRepository
	aiProvider          ai.IAiClientProvider
	walrusProvider      walrus_pkg.IWalrusProvider
	redisCache          cache.IRedisCache
	clients             map[string]sui.ISuiAPI
	errLogger           *log.Logger
}

func initializePendingCampaignService(
	pendingCampaignRepo i_repository.IPendingCampaignRepository,
	bankProfileRepo i_repository.IBankProfileRepository,
	aiProvider ai.IAiClientProvider,
	walrusProvider walrus_pkg.IWalrusProvider,
	clients map[string]sui.ISuiAPI,
	errLogger *log.Logger,
) business.IPendingCampaignService {
	return &pendingCampaignService{
		pendingCampaignRepo: pendingCampaignRepo,
		bankProfileRepo:     bankProfileRepo,
		aiProvider:          aiProvider,
		walrusProvider:      walrusProvider,
		redisCache:          cache.InitializeRedisCache(),
		clients:             clients,
		errLogger:           errLogger,
	}
}

func GeneratePendingCampaignService() (business.IPendingCampaignService, error) {
	var errLogger = util.GetLogConfig(shared.ERROR_LEVEL)

	cnn, err := db.ConnectDB(errLogger, db.InitializePostgreSQL())
	if err != nil {
		return nil, err
	}

	return initializePendingCampaignService(
		repository.InitializePendingCampaignRepo(cnn, errLogger),
		repository.InitializeBankProfileRepository(cnn, errLogger),
		ai.InitializeAiProvider(errLogger),
		walrus_pkg.InitializeWalrusProvider(errLogger),
		_networkAliases,
		errLogger,
	), nil
}

// // ApprovePendingCampaign implements business.IPendingCampaignService.
// func (p *pendingCampaignService) ApprovePendingCampaign(id string, ctx context.Context) (response.BuildTransactionResponse, error) {
// 	campaign, err := p.pendingCampaignRepo.GetPendingCampaign(id, ctx)
// 	if err != nil {
// 		return response.BuildTransactionResponse{}, err
// 	}

// 	if campaign == nil {
// 		return response.BuildTransactionResponse{}, errors.New(noti.GENERIC_ERROR_WARN_MSG)
// 	}

// 	if campaign.ReviewedBy != nil || campaign.ReviewStatus != request_pending_status {
// 		return response.BuildTransactionResponse{}, errors.New(noti.REQUEST_REVIEWED_MESSAGE)
// 	}

// 	var client = p.clients[constant.SuiTestnet]
// 	var reviewer string = ctx.Value("address").(string)
// 	var manageObj entities.Manage
// 	if !p.redisCache.Get(manageObj.GetRedisKey(), &manageObj, ctx) {
// 		res, err := on_chain.GetOnChainObject[entities.Manage](on_chain.GetOnChainObjectRequest{
// 			Client:    client,
// 			ObjectId:  os.Getenv(env.MANAGE_OBJECT_ID),
// 			ErrLogger: p.errLogger,
// 		}, ctx)
// 		if err != nil {
// 			return response.BuildTransactionResponse{}, err
// 		}

// 		if res != nil {
// 			p.redisCache.Set(manageObj.GetRedisKey(), res, time.Minute, ctx)
// 			manageObj = *res
// 		}
// 	}

// 	if !slices.Contains(manageObj.AdminIds, reviewer) {
// 		return response.BuildTransactionResponse{}, errors.New(noti.GENERIC_RIGHT_ACCESS_WARN_MSG)
// 	}

// 	var campaignModule = on_chain.InitializeModuleCampaign()
// 	var function string
// 	var args []interface{}

// 	if campaign.PoolID != os.Getenv(env.POOL_ID) || campaign.PoolName != "Main Pool" {
// 		function = campaignModule.GetFunctionCreateCampaignForRegionPool()
// 		args = campaignModule.ToCreateCampaignForRegionPoolArguments(on_chain.CreateCampaignForRegionPoolArguments{
// 			LocalPoolID: campaign.PoolID,
// 			CreateCampaignForMainPoolArguments: on_chain.CreateCampaignForMainPoolArguments{
// 				Creator:     campaign.ActorAddress,
// 				Target:      campaign.Target,
// 				Description: campaign.Description,
// 				ProofBlobID: campaign.ProofBlobID,
// 			},
// 		})
// 	} else {
// 		function = campaignModule.GetFunctionCreateCampaignForMainPool()
// 		args = campaignModule.ToCreateCampaignForMainPoolArguments(on_chain.CreateCampaignForMainPoolArguments{
// 			Creator:     campaign.ActorAddress,
// 			Target:      campaign.Target,
// 			Description: campaign.Description,
// 			ProofBlobID: campaign.ProofBlobID,
// 		})
// 	}
// 	txBytes, err := on_chain.BuildTransaction(on_chain.BuildTransactionRequest{
// 		Client:    client,
// 		Sender:    reviewer,
// 		Module:    campaignModule.GetModule(),
// 		Function:  function,
// 		ErrLogger: p.errLogger,
// 		Arguments: args,
// 	}, ctx)
// 	if err != nil {
// 		return response.BuildTransactionResponse{}, err
// 	}

// 	campaign.ReviewedBy = &reviewer
// 	campaign.ReviewStatus = request_approved_status

// 	return response.BuildTransactionResponse{
// 		TxBytes: txBytes,
// 	}, p.pendingCampaignRepo.UpdatePendingCampaign(*campaign, ctx)
// }

// ApprovePendingCampaign implements business.IPendingCampaignService.
func (p *pendingCampaignService) ApprovePendingCampaign(id string, ctx context.Context) error {
	campaign, err := p.pendingCampaignRepo.GetPendingCampaign(id, ctx)
	if err != nil {
		return err
	}

	if campaign == nil {
		return errors.New(noti.GENERIC_ERROR_WARN_MSG)
	}

	if campaign.ReviewedBy != nil || campaign.ReviewStatus != request_pending_status {
		return errors.New(noti.REQUEST_REVIEWED_MESSAGE)
	}

	var client = p.clients[constant.SuiTestnet]
	manageObj, err := on_chain.GetOnChainObject[entities.Manage](on_chain.GetOnChainObjectRequest{
		Client:   client,
		ObjectId: os.Getenv(env.MANAGE_OBJECT_ID),
	}, ctx)
	if err != nil {
		return err
	}

	if manageObj == nil {
		return errors.New(noti.INTERNALL_ERR_MSG)
	}

	var reviewer string = ctx.Value("address").(string)
	if !slices.Contains(manageObj.AdminIds, reviewer) {
		return errors.New(noti.GENERIC_RIGHT_ACCESS_WARN_MSG)
	}

	var campaignModule = on_chain.InitializeModuleCampaign()
	var function string
	var args []interface{}
	if campaign.PoolID != os.Getenv(env.POOL_ID) || campaign.PoolName != "Main Pool" {
		function = campaignModule.GetFunctionCreateCampaignForRegionPool()
		args = campaignModule.ToCreateCampaignForRegionPoolArguments(on_chain.CreateCampaignForRegionPoolArguments{
			LocalPoolID: campaign.PoolID,
			CreateCampaignForMainPoolArguments: on_chain.CreateCampaignForMainPoolArguments{
				Creator:     campaign.ActorAddress,
				Target:      campaign.Target,
				Description: campaign.Description,
				ProofBlobID: campaign.ProofBlobID,
				Sender:      campaign.ActorAddress,
			},
		})
	} else {
		function = campaignModule.GetFunctionCreateCampaignForMainPool()
		args = campaignModule.ToCreateCampaignForMainPoolArguments(on_chain.CreateCampaignForMainPoolArguments{
			Creator:     campaign.ActorAddress,
			Target:      campaign.Target,
			Description: campaign.Description,
			ProofBlobID: campaign.ProofBlobID,
			Sender:      campaign.ActorAddress,
		})
	}

	campaign.ReviewedBy = &reviewer
	campaign.ReviewStatus = request_approved_status
	if err := p.pendingCampaignRepo.UpdatePendingCampaign(*campaign, ctx); err != nil {
		return err
	}

	_, errRes := on_chain.ExecuteTransactionV2(on_chain.ExecuteTransactionRequestV2{
		Client:    client,
		Module:    campaignModule.GetModule(),
		Function:  function,
		ErrLogger: p.errLogger,
		Arguments: args,
	}, ctx)

	return errRes
}

// CreatePendingCampaign implements business.IPendingCampaignService.
func (p *pendingCampaignService) CreatePendingCampaign(req request.CreatePendingCampaignRequest, ctx context.Context) (*entities.PendingCampaign, error) {
	var client = p.clients[constant.SuiTestnet]
	manage, err := on_chain.GetOnChainObject[entities.Manage](on_chain.GetOnChainObjectRequest{
		Client:    p.clients[constant.SuiTestnet],
		ObjectId:  os.Getenv(env.MANAGE_OBJECT_ID),
		ErrLogger: p.errLogger,
	}, ctx)
	if err != nil {
		return nil, err
	}

	var genericErr error = errors.New(noti.GENERIC_ERROR_WARN_MSG)
	var genericRightErr = errors.New(noti.GENERIC_RIGHT_ACCESS_WARN_MSG)
	var sender string = ctx.Value("address").(string)
	var poolId string
	if req.PoolName != "Main Pool" {
		pool, err := on_chain.GetOnChainObject[entities.MainPool](on_chain.GetOnChainObjectRequest{
			Client:    client,
			ObjectId:  os.Getenv(env.POOL_ID),
			ErrLogger: p.errLogger,
		}, ctx)
		if err != nil {
			return nil, err
		}

		localPools, err := on_chain.GetOnChainObjects[entities.LocalPool](on_chain.GetOnChainObjectsRequest{
			Client:    client,
			ObjectIds: pool.LocalPools,
			ErrLogger: p.errLogger,
		}, ctx)
		if err != nil {
			return nil, err
		}

		var foundPool entities.LocalPool
		for _, localPool := range localPools {
			if localPool.Region == req.PoolName {
				foundPool = localPool
				break
			}
		}

		if foundPool.ID.ID == "" {
			return nil, genericErr
		}

		// Not region leader
		if !slices.Contains(foundPool.Mods, sender) {
			if !slices.Contains(manage.AdminIds, sender) { // Not admin
				return nil, genericRightErr
			}
		}

		poolId = foundPool.ID.ID
	} else {
		if !slices.Contains(manage.AdminIds, sender) { // Not admin
			return nil, genericRightErr
		}

		poolId = os.Getenv(env.POOL_ID)
	}

	bankProfile, err := p.bankProfileRepo.GetBankProfileByOwner(sender, ctx)
	if err != nil {
		return nil, err
	}

	if bankProfile == nil {
		return nil, errors.New(noti.LEADER_NOT_UPLOAD_BANK_PROFILE_MESSAGE)
	}

	var curTime time.Time = time.Now()
	var campaign = entities.PendingCampaign{
		ID:             util.GenerateId(),
		PoolID:         poolId,
		PoolName:       req.PoolName,
		ActorProfileID: ctx.Value("sub").(string),
		ActorAddress:   sender,
		Target:         req.Target,
		Description:    strings.TrimSpace(req.Description),
		ProofBlobID:    req.ProofBlobID,
		CreatedAt:      curTime,
		UpdatedAt:      curTime,
	}

	return &campaign, p.pendingCampaignRepo.CreatePendingCampaign(campaign, ctx)
}

// GetPendingCampaign implements business.IPendingCampaignService.
func (p *pendingCampaignService) GetPendingCampaign(id string, ctx context.Context) (*entities.PendingCampaign, error) {
	return p.pendingCampaignRepo.GetPendingCampaign(id, ctx)
}

// GetPendingCampaigns implements business.IPendingCampaignService.
func (p *pendingCampaignService) GetPendingCampaigns(req request.GetPendingCampaignsRequest, ctx context.Context) (response.PaginationDataResponse, error) {
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
	// var redisKey string = p.getGetPendingCampaignsRedisKey(req)
	// if p.redisCache.Get(redisKey, &res, ctx) {
	// 	return res, nil
	// }

	data, pages, err := p.pendingCampaignRepo.GetPendingCampaigns(req, ctx)
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

// RefusePendingCampaign implements business.IPendingCampaignService.
func (p *pendingCampaignService) RefusePendingCampaign(id string, ctx context.Context) error {
	campaign, err := p.pendingCampaignRepo.GetPendingCampaign(id, ctx)
	if err != nil {
		return err
	}

	if campaign == nil {
		return errors.New(noti.GENERIC_ERROR_WARN_MSG)
	}

	if campaign.ReviewedBy != nil || campaign.ReviewStatus != request_pending_status {
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

	campaign.ReviewedBy = &reviewer
	campaign.ReviewStatus = request_refused_status
	campaign.UpdatedAt = time.Now()
	return p.pendingCampaignRepo.UpdatePendingCampaign(*campaign, ctx)
}

func (p *pendingCampaignService) getGetPendingCampaignsRedisKey(req request.GetPendingCampaignsRequest) string {
	var keyword string = "empty"
	if req.Keyword != "" {
		keyword = req.Keyword
	}

	var pool string = "empty"
	if req.PoolName != "" {
		pool = req.PoolName
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

	return fmt.Sprintf("pending_campaign:kw:%s:pool:%s:status:%s:of:%s:reviewed:%s:min:%s:max:%s:sc:%s:o:%s:s:%d:p:%d",
		keyword, pool, status, creator, reviewer, minAmount, maxAmount, sortCriteria, req.SortOrder, req.PageSize, req.Page)
}
