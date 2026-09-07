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
	"raise-child/model/dtos/request"
	"raise-child/model/entities"
	"time"
)

type paymentRepo struct {
	db        *sql.DB
	errLogger *log.Logger
}

const (
	payment_table        string = "payments"
	payment_limit_record int    = 10
)

func InitializePaymentRepository(db *sql.DB, errLogger *log.Logger) repository.IPaymentRepository {
	return &paymentRepo{
		db:        db,
		errLogger: errLogger,
	}
}

// CreatePayment implements repository.IPaymentRepository.
func (p *paymentRepo) CreatePayment(payment entities.Payment, ctx context.Context) error {
	var query string = "INSERT INTO " + payment_table +
		" (id, actor, profile_id, proposal_id, donation_id, is_donate_tx, transaction_id, " +
		"amount, currency, status, method, cancel_reason, " +
		"message, expired_at, created_at, updated_at) " +
		"values ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16)"
	var errLogMsg string = fmt.Sprintf(noti.REPO_ERR_MSG, shared.PAYMENT_REPOSITORY) + "CreatePayment - "

	if _, err := p.db.ExecContext(ctx, query, payment.ID, payment.Actor, payment.ProfileID, payment.ProposalID, payment.DonationID, payment.IsDonateTx, payment.TransactionId,
		payment.Amount, payment.Currency, payment.Status, payment.Method, payment.CancelReason,
		payment.Message, payment.ExpiredAt, payment.CreatedAt, payment.UpdatedAt); err != nil {

		p.errLogger.Println(errLogMsg + err.Error())
		return errors.New(noti.INTERNALL_ERR_MSG)
	}

	return nil
}

// GetPaymentById implements repository.IPaymentRepository.
func (p *paymentRepo) GetPaymentById(id string, ctx context.Context) (*entities.Payment, error) {
	var query string = "SELECT * FROM " + payment_table + " WHERE id = $1"
	var errLogMsg string = fmt.Sprintf(noti.REPO_ERR_MSG, shared.PAYMENT_REPOSITORY) + "GetPaymentById - "

	var res entities.Payment
	if err := p.db.QueryRowContext(ctx, query, id).Scan(
		&res.ID, &res.Actor, &res.ProfileID, &res.ProposalID, &res.DonationID, &res.IsDonateTx, &res.TransactionId,
		&res.Amount, &res.Currency, &res.Status, &res.Method, &res.CancelReason, &res.Message, &res.ExpiredAt,
		&res.CreatedAt, &res.UpdatedAt, &res.ProofBlobID, &res.ReviewedBy, &res.ReviewStatus, &res.IsTransferred, &res.TransferredAt); err != nil {

		if err == sql.ErrNoRows {
			return nil, nil
		}

		p.errLogger.Println(errLogMsg + err.Error())
		return nil, errors.New(noti.INTERNALL_ERR_MSG)
	}

	return &res, nil
}

