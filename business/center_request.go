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

type centerRequestService struct {
	platformConfigRepo i_repository.IPlatformConfigRepository
	centerRequestRepo  i_repository.ICenterRequestRepository
	profileRepo        i_repository.IProfileRepository
	redisCache         cache.IRedisCache
	clients            map[string]sui.ISuiAPI
	errLogger          *log.Logger
}

const (
	default_page_size int = 10
)

var (
	min_region_staffs int64 = 2
)

func initializeCenterRequestService(
	platformConfigRepo i_repository.IPlatformConfigRepository,
	centerRequestRepo i_repository.ICenterRequestRepository,
	profileRepo i_repository.IProfileRepository,
	clients map[string]sui.ISuiAPI,
	errLogger *log.Logger,
) business.ICenterRequestService {
	return &centerRequestService{
		platformConfigRepo: platformConfigRepo,
		centerRequestRepo:  centerRequestRepo,
		profileRepo:        profileRepo,
		redisCache:         cache.InitializeRedisCache(),
		clients:            clients,
		errLogger:          errLogger,
	}
}

func GenerateCenterRequestService() (business.ICenterRequestService, error) {
	var errLogger = util.GetLogConfig(shared.ERROR_LEVEL)

	cnn, err := db.ConnectDB(errLogger, db.InitializePostgreSQL())
	if err != nil {
		return nil, err
	}

	return initializeCenterRequestService(
		repository.InitializePlatformConfigRepository(cnn, errLogger),
		repository.InitializeCenterRequestRepository(cnn, errLogger),
		repository.InitializeProfileRepository(cnn, errLogger),
		_networkAliases,
		errLogger,
	), nil
}

// ConfirmRequest implements business.ICenterRequestService.
func (c *centerRequestService) ConfirmRequest(id string, ctx context.Context) error {
	req, err := c.centerRequestRepo.GetRequest(id, ctx)
	if err != nil {
		return err
	}

	var genericErr error = errors.New(noti.GENERIC_ERROR_WARN_MSG)
	if req == nil {
		return genericErr
	}

	var client = c.clients[constant.SuiTestnet]
	manage, err := on_chain.GetOnChainObject[entities.Manage](on_chain.GetOnChainObjectRequest{
		Client:    client,
		ObjectId:  os.Getenv(env.MANAGE_OBJECT_ID),
		ErrLogger: c.errLogger,
	}, ctx)
	if err != nil {
		return err
	}

	var sender string = ctx.Value("address").(string)
	if req.CreatedBy != sender {
		return errors.New(noti.GENERIC_RIGHT_ACCESS_WARN_MSG)
	}

	var curTime time.Time = time.Now()
	if req.ClosedAt.After(curTime) {
		return errors.New(noti.STILL_PENDING_REQUEST_MESSAGE)
	}

	if len(req.Approvers) == 0 {
		return errors.New(noti.NOT_ARRPOVED_REQUEST_WARN_MSG)
	}

	var refuseCounts int
	if len(req.Refusers) > 0 {
		refuseCounts = len(req.Refusers)
	}

	var rate float32 = float32(len(req.Approvers)) / float32(len(req.Approvers)+refuseCounts)
	if rate < approve_rate_limit {
		req.Status = request_refused_status
	} else {
		req.Status = request_approved_status
	}

	req.UpdatedAt = curTime
	if err := c.centerRequestRepo.UpdateRegistrationRequest(*req, ctx); err != nil {
		return err
	}

	if req.Status == request_refused_status {
		return nil
	}

	var internalErr error = errors.New(noti.INTERNALL_ERR_MSG)
	if manage == nil {
		return internalErr
	}

	var startIdx int
	for i, region := range manage.LocalRegions {
		if region == req.Region {
			startIdx = i
			break
		}
	}

	leaders, err := on_chain.GetOnChainObjects[entities.StaffNft](on_chain.GetOnChainObjectsRequest{
		Client:    client,
		ObjectIds: manage.LocalLeaderNfts[startIdx:],
		ErrLogger: c.errLogger,
	}, ctx)
	if err != nil {
		return err
	}

	if leaders == nil {
		return internalErr
	}

	var regionLeaders []string
	for _, leader := range leaders {
		if leader.Region == req.Region {
			regionLeaders = append(regionLeaders, leader.Owner)
		}
	}

	var childModule = on_chain.InitializeModuleChild()
	var executeReq = on_chain.ExecuteTransactionRequestV2{
		Client:   client,
		Module:   childModule.GetModule(),
		Function: childModule.GetFunctionUploadCenter(),
		Arguments: childModule.ToCreateCenterArguments(on_chain.CreateCenterArguments{
			Region:      req.Region,
			Address:     req.Address,
			PhoneNumber: req.PhoneNumber,
			ImageBlobID: req.ImageBlobID,
			Leaders:     regionLeaders,
			Sender:      req.CreatedBy,
		}),
		ErrLogger: c.errLogger,
	}

	for i := 1; i <= 3; i++ {
		if _, err := on_chain.ExecuteTransactionV2(executeReq, ctx); err == nil {
			return nil
		}
	}

	return internalErr
}

