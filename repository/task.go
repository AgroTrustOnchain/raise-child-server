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

type taskRepo struct {
	db        *sql.DB
	errLogger *log.Logger
}

const task_table string = "tasks"

func InitializeTaskRepository(db *sql.DB, errLogger *log.Logger) repository.ITaskRepository {
	return &taskRepo{
		db:        db,
		errLogger: errLogger,
	}
}

// CreateTask implements repository.ITaskRepository.
func (t *taskRepo) CreateTask(task entities.Task, ctx context.Context) error {
	var query string = "INSERT INTO " + task_table +
		" (id, is_Child_task, child_task_id, created_by, region, description, start_period, end_period, created_at, updated_at) " +
		"values ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)"

	var errLogMsg string = fmt.Sprintf(noti.REPO_ERR_MSG, shared.TASK_REPOSITORY) + "CreateTask - "

	if _, err := t.db.ExecContext(ctx, query, task.ID, task.IsChildTask, task.ChildTaskDetailID, task.CreatedBy,
		task.Region, task.Description, task.StartPeriod, task.EndPeriod, task.CreatedAt, task.UpdatedAt); err != nil {

		t.errLogger.Println(errLogMsg + err.Error())
		return errors.New(noti.INTERNALL_ERR_MSG)
	}

	return nil
}

// GetTask implements repository.ITaskRepository.
func (t *taskRepo) GetTask(id string, ctx context.Context) (*entities.Task, error) {
	var query string = "SELECT * FROM " + task_table + " WHERE id = $1"
	var errLogMsg string = fmt.Sprintf(noti.REPO_ERR_MSG, shared.TASK_REPOSITORY) + "GetTask - "

	var res entities.Task
	if err := t.db.QueryRowContext(ctx, query, id).Scan(
		&res.ID, &res.IsChildTask, &res.ChildTaskDetailID, &res.CreatedBy, &res.AssignedProfileID,
		&res.AssignedStaff, &res.Region, &res.Description, &res.StartPeriod,
		&res.EndPeriod, &res.CreatedAt, &res.UpdatedAt); err != nil {

		if err == sql.ErrNoRows {
			return nil, nil
		}

		t.errLogger.Println(errLogMsg + err.Error())
		return nil, errors.New(noti.INTERNALL_ERR_MSG)
	}

	return &res, nil
}

// GetTasks implements repository.ITaskRepository.
func (t *taskRepo) GetTasks(req request.GetTasksRequest, ctx context.Context) ([]entities.Task, int, error) {
	var errLogMsg string = fmt.Sprintf(noti.REPO_ERR_MSG, shared.TASK_REPOSITORY) + "GetTasks - "
	var internalErr error = errors.New(noti.INTERNALL_ERR_MSG)

	var queryCondition string
	var isHavePreviosCondition bool = false
	if req.Keyword != "" {
		queryCondition += fmt.Sprintf("LOWER(description) LIKE LOWER('%s')", req.Keyword)
		isHavePreviosCondition = true
	}

	if req.Region != "" {
		if isHavePreviosCondition {
			queryCondition += " AND "
		}

		queryCondition += fmt.Sprintf("LOWER(region) = LOWER('%s')", req.Region)
		isHavePreviosCondition = true
	}

	if req.AssignedStaff != "" {
		if isHavePreviosCondition {
			queryCondition += " AND "
		}

		queryCondition += fmt.Sprintf("assgined_staff = '%s'", req.AssignedStaff)
		isHavePreviosCondition = true
	}

	var order string = "DESC"
	if req.SortOrder != "" {
		order = req.SortOrder
	}

	var query string = generateRetrieveQuery(generateRetrieveQueryRequest{
		table:       task_table,
		limitAmount: req.PageSize,
		condition:   queryCondition,
		order:       " ORDER BY created_at " + order,
		page:        req.Page,
		isGetCount:  false,
	})

	rows, err := t.db.QueryContext(ctx, query)
	if err != nil {
		t.errLogger.Println("Query get tasks:", query)
		t.errLogger.Println(errLogMsg + err.Error())
		return nil, 0, internalErr
	}
	defer rows.Close()

	var res []entities.Task
	for rows.Next() {
		var x entities.Task

		if err := rows.Scan(
			&x.ID, &x.IsChildTask, &x.ChildTaskDetailID, &x.CreatedBy, &x.AssignedProfileID,
			&x.AssignedStaff, &x.Region, &x.Description,
			&x.StartPeriod, &x.EndPeriod, &x.CreatedAt, &x.UpdatedAt); err != nil {

			t.errLogger.Println(errLogMsg + err.Error())
			return nil, 0, internalErr
		}

		res = append(res, x)
	}

	var totalRecords int
	t.db.QueryRowContext(ctx, generateCountTotalRecordsQuery(task_table, queryCondition)).Scan(&totalRecords)

	return res, caculateTotalPages(totalRecords, req.PageSize), nil
}

