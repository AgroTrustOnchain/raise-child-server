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

type pendingChildSpecialNeedProposalRepo struct {
	db        *sql.DB
	errLogger *log.Logger
}

const (
	pending_child_special_need_proposal_table string = "pending_child_special_need_proposals"
)

func InitializePendingChildSpecialNeedProposalRepo(db *sql.DB, errLogger *log.Logger) repository.IPendingChildSpecialNeedProposalRepository {
	return &pendingChildSpecialNeedProposalRepo{
		db:        db,
		errLogger: errLogger,
	}
}

// CreatePendingChildSpecialNeedProposal implements repository.IPendingChildSpecialNeedProposalRepository.
func (p *pendingChildSpecialNeedProposalRepo) CreatePendingChildSpecialNeedProposal(proposal entities.PendingChildSpecialNeedProposal, ctx context.Context) error {
	var query string = "INSERT INTO " + pending_child_special_need_proposal_table +
		" (id, child_id, region, actor_profile_id, actor_address, target, " +
		"description, proof_blob_id, ai_evaluation, created_at, updated_at) " +
		"values ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)"

	var errLogMsg string = fmt.Sprintf(noti.REPO_ERR_MSG, shared.PENDING_CHILD_SPECIAL_NEED_PROPOSAL_REPOSITORY) + "CreatePendingChildSpecialNeedProposal - "

	if _, err := p.db.ExecContext(ctx, query, proposal.ID, proposal.ChildID, proposal.Region, proposal.ActorProfileID, proposal.ActorAddress,
		proposal.Target, proposal.Description, proposal.ProofBlobID, proposal.AIEvaluation, proposal.CreatedAt, proposal.UpdatedAt,
	); err != nil {

		p.errLogger.Println(errLogMsg + err.Error())
		return errors.New(noti.INTERNALL_ERR_MSG)
	}

	return nil
}

// GetPendingChildSpecialNeedProposal implements repository.IPendingChildSpecialNeedProposalRepository.
func (p *pendingChildSpecialNeedProposalRepo) GetPendingChildSpecialNeedProposal(id string, ctx context.Context) (*entities.PendingChildSpecialNeedProposal, error) {
	var query string = "SELECT * FROM " + pending_child_special_need_proposal_table + " WHERE id = $1"
	var errLogMsg string = fmt.Sprintf(noti.REPO_ERR_MSG, shared.PENDING_CHILD_SPECIAL_NEED_PROPOSAL_REPOSITORY) + "GetPendingChildSpecialNeedProposal - "

	var res entities.PendingChildSpecialNeedProposal
	if err := p.db.QueryRowContext(ctx, query, id).Scan(
		&res.ID, &res.ChildID, &res.Region, &res.ActorProfileID, &res.ActorAddress,
		&res.Target, &res.Description, &res.ProofBlobID, &res.AIEvaluation,
		&res.ReviewStatus, &res.ReviewedBy, &res.CreatedAt, &res.UpdatedAt,
	); err != nil {

		if err == sql.ErrNoRows {
			return nil, nil
		}

		p.errLogger.Println(errLogMsg + err.Error())
		return nil, errors.New(noti.INTERNALL_ERR_MSG)
	}

	return &res, nil
}

