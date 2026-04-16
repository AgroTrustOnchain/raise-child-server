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
	"raise-child/util/ai"
	"raise-child/util/cache"
	"raise-child/util/db"
	on_chain "raise-child/util/on_chain"
	"sort"
	"strconv"
	"strings"
	"time"

	walrus_pkg "raise-child/util/walrus_pkg"

	"github.com/block-vision/sui-go-sdk/constant"
	"github.com/block-vision/sui-go-sdk/sui"
	"github.com/payOSHQ/payos-lib-golang"
)

type campaignService struct {
	profileRepo                 i_repository.IProfileRepository
	paymentRepo                 i_repository.IPaymentRepository
	pendingWithdrawProposalRepo i_repository.IPendingWithdrawProposalRepository
	donationRepo                i_repository.IOffChainDonationRepository
	withdrawRepo                i_repository.IOffChainWithdrawProposalRepository
	aiProvider                  ai.IAiClientProvider
	walrusProvider              walrus_pkg.IWalrusProvider
	redisCache                  cache.IRedisCache
	clients                     map[string]sui.ISuiAPI
	errLogger                   *log.Logger
}

func initializeCampaignService(
	profileRepo i_repository.IProfileRepository,
	paymentRepo i_repository.IPaymentRepository,
	pendingWithdrawProposalRepo i_repository.IPendingWithdrawProposalRepository,
	donationRepo i_repository.IOffChainDonationRepository,
	withdrawRepo i_repository.IOffChainWithdrawProposalRepository,
	aiProvider ai.IAiClientProvider,
	walrusProvider walrus_pkg.IWalrusProvider,
	clients map[string]sui.ISuiAPI,
	errLogger *log.Logger,
) business.ICampaignService {
	return &campaignService{
		profileRepo:                 profileRepo,
		paymentRepo:                 paymentRepo,
		pendingWithdrawProposalRepo: pendingWithdrawProposalRepo,
		donationRepo:                donationRepo,
		withdrawRepo:                withdrawRepo,
		aiProvider:                  aiProvider,
		walrusProvider:              walrusProvider,
		redisCache:                  cache.InitializeRedisCache(),
		clients:                     clients,
		errLogger:                   errLogger,
	}
}

func GenerateCampaignService() (business.ICampaignService, error) {
	var errLogger = util.GetLogConfig(shared.ERROR_LEVEL)

	cnn, err := db.ConnectDB(errLogger, db.InitializePostgreSQL())
	if err != nil {
		return nil, err
	}

	return initializeCampaignService(
		repository.InitializeProfileRepository(cnn, errLogger),
		repository.InitializePaymentRepository(cnn, errLogger),
		repository.InitializePendingWithdrawProposalRepo(cnn, errLogger),
		repository.InitializeOffChainDonationRepository(cnn, errLogger),
		repository.InitializeOffChainWithdrawProposalRepository(cnn, errLogger),
		ai.InitializeAiProvider(nil, errLogger),
		walrus_pkg.InitializeWalrusProvider(errLogger),
		_networkAliases,
		errLogger,
	), nil
}

