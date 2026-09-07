package repository

import (
	"context"
	"raise-child/interfaces/repository"
	"raise-child/model/dtos/request"
	"raise-child/model/entities"

	"github.com/stretchr/testify/mock"
)

type leaderNotiMockRepo struct {
	mock.Mock
}

// GetNotiByNeed implements repository.ILeaderNotiRepository.
func (l *leaderNotiMockRepo) GetNotiByNeed(id string, ctx context.Context) (*entities.LeaderNoti, error) {
	panic("unimplemented")
}

func InializeLeaderNotiMockRepo() repository.ILeaderNotiRepository {
	return &leaderNotiMockRepo{}
}

// AssignLeader implements repository.ILeaderNotiRepository.
func (l *leaderNotiMockRepo) AssignLeader(leader string, region string, ctx context.Context) error {
	var mockData = l.Called(leader, region, ctx)

	if mockFunc, ok := mockData.Get(0).(func(string, string, context.Context) error); ok {
		return mockFunc(leader, region, ctx)
	}

	if err, ok := mockData.Error(0).(error); ok {
		return err
	}

	return nil
}

// CreateNoti implements repository.ILeaderNotiRepository.
func (l *leaderNotiMockRepo) CreateNoti(noti entities.LeaderNoti, ctx context.Context) error {
	var mockData = l.Called(noti, ctx)

	if mockFunc, ok := mockData.Get(0).(func(entities.LeaderNoti, context.Context) error); ok {
		return mockFunc(noti, ctx)
	}

	if err, ok := mockData.Error(0).(error); ok {
		return err
	}

	return nil
}

// GetCurrentLeaderNotis implements repository.ILeaderNotiRepository.
func (l *leaderNotiMockRepo) GetCurrentLeaderNotis(req request.GetNotisRequest, leader string, ctx context.Context) ([]entities.LeaderNoti, error) {
	var mockData = l.Called(req, leader, ctx)

	var res1 []entities.LeaderNoti
	if mockFunc, ok := mockData.Get(0).(func(request.GetNotisRequest, string, context.Context) []entities.LeaderNoti); ok {
		res1 = mockFunc(req, leader, ctx)
	} else {
		res1 = mockData.Get(0).([]entities.LeaderNoti)
	}

	var res2 error
	if mockFunc, ok := mockData.Get(2).(func(request.GetNotisRequest, string, context.Context) error); ok {
		res2 = mockFunc(req, leader, ctx)
	} else {
		res2 = mockData.Error(2)
	}

	return res1, res2
}

// GetNoti implements repository.ILeaderNotiRepository.
func (l *leaderNotiMockRepo) GetNoti(id string, ctx context.Context) (*entities.LeaderNoti, error) {
	var mockData = l.Called(id, ctx)

	var res1 *entities.LeaderNoti
	if mockFunc, ok := mockData.Get(0).(func(string, context.Context) *entities.LeaderNoti); ok {
		res1 = mockFunc(id, ctx)
	} else {
		res1 = mockData.Get(0).(*entities.LeaderNoti)
	}

	var res2 error
	if mockFunc, ok := mockData.Get(1).(func(string, context.Context) error); ok {
		res2 = mockFunc(id, ctx)
	} else {
		res2 = mockData.Error(1)
	}

	return res1, res2
}

// GetNotiByMealNeed implements repository.ILeaderNotiRepository.
func (l *leaderNotiMockRepo) GetNotiByMealNeed(id string, ctx context.Context) (*entities.LeaderNoti, error) {
	var mockData = l.Called(id, ctx)

	var res1 *entities.LeaderNoti
	if mockFunc, ok := mockData.Get(0).(func(string, context.Context) *entities.LeaderNoti); ok {
		res1 = mockFunc(id, ctx)
	} else {
		res1 = mockData.Get(0).(*entities.LeaderNoti)
	}

	var res2 error
	if mockFunc, ok := mockData.Get(1).(func(string, context.Context) error); ok {
		res2 = mockFunc(id, ctx)
	} else {
		res2 = mockData.Error(1)
	}

	return res1, res2
}

// UpdateNoti implements repository.ILeaderNotiRepository.
func (l *leaderNotiMockRepo) UpdateNoti(noti entities.LeaderNoti, ctx context.Context) error {
	var mockData = l.Called(noti, ctx)

	if mockFunc, ok := mockData.Get(0).(func(entities.LeaderNoti, context.Context) error); ok {
		return mockFunc(noti, ctx)
	}

	if err, ok := mockData.Error(0).(error); ok {
		return err
	}

	return nil
}
