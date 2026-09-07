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

type mealSupportDurationRepo struct {
	db        *sql.DB
	errLogger *log.Logger
}

const meal_support_duration_table string = "meal_support_durations"

func InitializeMealSupportDurationRepository(db *sql.DB, errLogger *log.Logger) repository.IMealSupportDurationRepository {
	return &mealSupportDurationRepo{
		db:        db,
		errLogger: errLogger,
	}
}

// CreateMealSupportDuration implements repository.IMealSupportDurationRepository.
func (m *mealSupportDurationRepo) CreateMealSupportDuration(duration entities.OffChainMealSupportDuration, ctx context.Context) error {
	var query string = "INSERT INTO " + meal_support_duration_table +
		" (id, start_period, end_period) values ($1, $2, $3)"
	var errLogMsg string = fmt.Sprintf(noti.REPO_ERR_MSG, shared.MEAL_SUPPORT_DURATION_REPOSITORY) + "CreateMealSupportDuration - "

	if _, err := m.db.ExecContext(ctx, query, duration.ID, duration.StartPeriod, duration.EndPeriod); err != nil {

		m.errLogger.Println(errLogMsg + err.Error())
		return errors.New(noti.INTERNALL_ERR_MSG)
	}

	return nil
}

// GetMealSupportDuration implements repository.IMealSupportDurationRepository.
func (m *mealSupportDurationRepo) GetMealSupportDuration(id string, ctx context.Context) (*entities.OffChainMealSupportDuration, error) {
	var query string = "SELECT * FROM " + meal_support_duration_table + " WHERE id = $1"
	var errLogMsg string = fmt.Sprintf(noti.REPO_ERR_MSG, shared.MEAL_SUPPORT_DURATION_REPOSITORY) + "GetMealSupportDuration - "

	var res entities.OffChainMealSupportDuration
	if err := m.db.QueryRowContext(ctx, query, id).Scan(
		&res.ID, &res.StartPeriod, &res.EndPeriod); err != nil {

		if err == sql.ErrNoRows {
			return nil, nil
		}

		m.errLogger.Println(errLogMsg + err.Error())
		return nil, errors.New(noti.INTERNALL_ERR_MSG)
	}

	return &res, nil
}
