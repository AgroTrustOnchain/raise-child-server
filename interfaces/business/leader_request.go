package business

import (
	"context"
	"raise-child/model/dtos/request"
	"raise-child/model/dtos/response"
	"raise-child/model/entities"
)

type ILocalLeaderRequestService interface {
	GetRequests(req request.GetNormalStaffRegistrationRequests, ctx context.Context) (response.PaginationDataResponse, error)
	GetWalletRequests(id string, ctx context.Context) ([]entities.LocalLeaderRegistrationRequest, error)
	GetRequest(id string, ctx context.Context) (*entities.LocalLeaderRegistrationRequest, error)
	CreateRequest(req request.CreateRegistrationRequest, ctx context.Context) (*entities.LocalLeaderRegistrationRequest, error)
	VoteRequest(id string, req request.VoteRequest, ctx context.Context) error
	ConfirmRequest(id string, ctx context.Context) (response.BuildTransactionResponse, error)
}
