package repository

import (
	"context"
	"raise-child/interfaces/repository"
	"raise-child/model/dtos/request"
	"raise-child/model/entities"

	"github.com/stretchr/testify/mock"
)

type paymentMockRepo struct {
	mock.Mock
}

// GetPayment implements repository.IPaymentRepository.
func (p *paymentMockRepo) GetPayment(id string, ctx context.Context) (*entities.Payment, error) {
	panic("unimplemented")
}

func InializePaymentMockRepo() repository.IPaymentRepository {
	return &paymentMockRepo{}
}

// CreatePayment implements repository.IPaymentRepository.
func (p *paymentMockRepo) CreatePayment(payment entities.Payment, ctx context.Context) error {
	var mockData = p.Called(payment, ctx)

	if mockFunc, ok := mockData.Get(0).(func(entities.Payment, context.Context) error); ok {
		return mockFunc(payment, ctx)
	}

	if err, ok := mockData.Error(0).(error); ok {
		return err
	}

	return nil
}

// GetPaymentById implements repository.IPaymentRepository.
func (p *paymentMockRepo) GetPaymentById(id string, ctx context.Context) (*entities.Payment, error) {
	var mockData = p.Called(id, ctx)

	var res1 *entities.Payment
	if mockFunc, ok := mockData.Get(0).(func(string, context.Context) *entities.Payment); ok {
		res1 = mockFunc(id, ctx)
	} else {
		res1 = mockData.Get(0).(*entities.Payment)
	}

	var res2 error
	if mockFunc, ok := mockData.Get(1).(func(string, context.Context) error); ok {
		res2 = mockFunc(id, ctx)
	} else {
		res2 = mockData.Error(1)
	}

	return res1, res2
}

// GetPayments implements repository.IPaymentRepository.
func (p *paymentMockRepo) GetPayments(req request.GetPaymentsRequest, ctx context.Context) ([]entities.Payment, int, error) {
	var mockData = p.Called(ctx)

	var res1 []entities.Payment
	if mockFunc, ok := mockData.Get(0).(func(request.GetPaymentsRequest, context.Context) []entities.Payment); ok {
		res1 = mockFunc(req, ctx)
	} else {
		res1 = mockData.Get(0).([]entities.Payment)
	}

	var res2 int
	if mockFunc, ok := mockData.Get(1).(func(request.GetPaymentsRequest, context.Context) int); ok {
		res2 = mockFunc(req, ctx)
	} else {
		res2 = mockData.Get(1).(int)
	}

	var res3 error
	if mockFunc, ok := mockData.Get(2).(func(request.GetPaymentsRequest, context.Context) error); ok {
		res3 = mockFunc(req, ctx)
	} else {
		res3 = mockData.Error(2)
	}

	return res1, res2, res3
}

// IsWithdrawalPaymentInProcess implements repository.IPaymentRepository.
func (p *paymentMockRepo) IsWithdrawalPaymentInProcess(id string, ctx context.Context) (bool, error) {
	var mockData = p.Called(id, ctx)

	var res1 bool
	if mockFunc, ok := mockData.Get(0).(func(string, context.Context) bool); ok {
		res1 = mockFunc(id, ctx)
	} else {
		res1 = mockData.Get(0).(bool)
	}

	var res2 error
	if mockFunc, ok := mockData.Get(1).(func(string, context.Context) error); ok {
		res2 = mockFunc(id, ctx)
	} else {
		res2 = mockData.Error(1)
	}

	return res1, res2
}

// UpdatePayment implements repository.IPaymentRepository.
func (p *paymentMockRepo) UpdatePayment(payment entities.Payment, ctx context.Context) error {
	var mockData = p.Called(payment, ctx)

	if mockFunc, ok := mockData.Get(0).(func(entities.Payment, context.Context) error); ok {
		return mockFunc(payment, ctx)
	}

	if err, ok := mockData.Error(0).(error); ok {
		return err
	}

	return nil
}
