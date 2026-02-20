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
	"raise-child/util"
	"time"
)

type leaderNotiRepo struct {
	db        *sql.DB
	errLogger *log.Logger
}

const leader_noti_table string = "leader_notis"

func InitializeLeaderNotiRepository(db *sql.DB, errLogger *log.Logger) repository.ILeaderNotiRepository {
	return &leaderNotiRepo{
		db:        db,
		errLogger: errLogger,
	}
}

// AssignLeader implements repository.ILeaderNotiRepository.
func (l *leaderNotiRepo) AssignLeader(leader string, region string, ctx context.Context) error {
	var errLogMsg string = fmt.Sprintf(noti.REPO_ERR_MSG, shared.LEADER_NOTI_REPOSITORY) + "AssignLeader - "
	var query string = "UPDATE " + leader_noti_table + " SET assgined_leaders = array_append(assgined_leaders, $1) WHERE region = $2"

	res, err := l.db.ExecContext(ctx, query, leader, region)
	var internalErr error = errors.New(noti.INTERNALL_ERR_MSG)

	if err != nil {
		l.errLogger.Println(errLogMsg + err.Error())
		return internalErr
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		l.errLogger.Println(errLogMsg + err.Error())
		return internalErr
	}

	if rowsAffected == 0 {
		return errors.New(fmt.Sprintf(noti.UNDEFINED_OBJECT_WARN_MSG, leader_noti_table))
	}

	return nil
}

// CreateNoti implements repository.ILeaderNotiRepository.
func (l *leaderNotiRepo) CreateNoti(notification entities.LeaderNoti, ctx context.Context) error {
	var query string = "INSERT INTO " + volunteer_noti_table +
		" (id, meal_need_id, child_id, region, assgined_leaders, " +
		"expected_withdraw_periods, content) " +
		"values ($1, $2, $3, $4, $5, $6, $7)"

	var errLogMsg string = fmt.Sprintf(noti.REPO_ERR_MSG, shared.LEADER_NOTI_REPOSITORY) + "CreateNoti - "

	if _, err := l.db.ExecContext(ctx, query, notification.ID, notification.MealNeedID, notification.ChildID, notification.Region, notification.AssignedLeaders,
		notification.ExpectedWithdrawPeriods, notification.Content); err != nil {

		l.errLogger.Println(errLogMsg + err.Error())
		return errors.New(noti.INTERNALL_ERR_MSG)
	}

	return nil
}

// GetCurrentLeaderNotis implements repository.ILeaderNotiRepository.
func (l *leaderNotiRepo) GetCurrentLeaderNotis(req request.GetNotisRequest, leader string, ctx context.Context) ([]entities.LeaderNoti, error) {
	var rawCurTime string = util.TimeToRawDate(time.Now())
	var query string = generateRetrieveQuery(generateRetrieveQueryRequest{
		table:       leader_noti_table,
		limitAmount: req.PageSize,
		condition:   fmt.Sprintf("%s = ANY(assgined_leaders) AND %s = ANY(expected_withdraw_periods) ORDER BY created_at DESC", leader, rawCurTime),
		page:        req.Page,
		isGetCount:  false,
	})

	var errLogMsg string = fmt.Sprintf(noti.REPO_ERR_MSG, shared.VOLUNTEER_NOTI_REPOSITORY) + "GetVolunteerNotis - "
	var internalErr error = errors.New(noti.INTERNALL_ERR_MSG)
	rows, err := l.db.QueryContext(ctx, query)
	if err != nil {
		l.errLogger.Println(errLogMsg + err.Error())
		return nil, internalErr
	}

	var res []entities.LeaderNoti
	for rows.Next() {
		var x entities.LeaderNoti
		if err := rows.Scan(&x.ID, &x.MealNeedID, &x.ChildID, &x.Region, &x.AssignedLeaders,
			&x.ExpectedWithdrawPeriods, &x.Content, &x.CreatedAt, &x.UpdatedAt); err != nil {
			l.errLogger.Println(errLogMsg + err.Error())
			return nil, internalErr
		}

		res = append(res, x)
	}

	return res, nil
}

// GetNoti implements repository.ILeaderNotiRepository.
func (l *leaderNotiRepo) GetNoti(id string, ctx context.Context) (*entities.LeaderNoti, error) {
	panic("unimplemented")
}

// GetNotiByMealNeed implements repository.ILeaderNotiRepository.
func (l *leaderNotiRepo) GetNotiByMealNeed(id string, ctx context.Context) (*entities.LeaderNoti, error) {
	var query string = "SELECT * FROM " + leader_noti_table + " WHERE meal_need_id = $1"
	var errLogMsg string = fmt.Sprintf(noti.REPO_ERR_MSG, shared.LEADER_NOTI_REPOSITORY) + "GetNotiByMealNeed - "

	var res entities.LeaderNoti
	if err := l.db.QueryRowContext(ctx, query, id).Scan(res.ID, res.MealNeedID, res.ChildID, res.Region, res.AssignedLeaders,
		res.ExpectedWithdrawPeriods, res.Content, res.CreatedAt, res.UpdatedAt); err != nil {

		if err == sql.ErrNoRows {
			return nil, nil
		}

		l.errLogger.Println(errLogMsg + err.Error())
		return nil, errors.New(noti.INTERNALL_ERR_MSG)
	}

	return &res, nil
}

// UpdateNoti implements repository.ILeaderNotiRepository.
func (l *leaderNotiRepo) UpdateNoti(notification entities.LeaderNoti, ctx context.Context) error {
	var errLogMsg string = fmt.Sprintf(noti.REPO_ERR_MSG, shared.LEADER_NOTI_REPOSITORY) + "AssignLeader - "
	var query string = "UPDATE " + leader_noti_table + " SET assgined_leaders = $1, expected_withdraw_periods = $2  WHERE id = $3"

	res, err := l.db.ExecContext(ctx, query, notification.AssignedLeaders, notification.ExpectedWithdrawPeriods)
	var internalErr error = errors.New(noti.INTERNALL_ERR_MSG)

	if err != nil {
		l.errLogger.Println(errLogMsg + err.Error())
		return internalErr
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		l.errLogger.Println(errLogMsg + err.Error())
		return internalErr
	}

	if rowsAffected == 0 {
		return errors.New(fmt.Sprintf(noti.UNDEFINED_OBJECT_WARN_MSG, leader_noti_table))
	}

	return nil
}
