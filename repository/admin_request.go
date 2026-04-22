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
)

type adminRequestRepo struct {
	db        *sql.DB
	errLogger *log.Logger
}

const (
	admin_request_table        string = "admin_requests"
	admin_request_limit_record int    = 10
)

func InitializeAdminRequestRepository(db *sql.DB, errLogger *log.Logger) repository.IAdminRequestRepository {
	return &adminRequestRepo{
		db:        db,
		errLogger: errLogger,
	}
}

// CreateRegistrationRequest implements repository.IAdminRequestRepository.
func (a *adminRequestRepo) CreateRegistrationRequest(req entities.AdminRegistrationRequest, ctx context.Context) error {
	var query string = "INSERT INTO " + admin_request_table +
		" (id, identity_code, identity_card_blob_id, avatar_blob_id, " +
		"first_name, last_name, gender, date_of_birth, phone_number, email, " +
		"approvers, refusers, refuse_reasons, status, is_available_to_confirm, is_confirm_register, " +
		"created_by, created_at, updated_at, closed_at, is_closed) " +
		"values ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, " +
		"$11, $12, $13, $14, $15, $16, $17, $18, $19, $20)"

	var errLogMsg string = fmt.Sprintf(noti.REPO_ERR_MSG, shared.ADMIN_REQUEST_REPOSITORY) + "CreateRegistrationRequest - "

	if _, err := a.db.ExecContext(ctx, query, req.ID, req.IdentityCode, req.IdentityCardBlobID,
		req.AvatarBlobID, req.FirstName, req.LastName, req.Gender,
		req.DateOfBirth, req.PhoneNumber, req.Email, req.Approvers, req.Refusers,
		req.RefuseReasons, req.Status, req.IsAvailableToConfirm, req.IsConfirmRegister,
		req.CreatedBy, req.CreatedAt, req.UpdatedAt, req.ClosedAt); err != nil {

		a.errLogger.Println(errLogMsg + err.Error())
		return errors.New(noti.INTERNALL_ERR_MSG)
	}

	return nil
}