// CreateCampaignWithdrawProposal implements business.ICampaignService.
func (c *campaignService) CreateCampaignWithdrawProposal(req request.CreateCampaignWithdrawProposalRequest, ctx context.Context) (*entities.PendingWithdrawProposal, error) {
	var genericErr error = errors.New(noti.GENERIC_ERROR_WARN_MSG)
	if !util.IsValidSuiAddressStrict(req.CampaignID) {
		return nil, genericErr
	}

	var client = c.clients[constant.SuiTestnet]
	campaign, err := on_chain.GetOnChainObject[entities.OnChainCampaign](on_chain.GetOnChainObjectRequest{
		Client:    client,
		ObjectId:  req.CampaignID,
		ErrLogger: c.errLogger,
	}, ctx)
	if err != nil {
		return nil, err
	}

	if campaign == nil {
		return nil, genericErr
	}

	var sender string = ctx.Value("address").(string)
	if campaign.Creator != sender {
		return nil, errors.New(noti.GENERIC_RIGHT_ACCESS_WARN_MSG)
	}

	totalWithdrawAmount, _ := strconv.ParseInt(campaign.WithdrawAmount, 10, 64)
	totalDonation, _ := strconv.ParseInt(campaign.TotalDonated, 10, 64)
	if req.Amount > totalDonation-totalWithdrawAmount {
		return nil, errors.New(noti.CURRENT_BUDGET_NOT_ENOUGH_MESSAGE)
	}

	var description string = strings.TrimSpace(req.Description)
	var purpose string = string(entities.CAMPAIGN_PURPOSE)
	var aiEvaluation string
	if req.ProofBlobID != nil {
		proofBytes, _ := c.walrusProvider.FetchBytesImage(*req.ProofBlobID)
		if proofBytes != nil {
			aiEvaluation = c.aiProvider.ValidateWithdrawProposal(ai.ValidateWithdrawProposal{
				Purpose:         purpose,
				WithdrawAmount:  req.Amount,
				Description:     description,
				ProofBytesImage: proofBytes,
			}, ctx)
		}
	}

	var curTime time.Time = time.Now()
	var res = entities.PendingWithdrawProposal{
		ID:             util.GenerateId(),
		ProfileID:      ctx.Value("sub").(string),
		Creator:        sender,
		PoolID:         campaign.PoolID,
		PoolName:       campaign.PoolName,
		Purpose:        purpose,
		Target:         req.CampaignID,
		WithdrawAmount: req.Amount,
		ProofBlobID:    req.ProofBlobID,
		Description:    description,
		Status:         request_pending_status,
		AIEvaluation:   aiEvaluation,
		CreatedAt:      curTime,
		UpdatedAt:      curTime,
	}

	return &res, c.pendingWithdrawProposalRepo.CreatePendingWithdrawProposal(res, ctx)
}

// GetCampaign implements business.ICampaignService.
func (c *campaignService) GetCampaign(id string, ctx context.Context) (response.OnChainCampaignResponse, error) {
	var genericErr error = errors.New(noti.GENERIC_ERROR_WARN_MSG)
	if !util.IsValidSuiAddressStrict(id) {
		return response.OnChainCampaignResponse{}, genericErr
	}

	res, err := on_chain.GetOnChainObject[entities.OnChainCampaign](on_chain.GetOnChainObjectRequest{
		Client:    c.clients[constant.SuiTestnet],
		ObjectId:  id,
		ErrLogger: c.errLogger,
	}, ctx)
	if err != nil {
		return response.OnChainCampaignResponse{}, err
	}

	if res == nil {
		return response.OnChainCampaignResponse{}, genericErr
	}

	return res.ToOnChainCampaignResponse(), nil
}

