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
	profileRepo         i_repository.IProfileRepository
	childTaskDetailRepo i_repository.IChildTaskDetailRepository
	taskRepo            i_repository.ITaskRepository
	redisCache          cache.IRedisCache
	clients             map[string]sui.ISuiAPI
	errLogger           *log.Logger
}

func initializeTaskService(
	profileRepo i_repository.IProfileRepository,
	childTaskDetailRepo i_repository.IChildTaskDetailRepository,
	taskRepo i_repository.ITaskRepository,
	clients map[string]sui.ISuiAPI,
	errLogger *log.Logger,
) business.ITaskService {
	return &taskService{
		profileRepo:         profileRepo,
		childTaskDetailRepo: childTaskDetailRepo,
		taskRepo:            taskRepo,
		redisCache:          cache.InitializeRedisCache(),
		clients:             clients,
		errLogger:           errLogger,
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
		repository.InitializeChildTaskDetailRepository(cnn, errLogger),
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

	var client = t.clients[constant.SuiTestnet]
	manage, err := on_chain.GetOnChainObject[entities.Manage](on_chain.GetOnChainObjectRequest{
		Client:    client,
		ObjectId:  os.Getenv(env.MANAGE_OBJECT_ID),
		ErrLogger: t.errLogger,
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

	var sender string = ctx.Value("address").(string)
	var staffId string
	for i := 0; i < searchLen; i++ {
		if i < len(manage.VolunteerIds) {
			if sender == manage.VolunteerIds[i] {
				staffId = manage.VolunteerNfts[i]
				break
			}
		}

		if i < len(manage.LocalLeaderIds) {
			if sender == manage.LocalLeaderIds[i] {
				staffId = manage.LocalLeaderNfts[i]
				break
			}
		}
	}

	var genericRightErr error = errors.New(noti.GENERIC_RIGHT_ACCESS_WARN_MSG)
	if staffId == "" {
		return genericRightErr
	}

	nft, err := on_chain.GetOnChainObject[entities.StaffNft](on_chain.GetOnChainObjectRequest{
		Client:    client,
		ObjectId:  staffId,
		ErrLogger: t.errLogger,
	}, ctx)
	if err != nil {
		return err
	}

	if nft == nil {
		return internalErr
	}

	if nft.Region != task.Region {
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
	manage, err := on_chain.GetOnChainObject[entities.Manage](on_chain.GetOnChainObjectRequest{
		Client:    client,
		ObjectId:  os.Getenv(env.MANAGE_OBJECT_ID),
		ErrLogger: t.errLogger,
	}, ctx)
	if err != nil {
		return nil, err
	}

	var internalErr error = errors.New(noti.INTERNALL_ERR_MSG)
	if manage == nil {
		return nil, internalErr
	}

	var isRegionEstablished bool = false
	for i, localRegion := range manage.LocalRegions {
		if localRegion == req.Region {
			isRegionEstablished = manage.CenterConfirmStatuses[i]
			break
		}
	}

	if !isRegionEstablished {
		return nil, errors.New(noti.REGION_NOT_ESTABLISHED_MESSAGE)
	}

	var sender string = ctx.Value("address").(string)
	if !slices.Contains(manage.AdminIds, sender) {
		var foundIdx int = -1
		for i, leader := range manage.LocalLeaderIds {
			if leader == sender {
				foundIdx = i
				break
			}
		}

		var genericRightErr error = errors.New(noti.GENERIC_RIGHT_ACCESS_WARN_MSG)
		if foundIdx == -1 {
			return nil, genericRightErr
		}

		nft, err := on_chain.GetOnChainObject[entities.StaffNft](on_chain.GetOnChainObjectRequest{
			Client:    client,
			ObjectId:  manage.LocalLeaderNfts[foundIdx],
			ErrLogger: t.errLogger,
		}, ctx)
		if err != nil {
			return nil, err
		}

		if nft == nil {
			return nil, err
		}

		if nft.Region != req.Region {
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
	var childTaskId *string
	var curTime time.Time = time.Now()
	if req.IsChildTask != nil {
		if *req.IsChildTask {
			if req.ChildID == nil || req.NeedID == nil {
				return nil, genericErr
			}

			if !util.IsValidSuiAddressStrict(*req.ChildID) || !util.IsValidSuiAddressStrict(*req.NeedID) {
				return nil, genericErr
			}

			child, err := on_chain.GetOnChainObject[entities.Child](on_chain.GetOnChainObjectRequest{
				Client:    client,
				ObjectId:  *req.ChildID,
				ErrLogger: t.errLogger,
			}, ctx)
			if err != nil {
				return nil, err
			}

			if child == nil {
				return nil, genericErr
			}

			if child.Region != req.Region {
				return nil, errors.New(noti.NOT_STAFF_OF_REGION_MESSAGE)
			}

			var taskPurpose string
			if slices.Contains(child.BooksNeeds, *req.NeedID) {
				taskPurpose = string(entities.BOOKS_NEED_PURPOSE)
			} else if child.HealthInsuranceNeed == *req.NeedID {
				taskPurpose = string(entities.HEALTH_INSURANCE_NEED_PURPOSE)
			}

			if taskPurpose == "" {
				return nil, genericErr
			}

			var childTaskDetailId string = util.GenerateId()
			if err := t.childTaskDetailRepo.CreateChildTaskDetail(entities.ChildTaskDetail{
				ID:        childTaskDetailId,
				ChildID:   *req.ChildID,
				Purpose:   taskPurpose,
				Target:    *req.NeedID,
				CreatedAt: curTime,
			}, ctx); err != nil {
				return nil, err
			}

			childTaskId = &childTaskDetailId
		}
	}

	var startPeriod time.Time = util.ToStartOfDate(util.RawDateToTime(strings.TrimSpace(req.StartPeriod)))
	if startPeriod.IsZero() {
		return nil, genericErr
	}

	var endPeriod time.Time = util.ToStartOfDate(util.RawDateToTime(strings.TrimSpace(req.EndPeriod)))
	if endPeriod.IsZero() || endPeriod.Before(startPeriod) {
		return nil, genericErr
	}

	var task = entities.Task{
		ID:                util.GenerateId(),
		IsChildTask:       *req.IsChildTask,
		ChildTaskDetailID: childTaskId,
		CreatedBy:         sender,
		Region:            req.Region,
		Description:       strings.TrimSpace(req.Description),
		StartPeriod:       startPeriod,
		EndPeriod:         endPeriod,
		CreatedAt:         curTime,
		UpdatedAt:         curTime,
	}

	return &task, t.taskRepo.CreateTask(task, ctx)
}

// GetTask implements business.ITaskService.
func (t *taskService) GetTask(id string, ctx context.Context) (*entities.Task, error) {
	return t.taskRepo.GetTask(id, ctx)
}

// GetTasks implements business.ITaskService.
func (t *taskService) GetTasks(req request.GetTasksRequest, ctx context.Context) (response.PaginationDataResponse, error) {
	req.SortOrder = util.StandardizeSortOrder(req.SortOrder)
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

	var res response.PaginationDataResponse
	// var redisKey string = t.getGetTasksRedisKey(req)
	// if t.redisCache.Get(redisKey, &res, ctx) {
	// 	return res, nil
	// }

	data, pages, err := t.taskRepo.GetTasks(req, ctx)
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

	// t.redisCache.Set(redisKey, res, time.Minute*5, ctx)

	return res, nil
}

// GetUserTasks implements business.ITaskService.
func (t *taskService) GetUserTasks(wallet string, req request.GetTasksRequest, ctx context.Context) (response.PaginationDataResponse, error) {
	if !util.IsValidSuiAddressStrict(wallet) {
		return response.PaginationDataResponse{}, errors.New(noti.GENERIC_ERROR_WARN_MSG)
	}

	req.SortOrder = util.StandardizeSortOrder(req.SortOrder)
	req.Keyword = strings.TrimSpace(req.Keyword)
	if req.Page < 1 {
		req.Page = 1
	}

	if req.PageSize < 1 {
		req.PageSize = default_page_size
	}

	req.AssignedStaff = wallet
	data, pages, err := t.taskRepo.GetTasksOfUser(req, ctx)
	if err != nil {
		return response.PaginationDataResponse{}, err
	}

	var amount int
	if len(data) == 0 {
		amount = 0
	} else {
		amount = len(data)
	}

	return response.PaginationDataResponse{
		Data:       data,
		Amount:     amount,
		Page:       req.Page,
		TotalPages: pages,
	}, nil
}

// GetTasksOfRegionOnUser implements business.ITaskService.
func (t *taskService) GetTasksOfRegionOnUser(wallet string, req request.GetTasksRequest, ctx context.Context) (response.PaginationDataResponse, error) {
	var genericErr error = errors.New(noti.GENERIC_ERROR_WARN_MSG)
	if !util.IsValidSuiAddressStrict(wallet) {
		return response.PaginationDataResponse{}, genericErr
	}

	var client = t.clients[constant.SuiTestnet]
	manage, err := on_chain.GetOnChainObject[entities.Manage](on_chain.GetOnChainObjectRequest{
		Client:    client,
		ObjectId:  os.Getenv(env.MANAGE_OBJECT_ID),
		ErrLogger: t.errLogger,
	}, ctx)
	if err != nil {
		return response.PaginationDataResponse{}, err
	}

	var searchLen int
	if len(manage.VolunteerIds) > len(manage.LocalLeaderIds) {
		searchLen = len(manage.VolunteerIds)
	} else {
		searchLen = len(manage.LocalLeaderIds)
	}

	var foundIdx int = -1
	var isVolunteer bool = true
	for i := 0; i < searchLen; i++ {
		if i < len(manage.VolunteerIds) {
			if manage.VolunteerIds[i] == wallet {
				foundIdx = i
				break
			}
		}

		if i < len(manage.LocalLeaderIds) {
			if manage.LocalLeaderIds[i] == wallet {
				foundIdx = i
				isVolunteer = false
				break
			}
		}
	}

	if foundIdx == -1 {
		return response.PaginationDataResponse{}, nil
	}

	var nftId string
	if isVolunteer {
		nftId = manage.VolunteerNfts[foundIdx]
	} else {
		nftId = manage.LocalLeaderNfts[foundIdx]
	}

	nft, err := on_chain.GetOnChainObject[entities.StaffNft](on_chain.GetOnChainObjectRequest{
		Client:    client,
		ObjectId:  nftId,
		ErrLogger: t.errLogger,
	}, ctx)
	if err != nil {
		return response.PaginationDataResponse{}, err
	}

	req.SortOrder = util.StandardizeSortOrder(req.SortOrder)
	req.Keyword = strings.TrimSpace(req.Keyword)
	if req.Page < 1 {
		req.Page = 1
	}

	if req.PageSize < 1 {
		req.PageSize = default_page_size
	}

	if req.AssignedStaff != "" {
		if !util.IsValidSuiAddressStrict(req.AssignedStaff) {
			return response.PaginationDataResponse{}, genericErr
		}
	}

	req.Region = nft.Region

	data, pages, err := t.taskRepo.GetTasks(req, ctx)
	if err != nil {
		return response.PaginationDataResponse{}, err
	}

	var amount int
	if len(data) == 0 {
		amount = 0
	} else {
		amount = len(data)
	}

	return response.PaginationDataResponse{
		Data:       data,
		Amount:     amount,
		Page:       req.Page,
		TotalPages: pages,
	}, nil
}

// // ReviewAssignedProfileOfTask implements business.ITaskService.
// func (t *taskService) ReviewAssignedProfileOfTask(id string, req request.VoteRequest, ctx context.Context) error {
// 	task, err := t.taskRepo.GetTask(id, ctx)
// 	if err != nil {
// 		return err
// 	}

// 	var genericErr error = errors.New(noti.GENERIC_ERROR_WARN_MSG)
// 	if task == nil {
// 		return genericErr
// 	}

// 	if time.Now().After(task.EndPeriod) {
// 		return errors.New(noti.TASK_ENDED_MESSAGE)
// 	}

// 	if task.AssignedStaff == nil {
// 		return errors.New(noti.TASK_NOT_CLAIMED_MESSAGE)
// 	}

// 	var sender string = ctx.Value("address").(string)
// 	var staffModule = on_chain.InitializeModuleStaff()
// 	var client = t.clients[constant.SuiTestnet]
// 	staffNfts, err := on_chain.GetOnChainOwnedObjects[entities.StaffNft](on_chain.GetOnChainOwnedObjectsRequest{
// 		Client:       client,
// 		OwnerAddress: sender,
// 		StructType:   fmt.Sprintf("%s::%s::%s", os.Getenv(env.PACKAGE_ID), staffModule.GetModule, staffModule.GetStaffNftObjectStruct()),
// 		ErrLogger:    t.errLogger,
// 	}, ctx)
// 	if err != nil {
// 		return err
// 	}

// 	var genericRightErr error = errors.New(noti.GENERIC_RIGHT_ACCESS_WARN_MSG)
// 	if staffNfts == nil || len(staffNfts) == 0 {
// 		return genericRightErr
// 	}

// 	var isLeaderOfRegion bool = false
// 	for _, nft := range staffNfts {
// 		if nft.Region == task.Region && nft.Role == local_leader_role {
// 			isLeaderOfRegion = true
// 			break
// 		}
// 	}

// 	if !isLeaderOfRegion {
// 		var manageObj entities.Manage
// 		if !t.redisCache.Get(manageObj.GetRedisKey(), &manageObj, ctx) {
// 			res, err := on_chain.GetOnChainObject[entities.Manage](on_chain.GetOnChainObjectRequest{
// 				Client:    t.clients[constant.SuiTestnet],
// 				ObjectId:  os.Getenv(env.MANAGE_OBJECT_ID),
// 				ErrLogger: t.errLogger,
// 			}, ctx)
// 			if err != nil {
// 				return err
// 			}

// 			if res != nil {
// 				t.redisCache.Set(manageObj.GetRedisKey(), res, time.Minute, ctx)
// 				manageObj = *res
// 			}
// 		}

// 		if !slices.Contains(manageObj.AdminIds, sender) {
// 			return genericRightErr
// 		}
// 	}

// 	profile, err := t.profileRepo.GetProfile(ctx.Value("sub").(string), ctx)
// 	if err != nil {
// 		return err
// 	}

// 	if profile.Status == "Suspended" {
// 		return errors.New(noti.CURRENTLY_SUSPENDED_MESSAGE)
// 	}

// 	var reviewStatus string
// 	if req.IsVoteYes {
// 		reviewStatus = request_approved_status
// 	} else {
// 		reviewStatus = request_refused_status
// 	}

// 	task.ReviewedBy = &sender
// 	task.ReviewProfileStatus = reviewStatus

// 	return t.taskRepo.UpdateTask(*task, ctx)
// }

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

	return fmt.Sprintf("task:kw:%s:r:%s:status:%s:assigned:%s:o:%s:s:%d:p:%d",
		keyword, region, status, assignedStaff, req.SortOrder, req.PageSize, req.Page)
}
