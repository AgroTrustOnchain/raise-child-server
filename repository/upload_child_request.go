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
		" (id, profile_id, identity_code, avatar_blob_id, home_blob_id, " +
		"region, first_name, last_name, gender, date_of_birth, home_address, " +
		"first_guardian_name, first_guardian_phone, first_guardian_relation, first_guardian_identity_card_blob_id, " +
		"second_guardian_name, second_guardian_phone, second_guardian_relation, second_guardian_identity_card_blob_id, " +
		"ai_evaluation, status, created_by, created_at, updated_at) " +
		"values ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, " +
		"$13, $14, $15, $16, $17, $18, $19, $20, $21, $22, $23, $24)"

	var errLogMsg string = fmt.Sprintf(noti.REPO_ERR_MSG, shared.UPLOAD_CHILD_REQUEST_REPOSITORY) + "CreateUploadChildRequest - "
	var secondGuardianName, secondGuardianPhone, secondGuardianRelation, secondGuardianIdentityBlob *string
	if req.SecondGuardianProfile != nil {
		secondGuardianName = &req.SecondGuardianProfile.FullName
		secondGuardianPhone = &req.SecondGuardianProfile.PhoneNumber
		secondGuardianRelation = &req.SecondGuardianProfile.Relation
		secondGuardianIdentityBlob = &req.SecondGuardianProfile.IdentityCardBlobID
	}

	if _, err := u.db.ExecContext(ctx, query, req.ID, req.ProfileID, req.IdentityCode, req.AvatarBlobId, req.HomeBlobID,
		req.Region, req.FirstName, req.LastName, req.Gender, req.DateOfBirth, req.HomeAddress,
		req.FirstGuardianProfile.FullName, req.FirstGuardianProfile.PhoneNumber, req.FirstGuardianProfile.Relation, req.FirstGuardianProfile.IdentityCardBlobID,
		secondGuardianName, secondGuardianPhone, secondGuardianRelation, secondGuardianIdentityBlob,
		req.AIEvaluation, req.Status, req.CreatedBy, req.UpdatedAt, req.ClosedAt); err != nil {

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
	var secondGuardianName, secondGuardianPhone, secondGuardianRelation, secondGuardianIdentityBlob *string

	if err := u.db.QueryRowContext(ctx, query, id).Scan(
		&res.ID, &res.ProfileID, &res.IdentityCode, &res.AvatarBlobId, &res.HomeBlobID,
		&res.Region, &res.FirstName, &res.LastName, &res.Gender, &res.DateOfBirth, &res.HomeAddress,
		&res.FirstGuardianProfile.FullName, &res.FirstGuardianProfile.PhoneNumber, &res.FirstGuardianProfile.Relation, &res.FirstGuardianProfile.IdentityCardBlobID,
		&secondGuardianName, &secondGuardianPhone, &secondGuardianRelation, &secondGuardianIdentityBlob,
		pq.Array(&res.Approvers), pq.Array(&res.Refusers), pq.Array(&res.RefuseReasons),
		&res.AIEvaluation, &res.Status, &res.ReviewStatus, &res.IsConfirmUpload,
		&res.CreatedBy, &res.ReviewedBy, &res.CreatedAt, &res.UpdatedAt, &res.ClosedAt); err != nil {

		if err == sql.ErrNoRows {
			return nil, nil
		}

		u.errLogger.Println(errLogMsg + err.Error())
		return nil, errors.New(noti.INTERNALL_ERR_MSG)
	}

	if secondGuardianName != nil {
		res.SecondGuardianProfile = &entities.ChildGuardianProfile{
			FullName:           *secondGuardianName,
			PhoneNumber:        *secondGuardianPhone,
			Relation:           *secondGuardianRelation,
			IdentityCardBlobID: *secondGuardianIdentityBlob,
		}
	}

	return &res, nil
}

