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
	"time"
)

type profileRepo struct {
	db        *sql.DB
	errLogger *log.Logger
}

const profile_table string = "profiles"

func InitializeProfileRepository(db *sql.DB, errLogger *log.Logger) repository.IProfileRepository {
	return &profileRepo{
		db:        db,
		errLogger: errLogger,
	}
}

// IsPersonalInfoExist implements repository.IProfileRepository.
func (p *profileRepo) IsPersonalInfoExist(identityCode string, phoneNumber string, email string, ctx context.Context) (bool, error) {
	var query string = "SELECT id FROM " + profile_table + " WHERE LOWER(email) = LOWER(" + email + ") OR phone_number = " + phoneNumber + " OR identity_code = " + identityCode + " LIMIT 1"

	var id string
	if err := p.db.QueryRowContext(ctx, query).Scan(&id); err != nil {
		if err == sql.ErrNoRows {
			return false, nil
		}

		p.errLogger.Println(fmt.Sprintf(noti.REPO_ERR_MSG, shared.PROFILE_REPOSITORY) + "IsPersonalInfoExist - " + err.Error())
		return false, errors.New(noti.INTERNALL_ERR_MSG)
	}

	return id != "", nil
}

// IsEmailRegistered implements repository.IProfileRepository.
func (p *profileRepo) IsEmailRegistered(email string, ctx context.Context) (bool, error) {
	var errLogMsg string = fmt.Sprintf(noti.REPO_ERR_MSG, shared.PROFILE_REPOSITORY) + "IsEmailRegistered - "
	var query string = "SELECT id FROM " + profile_table + " WHERE LOWER(email) = LOWER(" + email + ") LIMIT 1"

	var id string
	if err := p.db.QueryRowContext(ctx, query).Scan(&id); err != nil {
		if err == sql.ErrNoRows {
			return false, nil
		}

		p.errLogger.Println(errLogMsg + err.Error())
		return false, errors.New(noti.INTERNALL_ERR_MSG)
	}

	return id != "", nil
}

// IsPhoneNumberRegistered implements repository.IProfileRepository.
func (p *profileRepo) IsPhoneNumberRegistered(phoneNumber string, ctx context.Context) (bool, error) {
	var errLogMsg string = fmt.Sprintf(noti.REPO_ERR_MSG, shared.PROFILE_REPOSITORY) + "IsPhoneNumberRegistered - "
	var query string = "SELECT id FROM " + profile_table + " WHERE phone_number = " + phoneNumber + " LIMIT 1"

	var id string
	if err := p.db.QueryRowContext(ctx, query).Scan(&id); err != nil {
		if err == sql.ErrNoRows {
			return false, nil
		}

		p.errLogger.Println(errLogMsg + err.Error())
		return false, errors.New(noti.INTERNALL_ERR_MSG)
	}

	return id != "", nil
}

// UploadProfile implements repository.IProfileRepository.
func (p *profileRepo) UploadProfile(pfl entities.Profile, ctx context.Context) error {
	var errLogMsg string = fmt.Sprintf(noti.REPO_ERR_MSG, shared.PROFILE_REPOSITORY) + "UploadProfile - "
	var query string = "UPDATE " + profile_table + " SET identity_code = $1, first_name = $2, last_name = $3, gender = $4, date_of_birth = $5, phone_number = $6, email = $7, updated_at = $8 WHERE id = $9"
	var internalErr error = errors.New(noti.INTERNALL_ERR_MSG)

	res, err := p.db.ExecContext(ctx, query, pfl.IdentityCode, pfl.FirstName, pfl.LastName,
		pfl.Gender, pfl.DateOfBirth, pfl.PhoneNumber, pfl.Email, pfl.UpdatedAt, pfl.ID)
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
		return errors.New(fmt.Sprintf(noti.UNDEFINED_OBJECT_WARN_MSG, profile_table))
	}

	return nil
}

// Login implements repository.IProfileRepository.
func (p *profileRepo) Login(id string, token string, ctx context.Context) error {
	var errLogMsg string = fmt.Sprintf(noti.REPO_ERR_MSG, shared.PROFILE_REPOSITORY) + "Login - "
	var query string = "UPDATE " + profile_table + " SET token = $1, updated_at = $2 WHERE id = $3"
	var internalErr error = errors.New(noti.INTERNALL_ERR_MSG)

	res, err := p.db.ExecContext(ctx, query, token, time.Now(), id)
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
		return errors.New(fmt.Sprintf(noti.UNDEFINED_OBJECT_WARN_MSG, profile_table))
	}

	return nil
}

