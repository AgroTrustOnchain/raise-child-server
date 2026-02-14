package repository

import (
	"context"
	"raise-child/model/entities"
)

type IOffChainWithdrawProposalRepository interface {
	GetOffChainWithdrawProposal(id string, ctx context.Context) (*entities.OffChainWithdrawProposal, error)
	GetOffChainWithdrawProposalByProposal(id string, ctx context.Context) (*entities.OffChainWithdrawProposal, error)
	CreateOffChainWithdrawProposal(proposal entities.OffChainWithdrawProposal, ctx context.Context) error
	SetOnChainProposalIdAfterExecuteTx(id, proposalId string, ctx context.Context) error
}
