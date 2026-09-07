package repository

import (
	"context"
	"raise-child/interfaces/repository"
	"raise-child/model/dtos/request"
	"raise-child/model/entities"
	"time"

	"github.com/stretchr/testify/mock"
)

type uploadChildRequestMockRepo struct {
	mock.Mock
}

// SetReviewStatus implements repository.IUploadChildRequestRepository.
func (u *uploadChildRequestMockRepo) SetReviewStatus(id string, reviewStatus string, reviewer string, closedAt *time.Time, ctx context.Context) error {
	panic("unimplemented")
}

func InializeUploadChildRequestMockRepo() repository.IUploadChildRequestRepository {
	return &uploadChildRequestMockRepo{}
}

// CreateUploadChildRequest implements repository.IUploadChildRequestRepository.
func (u *uploadChildRequestMockRepo) CreateUploadChildRequest(req entities.UploadChildRequest, ctx context.Context) error {
	var mockData = u.Called(req, ctx)

	if mockFunc, ok := mockData.Get(0).(func(entities.UploadChildRequest, context.Context) error); ok {
		return mockFunc(req, ctx)
	}

	if err, ok := mockData.Error(0).(error); ok {
		return err
	}

	return nil
}

// GetPendingRequests implements repository.IUploadChildRequestRepository.
func (u *uploadChildRequestMockRepo) GetPendingRequests(ctx context.Context) ([]entities.BackgroundRecord, []entities.BackgroundRecord, error) {
	var mockData = u.Called(ctx)

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

// GetUploadChildRequest implements repository.IUploadChildRequestRepository.
func (u *uploadChildRequestMockRepo) GetUploadChildRequest(id string, ctx context.Context) (*entities.UploadChildRequest, error) {
	var mockData = u.Called(id, ctx)

	var res1 *entities.UploadChildRequest
	if mockFunc, ok := mockData.Get(0).(func(string, context.Context) *entities.UploadChildRequest); ok {
		res1 = mockFunc(id, ctx)
	} else {
		res1 = mockData.Get(0).(*entities.UploadChildRequest)
	}

	var res2 error
	if mockFunc, ok := mockData.Get(1).(func(string, context.Context) error); ok {
		res2 = mockFunc(id, ctx)
	} else {
		res2 = mockData.Error(1)
	}

	return res1, res2
}

// GetUploadChildRequests implements repository.IUploadChildRequestRepository.
func (u *uploadChildRequestMockRepo) GetUploadChildRequests(req request.GetUploadChildRequests, ctx context.Context) ([]entities.UploadChildRequest, int, error) {
	var mockData = u.Called(ctx)

	var res1 []entities.UploadChildRequest
	if mockFunc, ok := mockData.Get(0).(func(request.GetUploadChildRequests, context.Context) []entities.UploadChildRequest); ok {
		res1 = mockFunc(req, ctx)
	} else {
		res1 = mockData.Get(0).([]entities.UploadChildRequest)
	}

	var res2 int
	if mockFunc, ok := mockData.Get(1).(func(request.GetUploadChildRequests, context.Context) int); ok {
		res2 = mockFunc(req, ctx)
	} else {
		res2 = mockData.Get(1).(int)
	}

	var res3 error
	if mockFunc, ok := mockData.Get(2).(func(request.GetUploadChildRequests, context.Context) error); ok {
		res3 = mockFunc(req, ctx)
	} else {
		res3 = mockData.Error(2)
	}

	return res1, res2, res3
}

// GetWalletUploadChildRequests implements repository.IUploadChildRequestRepository.
func (u *uploadChildRequestMockRepo) GetWalletUploadChildRequests(id string, page int, ctx context.Context) ([]entities.UploadChildRequest, int, error) {
	var mockData = u.Called(id, page, ctx)

	var res1 []entities.UploadChildRequest
	if mockFunc, ok := mockData.Get(0).(func(string, int, context.Context) []entities.UploadChildRequest); ok {
		res1 = mockFunc(id, page, ctx)
	} else {
		res1 = mockData.Get(0).([]entities.UploadChildRequest)
	}

	var res2 int
	if mockFunc, ok := mockData.Get(1).(func(string, int, context.Context) int); ok {
		res2 = mockFunc(id, page, ctx)
	} else {
		res2 = mockData.Get(1).(int)
	}

	var res3 error
	if mockFunc, ok := mockData.Get(1).(func(string, int, context.Context) error); ok {
		res3 = mockFunc(id, page, ctx)
	} else {
		res3 = mockData.Error(1)
	}

	return res1, res2, res3
}

// IsChildRequested implements repository.IUploadChildRequestRepository.
func (u *uploadChildRequestMockRepo) IsChildRequested(identityCode string, ctx context.Context) (bool, error) {
	var mockData = u.Called(identityCode, ctx)

	var res1 bool
	if mockFunc, ok := mockData.Get(0).(func(string, context.Context) bool); ok {
		res1 = mockFunc(identityCode, ctx)
	} else {
		res1 = mockData.Get(0).(bool)
	}

	var res2 error
	if mockFunc, ok := mockData.Get(1).(func(string, context.Context) error); ok {
		res2 = mockFunc(identityCode, ctx)
	} else {
		res2 = mockData.Error(1)
	}

	return res1, res2
}

// SetApprovedStatuses implements repository.IUploadChildRequestRepository.
func (u *uploadChildRequestMockRepo) SetApprovedStatuses(reqs []entities.BackgroundRecord, ctx context.Context) error {
	var mockData = u.Called(reqs, ctx)

	if mockFunc, ok := mockData.Get(0).(func([]entities.BackgroundRecord, context.Context) error); ok {
		return mockFunc(reqs, ctx)
	}

	if err, ok := mockData.Error(0).(error); ok {
		return err
	}

	return nil
}

// SetRefusedStatuses implements repository.IUploadChildRequestRepository.
func (u *uploadChildRequestMockRepo) SetRefusedStatuses(reqs []entities.BackgroundRecord, ctx context.Context) error {
	var mockData = u.Called(reqs, ctx)

	if mockFunc, ok := mockData.Get(0).(func([]entities.BackgroundRecord, context.Context) error); ok {
		return mockFunc(reqs, ctx)
	}

	if err, ok := mockData.Error(0).(error); ok {
		return err
	}

	return nil
}

// UpdateUploadChildRequest implements repository.IUploadChildRequestRepository.
func (u *uploadChildRequestMockRepo) UpdateUploadChildRequest(req entities.UploadChildRequest, ctx context.Context) error {
	var mockData = u.Called(req, ctx)

	if mockFunc, ok := mockData.Get(0).(func(entities.UploadChildRequest, context.Context) error); ok {
		return mockFunc(req, ctx)
	}

	if err, ok := mockData.Error(0).(error); ok {
		return err
	}

	return nil
}
