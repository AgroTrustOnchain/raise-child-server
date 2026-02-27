package repository

import (
	"context"
	"raise-child/model/dtos/request"
	"raise-child/model/entities"
)

type ISupportedRegionProposalRepository interface {
	GetSupportedRegionProposals(req request.GetSupportedRegionProposalsRequest, ctx context.Context) ([]entities.SupportedRegionProposal, int, error)
	GetSupportedRegionProposal(id string, ctx context.Context) (*entities.SupportedRegionProposal, error)
	CreateSupportedRegionProposal(proposal entities.SupportedRegionProposal, ctx context.Context) error
	UpdateSupportedRegionProposal(proposal entities.SupportedRegionProposal, ctx context.Context) error
	IsRegionRequested(region string, ctx context.Context) (bool, error)
}
