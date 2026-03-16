package repository

import (
	"context"
	"raise-child/interfaces/repository"
	"raise-child/model/entities"

	"github.com/stretchr/testify/mock"
)

type profileMockRepo struct {
	mock.Mock
}

// GetProfileOfFirsts implements repository.IProfileRepository.
func (p *profileMockRepo) GetProfileOfFirsts(position int, ctx context.Context) (*entities.Profile, error) {
	panic("unimplemented")
}

func InializeProfileMockRepo() repository.IProfileRepository {
	return &profileMockRepo{}
}

// CreateProfile implements repository.IProfileRepository.
func (p *profileMockRepo) CreateProfile(pfl entities.Profile, ctx context.Context) error {
	var mockData = p.Called(pfl, ctx)

	if mockFunc, ok := mockData.Get(0).(func(entities.Profile, context.Context) error); ok {
		return mockFunc(pfl, ctx)
	}

	if err, ok := mockData.Error(0).(error); ok {
		return err
	}

	return nil
}

// GetFirstProfile implements repository.IProfileRepository.
func (p *profileMockRepo) GetFirstProfile(ctx context.Context) (*entities.Profile, error) {
	var mockData = p.Called(ctx)

	var res1 *entities.Profile
	if mockFunc, ok := mockData.Get(0).(func(context.Context) *entities.Profile); ok {
		res1 = mockFunc(ctx)
	} else {
		res1 = mockData.Get(0).(*entities.Profile)
	}

	var res2 error
	if mockFunc, ok := mockData.Get(1).(func(context.Context) error); ok {
		res2 = mockFunc(ctx)
	} else {
		res2 = mockData.Error(1)
	}

	return res1, res2
}

// GetProfile implements repository.IProfileRepository.
func (p *profileMockRepo) GetProfile(id string, ctx context.Context) (*entities.Profile, error) {
	var mockData = p.Called(id, ctx)

	var res1 *entities.Profile
	if mockFunc, ok := mockData.Get(0).(func(string, context.Context) *entities.Profile); ok {
		res1 = mockFunc(id, ctx)
	} else {
		res1 = mockData.Get(0).(*entities.Profile)
	}

	var res2 error
	if mockFunc, ok := mockData.Get(1).(func(string, context.Context) error); ok {
		res2 = mockFunc(id, ctx)
	} else {
		res2 = mockData.Error(1)
	}

	return res1, res2
}

// IsEmailRegistered implements repository.IProfileRepository.
func (p *profileMockRepo) IsEmailRegistered(email string, ctx context.Context) (bool, error) {
	var mockData = p.Called(email, ctx)

	var res1 bool
	if mockFunc, ok := mockData.Get(0).(func(string, context.Context) bool); ok {
		res1 = mockFunc(email, ctx)
	} else {
		res1 = mockData.Get(0).(bool)
	}

	var res2 error
	if mockFunc, ok := mockData.Get(1).(func(string, context.Context) error); ok {
		res2 = mockFunc(email, ctx)
	} else {
		res2 = mockData.Error(1)
	}

	return res1, res2
}

// IsPersonalInfoExist implements repository.IProfileRepository.
func (p *profileMockRepo) IsPersonalInfoExist(identityCode string, phoneNumber string, email string, ctx context.Context) (bool, error) {
	var mockData = p.Called(identityCode, phoneNumber, email, ctx)

	var res1 bool
	if mockFunc, ok := mockData.Get(0).(func(string, string, string, context.Context) bool); ok {
		res1 = mockFunc(identityCode, phoneNumber, email, ctx)
	} else {
		res1 = mockData.Get(0).(bool)
	}

	var res2 error
	if mockFunc, ok := mockData.Get(1).(func(string, string, string, context.Context) error); ok {
		res2 = mockFunc(identityCode, phoneNumber, email, ctx)
	} else {
		res2 = mockData.Error(1)
	}

	return res1, res2
}

// IsPhoneNumberRegistered implements repository.IProfileRepository.
func (p *profileMockRepo) IsPhoneNumberRegistered(phoneNumber string, ctx context.Context) (bool, error) {
	var mockData = p.Called(phoneNumber, ctx)

	var res1 bool
	if mockFunc, ok := mockData.Get(0).(func(string, context.Context) bool); ok {
		res1 = mockFunc(phoneNumber, ctx)
	} else {
		res1 = mockData.Get(0).(bool)
	}

	var res2 error
	if mockFunc, ok := mockData.Get(1).(func(string, context.Context) error); ok {
		res2 = mockFunc(phoneNumber, ctx)
	} else {
		res2 = mockData.Error(1)
	}

	return res1, res2
}

// Login implements repository.IProfileRepository.
func (p *profileMockRepo) Login(id string, token string, ctx context.Context) error {
	var mockData = p.Called(id, token, ctx)

	if mockFunc, ok := mockData.Get(0).(func(string, string, context.Context) error); ok {
		return mockFunc(id, token, ctx)
	}

	if err, ok := mockData.Error(0).(error); ok {
		return err
	}

	return nil
}

// Logout implements repository.IProfileRepository.
func (p *profileMockRepo) Logout(id string, ctx context.Context) error {
	var mockData = p.Called(id, ctx)

	if mockFunc, ok := mockData.Get(0).(func(string, context.Context) error); ok {
		return mockFunc(id, ctx)
	}

	if err, ok := mockData.Error(0).(error); ok {
		return err
	}

	return nil
}

// UploadProfile implements repository.IProfileRepository.
func (p *profileMockRepo) UploadProfile(pfl entities.Profile, ctx context.Context) error {
	var mockData = p.Called(pfl, ctx)

	if mockFunc, ok := mockData.Get(0).(func(entities.Profile, context.Context) error); ok {
		return mockFunc(pfl, ctx)
	}

	if err, ok := mockData.Error(0).(error); ok {
		return err
	}

	return nil
}
