package business

import (
	"context"
	"raise-child/model/dtos/request"
	"raise-child/model/dtos/response"
)

type ISponsorService interface {
	GetSponsors(req request.GetSponsorsRequest, ctx context.Context) (response.PaginationDataResponse, error)
	GetSponsor(id string, ctx context.Context) (response.SponsorResponse, error)
}
