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

type pendingWithdrawProposalRepo struct {
	db        *sql.DB
	errLogger *log.Logger
}

const (
	pending_withdraw_proposal_table string = "pending_withdraw_proposals"
)

func InitializePendingWithdrawProposalRepo(db *sql.DB, errLogger *log.Logger) repository.IPendingWithdrawProposalRepository {
	return &pendingWithdrawProposalRepo{
		db:        db,
		errLogger: errLogger,
	}
}

// CreatePendingWithdrawProposal implements repository.IPendingWithdrawProposalRepository.
func (p *pendingWithdrawProposalRepo) CreatePendingWithdrawProposal(proposal entities.PendingWithdrawProposal, ctx context.Context) error {
	var query string = "INSERT INTO " + pending_withdraw_proposal_table +
		" (id, profile_id, creator, pool_id, pool_name, " +
		"purpose, target, withdraw_amount, proof_blob_id, description, " +
		"status, ai_evaluation, created_at, updated_at) " +
		"values ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14)"

	var errLogMsg string = fmt.Sprintf(noti.REPO_ERR_MSG, shared.PENDING_WITHDRAW_PROPOSAL_REPOSITORY) + "CreatePendingWithdrawProposal - "
	if _, err := p.db.ExecContext(ctx, query, proposal.ID, proposal.ProfileID, proposal.Creator, proposal.PoolID, proposal.PoolName,
		proposal.Purpose, proposal.Target, proposal.WithdrawAmount, proposal.ProofBlobID, proposal.Description,
		proposal.Status, proposal.AIEvaluation, proposal.CreatedAt, proposal.UpdatedAt); err != nil {

		p.errLogger.Println(errLogMsg + err.Error())
		return errors.New(noti.INTERNALL_ERR_MSG)
	}

	return nil
}

// GetPendingWithdrawProposal implements repository.IPendingWithdrawProposalRepository.
func (p *pendingWithdrawProposalRepo) GetPendingWithdrawProposal(id string, ctx context.Context) (*entities.PendingWithdrawProposal, error) {
	var query string = "SELECT * FROM " + pending_withdraw_proposal_table + " WHERE id = $1"
	var errLogMsg string = fmt.Sprintf(noti.REPO_ERR_MSG, shared.PENDING_WITHDRAW_PROPOSAL_REPOSITORY) + "GetPendingWithdrawProposal - "

	var res entities.PendingWithdrawProposal
	if err := p.db.QueryRowContext(ctx, query, id).Scan(
		&res.ID, &res.ProfileID, &res.Creator, &res.PoolID, &res.PoolName,
		&res.Purpose, &res.Target, &res.WithdrawAmount, &res.ProofBlobID, &res.Description,
		&res.Status, &res.AIEvaluation, &res.ReviewedBy, &res.CreatedAt, &res.UpdatedAt,
	); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}

		p.errLogger.Println(errLogMsg + err.Error())
		return nil, errors.New(noti.INTERNALL_ERR_MSG)
	}

	return &res, nil
}

