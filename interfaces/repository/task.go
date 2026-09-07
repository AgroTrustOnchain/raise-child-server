package repository

import (
	"context"
	"raise-child/model/dtos/request"
	"raise-child/model/entities"
)

type ITaskRepository interface {
	GetTasks(req request.GetTasksRequest, ctx context.Context) ([]entities.Task, int, error)
	GetTasksOfUser(req request.GetTasksRequest, ctx context.Context) ([]entities.TaskV2, int, error)
	GetTask(id string, ctx context.Context) (*entities.Task, error)
	CreateTask(task entities.Task, ctx context.Context) error
	UpdateTask(task entities.Task, ctx context.Context) error
}

type ITaskProofRepository interface {
	GetTaskProofs(req request.GetTaskProofsRequest, ctx context.Context) ([]entities.TaskProof, int, error)
	GetTaskProofsV2(req request.GetTaskProofsRequest, ctx context.Context) ([]entities.TaskProof, int, error)
	GetTaskProofsWithIsChildTask(req request.GetTaskProofsRequest, ctx context.Context) ([]entities.TaskProofWithIsChildTask, int, error)
	GetTaskProof(id string, ctx context.Context) (*entities.TaskProof, error)
	CreateTaskProof(proof entities.TaskProof, ctx context.Context) error
	UpdateTaskProof(proof entities.TaskProof, ctx context.Context) error
	IsTaskProofSumittedWithDetail(taskId, description, actorAddress, rawSubmitDate string, ctx context.Context) (bool, error)
}
