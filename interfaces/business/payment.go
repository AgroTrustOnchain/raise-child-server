package business

import (
	"context"
	"raise-child/model/dtos/request"
	"raise-child/model/dtos/response"
)

type IPaymentService interface {
	Donate(req request.DonateRequest, ctx context.Context) (response.UrlAPIResponse, error)
	Callback(id string, ctx context.Context) (string, error)
	CallbackWithAuth(id, capturedImgBlobId string, ctx context.Context) (string, error)
}
