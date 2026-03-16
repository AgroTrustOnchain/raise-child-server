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

type registratioRequestRepo struct {
	db        *sql.DB
	errLogger *log.Logger
}

const (
	registraion_request_table        string = "registration_requests"
	registraion_request_limit_record int    = 10
)

func InitializeRegistrationRequestRepo(db *sql.DB, errLogger *log.Logger) repository.IRegistrationRequestRepository {
	return &registratioRequestRepo{
		db:        db,
		errLogger: errLogger,
	}
}

// GetRoleRegistrationRequests implements repository.IRegistrationRequestRepository.
func (r *registratioRequestRepo) GetRoleRegistrationRequests(role string, ctx context.Context) ([]entities.RegistrationRequest, error) {
	var query string = "SELECT * FROM " + registraion_request_table + " WHERE register_role = $1"
	var errLogMsg string = fmt.Sprintf(noti.REPO_ERR_MSG, shared.REGISTRAION_REQUEST_REPOSITORY) + "GetRoleRegistrationRequests - "
	var internalErr error = errors.New(noti.INTERNALL_ERR_MSG)

	rows, err := r.db.QueryContext(ctx, query, role)
	if err != nil {
		r.errLogger.Println(errLogMsg + err.Error())
		return nil, internalErr
	}

	var res []entities.RegistrationRequest
	for rows.Next() {
		var x entities.RegistrationRequest
		// Original
		// if err := rows.Scan(
		// 	&x.ID, &x.RegisterRole, &x.IdentityCode, &x.IdentityCardBlobID,
		// 	&x.AvatarBlobID, &x.Region, &x.FirstName, &x.LastName, &x.Gender,
		// 	&x.DateOfBirth, &x.PhoneNumber, &x.Email, &x.Approvers, &x.Refusers,
		// 	&x.RefuseReasons, &x.Status, &x.IsAvailableToConfirm, &x.IsConfirmRegister,
		// 	&x.CreatedBy, &x.CreatedAt, &x.UpdatedAt, &x.ClosedAt); err != nil {

		// 	r.errLogger.Println(errLogMsg + err.Error())
		// 	return nil, internalErr
		// }

		if err := rows.Scan(
			&x.ID, &x.ProfileID, &x.RegisterRole, &x.IdentityCode, &x.IdentityCardBlobID,
			&x.AvatarBlobID, &x.Region, &x.FirstName, &x.LastName, &x.Gender,
			&x.DateOfBirth, &x.PhoneNumber, &x.Email, pq.Array(&x.Approvers), pq.Array(&x.Refusers),
			pq.Array(&x.RefuseReasons), &x.Status, &x.IsAvailableToConfirm, &x.IsConfirmRegister,
			&x.CreatedBy, &x.CreatedAt, &x.UpdatedAt, &x.ClosedAt); err != nil {

			r.errLogger.Println(errLogMsg + err.Error())
			return nil, internalErr
		}

		res = append(res, x)
	}

	return res, nil
}

// CreateRegistrationRequest implements repository.IRegistrationRequestRepository.
func (r *registratioRequestRepo) CreateRegistrationRequest(req entities.RegistrationRequest, ctx context.Context) error {
	var query string = "INSERT INTO " + registraion_request_table +
		" (id, profile_id, register_role, identity_code, identity_card_blob_id, avatar_blob_id, " +
		"region, first_name, last_name, gender, date_of_birth, phone_number, email, " +
		"approvers, refusers, refuse_reasons, status, is_available_to_confirm, is_confirm_register, " +
		"created_by, created_at, updated_at, closed_at) " +
		"values ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, " +
		"$11, $12, $13, $14, $15, $16, $17, $18, $19, $20, $21, $ss, $23)"

	var errLogMsg string = fmt.Sprintf(noti.REPO_ERR_MSG, shared.REGISTRAION_REQUEST_REPOSITORY) + "CreateRegistrationRequest - "

	if _, err := r.db.ExecContext(ctx, query, req.ID, req.ProfileID, req.RegisterRole, req.IdentityCode, req.IdentityCardBlobID,
		req.AvatarBlobID, req.Region, req.FirstName, req.LastName, req.Gender,
		req.DateOfBirth, req.PhoneNumber, req.Email, pq.Array(req.Approvers), pq.Array(req.Refusers),
		pq.Array(req.RefuseReasons), req.Status, req.IsAvailableToConfirm, req.IsConfirmRegister,
		req.CreatedBy, req.CreatedAt, req.UpdatedAt, req.ClosedAt); err != nil {

		r.errLogger.Println(errLogMsg + err.Error())
		return errors.New(noti.INTERNALL_ERR_MSG)
	}

	return nil
}

