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

type offChainDonationRepo struct {
	db        *sql.DB
	errLogger *log.Logger
}

const (
	offchain_donation_table string = "donations"
)

func InitializeOffChainDonationRepository(db *sql.DB, errLogger *log.Logger) repository.IOffChainDonationRepository {
	return &offChainDonationRepo{
		db:        db,
		errLogger: errLogger,
	}
}

// CreateDonation implements repository.IOffChainDonationRepository.
func (o *offChainDonationRepo) CreateDonation(donation entities.OffChainDonation, ctx context.Context) error {
	var query string = "INSERT INTO " + offchain_donation_table +
		" (id, purpose, target, meal_duration_id, created_at) " +
		"values ($1, $2, $3, $4, $5)"

	var errLogMsg string = fmt.Sprintf(noti.REPO_ERR_MSG, shared.OFFCHAIN_DONATION_REPOSITORY) + "CreateDonation - "

	if _, err := o.db.ExecContext(ctx, query, donation.ID, donation.Purpose, donation.Target, donation.MealDurationID, donation.CreatedAt); err != nil {

		o.errLogger.Println(errLogMsg + err.Error())
		return errors.New(noti.INTERNALL_ERR_MSG)
	}

	return nil
}

// GetDonation implements repository.IOffChainDonationRepository.
func (o *offChainDonationRepo) GetDonation(id string, ctx context.Context) (*entities.OffChainDonation, error) {
	var query string = "SELECT * FROM " + offchain_donation_table + " WHERE id = $1"
	var errLogMsg string = fmt.Sprintf(noti.REPO_ERR_MSG, shared.OFFCHAIN_DONATION_REPOSITORY) + "GetDonation - "

	var res entities.OffChainDonation
	if err := o.db.QueryRowContext(ctx, query, id).Scan(
		&res.ID, &res.Purpose, &res.Target, &res.MealDurationID, &res.CreatedAt); err != nil {

		if err == sql.ErrNoRows {
			return nil, nil
		}

		o.errLogger.Println(errLogMsg + err.Error())
		return nil, errors.New(noti.INTERNALL_ERR_MSG)
	}

	return &res, nil
}
