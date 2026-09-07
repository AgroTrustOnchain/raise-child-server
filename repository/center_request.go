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

	"github.com/lib/pq"
)

type centerRequestRepo struct {
	db        *sql.DB
	errLogger *log.Logger
}

const (
	center_request_table        string = "center_requests"
	center_request_limit_record int    = 10
)

func InitializeCenterRequestRepository(db *sql.DB, errLogger *log.Logger) repository.ICenterRequestRepository {
	return &centerRequestRepo{
		db:        db,
		errLogger: errLogger,
	}
}

// CreateRegistrationRequest implements repository.ICenterRequestRepository.
func (c *centerRequestRepo) CreateRegistrationRequest(req entities.CenterRequest, ctx context.Context) error {
	var query string = "INSERT INTO " + center_request_table +
		" (id, profile_id, region, address, phone_number, image_blob_id, " +
		"approvers, refusers, refuse_reasons, status, " +
		"created_by, created_at, updated_at, closed_at) " +
		"values ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14)"

	var errLogMsg string = fmt.Sprintf(noti.REPO_ERR_MSG, shared.CENTER_REQUEST_REPOSITORY) + "CreateRegistrationRequest - "

	// Original approach
	// if _, err := c.db.ExecContext(ctx, query, req.ID, req.Region, req.Address, req.PhoneNumber, req.ImageBlobID,
	// 	req.Approvers, req.Refusers, req.RefuseReasons, req.Status,
	// 	req.IsAvailableToConfirm, req.IsConfirmRegister,
	// 	req.CreatedBy, req.CreatedAt, req.UpdatedAt, req.ClosedAt); err != nil {

	// 	c.errLogger.Println(errLogMsg + err.Error())
	// 	return errors.New(noti.INTERNALL_ERR_MSG)
	// }

	if _, err := c.db.ExecContext(ctx, query, req.ID, req.ProfileID, req.Region, req.Address, req.PhoneNumber, req.ImageBlobID,
		pq.Array(req.Approvers), pq.Array(req.Refusers), pq.Array(req.RefuseReasons), req.Status,
		req.CreatedBy, req.CreatedAt, req.UpdatedAt, req.ClosedAt); err != nil {

		c.errLogger.Println(errLogMsg + err.Error())
		return errors.New(noti.INTERNALL_ERR_MSG)
	}

	return nil
}

// GetRegistrationRequests implements repository.ICenterRequestRepository.
func (c *centerRequestRepo) GetRegistrationRequests(req request.GetCenterRequests, ctx context.Context) ([]entities.CenterRequest, int, error) {
	var errLogMsg string = fmt.Sprintf(noti.REPO_ERR_MSG, shared.CENTER_REQUEST_REPOSITORY) + "GetRegistrationRequests - "
	var internalErr error = errors.New(noti.INTERNALL_ERR_MSG)

	var queryCondition string
	var isHavePreviosCondition bool = false
	if req.Region != "" {
		queryCondition += fmt.Sprintf("LOWER(region) = LOWER('%s')", req.Region)
		isHavePreviosCondition = true
	}

	if req.Status != "" {
		if isHavePreviosCondition {
			queryCondition += " AND "
		}

		queryCondition += fmt.Sprintf("LOWER(status) = LOWER('%s')", req.Status)
		isHavePreviosCondition = true
	}

	if req.Keyword != "" {
		if isHavePreviosCondition {
			queryCondition += " AND "
		}

		queryCondition += fmt.Sprintf("(LOWER(address) LIKE LOWER('%s') OR phone_number LIKE '%s')", req.Keyword, req.Keyword)
		isHavePreviosCondition = true
	}

	if req.IsClosed != nil {
		if isHavePreviosCondition {
			queryCondition += " AND "
		}

		var operation string = ">"
		if *req.IsClosed {
			operation = "<="
		}

		queryCondition += fmt.Sprintf("closed_at %s NOW()", operation)
	}

	if isHavePreviosCondition {
		queryCondition += " "
	}

	var order string = "DESC"
	if req.SortOrder != "" {
		order = req.SortOrder
	}

	var query string = generateRetrieveQuery(generateRetrieveQueryRequest{
		table:       center_request_table,
		limitAmount: req.PageSize,
		condition:   queryCondition,
		order:       " ORDER BY created_at " + order,
		page:        req.Page,
		isGetCount:  false,
	})

	rows, err := c.db.QueryContext(ctx, query)
	if err != nil {
		c.errLogger.Println(errLogMsg + err.Error())
		return nil, 0, internalErr
	}
	defer rows.Close()

	var res []entities.CenterRequest
	for rows.Next() {
		// var x entities.CenterRequest
		// if err := rows.Scan(
		// 	&x.ID, &x.Region, &x.Address, &x.PhoneNumber, &x.ImageBlobID,
		// 	&x.Approvers, &x.Refusers, &x.RefuseReasons, &x.Status,
		// 	&x.IsAvailableToConfirm, &x.IsConfirmRegister,
		// 	&x.CreatedBy, &x.CreatedAt, &x.UpdatedAt, &x.ClosedAt); err != nil {

		// 	c.errLogger.Println(errLogMsg + err.Error())
		// 	return nil, 0, internalErr
		// }

		// res = append(res, x)

		var x entities.CenterRequest
		if err := rows.Scan(
			&x.ID, &x.ProfileID, &x.Region, &x.Address, &x.PhoneNumber, &x.ImageBlobID,
			pq.Array(&x.Approvers), pq.Array(&x.Refusers), pq.Array(&x.RefuseReasons), &x.Status,
			&x.CreatedBy, &x.CreatedAt, &x.UpdatedAt, &x.ClosedAt); err != nil {

			c.errLogger.Println(errLogMsg + err.Error())
			return nil, 0, internalErr
		}

		res = append(res, x)
	}

	var totalRecords int
	c.db.QueryRowContext(ctx, generateCountTotalRecordsQuery(center_request_table, queryCondition)).Scan(&totalRecords)

	return res, caculateTotalPages(totalRecords, req.PageSize), nil
}

