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

type childtaskDetailRepo struct {
	db        *sql.DB
	errLogger *log.Logger
}

const child_task_detail_table string = "child_task_details"

func InitializeChildTaskDetailRepository(db *sql.DB, errLogger *log.Logger) repository.IChildTaskDetailRepository {
	return &childtaskDetailRepo{
		db:        db,
		errLogger: errLogger,
	}
}

// CreateChildTaskDetail implements repository.IChildTaskDetailRepository.
func (c *childtaskDetailRepo) CreateChildTaskDetail(detail entities.ChildTaskDetail, ctx context.Context) error {
	var query string = "INSERT INTO " + child_task_detail_table +
		" (id, child_id, purpose, target)  values ($1, $2, $3, $4)"

	var errLogMsg string = fmt.Sprintf(noti.REPO_ERR_MSG, shared.CHILD_TASK_DETAIL_REPOSITORY) + "CreateChildTaskDetail - "

	if _, err := c.db.ExecContext(ctx, query, detail.ID, detail.ChildID, detail.Purpose, detail.Target); err != nil {

		c.errLogger.Println(errLogMsg + err.Error())
		return errors.New(noti.INTERNALL_ERR_MSG)
	}

	return nil
}

// GetChildTaskDetail implements repository.IChildTaskDetailRepository.
func (c *childtaskDetailRepo) GetChildTaskDetail(id string, ctx context.Context) (*entities.ChildTaskDetail, error) {
	var query string = "SELECT id, child_id, purpose, target FROM " + child_task_detail_table + " WHERE id = $1"
	var errLogMsg string = fmt.Sprintf(noti.REPO_ERR_MSG, shared.CHILD_TASK_DETAIL_REPOSITORY) + "GetChildTaskDetail - "

	var res entities.ChildTaskDetail
	if err := c.db.QueryRowContext(ctx, query, id).Scan(
		&res.ID, &res.ChildID, &res.Purpose, &res.Target); err != nil {

		if err == sql.ErrNoRows {
			return nil, nil
		}

		c.errLogger.Println(errLogMsg + err.Error())
		return nil, errors.New(noti.INTERNALL_ERR_MSG)
	}

	return &res, nil
}
