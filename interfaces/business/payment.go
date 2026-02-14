package business

import (
	"context"
	"raise-child/model/dtos/request"
	"raise-child/model/dtos/response"
)

type IPaymentService interface {
	Donate(req request.DonateRequest, ctx context.Context) (response.UrlAPIResponse, error)
	ConfirmWithdrawProposal(id string, ctx context.Context) (map[string]interface{}, error)
	//CallbackTx(id string, ctx context.Context) (response.BuildTransactionResponse, error)
	Callback(id string, ctx context.Context) (string, error)
	CallbackWithAuth(id, capturedImgBlobId string, ctx context.Context) (string, error)
}
