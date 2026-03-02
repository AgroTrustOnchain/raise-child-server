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
	"time"

	"slices"

	"github.com/block-vision/sui-go-sdk/constant"
	"github.com/block-vision/sui-go-sdk/sui"
)

type regionService struct {
	regionRepo i_repository.ISupportedRegionProposalRepository
	redisCache cache.IRedisCache
	clients    map[string]sui.ISuiAPI
	regions    []string
	errLogger  *log.Logger
}

var _regions []string

func initalizeRegionService(regionRepo i_repository.ISupportedRegionProposalRepository,
	clients map[string]sui.ISuiAPI,
	regions []string,
	errLogger *log.Logger) business.IRegionService {
	return &regionService{
		regionRepo: regionRepo,
		redisCache: cache.InitializeRedisCache(),
		clients:    clients,
		regions:    regions,
		errLogger:  errLogger,
	}
}

func GenerateRegionService() (business.IRegionService, error) {
	var errLogger = util.GetLogConfig(shared.ERROR_LEVEL)

	cnn, err := db.ConnectDB(errLogger, db.InitializePostgreSQL())
	if err != nil {
		return nil, err
	}

	if _regions == nil || len(_regions) == 0 {
		_regions = []string{
			shared.TUYEN_QUANG_REGION,
			shared.LAO_CAI_REGION,
			shared.THAI_NGUYEN_REGION,
			shared.PHU_THO_REGION,
			shared.BAC_BINH_REGION,
			shared.HUNG_YEN_REGION,
			shared.HAI_PHONG_REGION,
			shared.NINH_BINH_REGION,
			shared.QUANG_TRI_REGION,
			shared.DA_NANG_REGION,
			shared.QUANG_NGAI_REGION,
			shared.GIA_LAI_REGION,
			shared.KHANH_HOA_REGION,
			shared.LAM_DONG_REGION,
			shared.DAK_LAK_REGION,
			shared.HO_CHI_MINH_REGION,
			shared.DONG_NAI_REGION,
			shared.TAY_NINH_REGION,
			shared.CAN_THO_REGION,
			shared.VINH_LONG_REGION,
			shared.DONG_THAP_REGION,
			shared.CA_MAU_REGION,
			shared.AN_GIANG_REGION,
			shared.HA_NOI_REGION,
			shared.HUE_REGION,
			shared.LAI_CHAU_REGION,
			shared.DIEN_BIEN_REGION,
			shared.SON_LA_REGION,
			shared.QUANG_NINH_REGION,
			shared.THANH_HOA_REGION,
			shared.NGHE_AN_REGION,
			shared.HA_TINH_REGION,
			shared.CAO_BANG_REGION,
		}
	}

	return initalizeRegionService(repository.InitializeSupportedRegionProposalRepository(cnn, errLogger), _networkAliases, _regions, errLogger), nil
}