// GetPendingWithdrawProposals implements repository.IPendingWithdrawProposalRepository.
func (p *pendingWithdrawProposalRepo) GetPendingWithdrawProposals(req request.GetPendingWithdrawProposalsRequest, ctx context.Context) ([]entities.PendingWithdrawProposal, int, error) {
	var errLogMsg string = fmt.Sprintf(noti.REPO_ERR_MSG, shared.PENDING_WITHDRAW_PROPOSAL_REPOSITORY) + "GetPendingWithdrawProposals - "
	var internalErr error = errors.New(noti.INTERNALL_ERR_MSG)

	var queryCondition string
	var isHavePreviosCondition bool = false
	if req.Keyword != "" {
		queryCondition += fmt.Sprintf("(LOWER(pool_name) LIKE LOWER('%%%s%%') OR LOWER(description) LIKE LOWER('%%%s%%'))", req.Keyword, req.Keyword)
		isHavePreviosCondition = true
	}

	if req.Creator != "" {
		if isHavePreviosCondition {
			queryCondition += " AND "
		}

		queryCondition += fmt.Sprintf("creator = '%s'", req.Creator)
		isHavePreviosCondition = true
	}

	if req.Reviewer != "" {
		if isHavePreviosCondition {
			queryCondition += " AND "
		}

		queryCondition += fmt.Sprintf("reviewed_by = '%s'", req.Reviewer)
		isHavePreviosCondition = true
	}

	if req.MinAmount != nil {
		if isHavePreviosCondition {
			queryCondition += " AND "
		}

		queryCondition += fmt.Sprintf("withdraw_amount >= %d", *req.MinAmount)
		isHavePreviosCondition = true
	}

	if req.MaxAmount != nil {
		if isHavePreviosCondition {
			queryCondition += " AND "
		}

		queryCondition += fmt.Sprintf("withdraw_amount <= %d", *req.MaxAmount)
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

	if req.SortCriteria == "" {
		req.SortCriteria = "created_at"
	}

	if req.SortOrder == "" {
		req.SortOrder = "DESC"
	}

	var query string = generateRetrieveQuery(generateRetrieveQueryRequest{
		table:       pending_withdraw_proposal_table,
		limitAmount: req.PageSize,
		condition:   queryCondition,
		order:       fmt.Sprintf(" ORDER BY %s %s", req.SortCriteria, req.SortOrder),
		page:        req.Page,
		isGetCount:  false,
	})

	rows, err := p.db.QueryContext(ctx, query)
	if err != nil {
		p.errLogger.Println(errLogMsg + err.Error())
		return nil, 0, internalErr
	}

	var res []entities.PendingWithdrawProposal
	for rows.Next() {
		var x entities.PendingWithdrawProposal
		if err := rows.Scan(
			&x.ID, &x.ProfileID, &x.Creator, &x.PoolID, &x.PoolName,
			&x.Purpose, &x.Target, &x.WithdrawAmount, &x.ProofBlobID, &x.Description,
			&x.Status, &x.AIEvaluation, &x.ReviewedBy, &x.CreatedAt, &x.UpdatedAt,
		); err != nil {

			p.errLogger.Println(errLogMsg + err.Error())
			return nil, 0, internalErr
		}

		res = append(res, x)
	}

	var totalRecords int
	p.db.QueryRowContext(ctx, generateCountTotalRecordsQuery(pending_withdraw_proposal_table, queryCondition)).Scan(&totalRecords)

	return res, caculateTotalPages(totalRecords, req.PageSize), nil
}

// UpdatePendingWithdrawProposal implements repository.IPendingWithdrawProposalRepository.
func (p *pendingWithdrawProposalRepo) UpdatePendingWithdrawProposal(proposal entities.PendingWithdrawProposal, ctx context.Context) error {
	var query string = "UPDATE " + pending_withdraw_proposal_table + " SET " +
		"withdraw_amount = $1, proof_blob_id = $2, description = $3, status = $4, " +
		"reviewed_by = $5 WHERE id = $6"

	var errLogMsg string = fmt.Sprintf(noti.REPO_ERR_MSG, shared.PENDING_WITHDRAW_PROPOSAL_REPOSITORY) + "UpdatePendingWithdrawProposal - "
	var internalErr error = errors.New(noti.INTERNALL_ERR_MSG)

	res, err := p.db.ExecContext(ctx, query, proposal.WithdrawAmount, proposal.ProofBlobID, proposal.Description, proposal.Status, proposal.ReviewedBy, proposal.ID)
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
		return fmt.Errorf(noti.UNDEFINED_OBJECT_WARN_MSG, pending_withdraw_proposal_table)
	}

	return nil
}

// IsPendingWithdrawProposalProposedWithSpecificInfo implements repository.IPendingWithdrawProposalRepository.
func (p *pendingWithdrawProposalRepo) IsPendingWithdrawProposalProposedWithSpecificInfo(purpose string, target string, description string, withdrawAmount int64, ctx context.Context) (bool, error) {
	var query string = "SELECT id FROM " + leader_noti_table + " WHERE purpose = $1 AND target = $2 AND description = $3 AND withdraw_amount = $4 AND (status = 'Pending' OR status = 'Approved') LIMIT 1"

	var id string
	if err := p.db.QueryRowContext(ctx, query, purpose, target, description, withdrawAmount).Scan(&id); err != nil {
		if err == sql.ErrNoRows {
			return false, nil
		}

		p.errLogger.Println(fmt.Sprintf(noti.REPO_ERR_MSG, shared.PENDING_WITHDRAW_PROPOSAL_REPOSITORY) + "IsPendingWithdrawProposalProposedWithSpecificInfo - " + err.Error())
		return false, errors.New(noti.INTERNALL_ERR_MSG)
	}

	return id != "", nil
}
