package business

import (
	"context"
	"raise-child/model/dtos/request"
)

type IOnChainService interface {
	ExecuteTransaction(req request.ExecuteTransactionRequest, ctx context.Context) error
}