// GetPendingChildSpecialNeedProposals implements repository.IPendingChildSpecialNeedProposalRepository.
func (p *pendingChildSpecialNeedProposalRepo) GetPendingChildSpecialNeedProposals(req request.GetPendingChildSpecialNeedProposalsRequest, ctx context.Context) ([]entities.PendingChildSpecialNeedProposal, int, error) {
	var errLogMsg string = fmt.Sprintf(noti.REPO_ERR_MSG, shared.PENDING_CHILD_SPECIAL_NEED_PROPOSAL_REPOSITORY) + "GetPendingChildSpecialNeedProposals - "
	var internalErr error = errors.New(noti.INTERNALL_ERR_MSG)

	var queryCondition string
	var isHavePreviosCondition bool = false
	if req.Keyword != "" {
		queryCondition += fmt.Sprintf("LOWER(description) LIKE LOWER('%%%s%%')", req.Keyword)
		isHavePreviosCondition = true
	}

	if req.Region != "" {
		if isHavePreviosCondition {
			queryCondition += " AND "
		}

		queryCondition += fmt.Sprintf("LOWER(region) = LOWER('%%%s%%')", req.Region)
		isHavePreviosCondition = true
	}

	if req.Status != "" {
		if isHavePreviosCondition {
			queryCondition += " AND "
		}

		queryCondition += fmt.Sprintf("LOWER(review_status) = LOWER('%s')", req.Status)
		isHavePreviosCondition = true
	}

	if req.Creator != "" {
		if isHavePreviosCondition {
			queryCondition += " AND "
		}

		queryCondition += fmt.Sprintf("actor_address = '%s'", req.Creator)
		isHavePreviosCondition = true
	}

	if req.Reviewer != "" {
		if isHavePreviosCondition {
			queryCondition += " AND "
		}

		queryCondition += fmt.Sprintf("reviewed_bdy = '%s'", req.Creator)
		isHavePreviosCondition = true
	}

	if req.MinAmount != nil {
		if isHavePreviosCondition {
			queryCondition += " AND "
		}

		queryCondition += fmt.Sprintf("target >= %d", *req.MinAmount)
		isHavePreviosCondition = true
	}

	if req.MaxAmount != nil {
		if isHavePreviosCondition {
			queryCondition += " AND "
		}

		queryCondition += fmt.Sprintf("target <= %d", *req.MaxAmount)
		isHavePreviosCondition = true
	}

	if req.SortCriteria == "" {
		req.SortCriteria = "created_at"
	}

	if req.SortOrder == "" {
		req.SortOrder = "DESC"
	}

	if isHavePreviosCondition {
		queryCondition += " "
	}

	var query string = generateRetrieveQuery(generateRetrieveQueryRequest{
		table:       pending_child_special_need_proposal_table,
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

	var res []entities.PendingChildSpecialNeedProposal
	for rows.Next() {
		var x entities.PendingChildSpecialNeedProposal

		if err := rows.Scan(
			&x.ID, &x.ChildID, &x.Region, &x.ActorProfileID, &x.ActorAddress,
			&x.Target, &x.Description, &x.ProofBlobID, &x.AIEvaluation,
			&x.ReviewStatus, &x.ReviewedBy, &x.CreatedAt, &x.UpdatedAt,
		); err != nil {
			p.errLogger.Println(errLogMsg + err.Error())
			return nil, 0, internalErr
		}

		res = append(res, x)
	}

	var totalRecords int
	p.db.QueryRowContext(ctx, generateCountTotalRecordsQuery(pending_child_special_need_proposal_table, queryCondition)).Scan(&totalRecords)

	return res, caculateTotalPages(totalRecords, req.PageSize), nil
}

// UpdatePendingChildSpecialNeedProposal implements repository.IPendingChildSpecialNeedProposalRepository.
func (p *pendingChildSpecialNeedProposalRepo) UpdatePendingChildSpecialNeedProposal(proposal entities.PendingChildSpecialNeedProposal, ctx context.Context) error {
	var query string = "UPDATE " + pending_child_special_need_proposal_table + " SET " +
		"child_id = $1, target = $2, description = $3, proof_blob_id = $4, " +
		"review_status = $5, reviewed_by = $6 WHERE id = $7"

	var errLogMsg string = fmt.Sprintf(noti.REPO_ERR_MSG, shared.PENDING_CHILD_SPECIAL_NEED_PROPOSAL_REPOSITORY) + "UpdatePendingChildSpecialNeedProposal - "
	var internalErr error = errors.New(noti.INTERNALL_ERR_MSG)

	res, err := p.db.ExecContext(ctx, query, proposal.ChildID, proposal.Target, proposal.Description,
		proposal.ProofBlobID, proposal.ReviewStatus, proposal.ReviewedBy, proposal.ID)
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
		return fmt.Errorf(noti.UNDEFINED_OBJECT_WARN_MSG, pending_child_special_need_proposal_table)
	}

	return nil
}
