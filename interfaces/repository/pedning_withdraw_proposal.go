package repository

import (
	"context"
	"raise-child/model/dtos/request"
	"raise-child/model/entities"
)

type IPendingWithdrawProposalRepository interface {
	GetPendingWithdrawProposal(id string, ctx context.Context) (*entities.PendingWithdrawProposal, error)
	GetPendingWithdrawProposals(req request.GetPendingWithdrawProposalsRequest, ctx context.Context) ([]entities.PendingWithdrawProposal, int, error)
	CreatePendingWithdrawProposal(proposal entities.PendingWithdrawProposal, ctx context.Context) error
	UpdatePendingWithdrawProposal(proposal entities.PendingWithdrawProposal, ctx context.Context) error
	IsPendingWithdrawProposalProposedWithSpecificInfo(purpose, target, description string, withdrawAmount int64, ctx context.Context) (bool, error)
}