// CreateRequest implements business.ICenterRequestService.
func (c *centerRequestService) CreateRequest(req request.CreateCenterRequest, ctx context.Context) (*entities.CenterRequest, error) {
	var client = c.clients[constant.SuiTestnet]
	manage, err := on_chain.GetOnChainObject[entities.Manage](on_chain.GetOnChainObjectRequest{
		Client:    client,
		ObjectId:  os.Getenv(env.MANAGE_OBJECT_ID),
		ErrLogger: c.errLogger,
	}, ctx)
	if err != nil {
		return nil, err
	}

	var internalErr error = errors.New(noti.INTERNALL_ERR_MSG)
	if manage == nil {
		return nil, internalErr
	}

	var sender string = ctx.Value("address").(string)
	var regionLeaderNft string
	for i, leader := range manage.LocalLeaderIds {
		if leader == sender {
			regionLeaderNft = manage.LocalLeaderNfts[i]
			break
		}
	}

	var genericRightErr error = errors.New(noti.GENERIC_RIGHT_ACCESS_WARN_MSG)
	if regionLeaderNft == "" {
		return nil, genericRightErr
	}

	senderLeaderNft, err := on_chain.GetOnChainObject[entities.StaffNft](on_chain.GetOnChainObjectRequest{
		Client:    client,
		ObjectId:  regionLeaderNft,
		ErrLogger: c.errLogger,
	}, ctx)
	if err != nil {
		return nil, err
	}

	if senderLeaderNft == nil {
		return nil, internalErr
	}

	if senderLeaderNft.Region != req.Region {
		return nil, errors.New(noti.NOT_STAFF_OF_REGION_MESSAGE)
	}

	staffs, err := on_chain.GetOnChainObjects[entities.StaffNft](on_chain.GetOnChainObjectsRequest{
		Client:    client,
		ObjectIds: manage.VolunteerNfts,
		ErrLogger: c.errLogger,
	}, ctx)
	if err != nil {
		return nil, err
	}

	if staffs == nil {
		return nil, internalErr
	}

	var regionStaffCounts int
	for _, staff := range staffs {
		if staff.Region == req.Region {
			regionStaffCounts += 1
			if regionStaffCounts == int(min_region_staffs) {
				break
			}
		}
	}

	regionStaffCounts += 1

	var minRegionStaffsAccepted int64
	var tmp *entities.PlatformConfig
	var numericTmp *entities.NumericConfig
	if config, _ := c.platformConfigRepo.GetConfigByKey(numericTmp.GetTable(), tmp.GetMinRegionStaffsConfigKey(), ctx); config != nil {
		minRegionStaffsAccepted = config.Value.(int64)
	} else {
		minRegionStaffsAccepted = min_region_staffs
	}

	if regionStaffCounts < int(minRegionStaffsAccepted) {
		return nil, errors.New("Not enough staffs")
	}

	// todo: AI validate address if matched with region
	var address string = strings.TrimSpace(req.Address)
	var phoneNumber string = strings.TrimSpace(req.PhoneNumber)
	var curTime time.Time = time.Now()
	var request = entities.CenterRequest{
		ID:          util.GenerateId(),
		ProfileID:   ctx.Value("sub").(string),
		Region:      req.Region,
		Address:     address,
		PhoneNumber: phoneNumber,
		ImageBlobID: req.ImageBlobID,
		Status:      request_pending_status,
		CreatedBy:   sender,
		CreatedAt:   curTime,
		UpdatedAt:   curTime,
		ClosedAt:    util.GetRequestDuration(),
	}

	return &request, c.centerRequestRepo.CreateRegistrationRequest(request, ctx)
}