// GetUploadChildRequests implements repository.IUploadChildRequestRepository.
func (u *uploadChildRepo) GetUploadChildRequests(req request.GetUploadChildRequests, ctx context.Context) ([]entities.UploadChildRequest, int, error) {
	var errLogMsg string = fmt.Sprintf(noti.REPO_ERR_MSG, shared.UPLOAD_CHILD_REQUEST_REPOSITORY) + "GetUploadChildRequests - "
	var internalErr error = errors.New(noti.INTERNALL_ERR_MSG)

	var queryCondition string
	var isHavePreviosCondition bool = false
	if req.Keyword != "" {
		queryCondition += fmt.Sprintf("(identity_code LIKE '%s' OR LOWER(first_name) LIKE LOWER('%%%%%s%%%%') OR LOWER(last_name) LIKE LOWER('%%%%%s%%%%') OR date_of_birth LIKE '%%%%%s%%%%' OR LOWER(home_address) LIKE LOWER('%%%%%s%%%%') OR LOWER(first_guardian_name) LIKE LOWER('%%%%%s%%%%')  OR LOWER(first_guardian_phone) LIKE LOWER('%%%%%s%%%%') LOWER(second_guardian_name) LIKE LOWER('%%%%%s%%%%') OR LOWER(second_guardian_phone) LIKE LOWER('%%%%%s%%%%'))",
			req.Keyword, req.Keyword, req.Keyword, req.Keyword, req.Keyword, req.Keyword, req.Keyword, req.Keyword)
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
		limitAmount: req.PageSize,
		condition:   queryCondition,
		page:        req.Page,
		isGetCount:  false,
	})

	rows, err := u.db.QueryContext(ctx, query)
	if err != nil {
		u.errLogger.Println(errLogMsg + err.Error())
		return nil, 0, internalErr
	}

	var res []entities.UploadChildRequest
	for rows.Next() {
		var x entities.UploadChildRequest
		var secondGuardianName, secondGuardianPhone, secondGuardianRelation, secondGuardianIdentityBlob *string

		if err := rows.Scan(
			&x.ID, &x.ProfileID, &x.IdentityCode, &x.AvatarBlobId, &x.HomeBlobID,
			&x.Region, &x.FirstName, &x.LastName, &x.Gender, &x.DateOfBirth, &x.HomeAddress,
			&x.FirstGuardianProfile.FullName, &x.FirstGuardianProfile.PhoneNumber, &x.FirstGuardianProfile.Relation, &x.FirstGuardianProfile.IdentityCardBlobID,
			&secondGuardianName, &secondGuardianPhone, &secondGuardianRelation, &secondGuardianIdentityBlob,
			pq.Array(&x.Approvers), pq.Array(&x.Refusers), pq.Array(&x.RefuseReasons),
			&x.AIEvaluation, &x.Status, &x.ReviewStatus, &x.IsConfirmUpload,
			&x.CreatedBy, &x.ReviewedBy, &x.CreatedAt, &x.UpdatedAt, &x.ClosedAt); err != nil {

			u.errLogger.Println(errLogMsg + err.Error())
			return nil, 0, internalErr
		}

		if secondGuardianName != nil {
			x.SecondGuardianProfile = &entities.ChildGuardianProfile{
				FullName:           *secondGuardianName,
				PhoneNumber:        *secondGuardianPhone,
				Relation:           *secondGuardianRelation,
				IdentityCardBlobID: *secondGuardianIdentityBlob,
			}
		}

		res = append(res, x)
	}

	var totalRecords int
	u.db.QueryRowContext(ctx, generateCountTotalRecordsQuery(upload_child_request_table, queryCondition)).Scan(&totalRecords)

	return res, caculateTotalPages(totalRecords, req.PageSize), nil
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

	rows, err := u.db.QueryContext(ctx, query)
	if err != nil {
		u.errLogger.Println(errLogMsg + err.Error())
		return nil, 0, internalErr
	}

	var res []entities.UploadChildRequest
	for rows.Next() {
		var x entities.UploadChildRequest
		var secondGuardianName, secondGuardianPhone, secondGuardianRelation, secondGuardianIdentityBlob *string

		if err := rows.Scan(
			&x.ID, &x.ProfileID, &x.IdentityCode, &x.AvatarBlobId, &x.HomeBlobID,
			&x.Region, &x.FirstName, &x.LastName, &x.Gender, &x.DateOfBirth, &x.HomeAddress,
			&x.FirstGuardianProfile.FullName, &x.FirstGuardianProfile.PhoneNumber, &x.FirstGuardianProfile.Relation, &x.FirstGuardianProfile.IdentityCardBlobID,
			&secondGuardianName, &secondGuardianPhone, &secondGuardianRelation, &secondGuardianIdentityBlob,
			pq.Array(&x.Approvers), pq.Array(&x.Refusers), pq.Array(&x.RefuseReasons),
			&x.AIEvaluation, &x.Status, &x.ReviewStatus, &x.IsConfirmUpload,
			&x.CreatedBy, &x.ReviewedBy, &x.CreatedAt, &x.UpdatedAt, &x.ClosedAt); err != nil {

			u.errLogger.Println(errLogMsg + err.Error())
			return nil, 0, internalErr
		}

		if secondGuardianName != nil {
			x.SecondGuardianProfile = &entities.ChildGuardianProfile{
				FullName:           *secondGuardianName,
				PhoneNumber:        *secondGuardianPhone,
				Relation:           *secondGuardianRelation,
				IdentityCardBlobID: *secondGuardianIdentityBlob,
			}
		}

		res = append(res, x)
	}

	var totalRecords int
	u.db.QueryRowContext(ctx, generateCountTotalRecordsQuery(upload_child_request_table, queryCondition)).Scan(&totalRecords)

	return res, caculateTotalPages(totalRecords, upload_child_request_limit_record), nil
}

