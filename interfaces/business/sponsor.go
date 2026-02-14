package business

import (
	"context"
	"raise-child/model/dtos/request"
	"raise-child/model/dtos/response"
)

type IDonorService interface {
	GetDonors(req request.GetDonorsRequest, ctx context.Context) (response.PaginationDataResponse, error)
	GetDonor(id string, ctx context.Context) (response.DonorResponse, error)
}
