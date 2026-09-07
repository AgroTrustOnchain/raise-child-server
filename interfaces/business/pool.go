package business

import (
	"context"
	"raise-child/model/dtos/response"
)

type IPoolService interface {
	GetLeaderPool(id string, ctx context.Context) (response.PoolResponse, error)
}