// GetPayments implements repository.IPaymentRepository.
func (p *paymentRepo) GetPayments(req request.GetPaymentsRequest, ctx context.Context) ([]entities.Payment, int, error) {
	var errLogMsg string = fmt.Sprintf(noti.REPO_ERR_MSG, shared.PAYMENT_REPOSITORY) + "GetPayments - "
	var internalErr error = errors.New(noti.INTERNALL_ERR_MSG)

	var queryCondition string
	var isHavePreviosCondition bool = false
	if req.Status != "" {
		queryCondition += fmt.Sprintf("status = '%s'", req.Status)
		isHavePreviosCondition = true
	}

	if req.Keyword != "" {
		if isHavePreviosCondition {
			queryCondition += " AND "
		}

		queryCondition += fmt.Sprintf("LOWER(message) LIKE LOWER('%s')", req.Keyword)
		isHavePreviosCondition = true
	}

	if req.Method != "" {
		if isHavePreviosCondition {
			queryCondition += " AND "
		}

		queryCondition += fmt.Sprintf("method = '%s'", req.Method)
		isHavePreviosCondition = true
	}

	if req.Actor != "" {
		if isHavePreviosCondition {
			queryCondition += " AND "
		}

		queryCondition += fmt.Sprintf("actor = '%s'", req.Actor)
		isHavePreviosCondition = true
	}

	if req.MinAmount != nil {
		if isHavePreviosCondition {
			queryCondition += " AND "
		}

		queryCondition += fmt.Sprintf("amount >= %d", *req.MinAmount)
		isHavePreviosCondition = true
	}

	if req.MaxAmount != nil {
		if isHavePreviosCondition {
			queryCondition += " AND "
		}

		queryCondition += fmt.Sprintf("amount <= %d", *req.MaxAmount)
		isHavePreviosCondition = true
	}

	if req.IsDonatePayment != nil {
		if isHavePreviosCondition {
			queryCondition += " AND "
		}

		queryCondition += fmt.Sprintf("is_donate_tx = %v", req.IsDonatePayment)
		isHavePreviosCondition = true
	}

	if req.IsPaymentExpired != nil {
		if isHavePreviosCondition {
			queryCondition += " AND "
		}

		var operation string = ">"
		if *req.IsPaymentExpired {
			operation = "<="
		}

		queryCondition += fmt.Sprintf("expired_at %s NOW()", operation)
	}

	var filterProp string = req.FilterProp
	if filterProp == "" {
		filterProp = "created_at"
	}

	var sortOrder string = req.SortOrder
	if sortOrder == "" {
		sortOrder = "DESC"
	}

	if isHavePreviosCondition {
		queryCondition += " "
	}

	var query string = generateRetrieveQuery(generateRetrieveQueryRequest{
		table:       payment_table,
		limitAmount: req.PageSize,
		condition:   queryCondition,
		order:       fmt.Sprintf(" ORDER BY %s %s", filterProp, sortOrder),
		page:        req.Page,
		isGetCount:  false,
	})

	rows, err := p.db.QueryContext(ctx, query)
	if err != nil {
		p.errLogger.Println(errLogMsg + err.Error())
		return nil, 0, internalErr
	}
	defer rows.Close()

	var res []entities.Payment
	for rows.Next() {
		var x entities.Payment
		if err := rows.Scan(
			&x.ID, &x.Actor, &x.ProfileID, &x.ProposalID, &x.DonationID, &x.IsDonateTx, &x.TransactionId,
			&x.Amount, &x.Currency, &x.Status, &x.Method, &x.CancelReason, &x.Message, &x.ExpiredAt,
			&x.CreatedAt, &x.UpdatedAt, &x.ProofBlobID, &x.ReviewedBy, &x.ReviewStatus, &x.IsTransferred, &x.TransferredAt); err != nil {

			p.errLogger.Println(errLogMsg + err.Error())
			return nil, 0, internalErr
		}

		res = append(res, x)
	}

	// Track total records in table
	var totalRecords int
	p.db.QueryRowContext(ctx, generateCountTotalRecordsQuery(payment_table, queryCondition)).Scan(&totalRecords)

	return res, caculateTotalPages(totalRecords, req.PageSize), nil
}

// UpdatePayment implements repository.IPaymentRepository.
func (p *paymentRepo) UpdatePayment(payment entities.Payment, ctx context.Context) error {
	var errLogMsg string = fmt.Sprintf(noti.REPO_ERR_MSG, shared.PAYMENT_REPOSITORY) + "UpdatePayment - "
	var query string = "UPDATE " + payment_table + " SET status = $1, method = $2, cancel_reason = $3, updated_at = $4, proof_blob_id = $5, reviewed_by = $6, review_status = $7, is_transferred = $8, transferred_at = $9 WHERE id = $10"

	res, err := p.db.ExecContext(ctx, query, payment.Status, payment.Method, payment.CancelReason, payment.UpdatedAt,
		payment.ProfileID, payment.ReviewedBy, payment.ReviewStatus, payment.IsTransferred, payment.TransferredAt, payment.ID)

	var internalErr error = errors.New(noti.INTERNALL_ERR_MSG)

	if err != nil {
		p.errLogger.Println(errLogMsg + err.Error())
		return internalErr
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		p.errLogger.Println(errLogMsg + err.Error())
		return internalErr
	}

	if rowsAffected == 0 {
		return fmt.Errorf(noti.UNDEFINED_OBJECT_WARN_MSG, payment_table)
	}

	return nil
}

// IsWithdrawalPaymentInProcess implements repository.IPaymentRepository.
func (p *paymentRepo) IsWithdrawalPaymentInProcess(id string, ctx context.Context) (bool, error) {
	var query string = "SELECT status, expired_at FROM " + payment_table + " WHERE proposal_id = $1 AND is_transferred = FALSE AND (status = 'Pending' OR status = 'Success')"
	var errLogMsg string = fmt.Sprintf(noti.REPO_ERR_MSG, shared.PAYMENT_REPOSITORY) + "IsWithdrawalPaymentInProcess - "
	var internalErr error = errors.New(noti.INTERNALL_ERR_MSG)

	rows, err := p.db.QueryContext(ctx, query, id)
	if err != nil {
		p.errLogger.Println(errLogMsg + err.Error())
		return false, internalErr
	}

	var curTime time.Time = time.Now()
	for rows.Next() {
		var status string
		var expiredAt time.Time
		if err := rows.Scan(&status, &expiredAt); err != nil {
			p.errLogger.Println(errLogMsg + err.Error())
			return false, internalErr
		}

		if status == "Pending" {
			if expiredAt.After(curTime) {
				return true, nil
			}
		} else {
			return true, nil
		}
	}

	return false, nil
}