// Logout implements repository.IProfileRepository.
func (p *profileRepo) Logout(id string, ctx context.Context) error {
	var errLogMsg string = fmt.Sprintf(noti.REPO_ERR_MSG, shared.PROFILE_REPOSITORY) + "Logout - "
	var query string = "UPDATE " + profile_table + " SET token = '', updated_at = $1 WHERE id = $2"
	var internalErr error = errors.New(noti.INTERNALL_ERR_MSG)

	res, err := p.db.ExecContext(ctx, query, time.Now(), id)
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
		return errors.New(fmt.Sprintf(noti.UNDEFINED_OBJECT_WARN_MSG, profile_table))
	}

	return nil
}

// CreateProfile implements repository.IProfileRepository.
func (p *profileRepo) CreateProfile(pfl entities.Profile, ctx context.Context) error {
	var query string = "INSERT INTO " + profile_table +
		" (id, salt, created_at, updated_at) values ($1, $2, $3, $4)"

	var errLogMsg string = fmt.Sprintf(noti.REPO_ERR_MSG, shared.PROFILE_REPOSITORY) + "CreateProfile - "

	if _, err := p.db.ExecContext(ctx, query, pfl.ID, pfl.Salt, pfl.CreatedAt, pfl.UpdatedAt); err != nil {

		p.errLogger.Println(errLogMsg + err.Error())
		return errors.New(noti.INTERNALL_ERR_MSG)
	}

	return nil
}

// GetProfile implements repository.IProfileRepository.
func (p *profileRepo) GetProfile(id string, ctx context.Context) (*entities.Profile, error) {
	var query string = "SELECT * FROM " + profile_table + " WHERE id = $1"
	var errLogMsg string = fmt.Sprintf(noti.REPO_ERR_MSG, shared.PROFILE_REPOSITORY) + "GetProfile - "

	var res entities.Profile
	if err := p.db.QueryRowContext(ctx, query, id).Scan(
		&res.ID, &res.Salt, &res.Status, &res.IdentityCode, &res.FirstName, &res.LastName,
		&res.Gender, res.DateOfBirth, &res.PhoneNumber, &res.Email,
		&res.Token, &res.UpdatedAt, &res.CreatedAt); err != nil {

		if err == sql.ErrNoRows {
			return nil, nil
		}

		p.errLogger.Println(errLogMsg + err.Error())
		return nil, errors.New(noti.INTERNALL_ERR_MSG)
	}

	return &res, nil
}

// GetFirstProfile implements repository.IProfileRepository.
func (p *profileRepo) GetFirstProfile(ctx context.Context) (*entities.Profile, error) {
	var query string = "SELECT * FROM " + profile_table + " ORDER BY created_at ASC LIMIT 1"
	var errLogMsg string = fmt.Sprintf(noti.REPO_ERR_MSG, shared.PROFILE_REPOSITORY) + "GetFirstProfile - "

	var res entities.Profile
	if err := p.db.QueryRowContext(ctx, query).Scan(
		&res.ID, &res.Salt, &res.Status, &res.IdentityCode, &res.FirstName, &res.LastName,
		&res.Gender, res.DateOfBirth, &res.PhoneNumber, &res.Email,
		&res.Token, &res.UpdatedAt, &res.CreatedAt); err != nil {

		if err == sql.ErrNoRows {
			return nil, nil
		}

		p.errLogger.Println(errLogMsg + err.Error())
		return nil, errors.New(noti.INTERNALL_ERR_MSG)
	}

	return &res, nil
}

// GetProfileOfFirsts implements repository.IProfileRepository.
func (p *profileRepo) GetProfileOfFirsts(position int, ctx context.Context) (*entities.Profile, error) {
	var query string = fmt.Sprintf("SELECT * FROM %s ORDER BY created_at ASC OFFSET %d LIMIT 1", profile_table, (position - 1))
	var errLogMsg string = fmt.Sprintf(noti.REPO_ERR_MSG, shared.PROFILE_REPOSITORY) + "GetProfileOfFirsts - "

	var res entities.Profile
	if err := p.db.QueryRowContext(ctx, query).Scan(
		&res.ID, &res.Salt, &res.Status, &res.IdentityCode, &res.FirstName, &res.LastName,
		&res.Gender, res.DateOfBirth, &res.PhoneNumber, &res.Email,
		&res.Token, &res.UpdatedAt, &res.CreatedAt); err != nil {

		if err == sql.ErrNoRows {
			return nil, nil
		}

		p.errLogger.Println(errLogMsg + err.Error())
		return nil, errors.New(noti.INTERNALL_ERR_MSG)
	}

	return &res, nil
}