// UpdateRegistrationRequest implements repository.IRegistrationRequestRepository.
func (r *registratioRequestRepo) UpdateRegistrationRequest(req entities.RegistrationRequest, ctx context.Context) error {
	var query string = "UPDATE " + registraion_request_table + " SET " +
		"region = $1, first_name = $2, last_name = $3, gender = $4, " +
		"date_of_birth = $5, phone_number = $6, email = $7, " +
		"approvers = $8, refusers = $9, refuse_reasons = $10, " +
		"status = $11, is_confirm_register = $12, is_available_to_confirm = $13, updated_at = $14 WHERE id = $15"

	var errLogMsg string = fmt.Sprintf(noti.REPO_ERR_MSG, shared.REGISTRAION_REQUEST_REPOSITORY) + "UpdateRegistrationRequest - "
	var internalErr error = errors.New(noti.INTERNALL_ERR_MSG)

	res, err := r.db.ExecContext(ctx, query, req.Region, req.FirstName, req.LastName, req.Gender,
		req.DateOfBirth, req.PhoneNumber, req.Email,
		pq.Array(req.Approvers), pq.Array(req.Refusers), pq.Array(req.RefuseReasons),
		req.Status, req.IsConfirmRegister, req.IsAvailableToConfirm, req.UpdatedAt, req.ID)
	if err != nil {
		r.errLogger.Println(errLogMsg + err.Error())
		return internalErr
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		r.errLogger.Println(errLogMsg + err.Error())
		return internalErr
	}

	if rowsAffected == 0 {
		return errors.New(fmt.Sprintf(noti.UNDEFINED_OBJECT_WARN_MSG, registraion_request_table))
	}

	return nil
}

// GetRegistrationRequest implements repository.IRegistrationRequestRepository.
func (r *registratioRequestRepo) GetRegistrationRequest(id string, ctx context.Context) (*entities.RegistrationRequest, error) {
	var query string = "SELECT * FROM " + registraion_request_table + " WHERE id = $1"
	var errLogMsg string = fmt.Sprintf(noti.REPO_ERR_MSG, shared.REGISTRAION_REQUEST_REPOSITORY) + "GetRegistrationRequest - "

	var res entities.RegistrationRequest
	if err := r.db.QueryRowContext(ctx, query, id).Scan(
		&res.ID, &res.ProfileID, &res.RegisterRole, &res.IdentityCode, &res.IdentityCardBlobID,
		&res.AvatarBlobID, &res.Region, &res.FirstName, &res.LastName, &res.Gender,
		&res.DateOfBirth, &res.PhoneNumber, &res.Email, pq.Array(&res.Approvers), pq.Array(&res.Refusers),
		pq.Array(&res.RefuseReasons), &res.Status, &res.IsAvailableToConfirm, &res.IsConfirmRegister,
		&res.CreatedBy, &res.CreatedAt, &res.UpdatedAt, &res.ClosedAt); err != nil {

		if err == sql.ErrNoRows {
			return nil, nil
		}

		r.errLogger.Println(errLogMsg + err.Error())
		return nil, errors.New(noti.INTERNALL_ERR_MSG)
	}

	return &res, nil
}

