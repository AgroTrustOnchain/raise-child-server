package business

import (
	"context"
	"raise-child/model/dtos/request"
	"raise-child/model/dtos/response"
	"raise-child/model/entities"
)

type IVolunteerRequestService interface {
	GetRequests(req request.GetNormalStaffRegistrationRequests, ctx context.Context) (response.PaginationDataResponse, error)
	GetWalletRequests(id string, ctx context.Context) ([]entities.VolunteerRegistrationRequest, error)
	GetRequest(id string, ctx context.Context) (*entities.VolunteerRegistrationRequest, error)
	CreateRequest(req request.VolunteerRegistrationRequest, ctx context.Context) (*entities.VolunteerRegistrationRequest, error)
	VoteRequest(id string, req request.VoteRequest, ctx context.Context) error
	ConfirmRequest(id string, ctx context.Context) (response.BuildTransactionResponse, error)
}