// GetRegistrationRequests implements repository.IAdminRequestRepository.
func (a *adminRequestRepo) GetRegistrationRequests(req request.GetAdminRegistrationRequets, ctx context.Context) ([]entities.AdminRegistrationRequest, int, error) {
	var errLogMsg string = fmt.Sprintf(noti.REPO_ERR_MSG, shared.ADMIN_REQUEST_REPOSITORY) + "GetRegistrationRequests - "
	var internalErr error = errors.New(noti.INTERNALL_ERR_MSG)

	var queryCondition string
	var isHavePreviosCondition bool = false
	if req.Keyword != "" {
		queryCondition += fmt.Sprintf("(LOWER(identity_code) LIKE LOWER('%%%%%s%%%%') OR LOWER(first_name) LIKE LOWER('%%%%%s%%%%') OR LOWER(last_name) LIKE LOWER('%%%%%s%%%%') OR date_of_birth LIKE '%%%%%s%%%%' OR phone_number LIKE '%%%%%s%%%%' OR LOWER(email) LIKE LOWER('%%%%%s%%%%'))", req.Keyword, req.Keyword, req.Keyword, req.Keyword, req.Keyword, req.Keyword)
		isHavePreviosCondition = true
	}

	if req.Gender != "" {
		if isHavePreviosCondition {
			queryCondition += " AND "
		}

		queryCondition += fmt.Sprintf("LOWER(gender) = LOWER('%%%%%s%%%%')", req.Gender)
		isHavePreviosCondition = true
	}

	if req.Status != "" {
		if isHavePreviosCondition {
			queryCondition += " AND "
		}

		queryCondition += fmt.Sprintf("LOWER(status) = LOWER('%%%%%s%%%%')", req.Status)
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

	if req.IsConfirm != nil {
		if isHavePreviosCondition {
			queryCondition += " AND "
		}

		queryCondition += fmt.Sprintf("is_confirm_register = %v", *req.IsConfirm)
	}

	if isHavePreviosCondition {
		queryCondition += " "
	}

	var order string = "DESC"
	if req.SortOrder != "" {
		order = req.SortOrder
	}

	queryCondition += "ORDER BY created_at " + order

	var query string = generateRetrieveQuery(generateRetrieveQueryRequest{
		table:       registration_request_table,
		limitAmount: registration_request_limit_record,
		condition:   queryCondition,
		page:        req.Page,
		isGetCount:  false,
	})

	rows, err := a.db.QueryContext(ctx, query)
	if err != nil {
		a.errLogger.Println(errLogMsg + err.Error())
		return nil, 0, internalErr
	}

	var res []entities.AdminRegistrationRequest
	for rows.Next() {
		var x entities.AdminRegistrationRequest
		if err := rows.Scan(
			&x.ID, &x.IdentityCode, &x.IdentityCardBlobID,
			&x.AvatarBlobID, &x.FirstName, &x.LastName, &x.Gender,
			&x.DateOfBirth, &x.PhoneNumber, &x.Email, &x.Approvers, &x.Refusers,
			&x.RefuseReasons, &x.Status, &x.IsAvailableToConfirm, &x.IsConfirmRegister,
			&x.CreatedBy, &x.CreatedAt, &x.UpdatedAt, &x.ClosedAt); err != nil {

			a.errLogger.Println(errLogMsg + err.Error())
			return nil, 0, internalErr
		}

		res = append(res, x)
	}

	var totalRecords int
	a.db.QueryRowContext(ctx, generateCountTotalRecordsQuery(admin_request_table, queryCondition)).Scan(&totalRecords)

	return res, caculateTotalPages(totalRecords, admin_request_limit_record), nil
}

// GetRequest implements repository.IAdminRequestRepository.
func (a *adminRequestRepo) GetRequest(id string, ctx context.Context) (*entities.AdminRegistrationRequest, error) {
	var query string = "SELECT * FROM " + admin_request_table + " WHERE id = $1"
	var errLogMsg string = fmt.Sprintf(noti.REPO_ERR_MSG, shared.ADMIN_REQUEST_REPOSITORY) + "GetRequest - "

	var res entities.AdminRegistrationRequest
	if err := a.db.QueryRowContext(ctx, query, id).Scan(
		&res.ID, &res.IdentityCode, &res.IdentityCardBlobID,
		&res.AvatarBlobID, &res.FirstName, &res.LastName, &res.Gender,
		&res.DateOfBirth, &res.PhoneNumber, &res.Email, &res.Approvers, &res.Refusers,
		&res.RefuseReasons, &res.Status, &res.IsAvailableToConfirm, &res.IsConfirmRegister,
		&res.CreatedBy, &res.CreatedAt, &res.UpdatedAt, &res.ClosedAt); err != nil {

		if err == sql.ErrNoRows {
			return nil, nil
		}

		a.errLogger.Println(errLogMsg + err.Error())
		return nil, errors.New(noti.INTERNALL_ERR_MSG)
	}

	return &res, nil
}

// GetWalletRegistrationRequests implements repository.IAdminRequestRepository.
func (a *adminRequestRepo) GetWalletRegistrationRequests(id string, ctx context.Context) ([]entities.AdminRegistrationRequest, error) {
	var query string = "SELECT * FROM " + admin_request_table + " WHERE created_by = $1 ORDER BY created_at DESC"
	var errLogMsg string = fmt.Sprintf(noti.REPO_ERR_MSG, shared.ADMIN_REQUEST_REPOSITORY) + "GetWalletRegistrationRequests - "
	var internalErr error = errors.New(noti.INTERNALL_ERR_MSG)

	rows, err := a.db.QueryContext(ctx, query, id)
	if err != nil {
		a.errLogger.Println(errLogMsg + err.Error())
		return nil, internalErr
	}

	var res []entities.AdminRegistrationRequest
	for rows.Next() {
		var x entities.AdminRegistrationRequest
		if err := rows.Scan(
			&x.ID, &x.IdentityCode, &x.IdentityCardBlobID,
			&x.AvatarBlobID, &x.FirstName, &x.LastName, &x.Gender,
			&x.DateOfBirth, &x.PhoneNumber, &x.Email, &x.Approvers, &x.Refusers,
			&x.RefuseReasons, &x.Status, &x.IsAvailableToConfirm, &x.IsConfirmRegister,
			&x.CreatedBy, &x.CreatedAt, &x.UpdatedAt, &x.ClosedAt); err != nil {

			a.errLogger.Println(errLogMsg + err.Error())
			return nil, internalErr
		}

		res = append(res, x)
	}

	return res, nil
}

// UpdateRegistrationRequest implements repository.IAdminRequestRepository.
func (a *adminRequestRepo) UpdateRegistrationRequest(req entities.AdminRegistrationRequest, ctx context.Context) error {
	var query string = "UPDATE " + admin_request_table + " SET " +
		"approvers = $1, refusers = $2, refuse_reasons = $3, status = $4, is_confirm_register = $5, " +
		"is_available_to_confirm = $6, updated_at = $7, avatar_blob_id = $8 WHERE id = $9"

	var errLogMsg string = fmt.Sprintf(noti.REPO_ERR_MSG, shared.ADMIN_REQUEST_REPOSITORY) + "UpdateRegistrationRequest - "
	var internalErr error = errors.New(noti.INTERNALL_ERR_MSG)

	res, err := a.db.ExecContext(ctx, query, req.Approvers, req.Refusers, req.RefuseReasons, req.Status, req.IsConfirmRegister,
		req.IsAvailableToConfirm, req.UpdatedAt, req.AvatarBlobID, req.ID)
	if err != nil {
		a.errLogger.Println(errLogMsg + err.Error())
		return internalErr
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		a.errLogger.Println(errLogMsg + err.Error())
		return internalErr
	}

	if rowsAffected == 0 {
		return errors.New(fmt.Sprintf(noti.UNDEFINED_OBJECT_WARN_MSG, admin_request_table))
	}

	return nil
}
