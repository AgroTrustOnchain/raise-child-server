package repository

import (
	"context"
	"raise-child/model/entities"
)

type IBackgroundChildrenWithdrawProposalRequestRepository interface {
	GetCurrentPendingRequests(ctx context.Context) ([]entities.BackgroundChildrenWithdrawProposalsRequest, error)
	CreateRequest(req entities.BackgroundChildrenWithdrawProposalsRequest, ctx context.Context) error
	SetRequestsExecuted(reqs []entities.BackgroundChildrenWithdrawProposalsRequest, ctx context.Context) error
	IsRegionProposed(region string, ctx context.Context) (bool, error)
}
