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
	"strings"
	"time"

	"github.com/block-vision/sui-go-sdk/constant"
	"github.com/block-vision/sui-go-sdk/sui"
)

type taskProofService struct {
	profileRepo         i_repository.IProfileRepository
	taskProofRepo       i_repository.ITaskProofRepository
	childTaskDetailRepo i_repository.IChildTaskDetailRepository
	taskRepo            i_repository.ITaskRepository
	aiProvider          ai.IAiClientProvider
	walrusProvider      walrus_pkg.IWalrusProvider
	redisCache          cache.IRedisCache
	clients             map[string]sui.ISuiAPI
	errLogger           *log.Logger
}

func initializeTaskProofService(
	profileRepo i_repository.IProfileRepository,
	taskProofRepo i_repository.ITaskProofRepository,
	childTaskDetailRepo i_repository.IChildTaskDetailRepository,
	taskRepo i_repository.ITaskRepository,
	aiProvider ai.IAiClientProvider,
	walrusProvider walrus_pkg.IWalrusProvider,
	clients map[string]sui.ISuiAPI,
	errLogger *log.Logger,
) business.ITaskProofService {
	return &taskProofService{
		profileRepo:         profileRepo,
		taskProofRepo:       taskProofRepo,
		childTaskDetailRepo: childTaskDetailRepo,
		taskRepo:            taskRepo,
		aiProvider:          aiProvider,
		walrusProvider:      walrusProvider,
		redisCache:          cache.InitializeRedisCache(),
		clients:             clients,
		errLogger:           errLogger,
	}
}

func GenerateTaskProofService() (business.ITaskProofService, error) {
	var errLogger = util.GetLogConfig(shared.ERROR_LEVEL)

	cnn, err := db.ConnectDB(errLogger, db.InitializePostgreSQL())
	if err != nil {
		return nil, err
	}

	return initializeTaskProofService(
		repository.InitializeProfileRepository(cnn, errLogger),
		repository.InitializeTaskProofRepository(cnn, errLogger),
		repository.InitializeChildTaskDetailRepository(cnn, errLogger),
		repository.InitializeTaskRepository(cnn, errLogger),
		ai.InitializeAiProvider(errLogger),
		walrus_pkg.InitializeWalrusProvider(errLogger),
		_networkAliases,
		errLogger,
	), nil
}

// GetTaskProof implements business.ITaskProofService.
func (t *taskProofService) GetTaskProof(id string, ctx context.Context) (*entities.TaskProof, error) {
	return t.taskProofRepo.GetTaskProof(id, ctx)
}

// GetTaskProofs implements business.ITaskProofService.
func (t *taskProofService) GetTaskProofs(req request.GetTaskProofsRequest, ctx context.Context) (response.PaginationDataResponse, error) {
	req.SortOrder = util.StandardizeSortOrder(req.SortOrder)
	req.Keyword = strings.TrimSpace(req.Keyword)
	if req.Page < 1 {
		req.Page = 1
	}

	if req.PageSize < 1 {
		req.PageSize = default_page_size
	}

	var genericErr error = errors.New(noti.GENERIC_ERROR_WARN_MSG)
	if req.ActorAddress != "" {
		if !util.IsValidSuiAddressStrict(req.ActorAddress) {
			return response.PaginationDataResponse{}, genericErr
		}
	}

	if req.ReviewedBy != "" {
		if !util.IsValidSuiAddressStrict(req.ReviewedBy) {
			return response.PaginationDataResponse{}, genericErr
		}
	}

	var res response.PaginationDataResponse
	// var redisKey string = t.getGetTaskProofsRedisKey(req)
	// if t.redisCache.Get(redisKey, &res, ctx) {
	// 	return res, nil
	// }

	data, pages, err := t.taskProofRepo.GetTaskProofsWithIsChildTask(req, ctx)
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

	//t.redisCache.Set(redisKey, res, time.Minute*5, ctx)

	return res, nil

	// ////////////////////////
	// // MOCK DATA
	// req.SortOrder = util.StandardizeSortOrder(req.SortOrder)
	// req.Keyword = strings.TrimSpace(req.Keyword)
	// if req.Page < 1 {
	// 	req.Page = 1
	// }

	// if req.PageSize < 1 {
	// 	req.PageSize = default_page_size
	// }

	// var res response.PaginationDataResponse
	// var redisKey string = t.getGetTaskProofsRedisKey(req)
	// if t.redisCache.Get(redisKey, &res, ctx) {
	// 	return res, nil
	// }

	// var data []entities.TaskProof = mockTaskProofs[(req.Page-1)*req.PageSize : req.Page*req.PageSize]
	// res = response.PaginationDataResponse{
	// 	Data:       data,
	// 	Amount:     len(data),
	// 	Page:       req.Page,
	// 	TotalPages: int(math.Ceil(float64(len(mockTaskProofs)) / float64(req.PageSize))),
	// }

	// t.redisCache.Set(redisKey, res, time.Minute*5, ctx)

	// return res, nil
}