// GetTasksOfUser implements repository.ITaskRepository.
func (t *taskRepo) GetTasksOfUser(req request.GetTasksRequest, ctx context.Context) ([]entities.TaskV2, int, error) {
	var rawCurTime string = util.TimeToRawDate(time.Now())
	var query string = `SELECT t.*, 
							EXISTS (
								SELECT 1 FROM task_proofs tp
								WHERE tp.task_id = t.id 
								AND tp.raw_submit_date = $1
								AND tp.actor_address = t.assgined_staff
								AND tp.review_status IN ('Pending', 'Approved')
							) AS is_submitted
						FROM tasks t
						WHERE t.assgined_staff = $2`

	var keywordCond string
	if req.Keyword != "" {
		keywordCond = fmt.Sprintf(" AND LOWER(t.description) LIKE LOWER('%s')", req.Keyword)
	}

	query += keywordCond

	var order string = "DESC"
	if req.SortOrder != "" {
		order = req.SortOrder
	}

	var offSet int = (req.Page - 1) * req.PageSize
	var paginationFilterCond string = fmt.Sprintf(" ORDER BY t.created_at %s LIMIT %d OFFSET %d", order, req.PageSize, offSet)
	query += paginationFilterCond

	var errLogMsg string = fmt.Sprintf(noti.REPO_ERR_MSG, shared.TASK_REPOSITORY) + "GetTasks - "
	var internalErr error = errors.New(noti.INTERNALL_ERR_MSG)
	rows, err := t.db.QueryContext(ctx, query, rawCurTime, req.AssignedStaff)
	if err != nil {
		t.errLogger.Println(errLogMsg + err.Error())
		return nil, 0, internalErr
	}
	defer rows.Close()

	var res []entities.TaskV2
	for rows.Next() {
		var x entities.TaskV2

		if err := rows.Scan(
			&x.ID, &x.IsChildTask, &x.ChildTaskDetailID, &x.CreatedBy, &x.AssignedProfileID,
			&x.AssignedStaff, &x.Region, &x.Description, &x.StartPeriod,
			&x.EndPeriod, &x.CreatedAt, &x.UpdatedAt, &x.IsSubmitted); err != nil {

			t.errLogger.Println(errLogMsg + err.Error())
			return nil, 0, internalErr
		}

		res = append(res, x)
	}

	var totalRecords int
	t.db.QueryRowContext(ctx, generateCountTotalRecordsQuery(task_table, fmt.Sprintf("assgined_staff = '%s' %s", req.AssignedStaff, keywordCond))).Scan(&totalRecords)

	return res, caculateTotalPages(totalRecords, req.PageSize), nil
}

// UpdateTask implements repository.ITaskRepository.
func (t *taskRepo) UpdateTask(task entities.Task, ctx context.Context) error {
	var query string = "UPDATE " + task_table + " SET " +
		"assigned_profile_id = $1, assgined_staff = $2, description = $3 WHERE id = $4"

	var errLogMsg string = fmt.Sprintf(noti.REPO_ERR_MSG, shared.TASK_REPOSITORY) + "UpdateTask - "
	var internalErr error = errors.New(noti.INTERNALL_ERR_MSG)

	res, err := t.db.ExecContext(ctx, query, task.AssignedProfileID, task.AssignedStaff, task.Description, task.ID)
	if err != nil {
		t.errLogger.Println(errLogMsg + err.Error())
		return internalErr
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		t.errLogger.Println(errLogMsg + err.Error())
		return internalErr
	}

	if rowsAffected == 0 {
		return fmt.Errorf(noti.UNDEFINED_OBJECT_WARN_MSG, task_table)
	}

	return nil
}
