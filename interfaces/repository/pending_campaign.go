package repository

import (
	"context"
	"raise-child/model/dtos/request"
	"raise-child/model/entities"
)

type IPendingCampaignRepository interface {
	GetPendingCampaigns(req request.GetPendingCampaignsRequest, ctx context.Context) ([]entities.PendingCampaign, int, error)
	GetPendingCampaign(id string, ctx context.Context) (*entities.PendingCampaign, error)
	CreatePendingCampaign(proposal entities.PendingCampaign, ctx context.Context) error
	UpdatePendingCampaign(proposal entities.PendingCampaign, ctx context.Context) error
}