// CreateSupportedRegionProposal implements business.IRegionService.
func (r *regionService) CreateSupportedRegionProposal(req request.CreateSupportedRegionProposalsRequest, ctx context.Context) (*entities.SupportedRegionProposal, error) {
	var address string = ctx.Value("address").(string)
	if !util.IsValidSuiAddressStrict(address) {
		return nil, errors.New(noti.GENERIC_ERROR_WARN_MSG)
	}

	isRequested, err := r.regionRepo.IsRegionRequested(req.Region, ctx)
	if err != nil {
		return nil, err
	}

	if isRequested {
		return nil, errors.New(noti.SUPPORRT_REGION_REQUEST_MESSAGE)
	}

	var client = r.clients[constant.SuiTestnet]
	var packageId string = os.Getenv(env.PACKAGE_ID)
	var manageModule = on_chain.InitializeModuleManage()
	adminNfts, err := on_chain.GetOnChainOwnedObjects[entities.AdminNft](on_chain.GetOnChainOwnedObjectsRequest{
		Client:       client,
		OwnerAddress: address,
		StructType:   fmt.Sprintf("%s::%s::%s", packageId, manageModule.GetModule(), manageModule.GetAdminNftStruct()),
	}, ctx)
	if err != nil {
		return nil, err
	}

	// Not admin
	if adminNfts == nil || len(adminNfts) == 0 {
		var staffModule = on_chain.InitializeModuleStaff()
		staffNfts, err := on_chain.GetOnChainOwnedObjects[entities.StaffNft](on_chain.GetOnChainOwnedObjectsRequest{
			Client:       client,
			OwnerAddress: address,
			StructType:   fmt.Sprintf("%s::%s::%s", packageId, staffModule.GetModule(), staffModule.GetStaffNftObjectStruct()),
		}, ctx)
		if err != nil {
			return nil, err
		}

		var genericRightErr error = errors.New(noti.GENERIC_RIGHT_ACCESS_WARN_MSG)
		if staffNfts == nil || len(staffNfts) == 0 {
			return nil, genericRightErr
		}

		var isLeaderOfRegion bool = false
		for _, nft := range staffNfts {
			if nft.Region == req.Region && nft.Role == local_leader_role {
				isLeaderOfRegion = true
				break
			}
		}

		if !isLeaderOfRegion {
			return nil, genericRightErr
		}
	}

	var curTime time.Time = time.Now()
	var proposal = entities.SupportedRegionProposal{
		ID:        util.GenerateId(),
		Sub:       ctx.Value("sub").(string),
		Region:    req.Region,
		Content:   req.Content,
		CreatedBy: address,
		CreatedAt: curTime,
		UpdatedAt: curTime,
	}

	return &proposal, r.regionRepo.CreateSupportedRegionProposal(proposal, ctx)
}

// GetSupportedRegionProposal implements business.IRegionService.
func (r *regionService) GetSupportedRegionProposal(id string, ctx context.Context) (*entities.SupportedRegionProposal, error) {
	return r.regionRepo.GetSupportedRegionProposal(id, ctx)
}

// GetSupportedRegionProposals implements business.IRegionService.
func (r *regionService) GetSupportedRegionProposals(req request.GetSupportedRegionProposalsRequest, ctx context.Context) (response.PaginationDataResponse, error) {
	if req.CreatedBy != "" {
		if !util.IsValidSuiAddressStrict(req.CreatedBy) {
			return response.PaginationDataResponse{}, errors.New(noti.GENERIC_ERROR_WARN_MSG)
		}
	}

	req.SortOrder = util.StanderizeSortOrder(req.SortOrder)
	if req.PageSize < 1 {
		req.PageSize = default_page_size
	}

	if req.Page < 1 {
		req.Page = 1
	}

	var res response.PaginationDataResponse
	var redisKey string = r.getGetSupportedRegionProposalsRedisKey(req)
	if r.redisCache.Get(redisKey, &res, ctx) {
		return res, nil
	}

	data, pages, err := r.regionRepo.GetSupportedRegionProposals(req, ctx)
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

	r.redisCache.Set(redisKey, res, time.Minute*5, ctx)

	return res, err
}

// GetRegions implements business.IRegionService.
func (r *regionService) GetRegions() response.RegionsResponse {
	return response.RegionsResponse{
		Regions: r.regions,
	}
}

func (r *regionService) getGetSupportedRegionProposalsRedisKey(req request.GetSupportedRegionProposalsRequest) string {
	var keyword string = "empty"
	if req.Keyword != "" {
		keyword = req.Keyword
	}

	var createdBy string = "empty"
	if req.CreatedBy != "" {
		createdBy = req.CreatedBy
	}

	return fmt.Sprintf("region_proposal:kw:%s:of:%s:o:%s:s:%d:p:%d",
		keyword, createdBy, req.SortOrder, req.PageSize, req.Page)
}

func isRegionExist(region string) bool {
	return slices.Contains(_regions, region)
}
