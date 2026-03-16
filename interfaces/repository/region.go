package repository

import (
	"context"
	"raise-child/model/dtos/request"
	"raise-child/model/entities"
)

type ISupportedRegionSuggestionRepository interface {
	GetSupportedRegionSuggestions(req request.GetSupportedRegionSuggestionsRequest, isGuestView bool, ctx context.Context) ([]entities.SupportedRegionSuggestion, int, error)
	GetSupportedRegionSuggestion(id string, ctx context.Context) (*entities.SupportedRegionSuggestion, error)
	CreateSupportedRegionSuggestion(proposal entities.SupportedRegionSuggestion, ctx context.Context) error
	UpdateSupportedRegionSuggestion(proposal entities.SupportedRegionSuggestion, ctx context.Context) error
	IsRegionRequested(region string, ctx context.Context) (bool, error)
}
