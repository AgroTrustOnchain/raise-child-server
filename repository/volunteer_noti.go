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

type volunteerNotiRepo struct {
	db        *sql.DB
	errLogger *log.Logger
}

const volunteer_noti_table string = "volunteer_notis"

func InitializeVolunteerNotiRepository(db *sql.DB, errLogger *log.Logger) repository.IVolunteerNotiRepository {
	return &volunteerNotiRepo{
		db:        db,
		errLogger: errLogger,
	}
}

// AssignVolunteer implements repository.IVolunteerNotiRepository.
func (v *volunteerNotiRepo) AssignVolunteer(volunteer string, region string, ctx context.Context) error {
	var errLogMsg string = fmt.Sprintf(noti.REPO_ERR_MSG, shared.VOLUNTEER_NOTI_REPOSITORY) + "AssignVolunteer - "
	var query string = "UPDATE " + volunteer_noti_table + " SET assgined_volunteers = array_append(assgined_volunteers, $1) WHERE region = $2 AND start_period > NOW()"

	res, err := v.db.ExecContext(ctx, query, volunteer, region)
	var internalErr error = errors.New(noti.INTERNALL_ERR_MSG)

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
		return errors.New(fmt.Sprintf(noti.UNDEFINED_OBJECT_WARN_MSG, volunteer_noti_table))
	}

	return nil
}

// CreateNoti implements repository.IVolunteerNotiRepository.
func (v *volunteerNotiRepo) CreateNoti(notification entities.VolunteerNoti, ctx context.Context) error {
	var query string = "INSERT INTO " + volunteer_noti_table +
		" (id, child_id, region, assgined_volunteers, " +
		"content, start_period, end_period) " +
		"values ($1, $2, $3, $4, $5, $6, $7)"

	var errLogMsg string = fmt.Sprintf(noti.REPO_ERR_MSG, shared.VOLUNTEER_NOTI_REPOSITORY) + "CreateNoti - "

	if _, err := v.db.ExecContext(ctx, query, notification.ID, notification.ChildID, notification.Region, notification.AssginedVolunteers,
		notification.Content, notification.StartPeriod, notification.EndPeriod); err != nil {

		v.errLogger.Println(errLogMsg + err.Error())
		return errors.New(noti.INTERNALL_ERR_MSG)
	}

	return nil
}

// GetNoti implements repository.IVolunteerNotiRepository.
func (v *volunteerNotiRepo) GetNoti(id string, ctx context.Context) (*entities.VolunteerNoti, error) {
	panic("unimplemented")
}

// GetVolunteerNotis implements repository.IVolunteerNotiRepository.
func (v *volunteerNotiRepo) GetCurrentVolunteerNotis(req request.GetNotisRequest, volunteer string, ctx context.Context) ([]entities.VolunteerNoti, error) {
	var query string = generateRetrieveQuery(generateRetrieveQueryRequest{
		table:       volunteer_noti_table,
		limitAmount: req.PageSize,
		condition:   fmt.Sprintf("%s = ANY(assigned_volunteers) AND NOW >= start_period AND NOW() <= end_period ORDER BY created_at DESC", volunteer),
		page:        req.Page,
		isGetCount:  false,
	})

	var errLogMsg string = fmt.Sprintf(noti.REPO_ERR_MSG, shared.VOLUNTEER_NOTI_REPOSITORY) + "GetVolunteerNotis - "
	var internalErr error = errors.New(noti.INTERNALL_ERR_MSG)
	rows, err := v.db.QueryContext(ctx, query)
	if err != nil {
		v.errLogger.Println(errLogMsg + err.Error())
		return nil, internalErr
	}

	var res []entities.VolunteerNoti
	for rows.Next() {
		var x entities.VolunteerNoti
		if err := rows.Scan(&x.ID, &x.ChildID, &x.Region, &x.AssginedVolunteers,
			&x.Content, &x.StartPeriod, &x.EndPeriod, &x.CreatedAt, &x.UpdatedAt); err != nil {
			v.errLogger.Println(errLogMsg + err.Error())
			return nil, internalErr
		}

		res = append(res, x)
	}

	return res, nil
}
