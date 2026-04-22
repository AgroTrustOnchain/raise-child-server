package business

import (
	"context"
	"errors"
	"fmt"
	"log"
	"math"
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
	"sort"
	"strconv"
	"strings"
	"time"

	"slices"

	"github.com/block-vision/sui-go-sdk/constant"
	"github.com/block-vision/sui-go-sdk/sui"
)

type regionService struct {
	regionRepo i_repository.ISupportedRegionSuggestionRepository
	redisCache cache.IRedisCache
	clients    map[string]sui.ISuiAPI
	regions    []string
	errLogger  *log.Logger
}

var _regions []string

func initalizeRegionService(
	regionRepo i_repository.ISupportedRegionSuggestionRepository,
	clients map[string]sui.ISuiAPI,
	regions []string,
	errLogger *log.Logger,
) business.IRegionService {
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
			shared.AGROTRUST_REGION,
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

	return initalizeRegionService(repository.InitializeSupportedRegionSuggestionRepository(cnn, errLogger), _networkAliases, _regions, errLogger), nil
}

// GetEstablishedRegions implements business.IRegionService.
func (r *regionService) GetEstablishedRegions(ctx context.Context) (response.RegionsResponse, error) {
	manage, err := on_chain.GetOnChainObject[entities.Manage](on_chain.GetOnChainObjectRequest{
		Client:    r.clients[constant.SuiTestnet],
		ObjectId:  os.Getenv(env.MANAGE_OBJECT_ID),
		ErrLogger: r.errLogger,
	}, ctx)
	if err != nil {
		return response.RegionsResponse{}, err
	}

	var res []string
	for i, region := range manage.LocalRegions {
		if manage.CenterConfirmStatuses[i] {
			res = append(res, region)
		}
	}

	return response.RegionsResponse{
		Regions: res,
	}, nil
}

