package repository

import (
	"context"
	"raise-child/model/dtos/request"
	"raise-child/model/entities"
)

type ILocalLeaderRequestRepository interface {
	GetRequest(id string, ctx context.Context) (*entities.LocalLeaderRegistrationRequest, error)
	GetRegistrationRequests(req request.GetNormalStaffRegistrationRequests, ctx context.Context) ([]entities.LocalLeaderRegistrationRequest, int, error)
	GetWalletRegistrationRequests(id string, ctx context.Context) ([]entities.LocalLeaderRegistrationRequest, error)
	CreateRegistrationRequest(req entities.LocalLeaderRegistrationRequest, ctx context.Context) error
	UpdateRegistrationRequest(req entities.LocalLeaderRegistrationRequest, ctx context.Context) error
}
