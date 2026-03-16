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

type taskService struct {
	profileRepo i_repository.IProfileRepository
	taskRepo    i_repository.ITaskRepository
	redisCache  cache.IRedisCache
	clients     map[string]sui.ISuiAPI
	errLogger   *log.Logger
}

func initializeTaskService(
	profileRepo i_repository.IProfileRepository,
	taskRepo i_repository.ITaskRepository,
	clients map[string]sui.ISuiAPI,
	errLogger *log.Logger,
) business.ITaskService {
	return &taskService{
		profileRepo: profileRepo,
		taskRepo:    taskRepo,
		redisCache:  cache.InitializeRedisCache(),
		clients:     clients,
		errLogger:   errLogger,
	}
}

func GenerateTaskService() (business.ITaskService, error) {
	var errLogger = util.GetLogConfig(shared.ERROR_LEVEL)

	cnn, err := db.ConnectDB(errLogger, db.InitializePostgreSQL())
	if err != nil {
		return nil, err
	}

	return initializeTaskService(
		repository.InitializeProfileRepository(cnn, errLogger),
		repository.InitializeTaskRepository(cnn, errLogger),
		_networkAliases,
		errLogger,
	), nil
}

// ClaimTask implements business.ITaskService.
func (t *taskService) ClaimTask(id string, ctx context.Context) error {
	task, err := t.taskRepo.GetTask(id, ctx)
	if err != nil {
		return err
	}

	var genericErr error = errors.New(noti.GENERIC_ERROR_WARN_MSG)
	if task == nil {
		return genericErr
	}

	if time.Now().After(task.EndPeriod) {
		return errors.New(noti.TASK_ENDED_MESSAGE)
	}

	if task.AssignedStaff != nil {
		return errors.New(noti.TASK_CLAIMED_MESSAGE)
	}

	var sender string = ctx.Value("address").(string)
	var staffModule = on_chain.InitializeModuleStaff()
	staffNfts, err := on_chain.GetOnChainOwnedObjects[entities.StaffNft](on_chain.GetOnChainOwnedObjectsRequest{
		Client:       t.clients[constant.SuiTestnet],
		OwnerAddress: sender,
		StructType:   fmt.Sprintf("%s::%s::%s", os.Getenv(env.PACKAGE_ID), staffModule.GetModule, staffModule.GetStaffNftObjectStruct()),
		ErrLogger:    t.errLogger,
	}, ctx)
	if err != nil {
		return err
	}

	if staffNfts == nil || len(staffNfts) == 0 {
		return errors.New(noti.GENERIC_RIGHT_ACCESS_WARN_MSG)
	}

	var isStaffOfRegion bool = false
	for _, nft := range staffNfts {
		if nft.Region == task.Region {
			isStaffOfRegion = true
			break
		}
	}

	if !isStaffOfRegion {
		return errors.New(noti.NOT_STAFF_OF_REGION_MESSAGE)
	}

	profile, err := t.profileRepo.GetProfile(ctx.Value("sub").(string), ctx)
	if err != nil {
		return err
	}

	if profile.Status == "Suspended" {
		return errors.New(noti.CURRENTLY_SUSPENDED_MESSAGE)
	}

	task.AssignedProfileID = &profile.ID
	task.AssignedStaff = &sender

	return t.taskRepo.UpdateTask(*task, ctx)
}

// CreateTask implements business.ITaskService.
func (t *taskService) CreateTask(req request.CreateTaskRequest, ctx context.Context) (*entities.Task, error) {
	var client = t.clients[constant.SuiTestnet]
	manageObj, err := on_chain.GetOnChainObject[entities.Manage](on_chain.GetOnChainObjectRequest{
		Client:    client,
		ObjectId:  os.Getenv(env.MANAGE_OBJECT_ID),
		ErrLogger: t.errLogger,
	}, ctx)
	if err != nil {
		return nil, err
	}

	var isRegionEstablished bool = false
	for i, localRegion := range manageObj.LocalRegions {
		if localRegion == req.Region {
			isRegionEstablished = manageObj.CenterConfirmStatuses[i]
			break
		}
	}

	if !isRegionEstablished {
		return nil, errors.New(noti.REGION_NOT_ESTABLISHED_MESSAGE)
	}

	var sender string = ctx.Value("address").(string)
	if !slices.Contains(manageObj.AdminIds, sender) {
		var staffModule = on_chain.InitializeModuleStaff()
		staffNfts, err := on_chain.GetOnChainOwnedObjects[entities.StaffNft](on_chain.GetOnChainOwnedObjectsRequest{
			Client:       client,
			OwnerAddress: sender,
			StructType:   fmt.Sprintf("%s::%s::%s", os.Getenv(env.PACKAGE_ID), staffModule.GetModule, staffModule.GetStaffNftObjectStruct()),
			ErrLogger:    t.errLogger,
		}, ctx)
		if err != nil {
			return nil, err
		}

		if staffNfts == nil || len(staffNfts) == 0 {
			return nil, errors.New(noti.GENERIC_RIGHT_ACCESS_WARN_MSG)
		}

		var isStaffOfRegion bool = false
		for _, nft := range staffNfts {
			if nft.Region == req.Region {
				isStaffOfRegion = true
				break
			}
		}

		if !isStaffOfRegion {
			return nil, errors.New(noti.NOT_STAFF_OF_REGION_MESSAGE)
		}
	}

	profile, err := t.profileRepo.GetProfile(ctx.Value("sub").(string), ctx)
	if err != nil {
		return nil, err
	}

	if profile.Status == "Suspended" {
		return nil, errors.New(noti.CURRENTLY_SUSPENDED_MESSAGE)
	}

	var genericErr error = errors.New(noti.GENERIC_ERROR_WARN_MSG)
	var startPeriod time.Time = util.ToStartOfDate(util.RawDateToTime(strings.TrimSpace(req.StartPeriod)))
	if startPeriod.IsZero() {
		return nil, genericErr
	}

	var endPeriod time.Time = util.ToStartOfDate(util.RawDateToTime(strings.TrimSpace(req.EndPeriod)))
	if endPeriod.IsZero() || endPeriod.Before(startPeriod) {
		return nil, genericErr
	}

	var curTime time.Time = time.Now()
	var task = entities.Task{
		ID:          util.GenerateId(),
		CreatedBy:   sender,
		Region:      req.Region,
		Description: strings.TrimSpace(req.Description),
		StartPeriod: startPeriod,
		EndPeriod:   endPeriod,
		CreatedAt:   curTime,
		UpdatedAt:   curTime,
	}

	return &task, t.taskRepo.CreateTask(task, ctx)
}

