package repository

import (
	"context"
	"raise-child/model/dtos/request"
	"raise-child/model/entities"
)

type IPendingChildSpecialNeedProposalRepository interface {
	GetPendingChildSpecialNeedProposals(req request.GetPendingChildSpecialNeedProposalsRequest, ctx context.Context) ([]entities.PendingChildSpecialNeedProposal, int, error)
	GetPendingChildSpecialNeedProposal(id string, ctx context.Context) (*entities.PendingChildSpecialNeedProposal, error)
	CreatePendingChildSpecialNeedProposal(proposal entities.PendingChildSpecialNeedProposal, ctx context.Context) error
	UpdatePendingChildSpecialNeedProposal(proposal entities.PendingChildSpecialNeedProposal, ctx context.Context) error
}
