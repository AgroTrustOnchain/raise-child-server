package repository

import (
	"context"
	"raise-child/model/entities"

	"github.com/stretchr/testify/mock"
)

type bankProfileMockRepo struct {
	mock.Mock
}

func InializeBankProfileMockRepo() *bankProfileMockRepo {
	return &bankProfileMockRepo{}
}

// CreateBankProfile implements repository.IBankProfileRepository.
func (b *bankProfileMockRepo) CreateBankProfile(bp entities.BankProfile, ctx context.Context) error {
	var mockData = b.Called(bp, ctx)

	if mockFunc, ok := mockData.Get(0).(func(entities.BankProfile, context.Context) error); ok {
		return mockFunc(bp, ctx)
	}

	if err, ok := mockData.Error(0).(error); ok {
		return err
	}

	return nil
}

// GetBankProfileById implements repository.IBankProfileRepository.
func (b *bankProfileMockRepo) GetBankProfileById(id string, ctx context.Context) (*entities.BankProfile, error) {
	var mockData = b.Called(id, ctx)

	var res1 *entities.BankProfile
	if mockFunc, ok := mockData.Get(0).(func(string, context.Context) *entities.BankProfile); ok {
		res1 = mockFunc(id, ctx)
	} else {
		res1 = mockData.Get(0).(*entities.BankProfile)
	}

	var res2 error
	if mockFunc, ok := mockData.Get(1).(func(string, context.Context) error); ok {
		res2 = mockFunc(id, ctx)
	} else {
		res2 = mockData.Error(1)
	}

	return res1, res2
}

// GetBankProfileByOwner implements repository.IBankProfileRepository.
func (b *bankProfileMockRepo) GetBankProfileByOwner(owner string, ctx context.Context) (*entities.BankProfile, error) {
	var mockData = b.Called(owner, ctx)

	var res1 *entities.BankProfile
	if mockFunc, ok := mockData.Get(0).(func(string, context.Context) *entities.BankProfile); ok {
		res1 = mockFunc(owner, ctx)
	} else {
		res1 = mockData.Get(0).(*entities.BankProfile)
	}

	var res2 error
	if mockFunc, ok := mockData.Get(1).(func(string, context.Context) error); ok {
		res2 = mockFunc(owner, ctx)
	} else {
		res2 = mockData.Error(1)
	}

	return res1, res2
}

// IsBankWithSubExist implements repository.IBankProfileRepository.
func (b *bankProfileMockRepo) IsBankWithSubExist(sub string, ctx context.Context) (bool, error) {
	var mockData = b.Called(sub, ctx)

	var res1 bool
	if mockFunc, ok := mockData.Get(0).(func(string, context.Context) bool); ok {
		res1 = mockFunc(sub, ctx)
	} else {
		res1 = mockData.Get(0).(bool)
	}

	var res2 error
	if mockFunc, ok := mockData.Get(1).(func(string, context.Context) error); ok {
		res2 = mockFunc(sub, ctx)
	} else {
		res2 = mockData.Error(1)
	}

	return res1, res2
}

// UpdateBankProfile implements repository.IBankProfileRepository.
func (b *bankProfileMockRepo) UpdateBankProfile(bp entities.BankProfile, ctx context.Context) error {
	var mockData = b.Called(bp, ctx)

	if mockFunc, ok := mockData.Get(0).(func(entities.BankProfile, context.Context) error); ok {
		return mockFunc(bp, ctx)
	}

	if err, ok := mockData.Error(0).(error); ok {
		return err
	}

	return nil
}