// GetTask implements business.ITaskService.
func (t *taskService) GetTask(id string, ctx context.Context) (*entities.Task, error) {
	return t.taskRepo.GetTask(id, ctx)
}

// GetTasks implements business.ITaskService.
func (t *taskService) GetTasks(req request.GetTasksRequest, ctx context.Context) (response.PaginationDataResponse, error) {
	req.SortOrder = util.StanderizeSortOrder(req.SortOrder)
	req.Keyword = strings.TrimSpace(req.Keyword)
	if req.Page < 1 {
		req.Page = 1
	}

	if req.PageSize < 1 {
		req.PageSize = default_page_size
	}

	var genericErr error = errors.New(noti.GENERIC_ERROR_WARN_MSG)
	if req.AssignedStaff != "" {
		if !util.IsValidSuiAddressStrict(req.AssignedStaff) {
			return response.PaginationDataResponse{}, genericErr
		}
	}

	if req.ReviewedBy != "" {
		if !util.IsValidSuiAddressStrict(req.ReviewedBy) {
			return response.PaginationDataResponse{}, genericErr
		}
	}

	var res response.PaginationDataResponse
	var redisKey string = t.getGetTasksRedisKey(req)
	if t.redisCache.Get(redisKey, &res, ctx) {
		return res, nil
	}

	data, pages, err := t.taskRepo.GetTasks(req, ctx)
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

	t.redisCache.Set(redisKey, res, time.Minute*5, ctx)

	return res, nil
}

// ReviewAssignedProfileOfTask implements business.ITaskService.
func (t *taskService) ReviewAssignedProfileOfTask(id string, req request.VoteRequest, ctx context.Context) error {
	task, err := t.taskRepo.GetTask(id, ctx)
	if err != nil {
		return err
	}

	var genericErr error = errors.New(noti.GENERIC_ERROR_WARN_MSG)
	if task == nil {
		return genericErr
	}

	if time.Now().After(task.EndPeriod) {
		return errors.New(noti.TASK_ENDED_MESSAGE)
	}

	if task.AssignedStaff == nil {
		return errors.New(noti.TASK_NOT_CLAIMED_MESSAGE)
	}

	var sender string = ctx.Value("address").(string)
	var staffModule = on_chain.InitializeModuleStaff()
	var client = t.clients[constant.SuiTestnet]
	staffNfts, err := on_chain.GetOnChainOwnedObjects[entities.StaffNft](on_chain.GetOnChainOwnedObjectsRequest{
		Client:       client,
		OwnerAddress: sender,
		StructType:   fmt.Sprintf("%s::%s::%s", os.Getenv(env.PACKAGE_ID), staffModule.GetModule, staffModule.GetStaffNftObjectStruct()),
		ErrLogger:    t.errLogger,
	}, ctx)
	if err != nil {
		return err
	}

	var genericRightErr error = errors.New(noti.GENERIC_RIGHT_ACCESS_WARN_MSG)
	if staffNfts == nil || len(staffNfts) == 0 {
		return genericRightErr
	}

	var isLeaderOfRegion bool = false
	for _, nft := range staffNfts {
		if nft.Region == task.Region && nft.Role == local_leader_role {
			isLeaderOfRegion = true
			break
		}
	}

	if !isLeaderOfRegion {
		manageObj, err := on_chain.GetOnChainObject[entities.Manage](on_chain.GetOnChainObjectRequest{
			Client:    client,
			ObjectId:  os.Getenv(env.MANAGE_OBJECT_ID),
			ErrLogger: t.errLogger,
		}, ctx)
		if err != nil {
			return err
		}

		if !slices.Contains(manageObj.AdminIds, sender) {
			return genericRightErr
		}
	}

	profile, err := t.profileRepo.GetProfile(ctx.Value("sub").(string), ctx)
	if err != nil {
		return err
	}

	if profile.Status == "Suspended" {
		return errors.New(noti.CURRENTLY_SUSPENDED_MESSAGE)
	}

	var reviewStatus string
	if req.IsVoteYes {
		reviewStatus = request_approved_status
	} else {
		reviewStatus = request_refused_status
	}

	task.ReviewedBy = &sender
	task.ReviewProfileStatus = reviewStatus

	return t.taskRepo.UpdateTask(*task, ctx)
}

func (t *taskService) getGetTasksRedisKey(req request.GetTasksRequest) string {
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

	var assignedStaff string = "empty"
	if req.AssignedStaff != "" {
		assignedStaff = req.AssignedStaff
	}

	var reviewedBy string = "empty"
	if req.ReviewedBy != "" {
		reviewedBy = req.ReviewedBy
	}

	return fmt.Sprintf("task:kw:%s:r:%s:status:%s:assigned:%s:reviewed:%s:o:%s:s:%d:p:%d",
		keyword, region, status, assignedStaff, reviewedBy, req.SortOrder, req.PageSize, req.Page)
}
