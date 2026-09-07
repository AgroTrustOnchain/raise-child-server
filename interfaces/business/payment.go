package business

import (
	"context"
	"raise-child/model/dtos/request"
	"raise-child/model/dtos/response"
	"raise-child/model/entities"
)

// ApprovePayment(id string, ctx context.Context) (response.BuildTransactionResponse, error)

type IPaymentService interface {
	GetPayments(req request.GetPaymentsRequest, ctx context.Context) (response.PaginationDataResponse, error)
	GetPayment(id string, ctx context.Context) (*entities.Payment, error)
	ApprovePayment(id string, ctx context.Context) error
	RefusePayment(id string, ctx context.Context) error
	Donate(req request.DonateRequest, ctx context.Context) (response.UrlAPIResponse, error)
	DonateV2(req request.DonateRequest, ctx context.Context) (response.PaymentUrlResponse, error)
	Callback(id string, ctx context.Context) (string, error)
	CallbackV2(id string, ctx context.Context) error
	CallbackWithAuth(id, capturedImgBlobId string, ctx context.Context) (string, error)
	CallbackWithAuthV2(id string, req request.PaymentAuthCallbackRequest, ctx context.Context) error
}