// IsChildRequested implements repository.IUploadChildRequestRepository.
func (u *uploadChildRepo) IsChildRequested(identityCode string, ctx context.Context) (bool, error) {
	var query string = "SELECT id FROM " + upload_child_request_table + " WHERE identity_code = $1 AND (status = 'Pending' OR status = 'Approved') LIMIT 1"

	var id string
	if err := u.db.QueryRowContext(ctx, query, identityCode).Scan(&id); err != nil {
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
	var query string = "UPDATE " + upload_child_request_table + " SET " +
		"region = $1, first_name = $2, last_name = $3, gender = $4, " +
		"date_of_birth = $5, approvers = $6, refusers = $7, refuse_reasons = $8, " +
		"status = $9, is_confirm_upload = $10, updated_at = $11 WHERE id = $12"

	var errLogMsg string = fmt.Sprintf(noti.REPO_ERR_MSG, shared.UPLOAD_CHILD_REQUEST_REPOSITORY) + "UpdateUploadChildRequest - "
	var internalErr error = errors.New(noti.INTERNALL_ERR_MSG)

	res, err := u.db.ExecContext(ctx, query, req.Region, req.FirstName, req.LastName, req.Gender,
		req.DateOfBirth, pq.Array(req.Approvers), pq.Array(req.Refusers), pq.Array(req.RefuseReasons),
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

// GetPendingRequests implements repository.IUploadChildRequestRepository.
func (u *uploadChildRepo) GetPendingRequests(ctx context.Context) ([]entities.BackgroundRecord, []entities.BackgroundRecord, error) {
	var query string = "SELECT id, approvers, refusers, created_by, status FROM " + upload_child_request_table + " WHERE is_available_to_confirm = false AND closed_at <= NOW() AND (status = 'Pending' OR status = 'Approved')"
	var errLogMsg string = fmt.Sprintf(noti.REPO_ERR_MSG, shared.UPLOAD_CHILD_REQUEST_REPOSITORY) + "GetPendingRequests - "
	var internalErr error = errors.New(noti.INTERNALL_ERR_MSG)

	rows, err := u.db.QueryContext(ctx, query)
	if err != nil {
		u.errLogger.Println(errLogMsg + err.Error())
		return nil, nil, internalErr
	}

	var pendingRes, approvedRes []entities.BackgroundRecord
	for rows.Next() {
		var x entities.BackgroundRecord
		var status string
		if err := rows.Scan(
			&x.ID, pq.Array(&x.Approvers), pq.Array(&x.Refusers), &x.Sender, &status); err != nil {

			u.errLogger.Println(errLogMsg + err.Error())
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

// SetApprovedStatuses implements repository.IUploadChildRequestRepository.
func (u *uploadChildRepo) SetApprovedStatuses(reqs []entities.BackgroundRecord, ctx context.Context) error {
	var query string = "UPDATE " + upload_child_request_table + " SET status = 'Approved', is_available_to_confirm = true, updated_at = $1 WHERE "
	for i, req := range reqs {
		query += "id = " + req.ID
		if i < len(reqs)-1 {
			query += " OR "
		}
	}

	if _, err := u.db.ExecContext(ctx, query, time.Now()); err != nil {
		u.errLogger.Println(fmt.Sprintf(noti.REPO_ERR_MSG, shared.UPLOAD_CHILD_REQUEST_REPOSITORY) + "SetApprovedStatuses - " + err.Error())
		return errors.New(noti.INTERNALL_ERR_MSG)
	}

	return nil
}

// SetRefusedStatuses implements repository.IUploadChildRequestRepository.
func (u *uploadChildRepo) SetRefusedStatuses(reqs []entities.BackgroundRecord, ctx context.Context) error {
	var query string = "UPDATE " + upload_child_request_table + " SET status = 'Refused', updated_at = $1 WHERE "
	for i, req := range reqs {
		query += "id = " + req.ID
		if i < len(reqs)-1 {
			query += " OR "
		}
	}

	if _, err := u.db.ExecContext(ctx, query, time.Now()); err != nil {
		u.errLogger.Println(fmt.Sprintf(noti.REPO_ERR_MSG, shared.UPLOAD_CHILD_REQUEST_REPOSITORY) + "SetRefusedStatuses - " + err.Error())
		return errors.New(noti.INTERNALL_ERR_MSG)
	}

	return nil
}

// SetReviewStatus implements repository.IUploadChildRequestRepository.
func (u *uploadChildRepo) SetReviewStatus(id string, reviewStatus string, reviewer string, closedAt *time.Time, ctx context.Context) error {
	var query string = "UPDATE " + upload_child_request_table + " SET review_status = $1, reviewed_by = $2, closed_at = $3 WHERE id = $4"

	var errLogMsg string = fmt.Sprintf(noti.REPO_ERR_MSG, shared.UPLOAD_CHILD_REQUEST_REPOSITORY) + "SetReviewStatus - "
	var internalErr error = errors.New(noti.INTERNALL_ERR_MSG)

	res, err := u.db.ExecContext(ctx, query, reviewStatus, reviewer, closedAt, id)
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
