package repository

import (
	"context"
	"raise-child/interfaces/repository"
	"raise-child/model/dtos/request"
	"raise-child/model/entities"

	"github.com/stretchr/testify/mock"
)

type centerRequestMockRepo struct {
	mock.Mock
}

func InializeCenterRequestMockRepo() repository.ICenterRequestRepository {
	return &centerRequestMockRepo{}
}

// CreateRegistrationRequest implements repository.ICenterRequestRepository.
func (c *centerRequestMockRepo) CreateRegistrationRequest(req entities.CenterRequest, ctx context.Context) error {
	var mockData = c.Called(req, ctx)

	if mockFunc, ok := mockData.Get(0).(func(entities.CenterRequest, context.Context) error); ok {
		return mockFunc(req, ctx)
	}

	if err, ok := mockData.Error(0).(error); ok {
		return err
	}

	return nil
}

// GetPendingRequests implements repository.ICenterRequestRepository.
func (c *centerRequestMockRepo) GetPendingRequests(ctx context.Context) ([]entities.BackgroundRecord, []entities.BackgroundRecord, error) {
	var mockData = c.Called(ctx)

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

// GetRegistrationRequests implements repository.ICenterRequestRepository.
func (c *centerRequestMockRepo) GetRegistrationRequests(req request.GetCenterRequests, ctx context.Context) ([]entities.CenterRequest, int, error) {
	var mockData = c.Called(ctx)

	var res1 []entities.CenterRequest
	if mockFunc, ok := mockData.Get(0).(func(request.GetCenterRequests, context.Context) []entities.CenterRequest); ok {
		res1 = mockFunc(req, ctx)
	} else {
		res1 = mockData.Get(0).([]entities.CenterRequest)
	}

	var res2 int
	if mockFunc, ok := mockData.Get(1).(func(request.GetCenterRequests, context.Context) int); ok {
		res2 = mockFunc(req, ctx)
	} else {
		res2 = mockData.Get(1).(int)
	}

	var res3 error
	if mockFunc, ok := mockData.Get(2).(func(request.GetCenterRequests, context.Context) error); ok {
		res3 = mockFunc(req, ctx)
	} else {
		res3 = mockData.Error(2)
	}

	return res1, res2, res3
}

// GetRequest implements repository.ICenterRequestRepository.
func (c *centerRequestMockRepo) GetRequest(id string, ctx context.Context) (*entities.CenterRequest, error) {
	var mockData = c.Called(id, ctx)

	var res1 *entities.CenterRequest
	if mockFunc, ok := mockData.Get(0).(func(string, context.Context) *entities.CenterRequest); ok {
		res1 = mockFunc(id, ctx)
	} else {
		res1 = mockData.Get(0).(*entities.CenterRequest)
	}

	var res2 error
	if mockFunc, ok := mockData.Get(1).(func(string, context.Context) error); ok {
		res2 = mockFunc(id, ctx)
	} else {
		res2 = mockData.Error(1)
	}

	return res1, res2
}

// GetWalletRegistrationRequests implements repository.ICenterRequestRepository.
func (c *centerRequestMockRepo) GetWalletRegistrationRequests(id string, ctx context.Context) ([]entities.CenterRequest, error) {
	var mockData = c.Called(id, ctx)

	var res1 []entities.CenterRequest
	if mockFunc, ok := mockData.Get(0).(func(string, context.Context) []entities.CenterRequest); ok {
		res1 = mockFunc(id, ctx)
	} else {
		res1 = mockData.Get(0).([]entities.CenterRequest)
	}

	var res2 error
	if mockFunc, ok := mockData.Get(1).(func(string, context.Context) error); ok {
		res2 = mockFunc(id, ctx)
	} else {
		res2 = mockData.Error(1)
	}

	return res1, res2
}

// IsRegionRequested implements repository.ICenterRequestRepository.
func (c *centerRequestMockRepo) IsRegionRequested(region string, ctx context.Context) (bool, error) {
	var mockData = c.Called(region, ctx)

	var res1 bool
	if mockFunc, ok := mockData.Get(0).(func(string, context.Context) bool); ok {
		res1 = mockFunc(region, ctx)
	} else {
		res1 = mockData.Get(0).(bool)
	}

	var res2 error
	if mockFunc, ok := mockData.Get(1).(func(string, context.Context) error); ok {
		res2 = mockFunc(region, ctx)
	} else {
		res2 = mockData.Error(1)
	}

	return res1, res2
}

// SetApprovedStatuses implements repository.ICenterRequestRepository.
func (c *centerRequestMockRepo) SetApprovedStatuses(reqs []entities.BackgroundRecord, ctx context.Context) error {
	var mockData = c.Called(reqs, ctx)

	if mockFunc, ok := mockData.Get(0).(func([]entities.BackgroundRecord, context.Context) error); ok {
		return mockFunc(reqs, ctx)
	}

	if err, ok := mockData.Error(0).(error); ok {
		return err
	}

	return nil
}

// SetRefusedStatuses implements repository.ICenterRequestRepository.
func (c *centerRequestMockRepo) SetRefusedStatuses(reqs []entities.BackgroundRecord, ctx context.Context) error {
	var mockData = c.Called(reqs, ctx)

	if mockFunc, ok := mockData.Get(0).(func([]entities.BackgroundRecord, context.Context) error); ok {
		return mockFunc(reqs, ctx)
	}

	if err, ok := mockData.Error(0).(error); ok {
		return err
	}

	return nil
}

// UpdateRegistrationRequest implements repository.ICenterRequestRepository.
func (c *centerRequestMockRepo) UpdateRegistrationRequest(req entities.CenterRequest, ctx context.Context) error {
	var mockData = c.Called(req, ctx)

	if mockFunc, ok := mockData.Get(0).(func(entities.CenterRequest, context.Context) error); ok {
		return mockFunc(req, ctx)
	}

	if err, ok := mockData.Error(0).(error); ok {
		return err
	}

	return nil
}
