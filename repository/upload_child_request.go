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
		" (id, profile_id, identity_code, avatar_blob_id, home_blob_id, " +
		"region, first_name, last_name, gender, date_of_birth, home_address, " +
		"first_guardian_name, first_guardian_phone, first_guardian_relation, first_guardian_identity_card_blob_id, " +
		"second_guardian_name, second_guardian_phone, second_guardian_relation, second_guardian_identity_card_blob_id, " +
		"status, created_by, created_at, updated_at, birth_certificate_blob_id) " +
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
		req.Status, req.CreatedBy, req.CreatedAt, req.UpdatedAt, req.BirthCertificateBlobID); err != nil {

		u.errLogger.Println(errLogMsg + err.Error())
		return errors.New(noti.INTERNALL_ERR_MSG)
	}

	return nil
}

// GetUploadChildRequest implements repository.IUploadChildRequestRepository.
func (u *uploadChildRepo) GetUploadChildRequest(id string, ctx context.Context) (*entities.UploadChildRequest, error) {
	var query string = "SELECT * FROM " + upload_child_request_table + " WHERE id = $1"
	var errLogMsg string = fmt.Sprintf(noti.REPO_ERR_MSG, shared.UPLOAD_CHILD_REQUEST_REPOSITORY) + "GetUploadChildRequest - "

	var res entities.UploadChildRequest
	var secondGuardianName, secondGuardianPhone, secondGuardianRelation, secondGuardianIdentityBlob *string

	if err := u.db.QueryRowContext(ctx, query, id).Scan(
		&res.ID, &res.ProfileID, &res.IdentityCode, &res.AvatarBlobId, &res.HomeBlobID,
		&res.Region, &res.FirstName, &res.LastName, &res.Gender, &res.DateOfBirth, &res.HomeAddress,
		&res.FirstGuardianProfile.FullName, &res.FirstGuardianProfile.PhoneNumber, &res.FirstGuardianProfile.Relation, &res.FirstGuardianProfile.IdentityCardBlobID,
		&secondGuardianName, &secondGuardianPhone, &secondGuardianRelation, &secondGuardianIdentityBlob,
		&res.Status, &res.IsConfirmUpload, &res.CreatedBy, &res.ReviewedBy, &res.OnchainID,
		&res.CreatedAt, &res.UpdatedAt, &res.BirthCertificateBlobID); err != nil {

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

// GetUploadChildRequestByOnchainId implements repository.IUploadChildRequestRepository.
func (u *uploadChildRepo) GetUploadChildRequestByOnchainId(id string, ctx context.Context) (*entities.UploadChildRequest, error) {
	var query string = "SELECT * FROM " + upload_child_request_table + " WHERE on_chain_id = $1"
	var errLogMsg string = fmt.Sprintf(noti.REPO_ERR_MSG, shared.UPLOAD_CHILD_REQUEST_REPOSITORY) + "GetUploadChildRequestByOnchainId - "

	var res entities.UploadChildRequest
	var secondGuardianName, secondGuardianPhone, secondGuardianRelation, secondGuardianIdentityBlob *string

	if err := u.db.QueryRowContext(ctx, query, id).Scan(
		&res.ID, &res.ProfileID, &res.IdentityCode, &res.AvatarBlobId, &res.HomeBlobID,
		&res.Region, &res.FirstName, &res.LastName, &res.Gender, &res.DateOfBirth, &res.HomeAddress,
		&res.FirstGuardianProfile.FullName, &res.FirstGuardianProfile.PhoneNumber, &res.FirstGuardianProfile.Relation, &res.FirstGuardianProfile.IdentityCardBlobID,
		&secondGuardianName, &secondGuardianPhone, &secondGuardianRelation, &secondGuardianIdentityBlob,
		&res.Status, &res.IsConfirmUpload, &res.CreatedBy, &res.ReviewedBy, &res.OnchainID,
		&res.CreatedAt, &res.UpdatedAt, &res.BirthCertificateBlobID); err != nil {

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
		queryCondition = fmt.Sprintf(`
		(
			LOWER(identity_code) LIKE LOWER('%s') 
			OR LOWER(first_name) LIKE LOWER('%s') 
			OR LOWER(last_name) LIKE LOWER('%s') 
			OR date_of_birth LIKE '%s' 
			OR LOWER(home_address) LIKE LOWER('%s') 
			OR LOWER(first_guardian_name) LIKE LOWER('%s')  
			OR LOWER(first_guardian_phone) LIKE LOWER('%s') 
			OR LOWER(second_guardian_name) LIKE LOWER('%s') 
			OR LOWER(second_guardian_phone) LIKE LOWER('%s')
		)
		`,
			req.Keyword, req.Keyword, req.Keyword, req.Keyword, req.Keyword, req.Keyword, req.Keyword, req.Keyword, req.Keyword)
		isHavePreviosCondition = true
	}

	if req.Region != "" {
		if isHavePreviosCondition {
			queryCondition += " AND "
		}

		queryCondition += fmt.Sprintf("LOWER(region) = LOWER('%s')", req.Region)
		isHavePreviosCondition = true
	}

	if req.Gender != "" {
		if isHavePreviosCondition {
			queryCondition += " AND "
		}

		queryCondition += fmt.Sprintf("LOWER(gender) = LOWER('%s')", req.Gender)
		isHavePreviosCondition = true
	}

	if req.Status != "" {
		if isHavePreviosCondition {
			queryCondition += " AND "
		}

		queryCondition += fmt.Sprintf("LOWER(status) = LOWER('%s')", req.Status)
		isHavePreviosCondition = true
	}

	if isHavePreviosCondition {
		queryCondition += " "
	}

	var order string = "DESC"
	if req.SortOrder != "" {
		order = req.SortOrder
	}

	var query string = generateRetrieveQuery(generateRetrieveQueryRequest{
		table:       upload_child_request_table,
		limitAmount: req.PageSize,
		condition:   queryCondition,
		order:       " ORDER BY created_at " + order,
		page:        req.Page,
		isGetCount:  false,
	})

	rows, err := u.db.QueryContext(ctx, query)
	if err != nil {
		u.errLogger.Println(errLogMsg + err.Error())
		return nil, 0, internalErr
	}
	defer rows.Close()

	var res []entities.UploadChildRequest
	for rows.Next() {
		var x entities.UploadChildRequest
		var secondGuardianName, secondGuardianPhone, secondGuardianRelation, secondGuardianIdentityBlob *string

		if err := rows.Scan(
			&x.ID, &x.ProfileID, &x.IdentityCode, &x.AvatarBlobId, &x.HomeBlobID,
			&x.Region, &x.FirstName, &x.LastName, &x.Gender, &x.DateOfBirth, &x.HomeAddress,
			&x.FirstGuardianProfile.FullName, &x.FirstGuardianProfile.PhoneNumber, &x.FirstGuardianProfile.Relation, &x.FirstGuardianProfile.IdentityCardBlobID,
			&secondGuardianName, &secondGuardianPhone, &secondGuardianRelation, &secondGuardianIdentityBlob,
			&x.Status, &x.IsConfirmUpload, &x.CreatedBy, &x.ReviewedBy, &x.OnchainID,
			&x.CreatedAt, &x.UpdatedAt, &x.BirthCertificateBlobID); err != nil {

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
	defer rows.Close()

	var res []entities.UploadChildRequest
	for rows.Next() {
		var x entities.UploadChildRequest
		var secondGuardianName, secondGuardianPhone, secondGuardianRelation, secondGuardianIdentityBlob *string

		if err := rows.Scan(
			&x.ID, &x.ProfileID, &x.IdentityCode, &x.AvatarBlobId, &x.HomeBlobID,
			&x.Region, &x.FirstName, &x.LastName, &x.Gender, &x.DateOfBirth, &x.HomeAddress,
			&x.FirstGuardianProfile.FullName, &x.FirstGuardianProfile.PhoneNumber, &x.FirstGuardianProfile.Relation, &x.FirstGuardianProfile.IdentityCardBlobID,
			&secondGuardianName, &secondGuardianPhone, &secondGuardianRelation, &secondGuardianIdentityBlob,
			&x.Status, &x.IsConfirmUpload, &x.CreatedBy, &x.ReviewedBy, &x.OnchainID,
			&x.CreatedAt, &x.UpdatedAt, &x.BirthCertificateBlobID); err != nil {

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
		"region = $1, first_name = $2, last_name = $3, gender = $4, date_of_birth = $5, " +
		"status = $6, is_confirm_upload = $7, on_chain_id = $8, updated_at = $9 WHERE id = $10"

	var errLogMsg string = fmt.Sprintf(noti.REPO_ERR_MSG, shared.UPLOAD_CHILD_REQUEST_REPOSITORY) + "UpdateUploadChildRequest - "
	var internalErr error = errors.New(noti.INTERNALL_ERR_MSG)

	res, err := u.db.ExecContext(ctx, query, req.Region, req.FirstName, req.LastName, req.Gender,
		req.DateOfBirth, req.Status, req.IsConfirmUpload, req.OnchainID, req.UpdatedAt, req.ID)
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
		return fmt.Errorf(noti.UNDEFINED_OBJECT_WARN_MSG, upload_child_request_table)
	}

	return nil
}
