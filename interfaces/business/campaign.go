package business

import (
	"context"
	"raise-child/model/dtos/request"
	"raise-child/model/dtos/response"
	"raise-child/model/entities"
)

type ICampaignService interface {
	GetCampaigns(req request.GetCampaignsRequest, ctx context.Context) (response.PaginationDataResponse, error)
	GetCampaign(id string, ctx context.Context) (response.OnChainCampaignResponse, error)
	CreateCampaignWithdrawProposal(req request.CreateCampaignWithdrawProposalRequest, ctx context.Context) (*entities.PendingWithdrawProposal, error)
	SupportCampaign(id string, req request.SupportCampaignRequest, ctx context.Context) (response.PaymentUrlResponse, error)
}
