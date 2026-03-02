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
)

type bankProfileRepo struct {
	db        *sql.DB
	errLogger *log.Logger
}

const bank_profile_table string = "bank_profiles"

func InitializeBankProfileRepository(db *sql.DB, errLogger *log.Logger) repository.IBankProfileRepository {
	return &bankProfileRepo{
		db:        db,
		errLogger: errLogger,
	}
}

// CreateBankProfile implements repository.IBankProfileRepository.
func (b *bankProfileRepo) CreateBankProfile(bp entities.BankProfile, ctx context.Context) error {
	var query string = "INSERT INTO " + bank_profile_table +
		" (id, profile_id, owner, bank_org, bank_code, owner_name, " +
		"payos_client_id, payos_api_key, payos_check_sum_key, " +
		"created_at, updated_at) " +
		"values ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)"
	var errLogMsg string = fmt.Sprintf(noti.REPO_ERR_MSG, shared.BANK_PROFILE_REPOSITORY) + "CreatePayment - "

	if _, err := b.db.ExecContext(ctx, query, bp.ID, bp.ProfileID, bp.Owner, bp.BankOrg, bp.BankCode, bp.OwnerName,
		bp.PayosClientID, bp.PayosApiKey, bp.PayosCheckSumKey,
		bp.CreatedAt, bp.UpdatedAt); err != nil {

		b.errLogger.Println(errLogMsg + err.Error())
		return errors.New(noti.INTERNALL_ERR_MSG)
	}

	return nil
}

// GetBankProfile implements repository.IBankProfileRepository.
func (b *bankProfileRepo) GetBankProfileById(id string, ctx context.Context) (*entities.BankProfile, error) {
	var query string = "SELECT * FROM " + bank_profile_table + " WHERE id = $1"
	var errLogMsg string = fmt.Sprintf(noti.REPO_ERR_MSG, shared.BANK_PROFILE_REPOSITORY) + "GetPaymentById - "

	var res entities.BankProfile
	if err := b.db.QueryRowContext(ctx, query, id).Scan(
		&res.ID, &res.ProfileID, &res.Owner, &res.BankOrg, &res.BankCode, &res.OwnerName,
		&res.PayosClientID, &res.PayosApiKey, &res.PayosCheckSumKey,
		&res.CreatedAt, &res.UpdatedAt); err != nil {

		if err == sql.ErrNoRows {
			return nil, nil
		}

		b.errLogger.Println(errLogMsg + err.Error())
		return nil, errors.New(noti.INTERNALL_ERR_MSG)
	}

	return &res, nil
}

// GetBankProfileByOwner implements repository.IBankProfileRepository.
func (b *bankProfileRepo) GetBankProfileByOwner(owner string, ctx context.Context) (*entities.BankProfile, error) {
	var query string = "SELECT * FROM " + bank_profile_table + " WHERE owner = $1"
	var errLogMsg string = fmt.Sprintf(noti.REPO_ERR_MSG, shared.BANK_PROFILE_REPOSITORY) + "GetPaymentById - "

	var res entities.BankProfile
	if err := b.db.QueryRowContext(ctx, query, owner).Scan(
		&res.ID, &res.ProfileID, &res.Owner, &res.BankOrg, &res.BankCode, &res.OwnerName,
		&res.PayosClientID, &res.PayosApiKey, &res.PayosCheckSumKey,
		&res.CreatedAt, &res.UpdatedAt); err != nil {

		if err == sql.ErrNoRows {
			return nil, nil
		}

		b.errLogger.Println(errLogMsg + err.Error())
		return nil, errors.New(noti.INTERNALL_ERR_MSG)
	}

	return &res, nil
}

// UpdateBankProfile implements repository.IBankProfileRepository.
func (b *bankProfileRepo) UpdateBankProfile(bp entities.BankProfile, ctx context.Context) error {
	var errLogMsg string = fmt.Sprintf(noti.REPO_ERR_MSG, shared.BANK_PROFILE_REPOSITORY) + "UpdatePayment - "
	var query string = "UPDATE " + bank_profile_table + " SET bank_org = $1, bank_code = $2, owner_name = $3, payos_client_id = $4, payos_api_key = $5, payos_check_sum_key = $6, updated_at = $7 WHERE id = $8"

	res, err := b.db.ExecContext(ctx, query, bp.BankOrg, bp.BankCode, bp.OwnerName, bp.PayosClientID, bp.PayosApiKey, bp.PayosCheckSumKey, bp.UpdatedAt, bp.ID)

	var internalErr error = errors.New(noti.INTERNALL_ERR_MSG)

	if err != nil {
		b.errLogger.Println(errLogMsg + err.Error())
		return internalErr
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		b.errLogger.Println(errLogMsg + err.Error())
		return internalErr
	}

	if rowsAffected == 0 {
		return errors.New(fmt.Sprintf(noti.UNDEFINED_OBJECT_WARN_MSG, bank_profile_table))
	}

	return nil
}

// IsBankWithSubExist implements repository.IBankProfileRepository.
func (b *bankProfileRepo) IsBankWithSubExist(profile_id string, ctx context.Context) (bool, error) {
	var query string = "SELECT id FROM " + bank_profile_table + " WHERE profile_id = $1 LIMIT 1"
	var errLogMsg string = fmt.Sprintf(noti.REPO_ERR_MSG, shared.BANK_PROFILE_REPOSITORY) + "IsBankWithSubExist - "
	var internalErr error = errors.New(noti.INTERNALL_ERR_MSG)

	var id string
	if err := b.db.QueryRowContext(ctx, query, profile_id).Scan(&id); err != nil {
		b.errLogger.Println(errLogMsg + err.Error())
		return false, internalErr
	}

	return id != "", nil
}