// GetRequest implements business.ICenterRequestService.
func (c *centerRequestService) GetRequest(id string, ctx context.Context) (*entities.CenterRequest, error) {
	res, err := c.centerRequestRepo.GetRequest(id, ctx)
	if err != nil {
		return nil, err
	}

	return res, nil
}

// GetRequests implements business.ICenterRequestService.
func (c *centerRequestService) GetRequests(req request.GetCenterRequests, ctx context.Context) (response.PaginationDataResponse, error) {
	if req.Page < 1 {
		req.Page = 1
	}

	if req.PageSize < 1 {
		req.PageSize = default_page_size
	}

	var res response.PaginationDataResponse
	// var redisKey string = c.getGetCenterRequestsRedisKey(req)
	// if c.redisCache.Get(redisKey, &res, ctx) {
	// 	return res, nil
	// }

	data, pages, err := c.centerRequestRepo.GetRegistrationRequests(req, ctx)
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

	return res, err
}

// GetWalletRequests implements business.ICenterRequestService.
func (c *centerRequestService) GetWalletRequests(id string, ctx context.Context) ([]entities.CenterRequest, error) {
	if !util.IsValidSuiAddressStrict(id) {
		return nil, errors.New(noti.GENERIC_ERROR_WARN_MSG)
	}

	var res []entities.CenterRequest
	// var redisKey string = c.getGetWalletCenterRequestsRedisKey(id)
	// if c.redisCache.Get(redisKey, &res, ctx) {
	// 	return res, nil
	// }

	var errRes error
	res, errRes = c.centerRequestRepo.GetWalletRegistrationRequests(id, ctx)
	if errRes != nil {
		return nil, errRes
	}

	return res, errRes
}

