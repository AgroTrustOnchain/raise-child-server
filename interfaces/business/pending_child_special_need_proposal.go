package business

import (
	"context"
	"raise-child/model/dtos/request"
	"raise-child/model/dtos/response"
	"raise-child/model/entities"
)

type IPendingChildSpecialNeedProposalService interface {
	GetPendingChildSpecialNeedProposals(req request.GetPendingChildSpecialNeedProposalsRequest, ctx context.Context) (response.PaginationDataResponse, error)
	GetPendingChildSpecialNeedProposal(id string, ctx context.Context) (*entities.PendingChildSpecialNeedProposal, error)
	ApprovePendingChildSpecialNeedProposal(id string, ctx context.Context) (response.BuildTransactionResponse, error)
	RefusePendingChildSpecialNeedProposal(id string, ctx context.Context) error
}
