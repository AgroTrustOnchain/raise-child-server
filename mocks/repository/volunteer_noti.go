package repository

import (
	"context"
	"raise-child/interfaces/repository"
	"raise-child/model/dtos/request"
	"raise-child/model/entities"

	"github.com/stretchr/testify/mock"
)

type volunteerNotiMockRepo struct {
	mock.Mock
}

func InializeVolunteerNotiMockRepo() repository.IVolunteerNotiRepository {
	return &volunteerNotiMockRepo{}
}

// AssignVolunteer implements repository.IVolunteerNotiRepository.
func (v *volunteerNotiMockRepo) AssignVolunteer(volunteer string, region string, ctx context.Context) error {
	var mockData = v.Called(volunteer, region, ctx)

	if mockFunc, ok := mockData.Get(0).(func(string, string, context.Context) error); ok {
		return mockFunc(volunteer, region, ctx)
	}

	if err, ok := mockData.Error(0).(error); ok {
		return err
	}

	return nil
}

// CreateNoti implements repository.IVolunteerNotiRepository.
func (v *volunteerNotiMockRepo) CreateNoti(noti entities.VolunteerNoti, ctx context.Context) error {
	var mockData = v.Called(noti, ctx)

	if mockFunc, ok := mockData.Get(0).(func(entities.VolunteerNoti, context.Context) error); ok {
		return mockFunc(noti, ctx)
	}

	if err, ok := mockData.Error(0).(error); ok {
		return err
	}

	return nil
}

// GetCurrentVolunteerNotis implements repository.IVolunteerNotiRepository.
func (v *volunteerNotiMockRepo) GetCurrentVolunteerNotis(req request.GetNotisRequest, volunteer string, ctx context.Context) ([]entities.VolunteerNoti, error) {
	var mockData = v.Called(req, volunteer, ctx)

	var res1 []entities.VolunteerNoti
	if mockFunc, ok := mockData.Get(0).(func(request.GetNotisRequest, string, context.Context) []entities.VolunteerNoti); ok {
		res1 = mockFunc(req, volunteer, ctx)
	} else {
		res1 = mockData.Get(0).([]entities.VolunteerNoti)
	}

	var res2 error
	if mockFunc, ok := mockData.Get(2).(func(request.GetNotisRequest, string, context.Context) error); ok {
		res2 = mockFunc(req, volunteer, ctx)
	} else {
		res2 = mockData.Error(2)
	}

	return res1, res2
}

// GetNoti implements repository.IVolunteerNotiRepository.
func (v *volunteerNotiMockRepo) GetNoti(id string, ctx context.Context) (*entities.VolunteerNoti, error) {
	var mockData = v.Called(id, ctx)

	var res1 *entities.VolunteerNoti
	if mockFunc, ok := mockData.Get(0).(func(string, context.Context) *entities.VolunteerNoti); ok {
		res1 = mockFunc(id, ctx)
	} else {
		res1 = mockData.Get(0).(*entities.VolunteerNoti)
	}

	var res2 error
	if mockFunc, ok := mockData.Get(1).(func(string, context.Context) error); ok {
		res2 = mockFunc(id, ctx)
	} else {
		res2 = mockData.Error(1)
	}

	return res1, res2
}
