package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"raise-child/constants/noti"
	"raise-child/constants/shared"
	"raise-child/interfaces/repository"
	"raise-child/model/entities"
)

type transactionRepo struct {
	db        *sql.DB
	errLogger *log.Logger
}

const (
	tx_table         string = "transactions"
	tx_limit_records int    = 10
)

func InitializeTransactionRepository(db *sql.DB, errLogger *log.Logger) repository.ITransactionRepository {
	return &transactionRepo{
		db:        db,
		errLogger: errLogger,
	}
}

// CreateTransaction implements repository.ITransactionRepository.
func (t *transactionRepo) CreateTransaction(tx entities.Transaction, ctx context.Context) error {
	var query string = "INSERT INTO " + tx_table +
		" (id, actor_address, action_type, amount, message, coin_type, created_at) " +
		"values ($1, $2, $3, $4, $5, $6, $7)"

	var errLogMsg string = fmt.Sprintf(noti.REPO_ERR_MSG, shared.TX_REPOSITORY) + "CreateTransaction - "

	if _, err := t.db.ExecContext(ctx, query, tx.ID, tx.ActorAddress, tx.ActionType,
		tx.Amount, tx.Message, tx.CoinType, tx.CreatedAt); err != nil {

		t.errLogger.Println(errLogMsg + err.Error())
		return errors.New(noti.INTERNALL_ERR_MSG)
	}

	return nil
}

// GetTransactions implements repository.ITransactionRepository.
func (t *transactionRepo) GetTransactions(page int, actionType string, ctx context.Context) ([]entities.Transaction, int, error) {
	// var errLogMsg string = fmt.Sprintf(noti.REPO_ERR_MSG, shared.TX_REPOSITORY) + "GetTransactions - "

	// var retrieveRequest = generateRetrieveQueryRequest{
	// 	table:       tx_table,
	// 	limitAmount: tx_limit_records,
	// }

	panic("unimplemented")

}

// GetWalletTransactions implements repository.ITransactionRepository.
func (t *transactionRepo) GetWalletTransactions(page int, address string, actionType string, ctx context.Context) ([]entities.Transaction, int, error) {
	panic("unimplemented")
}

// GetTransaction implements repository.ITransactionRepository.
func (t *transactionRepo) GetTransactionById(id string, ctx context.Context) (*entities.Transaction, error) {
	var query string = "SELECT * FROM " + tx_table + " WHERE id = $1"
	var errLogMsg string = fmt.Sprintf(noti.REPO_ERR_MSG, shared.TX_REPOSITORY) + "GetTransactionById - "

	var res entities.Transaction
	if err := t.db.QueryRowContext(ctx, query, id).Scan(
		&res.ID, &res.ActorAddress, &res.ActionType, &res.Amount,
		&res.Message, &res.CoinType, &res.CreatedAt); err != nil {

		if err == sql.ErrNoRows {
			return nil, nil
		}

		t.errLogger.Println(errLogMsg + err.Error())
		return nil, errors.New(noti.INTERNALL_ERR_MSG)
	}

	return &res, nil
}