// GetRegionDetail implements business.IRegionService.
func (r *regionService) GetRegionDetail(region string, req request.GetChildrenFromRegionDetailRequest, ctx context.Context) (response.RegionDetailResponse, error) {
	req.Keyword = util.StandardizeString(req.Keyword)
	req.SortOrder = util.StandardizeSortOrder(req.SortOrder)
	if req.Page < 1 {
		req.Page = 1
	}

	if req.PageSize < 1 {
		req.PageSize = default_page_size
	}

	var res response.RegionDetailResponse
	var redisKey string = r.getGetRegionDetailRedisKey(region, req)
	if r.redisCache.Get(redisKey, &res, ctx) {
		return res, nil
	}

	var client = r.clients[constant.SuiTestnet]
	manage, err := on_chain.GetOnChainObject[entities.Manage](on_chain.GetOnChainObjectRequest{
		Client:    client,
		ObjectId:  os.Getenv(env.MANAGE_OBJECT_ID),
		ErrLogger: r.errLogger,
	}, ctx)
	if err != nil {
		return response.RegionDetailResponse{}, err
	}

	var centerId string
	for i, localRegion := range manage.LocalRegions {
		if localRegion == region {
			if manage.CenterConfirmStatuses[i] {
				centerId = manage.ChildrenCenters[i]
				break
			}
		}
	}

	if centerId == "" {
		return response.RegionDetailResponse{}, errors.New(noti.GENERIC_ERROR_WARN_MSG)
	}

	center, err := on_chain.GetOnChainObject[entities.Center](on_chain.GetOnChainObjectRequest{
		Client:    client,
		ObjectId:  centerId,
		ErrLogger: r.errLogger,
	}, ctx)
	if err != nil {
		return response.RegionDetailResponse{}, err
	}

	var skippedRecords int = (req.Page - 1) * req.PageSize
	var paginationChildrenResponse response.PaginationDataResponse
	if skippedRecords >= len(center.ChildIDs) {
		paginationChildrenResponse = response.PaginationDataResponse{
			Page: req.Page,
		}
	} else {
		children, err := on_chain.GetOnChainObjects[entities.Child](on_chain.GetOnChainObjectsRequest{
			Client:    client,
			ObjectIds: center.ChildIDs,
			ErrLogger: r.errLogger,
		}, ctx)
		if err != nil {
			return response.RegionDetailResponse{}, err
		}

		var filteredChildren []entities.Child
		for i := len(center.ChildIDs) - 1; i >= 0; i-- {
			var child = children[i]
			if req.Keyword != "" {
				var firstName string = util.StandardizeString(child.FirstName)
				var lastName string = util.StandardizeString(child.LastName)
				if !strings.Contains(firstName, req.Keyword) && !strings.Contains(lastName, req.Keyword) && !strings.Contains(child.IdentityCode, req.Keyword) { // Not matched
					continue
				}
			}

			filteredChildren = append(filteredChildren, child)
		}

		if len(filteredChildren) > skippedRecords {
			sort.Slice(filteredChildren, func(i, j int) bool {
				var name1 string = filteredChildren[i].LastName + " " + filteredChildren[i].FirstName
				var name2 string = filteredChildren[j].LastName + " " + filteredChildren[j].FirstName

				if req.SortOrder == "ASC" {
					return name1 < name2
				}

				return name2 > name1
			})

			var data []response.ChildCardMinimumResponse
			for i := skippedRecords; i < len(filteredChildren); i++ {
				data = append(data, filteredChildren[i].ToChildCardMinimumResponse())
				if len(data) == req.PageSize {
					break
				}
			}

			paginationChildrenResponse = response.PaginationDataResponse{
				Data:       data,
				Amount:     len(data),
				Page:       req.Page,
				TotalPages: int(math.Ceil(float64(len(filteredChildren)) / float64(req.PageSize))),
			}
		} else {
			paginationChildrenResponse = response.PaginationDataResponse{
				Page: req.Page,
			}
		}
	}

	pool, err := on_chain.GetOnChainObject[entities.MainPool](on_chain.GetOnChainObjectRequest{
		Client:    client,
		ObjectId:  os.Getenv(env.POOL_ID),
		ErrLogger: r.errLogger,
	}, ctx)
	if err != nil {
		return response.RegionDetailResponse{}, err
	}

	localPools, err := on_chain.GetOnChainObjects[entities.LocalPool](on_chain.GetOnChainObjectsRequest{
		Client:    client,
		ObjectIds: pool.LocalPools,
		ErrLogger: r.errLogger,
	}, ctx)
	if err != nil {
		return response.RegionDetailResponse{}, err
	}

	var total int64
	var poolId string
	for _, localPool := range localPools {
		if localPool.Region == region {
			totalDonated, _ := strconv.ParseInt(localPool.TotalAmount, 10, 64)
			total = totalDonated
			poolId = localPool.ID.ID
			break
		}
	}

	var centerBlobId string
	if center.ImageBlobIDs != nil && len(center.ImageBlobIDs) > 0 {
		centerBlobId = center.ImageBlobIDs[len(center.ImageBlobIDs)-1]
	}
	res = response.RegionDetailResponse{
		Region:            region,
		PoolID:            poolId,
		CenterPhoneNumber: center.CenterPhoneNumber,
		CenterAddress:     center.CenterAddress,
		CenterImageBlobID: centerBlobId,
		TotalDonated:      total,
		Children:          paginationChildrenResponse,
	}

	r.redisCache.Set(redisKey, res, time.Minute*5, ctx)

	return res, nil
}

// CreateSupportedRegionSuggestion implements business.IRegionService.
func (r *regionService) CreateSupportedRegionSuggestion(req request.CreateSupportedRegionSuggestionsRequest, ctx context.Context) (*entities.SupportedRegionSuggestion, error) {
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

	// var client = r.clients[constant.SuiTestnet]
	// var packageId string = os.Getenv(env.PACKAGE_ID)
	// var manageModule = on_chain.InitializeModuleManage()
	// adminNfts, err := on_chain.GetOnChainOwnedObjects[entities.AdminNft](on_chain.GetOnChainOwnedObjectsRequest{
	// 	Client:       client,
	// 	OwnerAddress: address,
	// 	StructType:   fmt.Sprintf("%s::%s::%s", packageId, manageModule.GetModule(), manageModule.GetAdminNftStruct()),
	// }, ctx)
	// if err != nil {
	// 	return nil, err
	// }

	// // Not admin
	// if adminNfts == nil || len(adminNfts) == 0 {
	// 	var staffModule = on_chain.InitializeModuleStaff()
	// 	staffNfts, err := on_chain.GetOnChainOwnedObjects[entities.StaffNft](on_chain.GetOnChainOwnedObjectsRequest{
	// 		Client:       client,
	// 		OwnerAddress: address,
	// 		StructType:   fmt.Sprintf("%s::%s::%s", packageId, staffModule.GetModule(), staffModule.GetStaffNftObjectStruct()),
	// 	}, ctx)
	// 	if err != nil {
	// 		return nil, err
	// 	}

	// 	var genericRightErr error = errors.New(noti.GENERIC_RIGHT_ACCESS_WARN_MSG)
	// 	if staffNfts == nil || len(staffNfts) == 0 {
	// 		return nil, genericRightErr
	// 	}

	// 	var isLeaderOfRegion bool = false
	// 	for _, nft := range staffNfts {
	// 		if nft.Region == req.Region && nft.Role == local_leader_role {
	// 			isLeaderOfRegion = true
	// 			break
	// 		}
	// 	}

	// 	if !isLeaderOfRegion {
	// 		return nil, genericRightErr
	// 	}
	// }

	var proposal = entities.SupportedRegionSuggestion{
		ID:        util.GenerateId(),
		ProfileID: ctx.Value("sub").(string),
		Region:    req.Region,
		Content:   req.Content,
		CreatedBy: address,
	}

	return &proposal, r.regionRepo.CreateSupportedRegionSuggestion(proposal, ctx)
}

