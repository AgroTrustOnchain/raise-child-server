package business

import (
	"context"
	"raise-child/model/dtos/request"
	"raise-child/model/dtos/response"
)

type IOnChainService interface {
	ExecuteTransaction(req request.ExecuteTransactionRequest, ctx context.Context) error
	BuildMoneyTransaction(req request.MoneyActionRequest, ctx context.Context) (response.BuildTransactionResponse, error)
}
