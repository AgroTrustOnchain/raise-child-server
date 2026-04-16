package repository

import (
	"context"
	"raise-child/model/entities"

	"github.com/stretchr/testify/mock"
)

type mealSupportDurationMockRepo struct {
	mock.Mock
}

func InitializeMealSupportDurationMockRepo() *mealSupportDurationMockRepo {
	return &mealSupportDurationMockRepo{}
}

// CreateMealSupportDuration implements repository.IMealSupportDurationRepository.
func (m *mealSupportDurationMockRepo) CreateMealSupportDuration(duration entities.OffChainMealSupportDuration, ctx context.Context) error {
	var mockData = m.Called(duration, ctx)

	if mockFunc, ok := mockData.Get(0).(func(entities.OffChainMealSupportDuration, context.Context) error); ok {
		return mockFunc(duration, ctx)
	}

	if err, ok := mockData.Error(0).(error); ok {
		return err
	}

	return nil
}

// GetMealSupportDuration implements repository.IMealSupportDurationRepository.
func (m *mealSupportDurationMockRepo) GetMealSupportDuration(id string, ctx context.Context) (*entities.OffChainMealSupportDuration, error) {
	var mockData = m.Called(id, ctx)

	var res1 *entities.OffChainMealSupportDuration
	if mockFunc, ok := mockData.Get(0).(func(string, context.Context) *entities.OffChainMealSupportDuration); ok {
		res1 = mockFunc(id, ctx)
	} else {
		res1 = mockData.Get(0).(*entities.OffChainMealSupportDuration)
	}

	var res2 error
	if mockFunc, ok := mockData.Get(1).(func(string, context.Context) error); ok {
		res2 = mockFunc(id, ctx)
	} else {
		res2 = mockData.Error(1)
	}

	return res1, res2
}
