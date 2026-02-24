package repository

import (
	"context"
	"raise-child/interfaces/repository"
	"raise-child/model/entities"

	"github.com/stretchr/testify/mock"
)

type offChainDonationMockRepo struct {
	mock.Mock
}

func InializeOffChainDonationMockRepo() repository.IOffChainDonationRepository {
	return &offChainDonationMockRepo{}
}

// CreateDonation implements repository.IOffChainDonationRepository.
func (o *offChainDonationMockRepo) CreateDonation(donation entities.OffChainDonation, ctx context.Context) error {
	var mockData = o.Called(donation, ctx)

	if mockFunc, ok := mockData.Get(0).(func(entities.OffChainDonation, context.Context) error); ok {
		return mockFunc(donation, ctx)
	}

	if err, ok := mockData.Error(0).(error); ok {
		return err
	}

	return nil
}

// GetDonation implements repository.IOffChainDonationRepository.
func (o *offChainDonationMockRepo) GetDonation(id string, ctx context.Context) (*entities.OffChainDonation, error) {
	var mockData = o.Called(id, ctx)

	var res1 *entities.OffChainDonation
	if mockFunc, ok := mockData.Get(0).(func(string, context.Context) *entities.OffChainDonation); ok {
		res1 = mockFunc(id, ctx)
	} else {
		res1 = mockData.Get(0).(*entities.OffChainDonation)
	}

	var res2 error
	if mockFunc, ok := mockData.Get(1).(func(string, context.Context) error); ok {
		res2 = mockFunc(id, ctx)
	} else {
		res2 = mockData.Error(1)
	}

	return res1, res2
}
