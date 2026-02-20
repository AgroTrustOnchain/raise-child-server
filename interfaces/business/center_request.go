package business

import (
	"context"
	"raise-child/model/dtos/request"
	"raise-child/model/dtos/response"
	"raise-child/model/entities"
)

type ICenterRequestService interface {
	GetRequests(req request.GetCenterRequests, ctx context.Context) (response.PaginationDataResponse, error)
	GetWalletRequests(id string, ctx context.Context) ([]entities.CenterRequest, error)
	GetRequest(id string, ctx context.Context) (*entities.CenterRequest, error)
	CreateRequest(req request.CreateCenterRequest, ctx context.Context) (*entities.CenterRequest, error)
	VoteRequest(id string, req request.VoteRequest, ctx context.Context) error
	ConfirmRequest(id string, ctx context.Context) (response.BuildTransactionResponse, error)
	EditStaffNumbersToRequestCenter(req request.EditStaffNumbersToCenterRequest, ctx context.Context) error
}