// GetRequest implements repository.ICenterRequestRepository.
func (c *centerRequestRepo) GetRequest(id string, ctx context.Context) (*entities.CenterRequest, error) {
	var query string = "SELECT * FROM " + center_request_table + " WHERE id = $1"
	var errLogMsg string = fmt.Sprintf(noti.REPO_ERR_MSG, shared.CENTER_REQUEST_REPOSITORY) + "GetRequest - "

	var res entities.CenterRequest
	if err := c.db.QueryRowContext(ctx, query, id).Scan(
		&res.ID, &res.ProfileID, &res.Region, &res.Address, &res.PhoneNumber, &res.ImageBlobID,
		pq.Array(&res.Approvers), pq.Array(&res.Refusers), pq.Array(&res.RefuseReasons), &res.Status,
		&res.CreatedBy, &res.CreatedAt, &res.UpdatedAt, &res.ClosedAt); err != nil {

		if err == sql.ErrNoRows {
			return nil, nil
		}

		c.errLogger.Println(errLogMsg + err.Error())
		return nil, errors.New(noti.INTERNALL_ERR_MSG)
	}

	return &res, nil
}

// GetWalletRegistrationRequests implements repository.ICenterRequestRepository.
func (c *centerRequestRepo) GetWalletRegistrationRequests(id string, ctx context.Context) ([]entities.CenterRequest, error) {
	var query string = "SELECT * FROM " + center_request_table + " WHERE created_by = $1 ORDER BY created_at DESC"
	var errLogMsg string = fmt.Sprintf(noti.REPO_ERR_MSG, shared.CENTER_REQUEST_REPOSITORY) + "GetWalletRegistrationRequests - "
	var internalErr error = errors.New(noti.INTERNALL_ERR_MSG)

	rows, err := c.db.QueryContext(ctx, query, id)
	if err != nil {
		c.errLogger.Println(errLogMsg + err.Error())
		return nil, internalErr
	}
	defer rows.Close()

	var res []entities.CenterRequest
	for rows.Next() {
		var x entities.CenterRequest
		if err := rows.Scan(
			&x.ID, &x.ProfileID, &x.Region, &x.Address, &x.PhoneNumber, &x.ImageBlobID,
			pq.Array(&x.Approvers), pq.Array(&x.Refusers), pq.Array(&x.RefuseReasons), &x.Status,
			&x.CreatedBy, &x.CreatedAt, &x.UpdatedAt, &x.ClosedAt); err != nil {

			c.errLogger.Println(errLogMsg + err.Error())
			return nil, internalErr
		}

		res = append(res, x)
	}

	return res, nil
}

// IsRegionRequested implements repository.ICenterRequestRepository.
func (c *centerRequestRepo) IsRegionRequested(region string, ctx context.Context) (bool, error) {
	var query string = "SELECT id FROM " + center_request_table + " WHERE LOWER(region) = LOWER($1) AND (status = 'Pending' OR status = 'Approved') LIMIT 1"
	var errLogMsg string = fmt.Sprintf(noti.REPO_ERR_MSG, shared.CENTER_REQUEST_REPOSITORY) + "IsRegionRequested - "
	var internalErr error = errors.New(noti.INTERNALL_ERR_MSG)

	var id string
	if err := c.db.QueryRowContext(ctx, query, region).Scan(&id); err != nil {
		c.errLogger.Println(errLogMsg + err.Error())
		return false, internalErr
	}

	return id != "", nil
}

