package repository

import (
	"context"
	"raise-child/model/dtos/request"
	"raise-child/model/entities"
)

type ICenterRequestRepository interface {
	GetRequest(id string, ctx context.Context) (*entities.CenterRequest, error)
	GetRegistrationRequests(req request.GetCenterRequests, ctx context.Context) ([]entities.CenterRequest, int, error)
	GetWalletRegistrationRequests(id string, ctx context.Context) ([]entities.CenterRequest, error)
	CreateRegistrationRequest(req entities.CenterRequest, ctx context.Context) error
	UpdateRegistrationRequest(req entities.CenterRequest, ctx context.Context) error
	IsRegionRequested(region string, ctx context.Context) (bool, error)
	GetPendingRequests(ctx context.Context) ([]entities.BackgroundRecord, []entities.BackgroundRecord, error)
	SetApprovedStatuses(reqs []entities.BackgroundRecord, ctx context.Context) error
	SetRefusedStatuses(reqs []entities.BackgroundRecord, ctx context.Context) error
}