// GetRegistrationRequests implements repository.IRegistrationRequestRepository.
func (r *registratioRequestRepo) GetRegistrationRequests(req request.GetRegistrationRequests, ctx context.Context) ([]entities.RegistrationRequest, int, error) {
	var errLogMsg string = fmt.Sprintf(noti.REPO_ERR_MSG, shared.REGISTRAION_REQUEST_REPOSITORY) + "GetRegistrationRequests - "
	var internalErr error = errors.New(noti.INTERNALL_ERR_MSG)

	var queryCondition string
	var isHavePreviosCondition bool = false
	if req.Keyword != "" {
		queryCondition += fmt.Sprintf("(LOWER(identity_code) LIKE LOWER('%%%%%s%%%%') OR LOWER(first_name) LIKE LOWER('%%%%%s%%%%') OR LOWER(last_name) LIKE LOWER('%%%%%s%%%%') OR date_of_birth LIKE '%%%%%s%%%%' OR phone_number LIKE '%%%%%s%%%%' OR LOWER(email) LIKE LOWER('%%%%%s%%%%'))", req.Keyword, req.Keyword, req.Keyword, req.Keyword, req.Keyword, req.Keyword)
		isHavePreviosCondition = true
	}

	if req.RegisterRole != "" {
		if isHavePreviosCondition {
			queryCondition += " AND "
		}

		queryCondition += fmt.Sprintf("LOWER(register_role) = LOWER('%%%%%s%%%%')", req.RegisterRole)
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

	if req.IsAvailableToConfirm != nil {
		if isHavePreviosCondition {
			queryCondition += " AND "
		}

		queryCondition += fmt.Sprintf("is_available_to_confirm = %v", *req.IsAvailableToConfirm)
		isHavePreviosCondition = true
	}

	if req.IsClosed != nil {
		if isHavePreviosCondition {
			queryCondition += " AND "
		}

		var operation string = ">"
		if *req.IsClosed {
			operation = "<"
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

	queryCondition += "ORDER BY created_at " + order

	var query string = generateRetrieveQuery(generateRetrieveQueryRequest{
		table:       registraion_request_table,
		limitAmount: req.PageSize,
		condition:   queryCondition,
		page:        req.Page,
		isGetCount:  false,
	})

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		r.errLogger.Println(errLogMsg + err.Error())
		return nil, 0, internalErr
	}

	var res []entities.RegistrationRequest
	for rows.Next() {
		var x entities.RegistrationRequest
		// Original
		// if err := rows.Scan(
		// 	&x.ID, &x.RegisterRole, &x.IdentityCode, &x.IdentityCardBlobID,
		// 	&x.AvatarBlobID, &x.Region, &x.FirstName, &x.LastName, &x.Gender,
		// 	&x.DateOfBirth, &x.PhoneNumber, &x.Email, &x.Approvers, &x.Refusers,
		// 	&x.RefuseReasons, &x.Status, &x.IsAvailableToConfirm, &x.IsConfirmRegister,
		// 	&x.CreatedBy, &x.CreatedAt, &x.UpdatedAt, &x.ClosedAt); err != nil {

		// 	r.errLogger.Println(errLogMsg + err.Error())
		// 	return nil, 0, internalErr
		// }

		if err := rows.Scan(
			&x.ID, &x.ProfileID, &x.RegisterRole, &x.IdentityCode, &x.IdentityCardBlobID,
			&x.AvatarBlobID, &x.Region, &x.FirstName, &x.LastName, &x.Gender,
			&x.DateOfBirth, &x.PhoneNumber, &x.Email, pq.Array(&x.Approvers), pq.Array(&x.Refusers),
			pq.Array(&x.RefuseReasons), &x.Status, &x.IsAvailableToConfirm, &x.IsConfirmRegister,
			&x.CreatedBy, &x.CreatedAt, &x.UpdatedAt, &x.ClosedAt); err != nil {

			r.errLogger.Println(errLogMsg + err.Error())
			return nil, 0, internalErr
		}

		res = append(res, x)
	}

	var totalRecords int
	r.db.QueryRowContext(ctx, generateCountTotalRecordsQuery(registraion_request_table, queryCondition)).Scan(&totalRecords)

	return res, caculateTotalPages(totalRecords, req.PageSize), nil
}

// GetWalletRegistrationRequests implements repository.IRegistrationRequestRepository.
func (r *registratioRequestRepo) GetWalletRegistrationRequests(id string, ctx context.Context) ([]entities.RegistrationRequest, error) {
	var query string = "SELECT * FROM " + registraion_request_table + " WHERE created_by = $1 ORDER BY created_at DESC"
	var errLogMsg string = fmt.Sprintf(noti.REPO_ERR_MSG, shared.REGISTRAION_REQUEST_REPOSITORY) + "GetWalletRegistrationRequests - "
	var internalErr error = errors.New(noti.INTERNALL_ERR_MSG)

	rows, err := r.db.QueryContext(ctx, query, id)
	if err != nil {
		r.errLogger.Println(errLogMsg + err.Error())
		return nil, internalErr
	}

	var res []entities.RegistrationRequest
	for rows.Next() {
		var x entities.RegistrationRequest
		if err := rows.Scan(
			&x.ID, &x.ProfileID, &x.RegisterRole, &x.IdentityCode, &x.IdentityCardBlobID,
			&x.AvatarBlobID, &x.Region, &x.FirstName, &x.LastName, &x.Gender,
			&x.DateOfBirth, &x.PhoneNumber, &x.Email, pq.Array(&x.Approvers), pq.Array(&x.Refusers),
			pq.Array(&x.RefuseReasons), &x.Status, &x.IsAvailableToConfirm, &x.IsConfirmRegister,
			&x.CreatedBy, &x.CreatedAt, &x.UpdatedAt, &x.ClosedAt); err != nil {

			r.errLogger.Println(errLogMsg + err.Error())
			return nil, internalErr
		}

		res = append(res, x)
	}

	return res, nil
}

// GetPendingRequests implements repository.IRegistrationRequestRepository.
func (r *registratioRequestRepo) GetPendingRequests(ctx context.Context) ([]entities.BackgroundRecord, []entities.BackgroundRecord, error) {
	var query string = "SELECT id, approvers, refusers, register_role, created_by, status FROM " + registraion_request_table + " WHERE is_available_to_confirm = false AND closed_at <= NOW() AND (status = 'Pending' OR status = 'Approved')"
	var errLogMsg string = fmt.Sprintf(noti.REPO_ERR_MSG, shared.REGISTRAION_REQUEST_REPOSITORY) + "GetPendingRequests - "
	var internalErr error = errors.New(noti.INTERNALL_ERR_MSG)

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		r.errLogger.Println(errLogMsg + err.Error())
		return nil, nil, internalErr
	}

	var pendingRes, approvedRes []entities.BackgroundRecord
	for rows.Next() {
		var x entities.BackgroundRecord
		var status string
		if err := rows.Scan(
			&x.ID, pq.Array(&x.Approvers), pq.Array(&x.Refusers), &x.Role, &x.Sender, &status); err != nil {

			r.errLogger.Println(errLogMsg + err.Error())
			return nil, nil, internalErr
		}

		if status == "Pending" {
			pendingRes = append(pendingRes, x)
		} else {
			approvedRes = append(approvedRes, x)
		}
	}

	return pendingRes, approvedRes, nil
}

// SetApprovedStatuses implements repository.IRegistrationRequestRepository.
func (r *registratioRequestRepo) SetApprovedStatuses(reqs []entities.BackgroundRecord, ctx context.Context) error {
	var query string = "UPDATE " + registraion_request_table + " SET status = 'Approved', is_available_to_confirm = true, updated_at = $1 WHERE "
	for i, req := range reqs {
		query += "id = " + req.ID
		if i < len(reqs)-1 {
			query += " OR "
		}
	}

	if _, err := r.db.ExecContext(ctx, query, time.Now()); err != nil {
		r.errLogger.Println(fmt.Sprintf(noti.REPO_ERR_MSG, shared.REGISTRAION_REQUEST_REPOSITORY) + "SetApprovedStatuses - " + err.Error())
		return errors.New(noti.INTERNALL_ERR_MSG)
	}

	return nil
}

// SetRefusedStatuses implements repository.IRegistrationRequestRepository.
func (r *registratioRequestRepo) SetRefusedStatuses(reqs []entities.BackgroundRecord, ctx context.Context) error {
	var query string = "UPDATE " + registraion_request_table + " SET status = 'Refused', updated_at = $1 WHERE "
	for i, req := range reqs {
		query += "id = " + req.ID
		if i < len(reqs)-1 {
			query += " OR "
		}
	}

	if _, err := r.db.ExecContext(ctx, query, time.Now()); err != nil {
		r.errLogger.Println(fmt.Sprintf(noti.REPO_ERR_MSG, shared.REGISTRAION_REQUEST_REPOSITORY) + "SetRefusedStatuses - " + err.Error())
		return errors.New(noti.INTERNALL_ERR_MSG)
	}

	return nil
}