// GetCampaigns implements business.ICampaignService.
func (c *campaignService) GetCampaigns(req request.GetCampaignsRequest, ctx context.Context) (response.PaginationDataResponse, error) {
	var genericErr error = errors.New(noti.GENERIC_ERROR_WARN_MSG)
	if req.Creator != "" {
		if !util.IsValidSuiAddressStrict(req.Creator) {
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

	req.SortOrder = util.StanderizeSortOrder(req.SortOrder)
	req.Keyword = util.StanderizeString(req.Keyword)
	req.SortCriteria = util.StanderizeSortCriteria(req.SortCriteria)
	if req.Page < 1 {
		req.Page = 1
	}

	if req.PageSize < 1 {
		req.PageSize = default_page_size
	}

	var res response.PaginationDataResponse
	var redisKey string = c.getGetCampaignsRedisKey(req)
	if c.redisCache.Get(redisKey, &res, ctx) {
		return res, nil
	}

	var client = c.clients[constant.SuiTestnet]
	pool, err := on_chain.GetOnChainObject[entities.MainPool](on_chain.GetOnChainObjectRequest{
		Client:    client,
		ObjectId:  os.Getenv(env.POOL_ID),
		ErrLogger: c.errLogger,
	}, ctx)
	if err != nil {
		return response.PaginationDataResponse{}, err
	}

	var camapaigns []entities.OnChainCampaign
	var errRes error
	if req.PoolName != "" {
		if req.PoolName == "Main Pool" {
			camapaigns, errRes = on_chain.GetOnChainObjects[entities.OnChainCampaign](on_chain.GetOnChainObjectsRequest{
				Client:    client,
				ObjectIds: pool.Campaigns,
				ErrLogger: c.errLogger,
			}, ctx)
			if errRes != nil {
				return response.PaginationDataResponse{}, errRes
			}
		} else {
			if !isRegionExist(req.PoolName) {
				return response.PaginationDataResponse{}, genericErr
			}

			localPools, err := on_chain.GetOnChainObjects[entities.LocalPool](on_chain.GetOnChainObjectsRequest{
				Client:    client,
				ObjectIds: pool.Campaigns,
				ErrLogger: c.errLogger,
			}, ctx)
			if err != nil {
				return response.PaginationDataResponse{}, err
			}

			var foundLocalPool entities.LocalPool
			for _, localPool := range localPools {
				if localPool.Region == req.PoolName {
					foundLocalPool = localPool
					break
				}
			}

			if foundLocalPool.ID.ID == "" {
				return response.PaginationDataResponse{}, genericErr
			}

			camapaigns, errRes = on_chain.GetOnChainObjects[entities.OnChainCampaign](on_chain.GetOnChainObjectsRequest{
				Client:    client,
				ObjectIds: foundLocalPool.Campaigns,
				ErrLogger: c.errLogger,
			}, ctx)
			if errRes != nil {
				return response.PaginationDataResponse{}, errRes
			}
		}
	} else {
		camapaigns, errRes = on_chain.GetOnChainObjects[entities.OnChainCampaign](on_chain.GetOnChainObjectsRequest{
			Client:    client,
			ObjectIds: pool.AllCampaigns,
			ErrLogger: c.errLogger,
		}, ctx)
		if errRes != nil {
			return response.PaginationDataResponse{}, errRes
		}
	}

	if camapaigns == nil || len(camapaigns) == 0 {
		return response.PaginationDataResponse{
			Page: req.Page,
		}, nil
	}

	var filteredCampaigns []entities.OnChainCampaign
	for i := len(camapaigns) - 1; i >= 0; i-- {
		var campaign entities.OnChainCampaign = camapaigns[i]
		if req.Creator != "" {
			if campaign.Creator != req.Creator { // Not matched
				continue
			}
		}

		if req.Keyword != "" {
			if !strings.Contains(campaign.Description, req.Keyword) && !strings.Contains(campaign.PoolName, req.Keyword) {
				continue
			}
		}

		target, _ := strconv.ParseInt(campaign.Target, 10, 64)
		if req.MinAmount != nil {
			if target < *req.MinAmount {
				continue
			}
		}

		if req.MaxAmount != nil {
			if target > *req.MaxAmount {
				continue
			}
		}

		filteredCampaigns = append(filteredCampaigns, campaign)
	}

	sort.Slice(filteredCampaigns, func(i, j int) bool {
		if req.SortCriteria == "target" {
			target1, _ := strconv.ParseInt(filteredCampaigns[i].Target, 10, 64)
			target2, _ := strconv.ParseInt(filteredCampaigns[j].Target, 10, 64)
			if req.SortOrder == "DESC" {
				return target2 > target1
			}

			return target2 < target1
		}

		if req.SortOrder == "ASC" {
			return false
		}

		return true
	})

	var skippedRecords int = (req.Page - 1) * req.PageSize
	if len(filteredCampaigns) <= skippedRecords {
		return response.PaginationDataResponse{}, nil
	}

	var data []response.OnChainCampaignResponse
	for i := skippedRecords; i < len(filteredCampaigns); i++ {
		data = append(data, filteredCampaigns[i].ToOnChainCampaignResponse())
		if len(data) == req.PageSize {
			break
		}
	}

	res = response.PaginationDataResponse{
		Data:       data,
		Amount:     len(data),
		Page:       req.Page,
		TotalPages: int(math.Ceil(float64(len(filteredCampaigns)) / float64(req.PageSize))),
	}

	c.redisCache.Set(redisKey, res, time.Minute*5, ctx)
	return res, nil
}

// SupportCampaign implements business.ICampaignService.
func (c *campaignService) SupportCampaign(id string, req request.SupportCampaignRequest, ctx context.Context) (response.UrlAPIResponse, error) {
	profile, err := c.profileRepo.GetProfile(ctx.Value("sub").(string), ctx)
	if err != nil {
		return response.UrlAPIResponse{}, err
	}

	if profile == nil || profile.IdentityCode == nil {
		return response.UrlAPIResponse{}, errors.New(noti.PROFILE_EMPTY_MESSAGE)
	}

	var genericErr error = errors.New(noti.GENERIC_ERROR_WARN_MSG)
	if !util.IsValidSuiAddressStrict(id) {
		return response.UrlAPIResponse{}, genericErr
	}

	var client = c.clients[constant.SuiTestnet]
	campaign, err := on_chain.GetOnChainObject[entities.OnChainCampaign](on_chain.GetOnChainObjectRequest{
		Client:    client,
		ObjectId:  id,
		ErrLogger: c.errLogger,
	}, ctx)
	if err != nil {
		return response.UrlAPIResponse{}, err
	}

	if campaign == nil {
		return response.UrlAPIResponse{}, genericErr
	}

	target, _ := strconv.ParseInt(campaign.Target, 10, 64)
	totalDonations, _ := strconv.ParseInt(campaign.TotalDonated, 10, 64)
	if req.Amount > target-totalDonations {
		return response.UrlAPIResponse{}, errors.New(noti.SUPPORT_SURPASS_CAMPAIGN_TARGET_MESSAGE)
	}

	var paymentId string = util.GenerateId()
	var orderCode int = util.GenerateNumber()
	var callbackUrl string = os.Getenv(payment.PAYMENT_CALLBACK_URL) + paymentId
	var description string = entities.POOL_CAMPAIGN_PAYMENT_DESCRIPTION.GenerateSupportPaymentDescription()
	data, err := payos.CreatePaymentLink(payos.CheckoutRequestType{
		OrderCode:   int64(orderCode),
		Amount:      int(req.Amount),
		Description: description,
		ReturnUrl:   callbackUrl,
		CancelUrl:   callbackUrl,
	})
	if err != nil {
		c.errLogger.Println("Err: ", err.Error())
		return response.UrlAPIResponse{}, errors.New(noti.INTERNALL_ERR_MSG)
	}

	var donationId string = util.GenerateId()
	var curTime time.Time = time.Now()
	if err := c.donationRepo.CreateDonation(entities.OffChainDonation{
		ID:        donationId,
		Purpose:   string(entities.CAMPAIGN_PURPOSE),
		Target:    id,
		CreatedAt: curTime,
	}, ctx); err != nil {
		return response.UrlAPIResponse{}, err
	}

	return response.UrlAPIResponse{
			Url: data.CheckoutUrl,
		}, c.paymentRepo.CreatePayment(entities.Payment{
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
			Message:       strings.TrimSpace(req.Description),
			ExpiredAt:     time.Unix(int64(*data.ExpiredAt), 0),
			CreatedAt:     curTime,
			UpdatedAt:     curTime,
		}, ctx)
}

func (c *campaignService) getGetCampaignsRedisKey(req request.GetCampaignsRequest) string {
	var keyword string = "empty"
	if req.Keyword != "" {
		keyword = req.Keyword
	}

	var pool string = "empty"
	if req.PoolName != "" {
		pool = req.PoolName
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

	var sortCriteria string = "empty"
	if req.SortCriteria != "" {
		sortCriteria = req.SortCriteria
	}

	return fmt.Sprintf("campaign:kw:%s:pool:%s:of:%s:min:%s:max:%s:sc:%s:o:%s:s:%d:p:%d",
		keyword, pool, creator, minAmount, maxAmount, sortCriteria, req.SortOrder, req.PageSize, req.Page)
}
