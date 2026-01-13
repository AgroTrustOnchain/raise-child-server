package repository

import (
	"context"
	"raise-child/model/entities"
)

type ITransactionRepository interface {
	GetTransactions(page int, actionType string, ctx context.Context) ([]entities.Transaction, int, error)
	GetWalletTransactions(page int, address string, actionType string, ctx context.Context) ([]entities.Transaction, int, error)
	GetTransactionById(id string, ctx context.Context) (*entities.Transaction, error)
	CreateTransaction(tx entities.Transaction, ctx context.Context) error
}
