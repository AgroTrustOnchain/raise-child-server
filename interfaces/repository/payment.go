package repository

import (
	"context"
	"raise-child/model/dtos/request"
	"raise-child/model/entities"
)

type IPaymentRepository interface {
	GetPayments(req request.GetPaymentsRequest, ctx context.Context) ([]entities.Payment, int, error)
	IsWithdrawalPaymentInProcess(id string, ctx context.Context) (bool, error)
	GetPaymentById(id string, ctx context.Context) (*entities.Payment, error)
	CreatePayment(payment entities.Payment, ctx context.Context) error
	UpdatePayment(payment entities.Payment, ctx context.Context) error
}
