package business

import (
	"context"
	"raise-child/model/dtos/request"
	"raise-child/model/dtos/response"
	"raise-child/model/entities"
)

type ITaskService interface {
	GetTasks(req request.GetTasksRequest, ctx context.Context) (response.PaginationDataResponse, error)
	GetUserTasks(wallet string, req request.GetTasksRequest, ctx context.Context) (response.PaginationDataResponse, error)
	GetTasksOfRegionOnUser(wallet string, req request.GetTasksRequest, ctx context.Context) (response.PaginationDataResponse, error)
	GetTask(id string, ctx context.Context) (*entities.Task, error)
	CreateTask(req request.CreateTaskRequest, ctx context.Context) (*entities.Task, error)
	ClaimTask(id string, ctx context.Context) error
	//ReviewAssignedProfileOfTask(id string, req request.VoteRequest, ctx context.Context) error
}

// ApproveTaskProof(id string, ctx context.Context) (response.BuildTransactionResponse, error)

type ITaskProofService interface {
	GetTaskProofs(req request.GetTaskProofsRequest, ctx context.Context) (response.PaginationDataResponse, error)
	GetTaskProof(id string, ctx context.Context) (*entities.TaskProof, error)
	SubmitTaskProof(id string, req request.SubmitTaskProofRequest, ctx context.Context) error
	ApproveTaskProof(id string, ctx context.Context) error
	RefuseTaskProof(id string, ctx context.Context) error
}
