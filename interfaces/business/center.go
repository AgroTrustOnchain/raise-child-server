package business

import (
	"context"
	"raise-child/model/dtos/request"
	"raise-child/model/dtos/response"
)

type ICenterService interface {
	GetCenters(req request.GetCentersRequest, ctx context.Context) (response.PaginationDataResponse, error)
}
