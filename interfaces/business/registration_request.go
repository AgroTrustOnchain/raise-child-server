package business

import (
	"context"
	"raise-child/model/dtos/request"
	"raise-child/model/dtos/response"
	"raise-child/model/entities"
)

type IRegistrationRequestService interface {
	GetRegistrationRequests(req request.GetRegistrationRequests, ctx context.Context) (response.PaginationDataResponse, error)
	GetWalletRegistrationRequests(id string, ctx context.Context) ([]entities.RegistrationRequest, error)
	GetRegistrationRequest(id string, ctx context.Context) (*entities.RegistrationRequest, error)
	CreateRegistrationRequest(req request.CreateRegistrationRequest, ctx context.Context) (*entities.RegistrationRequest, error)
	VoteRegistrationRequest(id string, req request.VoteRequest, ctx context.Context) error
	ConfirmRegistrationRequest(id string, ctx context.Context) (response.BuildTransactionResponse, error)
}
