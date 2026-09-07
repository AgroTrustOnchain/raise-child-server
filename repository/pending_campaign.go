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

type pendingCampaignRepo struct {
	db        *sql.DB
	errLogger *log.Logger
}

const (
	pending_campaign_table string = "pending_campaigns"
)

func InitializePendingCampaignRepo(db *sql.DB, errLogger *log.Logger) repository.IPendingCampaignRepository {
	return &pendingCampaignRepo{
		db:        db,
		errLogger: errLogger,
	}
}

// CreatePendingCampaign implements repository.IPendingCampaignRepository.
func (p *pendingCampaignRepo) CreatePendingCampaign(campaign entities.PendingCampaign, ctx context.Context) error {
	var query string = "INSERT INTO " + pending_campaign_table +
		" (id, actor_profile_id, actor_address, pool_id, pool_name, target, " +
		"description, proof_blob_id, ai_evaluation, created_at, updated_at) " +
		"values ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)"

	var errLogMsg string = fmt.Sprintf(noti.REPO_ERR_MSG, shared.PENDING_CAMPAIGN_REPOSITORY) + "CreatePendingCampaign - "

	if _, err := p.db.ExecContext(ctx, query, campaign.ID, campaign.ActorProfileID, campaign.ActorAddress, campaign.PoolID, campaign.PoolName,
		campaign.Target, campaign.Description, campaign.ProofBlobID, campaign.AIEvaluation, campaign.CreatedAt, campaign.UpdatedAt,
	); err != nil {

		p.errLogger.Println(errLogMsg + err.Error())
		return errors.New(noti.INTERNALL_ERR_MSG)
	}

	return nil
}

// GetPendingCampaign implements repository.IPendingCampaignRepository.
func (p *pendingCampaignRepo) GetPendingCampaign(id string, ctx context.Context) (*entities.PendingCampaign, error) {
	var query string = "SELECT * FROM " + pending_campaign_table + " WHERE id = $1"
	var errLogMsg string = fmt.Sprintf(noti.REPO_ERR_MSG, shared.PENDING_CAMPAIGN_REPOSITORY) + "GetPendingCampaign - "

	var res entities.PendingCampaign
	if err := p.db.QueryRowContext(ctx, query, id).Scan(
		&res.ID, &res.ActorProfileID, &res.ActorAddress, &res.PoolID, &res.PoolName,
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

// GetPendingCampaigns implements repository.IPendingCampaignRepository.
func (p *pendingCampaignRepo) GetPendingCampaigns(req request.GetPendingCampaignsRequest, ctx context.Context) ([]entities.PendingCampaign, int, error) {
	var errLogMsg string = fmt.Sprintf(noti.REPO_ERR_MSG, shared.PENDING_CAMPAIGN_REPOSITORY) + "GetPendingCampaigns - "
	var internalErr error = errors.New(noti.INTERNALL_ERR_MSG)

	var queryCondition string
	var isHavePreviosCondition bool = false
	if req.Keyword != "" {
		queryCondition += fmt.Sprintf("LOWER(description) LIKE LOWER('%%%s%%')", req.Keyword)
		isHavePreviosCondition = true
	}

	if req.PoolName != "" {
		if isHavePreviosCondition {
			queryCondition += " AND "
		}

		queryCondition += fmt.Sprintf("LOWER(pool_name) = LOWER('%%%s%%')", req.PoolName)
		isHavePreviosCondition = true
	}

	if req.Status != "" {
		if isHavePreviosCondition {
			queryCondition += " AND "
		}

		queryCondition += fmt.Sprintf("LOWER(review_status) = LOWER('%%%s%%')", req.Status)
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
		table:       pending_campaign_table,
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

	var res []entities.PendingCampaign
	for rows.Next() {
		var x entities.PendingCampaign

		if err := rows.Scan(
			&x.ID, &x.ActorProfileID, &x.ActorAddress, &x.PoolID, &x.PoolName,
			&x.Target, &x.Description, &x.ProofBlobID, &x.AIEvaluation,
			&x.ReviewStatus, &x.ReviewedBy, &x.CreatedAt, &x.UpdatedAt,
		); err != nil {
			p.errLogger.Println(errLogMsg + err.Error())
			return nil, 0, internalErr
		}

		res = append(res, x)
	}

	var totalRecords int
	p.db.QueryRowContext(ctx, generateCountTotalRecordsQuery(pending_campaign_table, queryCondition)).Scan(&totalRecords)

	return res, caculateTotalPages(totalRecords, req.PageSize), nil
}

// UpdatePendingCampaign implements repository.IPendingCampaignRepository.
func (p *pendingCampaignRepo) UpdatePendingCampaign(campaign entities.PendingCampaign, ctx context.Context) error {
	var query string = "UPDATE " + pending_campaign_table + " SET " +
		"target = $1, description = $2, proof_blob_id = $3, " +
		"review_status = $4, reviewed_by = $5, updated_at = $6 WHERE id = $7"

	var errLogMsg string = fmt.Sprintf(noti.REPO_ERR_MSG, shared.PENDING_CAMPAIGN_REPOSITORY) + "UpdatePendingCampaign - "
	var internalErr error = errors.New(noti.INTERNALL_ERR_MSG)

	res, err := p.db.ExecContext(ctx, query, campaign.Target, campaign.Description, campaign.ProofBlobID,
		campaign.ReviewStatus, campaign.ReviewedBy, campaign.UpdatedAt, campaign.ID)
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
		return fmt.Errorf(noti.UNDEFINED_OBJECT_WARN_MSG, pending_campaign_table)
	}

	return nil
}
