package repository

import (
	"context"
	"raise-child/model/entities"
)

type IChildTaskDetailRepository interface {
	GetChildTaskDetail(id string, ctx context.Context) (*entities.ChildTaskDetail, error)
	CreateChildTaskDetail(detail entities.ChildTaskDetail, ctx context.Context) error
}