// VoteRequest implements business.ICenterRequestService.
func (c *centerRequestService) VoteRequest(id string, req request.VoteRequest, ctx context.Context) error {
	request, err := c.centerRequestRepo.GetRequest(id, ctx)
	if err != nil {
		return err
	}

	var genericErr error = errors.New(noti.GENERIC_ERROR_WARN_MSG)
	if request == nil {
		return genericErr
	}

	var curTime time.Time = time.Now()
	if request.ClosedAt.Before(curTime) {
		return errors.New(noti.REQUEST_CLOSED_MESSAGE)
	}

	var voter string = ctx.Value("address").(string)
	if voter == request.CreatedBy {
		return errors.New(noti.OWNER_VOTE_WARN_MSG)
	}

	if slices.Contains(request.Approvers, voter) || slices.Contains(request.Refusers, voter) {
		return errors.New(noti.ALREADY_VOTE_MESSAGE)
	}

	var client = c.clients[constant.SuiTestnet]
	manage, err := on_chain.GetOnChainObject[entities.Manage](on_chain.GetOnChainObjectRequest{
		Client:    client,
		ObjectId:  os.Getenv(env.MANAGE_OBJECT_ID),
		ErrLogger: c.errLogger,
	}, ctx)
	if err != nil {
		return err
	}

	var internalErr error = errors.New(noti.INTERNALL_ERR_MSG)
	if manage == nil {
		return internalErr
	}

	var searchLen int
	if len(manage.VolunteerIds) > len(manage.LocalLeaderIds) {
		searchLen = len(manage.VolunteerIds)
	} else {
		searchLen = len(manage.LocalLeaderIds)
	}

	var staffId string
	for i := 0; i < searchLen; i++ {
		if i < len(manage.VolunteerIds) {
			if voter == manage.VolunteerIds[i] {
				staffId = manage.VolunteerNfts[i]
				break
			}
		}

		if i < len(manage.LocalLeaderIds) {
			if voter == manage.LocalLeaderIds[i] {
				staffId = manage.LocalLeaderNfts[i]
				break
			}
		}
	}

	if staffId == "" {
		return errors.New(noti.GENERIC_RIGHT_ACCESS_WARN_MSG)
	}

	nft, err := on_chain.GetOnChainObject[entities.StaffNft](on_chain.GetOnChainObjectRequest{
		Client:    client,
		ObjectId:  staffId,
		ErrLogger: c.errLogger,
	}, ctx)
	if err != nil {
		return err
	}

	if nft == nil {
		return internalErr
	}

	if nft.Region != request.Region {
		return errors.New(noti.NOT_STAFF_OF_REGION_MESSAGE)
	}

	if req.IsVoteYes {
		request.Approvers = append(request.Approvers, voter)
	} else {
		request.Refusers = append(request.Refusers, voter)
		if req.RefuseReason == "" {
			return errors.New(noti.FIELD_EMPTY_WARN_MSG)
		}

		request.RefuseReasons = append(request.RefuseReasons, strings.TrimSpace(req.RefuseReason))
	}

	request.UpdatedAt = curTime

	return c.centerRequestRepo.UpdateRegistrationRequest(*request, ctx)
}

// EditStaffNumbersToRequestCenter implements business.ICenterRequestService.
func (c *centerRequestService) EditStaffNumbersToRequestCenter(req request.EditStaffNumbersToCenterRequest, ctx context.Context) error {
	var module = on_chain.InitializeModuleManage()
	caps, err := on_chain.GetOnChainOwnedObjects[entities.Cap](on_chain.GetOnChainOwnedObjectsRequest{
		Client:       c.clients[constant.SuiTestnet],
		OwnerAddress: ctx.Value("address").(string),
		StructType:   fmt.Sprintf("%s::%s::%s", os.Getenv(env.PACKAGE_ID), module.GetModule(), module.GetAdminCapStruct()),
		ErrLogger:    c.errLogger,
	}, ctx)
	if err != nil {
		return err
	}

	if len(caps) == 0 {
		return errors.New(noti.GENERIC_RIGHT_ACCESS_WARN_MSG)
	}

	if req.MinStaffNumber != nil {
		if *req.MinStaffNumber > 0 {
			//min_region_staffs = *req.MinStaffNumber
		} else {

		}
	}

	return nil
}

func (c *centerRequestService) getGetCenterRequestsRedisKey(req request.GetCenterRequests) string {
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

	var isClosed string = "empty"
	if req.IsClosed != nil {
		isClosed = fmt.Sprintf("%v", *req.IsClosed)
	}

	return fmt.Sprintf("center_rq:kw:%s:r:%s:status:%s:close:%s:o:%s:s:%d:p:%d",
		keyword, region, status, isClosed, req.SortOrder, req.PageSize, req.Page)
}

func (c *centerRequestService) getGetWalletCenterRequestsRedisKey(wallet string) string {
	return fmt.Sprintf("center_rq:wallet:%s", wallet)
}

func (c *centerRequestService) getCenterRequestRedisKey(id string) string {
	return fmt.Sprintf("center_rq:%s", id)
}

var global_cur_time = time.Now()
