package repository

import (
	"context"
	"raise-child/model/dtos/request"
	"raise-child/model/entities"
)

type IRegistrationRequestRepository interface {
	GetRegistrationRequest(id string, ctx context.Context) (*entities.RegistrationRequest, error)
	GetRegistrationRequestByOnchainId(id string, ctx context.Context) (*entities.RegistrationRequest, error)
	GetRegistrationRequestWithDetail(detail entities.GetRegistrationRequestWithDetail, ctx context.Context) (*entities.RegistrationRequest, error)
	GetRegistrationRequests(req request.GetRegistrationRequests, ctx context.Context) ([]entities.RegistrationRequest, int, error)
	GetRoleRegistrationRequests(role string, ctx context.Context) ([]entities.RegistrationRequest, error)
	GetWalletRegistrationRequests(id string, ctx context.Context) ([]entities.RegistrationRequest, error)
	CreateRegistrationRequest(req entities.RegistrationRequest, ctx context.Context) error
	UpdateRegistrationRequest(req entities.RegistrationRequest, ctx context.Context) error
	GetPendingRequests(ctx context.Context) ([]entities.BackgroundRecord, []entities.BackgroundRecord, error)
	GetPendingRequestsV2(ctx context.Context) ([]entities.RegistrationRequest, error)
	SetApprovedStatuses(reqs []entities.BackgroundRecord, ctx context.Context) error
	SetRefusedStatuses(reqs []entities.BackgroundRecord, ctx context.Context) error
	SetApprovedStatusesV2(reqs []entities.RegistrationRequest, ctx context.Context) error
	SetRefusedStatusesV2(reqs []entities.RegistrationRequest, ctx context.Context) error
}
