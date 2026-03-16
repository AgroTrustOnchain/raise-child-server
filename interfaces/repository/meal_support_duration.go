package repository

import (
	"context"
	"raise-child/model/entities"
)

type IMealSupportDurationRepository interface {
	GetMealSupportDuration(id string, ctx context.Context) (*entities.OffChainMealSupportDuration, error)
	CreateMealSupportDuration(duration entities.OffChainMealSupportDuration, ctx context.Context) error
}
