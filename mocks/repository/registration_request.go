package repository

import (
	"context"
	"raise-child/interfaces/repository"
	"raise-child/model/dtos/request"
	"raise-child/model/entities"

	"github.com/stretchr/testify/mock"
)

type registrationRequestMockRepo struct {
	mock.Mock
}

func InializeRegistrationRequestMockRepo() repository.IRegistrationRequestRepository {
	return &registrationRequestMockRepo{}
}

// CreateRegistrationRequest implements repository.IRegistrationRequestRepository.
func (r *registrationRequestMockRepo) CreateRegistrationRequest(req entities.RegistrationRequest, ctx context.Context) error {
	var mockData = r.Called(req, ctx)

	if mockFunc, ok := mockData.Get(0).(func(entities.RegistrationRequest, context.Context) error); ok {
		return mockFunc(req, ctx)
	}

	if err, ok := mockData.Error(0).(error); ok {
		return err
	}

	return nil
}

// GetPendingRequests implements repository.IRegistrationRequestRepository.
func (r *registrationRequestMockRepo) GetPendingRequests(ctx context.Context) ([]entities.BackgroundRecord, []entities.BackgroundRecord, error) {
	var mockData = r.Called(ctx)

	var res1 []entities.BackgroundRecord
	if mockFunc, ok := mockData.Get(0).(func(context.Context) []entities.BackgroundRecord); ok {
		res1 = mockFunc(ctx)
	} else {
		res1 = mockData.Get(0).([]entities.BackgroundRecord)
	}

	var res2 []entities.BackgroundRecord
	if mockFunc, ok := mockData.Get(1).(func(context.Context) []entities.BackgroundRecord); ok {
		res2 = mockFunc(ctx)
	} else {
		res2 = mockData.Get(1).([]entities.BackgroundRecord)
	}

	var res3 error
	if mockFunc, ok := mockData.Get(2).(func(context.Context) error); ok {
		res3 = mockFunc(ctx)
	} else {
		res3 = mockData.Error(2)
	}

	return res1, res2, res3
}

// GetRegistrationRequest implements repository.IRegistrationRequestRepository.
func (r *registrationRequestMockRepo) GetRegistrationRequest(id string, ctx context.Context) (*entities.RegistrationRequest, error) {
	var mockData = r.Called(id, ctx)

	var res1 *entities.RegistrationRequest
	if mockFunc, ok := mockData.Get(0).(func(string, context.Context) *entities.RegistrationRequest); ok {
		res1 = mockFunc(id, ctx)
	} else {
		res1 = mockData.Get(0).(*entities.RegistrationRequest)
	}

	var res2 error
	if mockFunc, ok := mockData.Get(1).(func(string, context.Context) error); ok {
		res2 = mockFunc(id, ctx)
	} else {
		res2 = mockData.Error(1)
	}

	return res1, res2
}

// GetRegistrationRequests implements repository.IRegistrationRequestRepository.
func (r *registrationRequestMockRepo) GetRegistrationRequests(req request.GetRegistrationRequests, ctx context.Context) ([]entities.RegistrationRequest, int, error) {
	var mockData = r.Called(ctx)

	var res1 []entities.RegistrationRequest
	if mockFunc, ok := mockData.Get(0).(func(request.GetRegistrationRequests, context.Context) []entities.RegistrationRequest); ok {
		res1 = mockFunc(req, ctx)
	} else {
		res1 = mockData.Get(0).([]entities.RegistrationRequest)
	}

	var res2 int
	if mockFunc, ok := mockData.Get(1).(func(request.GetRegistrationRequests, context.Context) int); ok {
		res2 = mockFunc(req, ctx)
	} else {
		res2 = mockData.Get(1).(int)
	}

	var res3 error
	if mockFunc, ok := mockData.Get(2).(func(request.GetRegistrationRequests, context.Context) error); ok {
		res3 = mockFunc(req, ctx)
	} else {
		res3 = mockData.Error(2)
	}

	return res1, res2, res3
}

// GetRoleRegistrationRequests implements repository.IRegistrationRequestRepository.
func (r *registrationRequestMockRepo) GetRoleRegistrationRequests(role string, ctx context.Context) ([]entities.RegistrationRequest, error) {
	var mockData = r.Called(role, ctx)

	var res1 []entities.RegistrationRequest
	if mockFunc, ok := mockData.Get(0).(func(string, context.Context) []entities.RegistrationRequest); ok {
		res1 = mockFunc(role, ctx)
	} else {
		res1 = mockData.Get(0).([]entities.RegistrationRequest)
	}

	var res2 error
	if mockFunc, ok := mockData.Get(1).(func(string, context.Context) error); ok {
		res2 = mockFunc(role, ctx)
	} else {
		res2 = mockData.Error(1)
	}

	return res1, res2
}

// GetWalletRegistrationRequests implements repository.IRegistrationRequestRepository.
func (r *registrationRequestMockRepo) GetWalletRegistrationRequests(id string, ctx context.Context) ([]entities.RegistrationRequest, error) {
	panic("unimplemented")
}

// SetApprovedStatuses implements repository.IRegistrationRequestRepository.
func (r *registrationRequestMockRepo) SetApprovedStatuses(reqs []entities.BackgroundRecord, ctx context.Context) error {
	var mockData = r.Called(reqs, ctx)

	if mockFunc, ok := mockData.Get(0).(func([]entities.BackgroundRecord, context.Context) error); ok {
		return mockFunc(reqs, ctx)
	}

	if err, ok := mockData.Error(0).(error); ok {
		return err
	}

	return nil
}

// SetRefusedStatuses implements repository.IRegistrationRequestRepository.
func (r *registrationRequestMockRepo) SetRefusedStatuses(reqs []entities.BackgroundRecord, ctx context.Context) error {
	var mockData = r.Called(reqs, ctx)

	if mockFunc, ok := mockData.Get(0).(func([]entities.BackgroundRecord, context.Context) error); ok {
		return mockFunc(reqs, ctx)
	}

	if err, ok := mockData.Error(0).(error); ok {
		return err
	}

	return nil
}

// UpdateRegistrationRequest implements repository.IRegistrationRequestRepository.
func (r *registrationRequestMockRepo) UpdateRegistrationRequest(req entities.RegistrationRequest, ctx context.Context) error {
	var mockData = r.Called(req, ctx)

	if mockFunc, ok := mockData.Get(0).(func(entities.RegistrationRequest, context.Context) error); ok {
		return mockFunc(req, ctx)
	}

	if err, ok := mockData.Error(0).(error); ok {
		return err
	}

	return nil
}
