package repository

import (
	"context"
	"raise-child/model/dtos/request"
	"raise-child/model/entities"
)

type IVolunteerRequestRepository interface {
	GetRequest(id string, ctx context.Context) (*entities.VolunteerRegistrationRequest, error)
	GetRegistrationRequests(req request.GetNormalStaffRegistrationRequests, ctx context.Context) ([]entities.VolunteerRegistrationRequest, int, error)
	GetWalletRegistrationRequests(id string, ctx context.Context) ([]entities.VolunteerRegistrationRequest, error)
	CreateRegistrationRequest(req entities.VolunteerRegistrationRequest, ctx context.Context) error
	UpdateRegistrationRequest(req entities.VolunteerRegistrationRequest, ctx context.Context) error
}
