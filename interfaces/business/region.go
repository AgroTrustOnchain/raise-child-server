package business

import (
	"context"
	"raise-child/model/dtos/request"
	"raise-child/model/dtos/response"
	"raise-child/model/entities"
)

type IRegionService interface {
	GetRegions() response.RegionsResponse
	GetEstablishedRegions(ctx context.Context) (response.RegionsResponse, error)
	GetRegionDetail(region string, req request.GetChildrenFromRegionDetailRequest, ctx context.Context) (response.RegionDetailResponse, error)
	GetSupportedRegionSuggestions(req request.GetSupportedRegionSuggestionsRequest, ctx context.Context) (response.PaginationDataResponse, error)
	GetWalletSupportedRegionSuggestions(req request.GetSupportedRegionSuggestionsRequest, ctx context.Context) (response.PaginationDataResponse, error)
	AdminGetSupportedRegionSuggestions(req request.GetSupportedRegionSuggestionsRequest, ctx context.Context) (response.PaginationDataResponse, error)
	GetSupportedRegionSuggestion(id string, ctx context.Context) (*entities.SupportedRegionSuggestion, error)
	CreateSupportedRegionSuggestion(req request.CreateSupportedRegionSuggestionsRequest, ctx context.Context) (*entities.SupportedRegionSuggestion, error)
	ReviewRegionSuggestion(id string, req request.VoteRequest, ctx context.Context) error
}
