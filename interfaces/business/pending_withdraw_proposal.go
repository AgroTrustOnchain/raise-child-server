package business

import (
	"context"
	"raise-child/model/dtos/request"
	"raise-child/model/dtos/response"
	"raise-child/model/entities"
)

type IPendingWithdrawProposalService interface {
	GetPendingWithdrawProposal(id string, ctx context.Context) (*entities.PendingWithdrawProposal, error)
	GetPendingWithdrawProposals(req request.GetPendingWithdrawProposalsRequest, ctx context.Context) (response.PaginationDataResponse, error)
	CreatePendingWithdrawProposal(req request.CreatePendingWithdrawProposalRequest, ctx context.Context) (*entities.PendingWithdrawProposal, error)
	ApprovePendingWithdrawProposal(id string, ctx context.Context) (response.BuildTransactionResponse, error)
	RefusePendingWithdrawProposal(id string, ctx context.Context) error
}
