package repository

import (
	"context"
	"raise-child/model/dtos/request"
	"raise-child/model/entities"
)

type IAdminRequestRepository interface {
	GetRequest(id string, ctx context.Context) (*entities.AdminRegistrationRequest, error)
	GetRegistrationRequests(req request.GetAdminRegistrationRequets, ctx context.Context) ([]entities.AdminRegistrationRequest, int, error)
	GetWalletRegistrationRequests(id string, ctx context.Context) ([]entities.AdminRegistrationRequest, error)
	CreateRegistrationRequest(req entities.AdminRegistrationRequest, ctx context.Context) error
	UpdateRegistrationRequest(req entities.AdminRegistrationRequest, ctx context.Context) error
}
