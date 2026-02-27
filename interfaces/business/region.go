package business

import (
	"context"
	"raise-child/model/dtos/request"
	"raise-child/model/dtos/response"
	"raise-child/model/entities"
)

type IRegionService interface {
	GetRegions() response.RegionsResponse
	GetSupportedRegionProposals(req request.GetSupportedRegionProposalsRequest, ctx context.Context) (response.PaginationDataResponse, error)
	GetSupportedRegionProposal(id string, ctx context.Context) (*entities.SupportedRegionProposal, error)
	CreateSupportedRegionProposal(req request.CreateSupportedRegionProposalsRequest, ctx context.Context) (*entities.SupportedRegionProposal, error)
}
