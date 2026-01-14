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

type uploadChildRepo struct {
	db        *sql.DB
	errLogger *log.Logger
}

const (
	upload_child_request_table        string = "upload_child_requests"
	upload_child_request_limit_record int    = 10
)

func InitializeUploadChildRequestRepo(db *sql.DB, errLogger *log.Logger) repository.IUploadChildRequestRepository {
	return &uploadChildRepo{
		db:        db,
		errLogger: errLogger,
	}
}

// CreateUploadChildRequest implements repository.IUploadChildRequestRepository.
func (u *uploadChildRepo) CreateUploadChildRequest(req entities.UploadChildRequest, ctx context.Context) error {
	var query string = "INSERT INTO " + upload_child_request_table +
		" (id, identity_code, avatar_blob_id, " +
		"region, first_name, last_name, gender, date_of_birth, " +
		"approvers, refusers, refuse_reasons, status, is_confirm_upload, " +
		"created_by, created_at, updated_at, closed_at, is_closed) " +
		"values ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, " +
		"$11, $12, $13, $14, $15, $16, $17)"

	var errLogMsg string = fmt.Sprintf(noti.REPO_ERR_MSG, shared.UPLOAD_CHILD_REQUEST_REPOSITORY) + "CreateUploadChildRequest - "

	if _, err := u.db.Exec(query, req.ID, req.IdentityCode, req.AvatarBlobId,
		req.Region, req.FirstName, req.LastName, req.Gender, req.DateOfBirth,
		req.Aprrovers, req.Refusers, req.RefuseReasons, req.Status, req.IsConfirmUpload,
		req.CreatedBy, req.CreatedAt, req.UpdatedAt, req.ClosedAt); err != nil {

		u.errLogger.Println(errLogMsg + err.Error())
		return errors.New(noti.INTERNALL_ERR_MSG)
	}

	return nil
}

// GetUploadChildRequest implements repository.IUploadChildRequestRepository.
func (u *uploadChildRepo) GetUploadChildRequest(id string, ctx context.Context) (*entities.UploadChildRequest, error) {
	var query string = "SELECT * FROM " + registraion_request_table + " WHERE id = $1"
	var errLogMsg string = fmt.Sprintf(noti.REPO_ERR_MSG, shared.UPLOAD_CHILD_REQUEST_REPOSITORY) + "GetUploadChildRequest - "

	var res entities.UploadChildRequest
	if err := u.db.QueryRow(query, id).Scan(
		res.ID, res.IdentityCode, res.AvatarBlobId,
		res.Region, res.FirstName, res.LastName, res.Gender, res.DateOfBirth,
		res.Aprrovers, res.Refusers, res.RefuseReasons, res.Status, res.IsConfirmUpload,
		res.CreatedBy, res.CreatedAt, res.UpdatedAt, res.ClosedAt); err != nil {

		if err == sql.ErrNoRows {
			return nil, nil
		}

		u.errLogger.Println(errLogMsg + err.Error())
		return nil, errors.New(noti.INTERNALL_ERR_MSG)
	}

	return &res, nil
}

// GetUploadChildRequests implements repository.IUploadChildRequestRepository.
func (u *uploadChildRepo) GetUploadChildRequests(req request.GetUploadChildRequests, ctx context.Context) ([]entities.UploadChildRequest, int, error) {
	var errLogMsg string = fmt.Sprintf(noti.REPO_ERR_MSG, shared.UPLOAD_CHILD_REQUEST_REPOSITORY) + "GetWalletRegistrationRequests - "
	var internalErr error = errors.New(noti.INTERNALL_ERR_MSG)

	var queryCondition string
	var isHavePreviosCondition bool = false
	if req.Keyword != "" {
		queryCondition += fmt.Sprintf("(identity_code LIKE '%s' OR LOWER(first_name) LIKE LOWER('%%%%%s%%%%') OR LOWER(last_name) LIKE LOWER('%%%%%s%%%%') OR date_of_birth LIKE '%%%%%s%%%%')", req.Keyword, req.Keyword, req.Keyword, req.Keyword)
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
		table:       upload_child_request_table,
		limitAmount: upload_child_request_limit_record,
		condition:   queryCondition,
		page:        req.Page,
		isGetCount:  false,
	})

	rows, err := u.db.Query(query)
	if err != nil {
		u.errLogger.Println(errLogMsg + err.Error())
		return nil, 0, internalErr
	}

	var res []entities.UploadChildRequest
	for rows.Next() {
		var x entities.UploadChildRequest
		if err := rows.Scan(
			x.ID, x.IdentityCode, x.AvatarBlobId,
			x.Region, x.FirstName, x.LastName, x.Gender, x.DateOfBirth,
			x.Aprrovers, x.Refusers, x.RefuseReasons, x.Status, x.IsConfirmUpload,
			x.CreatedBy, x.CreatedAt, x.UpdatedAt, x.ClosedAt); err != nil {

			u.errLogger.Println(errLogMsg + err.Error())
			return nil, 0, internalErr
		}

		res = append(res, x)
	}

	var totalRecords int
	u.db.QueryRow(generateCountTotalRecordsQuery(upload_child_request_table, queryCondition)).Scan(&totalRecords)

	return res, caculateTotalPages(totalRecords, upload_child_request_limit_record), nil
}

