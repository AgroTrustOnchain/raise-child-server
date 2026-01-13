package business

import (
	"context"
	"raise-child/model/dtos/request"
	"raise-child/model/dtos/response"
)

type IPaymentService interface {
	Donate(req request.DonateRequest, ctx context.Context) (string, error)
	CallbackTx(id string, ctx context.Context) (response.BuildTransactionResponse, error)
}
