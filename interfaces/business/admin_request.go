package business

import (
	"context"
	"raise-child/model/dtos/request"
	"raise-child/model/dtos/response"
	"raise-child/model/entities"
)

type IAdminRequestService interface {
	GetRequests(req request.GetAdminRegistrationRequets, ctx context.Context) (response.PaginationDataResponse, error)
	GetWalletRequests(id string, ctx context.Context) ([]entities.AdminRegistrationRequest, error)
	GetRequest(id string, ctx context.Context) (*entities.AdminRegistrationRequest, error)
	CreateRequest(req request.AdminRegistrationRequest, ctx context.Context) (*entities.AdminRegistrationRequest, error)
	VoteRequest(id string, req request.VoteRequest, ctx context.Context) error
	ConfirmRequest(id string, ctx context.Context) (response.BuildTransactionResponse, error)
}