// UpdateRegistrationRequest implements repository.ICenterRequestRepository.
func (c *centerRequestRepo) UpdateRegistrationRequest(req entities.CenterRequest, ctx context.Context) error {
	var query string = "UPDATE " + center_request_table + " SET " +
		"address = $1, phone_number = $2, image_blob_id = $3, " +
		"approvers = $4, refusers = $5, refuse_reasons = $6, " +
		"status = $7, updated_at = $8 WHERE id = $9"

	var errLogMsg string = fmt.Sprintf(noti.REPO_ERR_MSG, shared.CENTER_REQUEST_REPOSITORY) + "UpdateRegistrationRequest - "
	var internalErr error = errors.New(noti.INTERNALL_ERR_MSG)

	res, err := c.db.ExecContext(ctx, query, req.Address, req.PhoneNumber, req.ImageBlobID,
		pq.Array(req.Approvers), pq.Array(req.Refusers), pq.Array(req.RefuseReasons),
		req.Status, req.UpdatedAt, req.ID)
	if err != nil {
		c.errLogger.Println(errLogMsg + err.Error())
		return internalErr
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		c.errLogger.Println(errLogMsg + err.Error())
		return internalErr
	}

	if rowsAffected == 0 {
		return fmt.Errorf(noti.UNDEFINED_OBJECT_WARN_MSG, center_request_table)
	}

	return nil
}

// GetPendingRequests implements repository.ICenterRequestRepository.
func (c *centerRequestRepo) GetPendingRequests(ctx context.Context) ([]entities.CenterRequest, error) {
	var query string = "SELECT * FROM " + center_request_table + " WHERE closed_at <= NOW() AND status = 'Pending'"
	var errLogMsg string = fmt.Sprintf(noti.REPO_ERR_MSG, shared.CENTER_REQUEST_REPOSITORY) + "GetPendingRequests - "
	var internalErr error = errors.New(noti.INTERNALL_ERR_MSG)

	rows, err := c.db.QueryContext(ctx, query)
	if err != nil {
		c.errLogger.Println(errLogMsg + err.Error())
		return nil, internalErr
	}
	defer rows.Close()

	var res []entities.CenterRequest
	for rows.Next() {
		var x entities.CenterRequest
		if err := rows.Scan(
			&x.ID, &x.ProfileID, &x.Region, &x.Address, &x.PhoneNumber, &x.ImageBlobID,
			pq.Array(&x.Approvers), pq.Array(&x.Refusers), pq.Array(&x.RefuseReasons), &x.Status,
			&x.CreatedBy, &x.CreatedAt, &x.UpdatedAt, &x.ClosedAt); err != nil {

			c.errLogger.Println(errLogMsg + err.Error())
			return nil, internalErr
		}

		res = append(res, x)
	}

	return res, nil
}

// SetApprovedStatuses implements repository.ICenterRequestRepository.
func (c *centerRequestRepo) SetApprovedStatuses(reqs []entities.CenterRequest, ctx context.Context) error {
	if len(reqs) == 0 {
		return nil
	}

	var query string = "UPDATE " + center_request_table + " SET status = 'Approved', updated_at = $1 WHERE "
	for i, req := range reqs {
		query += fmt.Sprintf("id = '%s'", req.ID)
		if i < len(reqs)-1 {
			query += " OR "
		}
	}

	if _, err := c.db.ExecContext(ctx, query, time.Now()); err != nil {
		c.errLogger.Println(fmt.Sprintf(noti.REPO_ERR_MSG, shared.CENTER_REQUEST_REPOSITORY) + "SetApprovedStatuses - " + err.Error())
		return errors.New(noti.INTERNALL_ERR_MSG)
	}

	return nil
}

// SetRefusedStatuses implements repository.ICenterRequestRepository.
func (c *centerRequestRepo) SetRefusedStatuses(reqs []entities.CenterRequest, ctx context.Context) error {
	if len(reqs) == 0 {
		return nil
	}

	var query string = "UPDATE " + center_request_table + " SET status = 'Refused', updated_at = $1 WHERE "
	for i, req := range reqs {
		query += fmt.Sprintf("id = '%s'", req.ID)
		if i < len(reqs)-1 {
			query += " OR "
		}
	}

	if _, err := c.db.ExecContext(ctx, query, time.Now()); err != nil {
		c.errLogger.Println(fmt.Sprintf(noti.REPO_ERR_MSG, shared.CENTER_REQUEST_REPOSITORY) + "SetRefusedStatuses - " + err.Error())
		return errors.New(noti.INTERNALL_ERR_MSG)
	}

	return nil
}
