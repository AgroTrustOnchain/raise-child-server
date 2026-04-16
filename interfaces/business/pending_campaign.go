package business

import (
	"context"
	"raise-child/model/dtos/request"
	"raise-child/model/dtos/response"
	"raise-child/model/entities"
)

type IPendingCampaignService interface {
	GetPendingCampaigns(req request.GetPendingCampaignsRequest, ctx context.Context) (response.PaginationDataResponse, error)
	GetPendingCampaign(id string, ctx context.Context) (*entities.PendingCampaign, error)
	CreatePendingCampaign(req request.CreatePendingCampaignRequest, ctx context.Context) (*entities.PendingCampaign, error)
	ApprovePendingCampaign(id string, ctx context.Context) (response.BuildTransactionResponse, error)
	RefusePendingCampaign(id string, ctx context.Context) error
}
