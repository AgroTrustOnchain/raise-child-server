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
	"raise-child/util"
	"time"
)

type backgroundChildrenWithdrawRequestRepo struct {
	db        *sql.DB
	errLogger *log.Logger
}

const (
	background_children_withdraw_request_table string = "background_children_withdraw_requests"
)

func InitializeBackgroundChildrenWithdrawRequestRepository(db *sql.DB, errLogger *log.Logger) repository.IBackgroundChildrenWithdrawProposalRequestRepository {
	return &backgroundChildrenWithdrawRequestRepo{
		db:        db,
		errLogger: errLogger,
	}
}

// CreateRequest implements repository.IBackgroundChildrenWithdrawProposalRequestRepository.
func (b *backgroundChildrenWithdrawRequestRepo) CreateRequest(req entities.BackgroundChildrenWithdrawProposalsRequest, ctx context.Context) error {
	var query string = "INSERT INTO " + background_children_withdraw_request_table +
		" (id, profile_id, actor_address, region, is_executed, raw_proposed_date , created_at, updated_at) " +
		"values ($1, $2, $3, $4, $5, $6, $7, $8)"

	var errLogMsg string = fmt.Sprintf(noti.REPO_ERR_MSG, shared.BACKGROUND_CHILDREN_WITHDRAW_REQUEST_REPOSITORY) + "CreateRequest - "

	if _, err := b.db.ExecContext(ctx, query, req.ID, req.ProfileID, req.ActorAddress, req.Region,
		req.IsExecuted, req.RawProposedDate, req.CreatedAt, req.UpdatedAt); err != nil {

		b.errLogger.Println(errLogMsg + err.Error())
		return errors.New(noti.INTERNALL_ERR_MSG)
	}

	return nil
}

// GetCurrentPendingRequests implements repository.IBackgroundChildrenWithdrawProposalRequestRepository.
func (b *backgroundChildrenWithdrawRequestRepo) GetCurrentPendingRequests(ctx context.Context) ([]entities.BackgroundChildrenWithdrawProposalsRequest, error) {
	var errLogMsg string = fmt.Sprintf(noti.REPO_ERR_MSG, shared.BACKGROUND_CHILDREN_WITHDRAW_REQUEST_REPOSITORY) + "GetCurrentPendingRequests - "
	var internalErr error = errors.New(noti.INTERNALL_ERR_MSG)

	var query string = "SELECT * FROM " + background_children_withdraw_request_table + " WHERE is_executed = FALSE AND created_at >= CURRENT_DATE AND created_at <= NOW()"
	rows, err := b.db.QueryContext(ctx, query)
	if err != nil {
		b.errLogger.Println(errLogMsg + err.Error())
		return nil, internalErr
	}
	defer rows.Close()

	var res []entities.BackgroundChildrenWithdrawProposalsRequest
	for rows.Next() {
		var x entities.BackgroundChildrenWithdrawProposalsRequest
		if err := rows.Scan(
			&x.ID, &x.ProfileID, &x.ActorAddress, &x.Region, &x.IsExecuted,
			&x.RawProposedDate, &x.CreatedAt, &x.UpdatedAt); err != nil {

			b.errLogger.Println(errLogMsg + err.Error())
			return nil, internalErr
		}

		res = append(res, x)
	}

	return res, nil
}

// SetRequestsExecuted implements repository.IBackgroundChildrenWithdrawProposalRequestRepository.
func (b *backgroundChildrenWithdrawRequestRepo) SetRequestsExecuted(reqs []entities.BackgroundChildrenWithdrawProposalsRequest, ctx context.Context) error {
	if len(reqs) == 0 {
		return nil
	}

	var query string = "UPDATE " + background_children_withdraw_request_table + " SET is_executed = TRUE, updated_at = $1 WHERE "
	for i, req := range reqs {
		query += fmt.Sprintf("id = '%s'", req.ID)
		if i < len(reqs)-1 {
			query += " OR "
		}
	}

	if _, err := b.db.ExecContext(ctx, query, time.Now()); err != nil {
		b.errLogger.Println(fmt.Sprintf(noti.REPO_ERR_MSG, shared.BACKGROUND_CHILDREN_WITHDRAW_REQUEST_REPOSITORY) + "SetRequestsExecuted - " + err.Error())
		return errors.New(noti.INTERNALL_ERR_MSG)
	}

	return nil
}

// IsRegionProposed implements repository.IBackgroundChildrenWithdrawProposalRequestRepository.
func (b *backgroundChildrenWithdrawRequestRepo) IsRegionProposed(region string, ctx context.Context) (bool, error) {
	var query string = "SELECT id FROM " + background_children_withdraw_request_table + " WHERE region = $1 AND raw_proposed_date = $2 AND is_executed = FALSE LIMIT 1"

	var id string
	if err := b.db.QueryRowContext(ctx, query, region, util.TimeToRawDate(time.Now())).Scan(&id); err != nil {
		if err == sql.ErrNoRows {
			return false, nil
		}

		b.errLogger.Println(fmt.Sprintf(noti.REPO_ERR_MSG, shared.BACKGROUND_CHILDREN_WITHDRAW_REQUEST_REPOSITORY) + "IsRegionProposed - " + err.Error())
		return false, errors.New(noti.INTERNALL_ERR_MSG)
	}

	return id != "", nil
}
