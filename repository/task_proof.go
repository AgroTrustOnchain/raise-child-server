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

type taskProofRepo struct {
	db        *sql.DB
	errLogger *log.Logger
}

const task_proof_table string = "task_proofs"

func InitializeTaskProofRepository(db *sql.DB, errLogger *log.Logger) repository.ITaskProofRepository {
	return &taskProofRepo{
		db:        db,
		errLogger: errLogger,
	}
}

// CreateTaskProof implements repository.ITaskProofRepository.
func (t *taskProofRepo) CreateTaskProof(proof entities.TaskProof, ctx context.Context) error {
	var query string = "INSERT INTO " + task_proof_table +
		" (id, task_id, description, actor_profile_id, actor_address, " +
		"image_blob_id, ai_evaluation, raw_submit_date, created_at, updated_at) " +
		"values ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)"

	var errLogMsg string = fmt.Sprintf(noti.REPO_ERR_MSG, shared.TASK_PROOF_REPOSITORY) + "CreateTaskProof - "

	if _, err := t.db.ExecContext(ctx, query, proof.ID, proof.TaskID, proof.Description, proof.ActorProfileID, proof.ActorAddress,
		proof.ImageBlobID, proof.AIEvaluation, proof.RawSubmitDate, proof.CreatedAt, proof.UpdatedAt); err != nil {

		t.errLogger.Println(errLogMsg + err.Error())
		return errors.New(noti.INTERNALL_ERR_MSG)
	}

	return nil
}

// GetTaskProof implements repository.ITaskProofRepository.
func (t *taskProofRepo) GetTaskProof(id string, ctx context.Context) (*entities.TaskProof, error) {
	var query string = "SELECT * FROM " + task_proof_table + " WHERE id = $1"
	var errLogMsg string = fmt.Sprintf(noti.REPO_ERR_MSG, shared.TASK_PROOF_REPOSITORY) + "GetTaskProof - "

	var res entities.TaskProof
	if err := t.db.QueryRowContext(ctx, query, id).Scan(
		&res.ID, &res.TaskID, &res.Description, &res.ActorProfileID, &res.ActorAddress,
		&res.ImageBlobID, &res.ReviewedBy, &res.AIEvaluation, &res.ReviewStatus,
		&res.RawSubmitDate, &res.CreatedAt, &res.UpdatedAt); err != nil {

		if err == sql.ErrNoRows {
			return nil, nil
		}

		t.errLogger.Println(errLogMsg + err.Error())
		return nil, errors.New(noti.INTERNALL_ERR_MSG)
	}

	return &res, nil
}

// GetTaskProofs implements repository.ITaskProofRepository.
func (t *taskProofRepo) GetTaskProofs(req request.GetTaskProofsRequest, ctx context.Context) ([]entities.TaskProof, int, error) {
	var errLogMsg string = fmt.Sprintf(noti.REPO_ERR_MSG, shared.TASK_PROOF_REPOSITORY) + "GetTaskProofs - "
	var internalErr error = errors.New(noti.INTERNALL_ERR_MSG)

	var queryCondition string
	var isHavePreviosCondition bool = false
	if req.Keyword != "" {
		queryCondition += fmt.Sprintf("LOWER(description) LIKE LOWER('%%%%%s%%%%'))", req.Keyword)
		isHavePreviosCondition = true
	}

	if req.Status != "" {
		if isHavePreviosCondition {
			queryCondition += " AND "
		}

		queryCondition += fmt.Sprintf("LOWER(status) = LOWER('%%%%%s%%%%')", req.Status)
		isHavePreviosCondition = true
	}

	if req.ActorAddress != "" {
		if isHavePreviosCondition {
			queryCondition += " AND "
		}

		queryCondition += fmt.Sprintf("actor_address = '%%%%%s%%%%'", req.ActorAddress)
		isHavePreviosCondition = true
	}

	if req.ReviewedBy != "" {
		if isHavePreviosCondition {
			queryCondition += " AND "
		}

		queryCondition += fmt.Sprintf("reviewed_by = '%%%%%s%%%%'", req.ReviewedBy)
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
		t.errLogger.Println(errLogMsg + err.Error())
		return nil, 0, internalErr
	}

	var res []entities.TaskProof
	for rows.Next() {
		var x entities.TaskProof

		if err := rows.Scan(
			&x.ID, &x.TaskID, &x.Description, &x.ActorProfileID, &x.ActorAddress,
			&x.ImageBlobID, &x.ReviewedBy, &x.AIEvaluation, &x.ReviewStatus,
			&x.RawSubmitDate, &x.CreatedAt, &x.UpdatedAt); err != nil {

			t.errLogger.Println(errLogMsg + err.Error())
			return nil, 0, internalErr
		}

		res = append(res, x)
	}

	var totalRecords int
	t.db.QueryRowContext(ctx, generateCountTotalRecordsQuery(task_proof_table, queryCondition)).Scan(&totalRecords)

	return res, caculateTotalPages(totalRecords, req.PageSize), nil
}

// IsTaskProofSumittedWithDetail implements repository.ITaskProofRepository.
func (t *taskProofRepo) IsTaskProofSumittedWithDetail(taskId string, description string, actorAddress string, rawSubmitDate string, ctx context.Context) (bool, error) {
	var query string = "SELECT id FROM " + task_proof_table + " WHERE task_id = $1 AND description = $2 AND actor_address = $3 AND raw_submit_date = $4 AND (review_status = 'Pending' OR review_status = 'Approved') LIMIT 1"

	var id string
	if err := t.db.QueryRowContext(ctx, query, taskId, description, actorAddress, rawSubmitDate).Scan(&id); err != nil {
		if err == sql.ErrNoRows {
			return false, nil
		}

		t.errLogger.Println(fmt.Sprintf(noti.REPO_ERR_MSG, shared.TASK_PROOF_REPOSITORY) + "IsTaskProofSumittedWithDetail - " + err.Error())
		return false, errors.New(noti.INTERNALL_ERR_MSG)
	}

	return id != "", nil
}

// UpdateTaskProof implements repository.ITaskProofRepository.
func (t *taskProofRepo) UpdateTaskProof(proof entities.TaskProof, ctx context.Context) error {
	var query string = "UPDATE " + task_proof_table + " SET " +
		"reviewed_by = $1, review_status = $2 WHERE id = $3"

	var errLogMsg string = fmt.Sprintf(noti.REPO_ERR_MSG, shared.TASK_PROOF_REPOSITORY) + "UpdateTaskProof - "
	var internalErr error = errors.New(noti.INTERNALL_ERR_MSG)

	res, err := t.db.ExecContext(ctx, query, proof.ReviewedBy, proof.ReviewStatus, proof.ID)
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
		return errors.New(fmt.Sprintf(noti.UNDEFINED_OBJECT_WARN_MSG, task_proof_table))
	}

	return nil
}
