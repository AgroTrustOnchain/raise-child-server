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

type volunteerRequestRepo struct {
	db        *sql.DB
	errLogger *log.Logger
}

const (
	volunteer_request_table        string = "volunteer_requests"
	volunteer_request_limit_record int    = 10
)

func InitializeVolunteerRequestRepository(db *sql.DB, errLogger *log.Logger) repository.IVolunteerRequestRepository {
	return &volunteerRequestRepo{
		db:        db,
		errLogger: errLogger,
	}
}

// CreateRegistrationRequest implements repository.IvolunteerRequestRepository.
func (v *volunteerRequestRepo) CreateRegistrationRequest(req entities.VolunteerRegistrationRequest, ctx context.Context) error {
	var query string = "INSERT INTO " + volunteer_request_table +
		" (id, identity_code, identity_card_blob_id, avatar_blob_id, region, " +
		"first_name, last_name, gender, date_of_birth, phone_number, email, " +
		"approvers, refusers, refuse_reasons, status, is_available_to_confirm, is_confirm_register, " +
		"created_by, created_at, updated_at, closed_at, is_closed) " +
		"values ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, " +
		"$11, $12, $13, $14, $15, $16, $17, $18, $19, $20, $21)"

	var errLogMsg string = fmt.Sprintf(noti.REPO_ERR_MSG, shared.VOLUNTEER_REQUEST_REPOSITORY) + "CreateRegistrationRequest - "

	if _, err := v.db.ExecContext(ctx, query, req.ID, req.IdentityCode, req.IdentityCardBlobID, req.Region,
		req.AvatarBlobID, req.FirstName, req.LastName, req.Gender,
		req.DateOfBirth, req.PhoneNumber, req.Email, req.Approvers, req.Refusers,
		req.RefuseReasons, req.Status, req.IsAvailableToConfirm, req.IsConfirmRegister,
		req.CreatedBy, req.CreatedAt, req.UpdatedAt, req.ClosedAt); err != nil {

		v.errLogger.Println(errLogMsg + err.Error())
		return errors.New(noti.INTERNALL_ERR_MSG)
	}

	return nil
}