// ReviewRegionSuggestion implements business.IRegionService.
func (r *regionService) ReviewRegionSuggestion(id string, req request.VoteRequest, ctx context.Context) error {
	suggestion, err := r.regionRepo.GetSupportedRegionSuggestion(id, ctx)
	if err != nil {
		return err
	}

	if suggestion == nil {
		return errors.New(noti.GENERIC_ERROR_WARN_MSG)
	}

	if suggestion.Status != request_pending_status || suggestion.ReviewedBy != nil {
		return errors.New(noti.REQUEST_REVIEWED_MESSAGE)
	}

	manageObj, err := on_chain.GetOnChainObject[entities.Manage](on_chain.GetOnChainObjectRequest{
		Client:    r.clients[constant.SuiTestnet],
		ObjectId:  os.Getenv(env.MANAGE_OBJECT_ID),
		ErrLogger: r.errLogger,
	}, ctx)
	if err != nil {
		return err
	}

	var sender string = ctx.Value("address").(string)
	if !slices.Contains(manageObj.AdminIds, sender) {
		return errors.New(noti.GENERIC_RIGHT_ACCESS_WARN_MSG)
	}

	suggestion.ReviewedBy = &sender
	if req.IsVoteYes {
		suggestion.Status = request_approved_status
	} else {
		suggestion.Status = request_refused_status
	}

	return r.regionRepo.UpdateSupportedRegionSuggestion(*suggestion, ctx)
}

// GetSupportedRegionSuggestion implements business.IRegionService.
func (r *regionService) GetSupportedRegionSuggestion(id string, ctx context.Context) (*entities.SupportedRegionSuggestion, error) {
	res, err := r.regionRepo.GetSupportedRegionSuggestion(id, ctx)
	if err != nil {
		return nil, err
	}

	if res.Status == request_pending_status || res.Status == request_refused_status {
		var addressValue = ctx.Value("address")
		var genericRightErr error = errors.New(noti.GENERIC_RIGHT_ACCESS_WARN_MSG)
		if addressValue == nil {
			return nil, genericRightErr
		}

		address, _ := addressValue.(string)
		if res.CreatedBy != address {
			manageObj, err := on_chain.GetOnChainObject[entities.Manage](on_chain.GetOnChainObjectRequest{
				Client:    r.clients[constant.SuiTestnet],
				ObjectId:  os.Getenv(env.MANAGE_OBJECT_ID),
				ErrLogger: r.errLogger,
			}, ctx)
			if err != nil {
				return nil, err
			}

			if !slices.Contains(manageObj.AdminIds, address) {
				return nil, genericRightErr
			}
		}
	}

	return res, nil
}

