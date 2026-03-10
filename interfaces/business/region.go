package business

import (
	"context"
	"raise-child/model/dtos/request"
	"raise-child/model/dtos/response"
	"raise-child/model/entities"
)

type IRegionService interface {
	GetRegions() response.RegionsResponse
	GetSupportedRegionSuggestions(req request.GetSupportedRegionSuggestionsRequest, ctx context.Context) (response.PaginationDataResponse, error)
	GetSupportedRegionSuggestion(id string, ctx context.Context) (*entities.SupportedRegionSuggestion, error)
	CreateSupportedRegionSuggestion(req request.CreateSupportedRegionSuggestionsRequest, ctx context.Context) (*entities.SupportedRegionSuggestion, error)
}