// ApproveTaskProof implements business.ITaskProofService.
func (t *taskProofService) ApproveTaskProof(id string, ctx context.Context) error {
	proof, err := t.taskProofRepo.GetTaskProof(id, ctx)
	if err != nil {
		return err
	}

	var genericErr error = errors.New(noti.GENERIC_ERROR_WARN_MSG)
	if proof == nil {
		return genericErr
	}

	if proof.ReviewStatus != request_pending_status {
		return errors.New(noti.TASK_PROOF_REVIEWED_MESSAGE)
	}

	var genericRightErr error = errors.New(noti.GENERIC_RIGHT_ACCESS_WARN_MSG)
	var sender string = ctx.Value("address").(string)
	if proof.ActorAddress == sender {
		return genericRightErr
	}

	task, err := t.taskRepo.GetTask(proof.TaskID, ctx)
	if err != nil {
		return err
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

	var foundIdx int = -1
	for i, leader := range manage.LocalLeaderIds {
		if leader == sender {
			foundIdx = i
			break
		}
	}

	if foundIdx == -1 {
		return genericRightErr
	}

	leaderNft, err := on_chain.GetOnChainObject[entities.StaffNft](on_chain.GetOnChainObjectRequest{
		Client:    client,
		ObjectId:  manage.LocalLeaderNfts[foundIdx],
		ErrLogger: t.errLogger,
	}, ctx)
	if err != nil {
		return err
	}

	if leaderNft.Region != task.Region {
		return genericRightErr
	}

	var args []interface{}
	var function string
	var childModule = on_chain.InitializeModuleChild()
	if task.IsChildTask {
		detail, err := t.childTaskDetailRepo.GetChildTaskDetail(*task.ChildTaskDetailID, ctx)
		if err != nil {
			return err
		}

		switch detail.Purpose {
		case string(entities.MEAL_NEED_PURPOSE):
			function = childModule.GetFunctionConfirmProvideMealForChildV2()
		case string(entities.BOOKS_NEED_PURPOSE):
			function = childModule.GetFunctionConfirmProvideBooksForChildV2()
		case string(entities.HEALTH_INSURANCE_NEED_PURPOSE):
			function = childModule.GetFunctionConfirmProvideHealthInsuranceForChildV2()
		}

		args = childModule.ToConfirmProvideNeedForChildArgumentsV2(on_chain.ConfirmProvideNeedForChildArgumentsV2{
			ChildID:     detail.ChildID,
			NeedID:      detail.Target,
			StaffNft:    leaderNft.ID.ID,
			ImageBlobID: proof.ImageWalrusBlobID,
			ProvideDate: proof.RawSubmitDate,
			Actor:       proof.ActorAddress,
			Sender:      sender,
		})
	} else {
		var centerId string
		for i, region := range manage.LocalRegions {
			if region == task.Region {
				centerId = manage.ChildrenCenters[i]
				break
			}
		}

		function = childModule.GetFunctionSubmitTask()
		args = childModule.ToSubmitTaskArguments(on_chain.SubmitTaskArguments{
			Center:      centerId,
			StaffNft:    leaderNft.ID.ID,
			Description: task.Description,
			ImageBlobID: proof.ImageWalrusBlobID,
			Actor:       proof.ActorAddress,
		})
	}

	proof.ReviewedBy = &sender
	proof.ReviewStatus = request_approved_status
	if err := t.taskProofRepo.UpdateTaskProof(*proof, ctx); err != nil {
		return err
	}

	var req = on_chain.ExecuteTransactionRequestV2{
		Client:    client,
		Module:    childModule.GetModule(),
		Function:  function,
		Arguments: args,
		ErrLogger: t.errLogger,
	}

	for i := 1; i <= 3; i++ {
		if _, err := on_chain.ExecuteTransactionV2(req, ctx); err == nil {
			return nil
		}
	}

	return internalErr
}

// RefuseTaskProof implements business.ITaskProofService.
func (t *taskProofService) RefuseTaskProof(id string, ctx context.Context) error {
	proof, err := t.taskProofRepo.GetTaskProof(id, ctx)
	if err != nil {
		return err
	}

	var genericErr error = errors.New(noti.GENERIC_ERROR_WARN_MSG)
	if proof == nil {
		return genericErr
	}

	if proof.ReviewStatus != request_pending_status {
		return errors.New(noti.TASK_PROOF_REVIEWED_MESSAGE)
	}

	task, err := t.taskRepo.GetTask(proof.TaskID, ctx)
	if err != nil {
		return err
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

	var sender string = ctx.Value("address").(string)
	var foundIdx int = -1
	for i, leader := range manage.LocalLeaderIds {
		if leader == sender {
			foundIdx = i
			break
		}
	}

	var genericRightErr error = errors.New(noti.GENERIC_RIGHT_ACCESS_WARN_MSG)
	if foundIdx == -1 {
		return genericRightErr
	}

	leaderNft, err := on_chain.GetOnChainObject[entities.StaffNft](on_chain.GetOnChainObjectRequest{
		Client:    client,
		ObjectId:  manage.LocalLeaderNfts[foundIdx],
		ErrLogger: t.errLogger,
	}, ctx)
	if err != nil {
		return err
	}

	if leaderNft.Region != task.Region {
		return genericRightErr
	}

	profile, err := t.profileRepo.GetProfile(ctx.Value("sub").(string), ctx)
	if err != nil {
		return err
	}

	if profile.Status == "Suspended" {
		return errors.New(noti.CURRENTLY_SUSPENDED_MESSAGE)
	}

	proof.ReviewedBy = &sender
	proof.ReviewStatus = request_refused_status

	return t.taskProofRepo.UpdateTaskProof(*proof, ctx)
}

// SubmitTaskProof implements business.ITaskProofService.
func (t *taskProofService) SubmitTaskProof(id string, req request.SubmitTaskProofRequest, ctx context.Context) error {
	task, err := t.taskRepo.GetTask(id, ctx)
	if err != nil {
		return err
	}

	var genericErr error = errors.New(noti.GENERIC_ERROR_WARN_MSG)
	if task == nil {
		return genericErr
	}

	var curTime time.Time = time.Now()
	if curTime.After(task.EndPeriod) {
		return errors.New(noti.TASK_ENDED_MESSAGE)
	}

	if curTime.Before(task.StartPeriod) || curTime.After(task.EndPeriod) {
		return errors.New(noti.TASK_NOT_SUBMITTED_DATE_MESSAGE)
	}

	var sender string = ctx.Value("address").(string)
	if task.AssignedStaff == nil {
		return errors.New(noti.TASK_NOT_CLAIMED_MESSAGE)
	} else {
		if *task.AssignedStaff != sender {
			return errors.New(noti.TASK_NOT_OF_STAFF_MESSAGE)
		}
	}

	var rawSubmitDate string = util.TimeToRawDate(curTime)
	isSubmitted, err := t.taskProofRepo.IsTaskProofSumittedWithDetail(id, task.Description, sender, rawSubmitDate, ctx)
	if err != nil {
		return err
	}

	if isSubmitted {
		return errors.New(noti.TASK_PROOF_SUBMITTED_MESSAGE)
	}

	// TODO: AI Evaluation
	taskProof := ai.ValidateTaskProof{
		TaskDescription: task.Description,
		ProofImageURL:   req.ImageUrl,
		CreatedAt:       curTime,
	}

	var aiResponse *ai.ValidateTaskProofResponse = &ai.ValidateTaskProofResponse{
		AIEvaluation: "uncertain",
		AIReason:     "AI validation temporarily unavailable, please wait for human review",
	}

	if task.IsChildTask {
		if childTaskDetail, err := t.childTaskDetailRepo.GetChildTaskDetail(*task.ChildTaskDetailID, ctx); err == nil {
			child, _ := on_chain.GetOnChainObject[entities.Child](on_chain.GetOnChainObjectRequest{
				Client:    t.clients[constant.SuiTestnet],
				ObjectId:  childTaskDetail.ChildID,
				ErrLogger: t.errLogger,
			}, ctx)

			if child != nil {
				aiResponse, err = t.aiProvider.ValidateTaskProof(taskProof, ctx)
				if err != nil {
					t.errLogger.Printf("Failed to evaluate task proof with AI: %v", err)
				}
			}
		}
	}

	return t.taskProofRepo.CreateTaskProof(entities.TaskProof{
		ID:                    util.GenerateId(),
		TaskID:                id,
		Description:           task.Description,
		ActorProfileID:        ctx.Value("sub").(string),
		ActorAddress:          sender,
		ImageWalrusBlobID:     t.walrusProvider.UploadImageUrlToWalrus(req.ImageUrl),
		ImageCloudinaryBlobID: req.ImageCloudinaryBlobID,
		AIEvaluation:          aiResponse.AIEvaluation,
		AIReason:              aiResponse.AIReason,
		RawSubmitDate:         rawSubmitDate,
		CreatedAt:             curTime,
		UpdatedAt:             curTime,
	}, ctx)
}

func (t *taskProofService) getGetTaskProofsRedisKey(req request.GetTaskProofsRequest) string {
	var keyword string = "empty"
	if req.Keyword != "" {
		keyword = req.Keyword
	}

	var status string = "empty"
	if req.Status != "" {
		status = req.Status
	}

	var actorAddress string = "empty"
	if req.ActorAddress != "" {
		actorAddress = req.ActorAddress
	}

	var reviewedBy string = "empty"
	if req.ReviewedBy != "" {
		reviewedBy = req.ReviewedBy
	}

	return fmt.Sprintf("task:kw:%s:status:%s:actor:%s:reviewed:%s:o:%s:s:%d:p:%d",
		keyword, status, actorAddress, reviewedBy, req.SortOrder, req.PageSize, req.Page)
}