// GetSupportedRegionSuggestions implements business.IRegionService.
func (r *regionService) GetSupportedRegionSuggestions(req request.GetSupportedRegionSuggestionsRequest, ctx context.Context) (response.PaginationDataResponse, error) {
	if req.CreatedBy != "" {
		if !util.IsValidSuiAddressStrict(req.CreatedBy) {
			return response.PaginationDataResponse{}, errors.New(noti.GENERIC_ERROR_WARN_MSG)
		}
	}

	req.SortOrder = util.StandardizeSortOrder(req.SortOrder)
	req.Keyword = strings.TrimSpace(req.Keyword)
	if req.PageSize < 1 {
		req.PageSize = default_page_size
	}

	if req.Page < 1 {
		req.Page = 1
	}

	var res response.PaginationDataResponse
	var redisKey string = r.getGetSupportedRegionSuggestionsRedisKey(req)
	if r.redisCache.Get(redisKey, &res, ctx) {
		return res, nil
	}

	data, pages, err := r.regionRepo.GetSupportedRegionSuggestions(req, true, ctx)
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

// AdminGetSupportedRegionSuggestions implements business.IRegionService.
func (r *regionService) AdminGetSupportedRegionSuggestions(req request.GetSupportedRegionSuggestionsRequest, ctx context.Context) (response.PaginationDataResponse, error) {
	//var address string = ctx.Value("address").(string)
	manageObj, err := on_chain.GetOnChainObject[entities.Manage](on_chain.GetOnChainObjectRequest{
		Client:    r.clients[constant.SuiTestnet],
		ObjectId:  os.Getenv(env.MANAGE_OBJECT_ID),
		ErrLogger: r.errLogger,
	}, ctx)
	if err != nil {
		return response.PaginationDataResponse{}, err
	}

	if !slices.Contains(manageObj.AdminIds, ctx.Value("address").(string)) {
		return response.PaginationDataResponse{}, errors.New(noti.GENERIC_RIGHT_ACCESS_WARN_MSG)
	}

	if req.CreatedBy != "" {
		if !util.IsValidSuiAddressStrict(req.CreatedBy) {
			return response.PaginationDataResponse{}, errors.New(noti.GENERIC_ERROR_WARN_MSG)
		}
	}

	req.SortOrder = util.StandardizeSortOrder(req.SortOrder)
	req.Keyword = strings.TrimSpace(req.Keyword)
	if req.PageSize < 1 {
		req.PageSize = default_page_size
	}

	if req.Page < 1 {
		req.Page = 1
	}

	var res response.PaginationDataResponse
	var redisKey string = r.getGetAuthenticatedSupportedRegionSuggestionsRedisKey(req)
	if r.redisCache.Get(redisKey, &res, ctx) {
		return res, nil
	}

	data, pages, err := r.regionRepo.GetSupportedRegionSuggestions(req, false, ctx)
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

// GetWalletSupportedRegionSuggestions implements business.IRegionService.
func (r *regionService) GetWalletSupportedRegionSuggestions(req request.GetSupportedRegionSuggestionsRequest, ctx context.Context) (response.PaginationDataResponse, error) {
	var address string = ctx.Value("address").(string)
	if address != req.CreatedBy {
		return response.PaginationDataResponse{}, errors.New(noti.GENERIC_RIGHT_ACCESS_WARN_MSG)
	}

	req.SortOrder = util.StandardizeSortOrder(req.SortOrder)
	req.Keyword = strings.TrimSpace(req.Keyword)
	if req.PageSize < 1 {
		req.PageSize = default_page_size
	}

	if req.Page < 1 {
		req.Page = 1
	}

	var res response.PaginationDataResponse
	var redisKey string = r.getGetAuthenticatedSupportedRegionSuggestionsRedisKey(req)
	if r.redisCache.Get(redisKey, &res, ctx) {
		return res, nil
	}

	data, pages, err := r.regionRepo.GetSupportedRegionSuggestions(req, false, ctx)
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

func (r *regionService) getGetSupportedRegionSuggestionsRedisKey(req request.GetSupportedRegionSuggestionsRequest) string {
	var keyword string = "empty"
	if req.Keyword != "" {
		keyword = req.Keyword
	}

	var createdBy string = "empty"
	if req.CreatedBy != "" {
		createdBy = req.CreatedBy
	}

	return fmt.Sprintf("region_suggestion:kw:%s:of:%s:o:%s:s:%d:p:%d",
		keyword, createdBy, req.SortOrder, req.PageSize, req.Page)
}

func (r *regionService) getGetAuthenticatedSupportedRegionSuggestionsRedisKey(req request.GetSupportedRegionSuggestionsRequest) string {
	var keyword string = "empty"
	if req.Keyword != "" {
		keyword = req.Keyword
	}

	var createdBy string = "empty"
	if req.CreatedBy != "" {
		createdBy = req.CreatedBy
	}

	return fmt.Sprintf("auth_region_suggestion:kw:%s:of:%s:o:%s:s:%d:p:%d",
		keyword, createdBy, req.SortOrder, req.PageSize, req.Page)
}

func (r *regionService) getGetRegionDetailRedisKey(region string, req request.GetChildrenFromRegionDetailRequest) string {
	var keyword string = "empty"
	if req.Keyword != "" {
		keyword = req.Keyword
	}

	return fmt.Sprintf("region_detail:region:%s:kw:%s:o:%s:s:%d:p:%d",
		region, keyword, req.SortOrder, req.PageSize, req.Page)
}

func isRegionExist(region string) bool {
	return slices.Contains(_regions, region)
}
