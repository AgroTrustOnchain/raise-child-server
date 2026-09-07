package business

import (
	"context"
	"raise-child/model/dtos/request"
	"raise-child/model/dtos/response"
	"raise-child/model/entities"
)

// VoteUploadChildRequest(id string, req request.VoteRequest, ctx context.Context) error
// ConfirmUploadChildRequest(id string, ctx context.Context) (response.BuildTransactionResponse, error)

type IUploadChildRequestService interface {
	GetUploadChildRequests(req request.GetUploadChildRequests, ctx context.Context) (response.PaginationDataResponse, error)
	GetWalletUploadChildRequests(id string, page int, ctx context.Context) (response.PaginationDataResponse, error)
	GetUploadChildRequest(id string, ctx context.Context) (response.UploadChildRequestResponse, error)
	CreateUploadChildRequest(req request.UploadChildRequest, ctx context.Context) (*entities.UploadChildRequest, error)
	ReviewUploadChildRequest(id string, req request.VoteRequest, ctx context.Context) error
}
