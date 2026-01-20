package business

import (
	"context"
	"raise-child/model/dtos/request"
	"raise-child/model/dtos/response"
)

type IAdminService interface {
	GetAdmins(req request.GetAdminsRequest, ctx context.Context) (response.PaginationDataResponse, error)
	UpdatePublisherInfo(req request.UpdatePublisherInfoRequest, ctx context.Context) (response.BuildTransactionResponse, error)
}