// GetRegistrationRequests implements repository.IvolunteerRequestRepository.
func (v *volunteerRequestRepo) GetRegistrationRequests(req request.GetNormalStaffRegistrationRequests, ctx context.Context) ([]entities.VolunteerRegistrationRequest, int, error) {
	var errLogMsg string = fmt.Sprintf(noti.REPO_ERR_MSG, shared.ADMIN_REQUEST_REPOSITORY) + "GetRegistrationRequests - "
	var internalErr error = errors.New(noti.INTERNALL_ERR_MSG)

	var queryCondition string
	var isHavePreviosCondition bool = false
	if req.Keyword != "" {
		queryCondition += fmt.Sprintf("(LOWER(identity_code) LIKE LOWER('%%%%%s%%%%') OR LOWER(first_name) LIKE LOWER('%%%%%s%%%%') OR LOWER(last_name) LIKE LOWER('%%%%%s%%%%') OR date_of_birth LIKE '%%%%%s%%%%' OR phone_number LIKE '%%%%%s%%%%' OR LOWER(email) LIKE LOWER('%%%%%s%%%%'))", req.Keyword, req.Keyword, req.Keyword, req.Keyword, req.Keyword, req.Keyword)
		isHavePreviosCondition = true
	}

	if req.Region != "" {
		if isHavePreviosCondition {
			queryCondition += " AND "
		}

		queryCondition += fmt.Sprintf("LOWER(region) = LOWER('%%%%%s%%%%')", req.Region)
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

	rows, err := v.db.QueryContext(ctx, query)
	if err != nil {
		v.errLogger.Println(errLogMsg + err.Error())
		return nil, 0, internalErr
	}

	var res []entities.VolunteerRegistrationRequest
	for rows.Next() {
		var x entities.VolunteerRegistrationRequest
		if err := rows.Scan(
			&x.ID, &x.IdentityCode, &x.IdentityCardBlobID, &x.Region,
			&x.AvatarBlobID, &x.FirstName, &x.LastName, &x.Gender,
			&x.DateOfBirth, &x.PhoneNumber, &x.Email, &x.Approvers, &x.Refusers,
			&x.RefuseReasons, &x.Status, &x.IsAvailableToConfirm, &x.IsConfirmRegister,
			&x.CreatedBy, &x.CreatedAt, &x.UpdatedAt, &x.ClosedAt); err != nil {

			v.errLogger.Println(errLogMsg + err.Error())
			return nil, 0, internalErr
		}

		res = append(res, x)
	}

	var totalRecords int
	v.db.QueryRowContext(ctx, generateCountTotalRecordsQuery(volunteer_request_table, queryCondition)).Scan(&totalRecords)

	return res, caculateTotalPages(totalRecords, volunteer_request_limit_record), nil
}

// GetRequest implements repository.IvolunteerRequestRepository.
func (v *volunteerRequestRepo) GetRequest(id string, ctx context.Context) (*entities.VolunteerRegistrationRequest, error) {
	var query string = "SELECT * FROM " + volunteer_request_table + " WHERE id = $1"
	var errLogMsg string = fmt.Sprintf(noti.REPO_ERR_MSG, shared.VOLUNTEER_REQUEST_REPOSITORY) + "GetRequest - "

	var res entities.VolunteerRegistrationRequest
	if err := v.db.QueryRowContext(ctx, query, id).Scan(
		&res.ID, &res.IdentityCode, &res.IdentityCardBlobID, &res.Region,
		&res.AvatarBlobID, &res.FirstName, &res.LastName, &res.Gender,
		&res.DateOfBirth, &res.PhoneNumber, &res.Email, &res.Approvers, &res.Refusers,
		&res.RefuseReasons, &res.Status, &res.IsAvailableToConfirm, &res.IsConfirmRegister,
		&res.CreatedBy, &res.CreatedAt, &res.UpdatedAt, &res.ClosedAt); err != nil {

		if err == sql.ErrNoRows {
			return nil, nil
		}

		v.errLogger.Println(errLogMsg + err.Error())
		return nil, errors.New(noti.INTERNALL_ERR_MSG)
	}

	return &res, nil
}

// GetWalletRegistrationRequests implements repository.IvolunteerRequestRepository.
func (v *volunteerRequestRepo) GetWalletRegistrationRequests(id string, ctx context.Context) ([]entities.VolunteerRegistrationRequest, error) {
	var query string = "SELECT * FROM " + volunteer_request_table + " WHERE created_by = $1 ORDER BY created_at DESC"
	var errLogMsg string = fmt.Sprintf(noti.REPO_ERR_MSG, shared.VOLUNTEER_REQUEST_REPOSITORY) + "GetWalletRegistrationRequests - "
	var internalErr error = errors.New(noti.INTERNALL_ERR_MSG)

	rows, err := v.db.QueryContext(ctx, query, id)
	if err != nil {
		v.errLogger.Println(errLogMsg + err.Error())
		return nil, internalErr
	}

	var res []entities.VolunteerRegistrationRequest
	for rows.Next() {
		var x entities.VolunteerRegistrationRequest
		if err := rows.Scan(
			&x.ID, &x.IdentityCode, &x.IdentityCardBlobID, &x.Region,
			&x.AvatarBlobID, &x.FirstName, &x.LastName, &x.Gender,
			&x.DateOfBirth, &x.PhoneNumber, &x.Email, &x.Approvers, &x.Refusers,
			&x.RefuseReasons, &x.Status, &x.IsAvailableToConfirm, &x.IsConfirmRegister,
			&x.CreatedBy, &x.CreatedAt, &x.UpdatedAt, &x.ClosedAt); err != nil {

			v.errLogger.Println(errLogMsg + err.Error())
			return nil, internalErr
		}

		res = append(res, x)
	}

	return res, nil
}

// UpdateRegistrationRequest implements repository.IvolunteerRequestRepository.
func (v *volunteerRequestRepo) UpdateRegistrationRequest(req entities.VolunteerRegistrationRequest, ctx context.Context) error {
	var query string = "UPDATE " + volunteer_request_table + " SET " +
		"approvers = $1, refusers = $2, refuse_reasons = $3, status = $4, is_confirm_register = $5, " +
		"is_available_to_confirm = $6, updated_at = $7, avatar_blob_id = $8, region = $9, WHERE id = $10"

	var errLogMsg string = fmt.Sprintf(noti.REPO_ERR_MSG, shared.VOLUNTEER_REQUEST_REPOSITORY) + "UpdateRegistrationRequest - "
	var internalErr error = errors.New(noti.INTERNALL_ERR_MSG)

	res, err := v.db.ExecContext(ctx, query, req.Approvers, req.Refusers, req.RefuseReasons, req.Status, req.IsConfirmRegister,
		req.IsAvailableToConfirm, req.UpdatedAt, req.AvatarBlobID, req.Region, req.ID)
	if err != nil {
		v.errLogger.Println(errLogMsg + err.Error())
		return internalErr
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		v.errLogger.Println(errLogMsg + err.Error())
		return internalErr
	}

	if rowsAffected == 0 {
		return fmt.Errorf(noti.UNDEFINED_OBJECT_WARN_MSG, volunteer_request_table)
	}

	return nil
}