// GetWalletUploadChildRequests implements repository.IUploadChildRequestRepository.
func (u *uploadChildRepo) GetWalletUploadChildRequests(id string, page int, ctx context.Context) ([]entities.UploadChildRequest, int, error) {
	var errLogMsg string = fmt.Sprintf(noti.REPO_ERR_MSG, shared.UPLOAD_CHILD_REQUEST_REPOSITORY) + "GetWalletUploadChildRequests - "
	var internalErr error = errors.New(noti.INTERNALL_ERR_MSG)

	var queryCondition string = fmt.Sprintf("created_by = %s ORDER BY created_at DESC", id)
	var query string = generateRetrieveQuery(generateRetrieveQueryRequest{
		table:       upload_child_request_table,
		limitAmount: upload_child_request_limit_record,
		condition:   queryCondition,
		page:        page,
		isGetCount:  false,
	})

	rows, err := u.db.Query(query)
	if err != nil {
		u.errLogger.Println(errLogMsg + err.Error())
		return nil, 0, internalErr
	}

	var res []entities.UploadChildRequest
	for rows.Next() {
		var x entities.UploadChildRequest
		if err := rows.Scan(
			x.ID, x.IdentityCode, x.AvatarBlobId,
			x.Region, x.FirstName, x.LastName, x.Gender, x.DateOfBirth,
			x.Aprrovers, x.Refusers, x.RefuseReasons, x.Status, x.IsConfirmUpload,
			x.CreatedBy, x.CreatedAt, x.UpdatedAt, x.ClosedAt); err != nil {

			u.errLogger.Println(errLogMsg + err.Error())
			return nil, 0, internalErr
		}

		res = append(res, x)
	}

	var totalRecords int
	u.db.QueryRow(generateCountTotalRecordsQuery(upload_child_request_table, queryCondition)).Scan(&totalRecords)

	return res, caculateTotalPages(totalRecords, upload_child_request_limit_record), nil
}

// IsChildRequested implements repository.IUploadChildRequestRepository.
func (u *uploadChildRepo) IsChildRequested(identityCode string, ctx context.Context) (bool, error) {
	var query string = "SELECT * FROM " + upload_child_request_table + " WHERE identity_code = $1 AND (status = 'Pending' OR status = 'Approved') LIMIT 1"

	var id string
	if err := u.db.QueryRow(query).Scan(&id); err != nil {
		if err == sql.ErrNoRows {
			return false, nil
		}

		u.errLogger.Println(fmt.Sprintf(noti.REPO_ERR_MSG, shared.UPLOAD_CHILD_REQUEST_REPOSITORY) + "IsChildRequested - " + err.Error())
		return false, errors.New(noti.INTERNALL_ERR_MSG)
	}

	return id != "", nil
}

// UpdateUploadChildRequest implements repository.IUploadChildRequestRepository.
func (u *uploadChildRepo) UpdateUploadChildRequest(req entities.UploadChildRequest, ctx context.Context) error {
	var query string = "UPDATE " + registraion_request_table + " SET " +
		"region = $1, first_name = $2, last_name = $3, gender = $4, " +
		"date_of_birth = $5, approvers = $6, refusers = $7, refuse_reasons = $8, " +
		"status = $9, is_confirm_upload = $10, updated_at = $11 WHERE id = $12"

	var errLogMsg string = fmt.Sprintf(noti.REPO_ERR_MSG, shared.UPLOAD_CHILD_REQUEST_REPOSITORY) + "UpdateUploadChildRequest - "
	var internalErr error = errors.New(noti.INTERNALL_ERR_MSG)

	res, err := u.db.Exec(query, req.Region, req.FirstName, req.LastName, req.Gender,
		req.DateOfBirth, req.Aprrovers, req.Refusers, req.RefuseReasons,
		req.Status, req.IsConfirmUpload, req.UpdatedAt, req.ID)
	if err != nil {
		u.errLogger.Println(errLogMsg + err.Error())
		return internalErr
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		u.errLogger.Println(errLogMsg + err.Error())
		return internalErr
	}

	if rowsAffected == 0 {
		return errors.New(fmt.Sprintf(noti.UNDEFINED_OBJECT_WARN_MSG, upload_child_request_table))
	}

	return nil
}
