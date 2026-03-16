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

type supportedRegionProposalRepo struct {
	db        *sql.DB
	errLogger *log.Logger
}

const supported_region_proposal_table string = "supported_region_suggestions"

func InitializeSupportedRegionSuggestionRepository(db *sql.DB, errLogger *log.Logger) repository.ISupportedRegionSuggestionRepository {
	return &supportedRegionProposalRepo{
		db:        db,
		errLogger: errLogger,
	}
}

// CreateSupportedRegionSuggestion implements repository.ISupportedRegionSuggestionRepository.
func (s *supportedRegionProposalRepo) CreateSupportedRegionSuggestion(proposal entities.SupportedRegionSuggestion, ctx context.Context) error {
	var query string = "INSERT INTO " + supported_region_proposal_table +
		" (id, profile_id, region, content, created_by) " +
		"values ($1, $2, $3, $4, $5)"

	var errLogMsg string = fmt.Sprintf(noti.REPO_ERR_MSG, shared.SUPPORTED_REGION_PROPOSAL) + "CreateSupportedRegionSuggestion - "

	if _, err := s.db.ExecContext(ctx, query, proposal.ID, proposal.ProfileID, proposal.Region, proposal.Content, proposal.CreatedBy); err != nil {

		s.errLogger.Println(errLogMsg + err.Error())
		return errors.New(noti.INTERNALL_ERR_MSG)
	}

	return nil
}

// GetSupportedRegionSuggestion implements repository.ISupportedRegionSuggestionRepository.
func (s *supportedRegionProposalRepo) GetSupportedRegionSuggestion(id string, ctx context.Context) (*entities.SupportedRegionSuggestion, error) {
	var query string = "SELECT * FROM " + supported_region_proposal_table + " WHERE id = $1"
	var errLogMsg string = fmt.Sprintf(noti.REPO_ERR_MSG, shared.SUPPORTED_REGION_PROPOSAL) + "GetSupportedRegionSuggestion - "

	var res entities.SupportedRegionSuggestion
	if err := s.db.QueryRowContext(ctx, query, id).Scan(
		&res.ID, &res.ProfileID, &res.Region, &res.Content, &res.Status,
		&res.CreatedBy, &res.ReviewedBy, &res.CreatedAt, &res.UpdatedAt); err != nil {

		if err == sql.ErrNoRows {
			return nil, nil
		}

		s.errLogger.Println(errLogMsg + err.Error())
		return nil, errors.New(noti.INTERNALL_ERR_MSG)
	}

	return &res, nil
}

// GetSupportedRegionSuggestions implements repository.ISupportedRegionSuggestionRepository.
func (s *supportedRegionProposalRepo) GetSupportedRegionSuggestions(req request.GetSupportedRegionSuggestionsRequest, isGuestView bool, ctx context.Context) ([]entities.SupportedRegionSuggestion, int, error) {
	var errLogMsg string = fmt.Sprintf(noti.REPO_ERR_MSG, shared.PAYMENT_REPOSITORY) + "GetPayments - "
	var internalErr error = errors.New(noti.INTERNALL_ERR_MSG)

	var queryCondition string
	var isHavePreviosCondition bool = false
	if isGuestView {
		queryCondition += "status = 'Approved'"
		isHavePreviosCondition = true
	}

	if req.Keyword != "" {
		if isHavePreviosCondition {
			queryCondition += " AND "
		}

		queryCondition += fmt.Sprintf("(LOWER(region) LIKE LOWER('%%%%%s%%%%') OR LOWER(content) LIKE LOWER('%%%%%s%%%%'))", req.Keyword, req.Keyword)
		isHavePreviosCondition = true
	}

	if req.CreatedBy != "" {
		if isHavePreviosCondition {
			queryCondition += " AND "
		}

		queryCondition += fmt.Sprintf("created_by = %s", req.CreatedBy)
		isHavePreviosCondition = true
	}

	var order string = "DESC"
	if req.SortOrder != "" {
		order = req.SortOrder
	}

	queryCondition += "ORDER BY created_at " + order
	var query string = generateRetrieveQuery(generateRetrieveQueryRequest{
		table:       supported_region_proposal_table,
		limitAmount: req.PageSize,
		condition:   queryCondition,
		page:        req.Page,
		isGetCount:  false,
	})

	rows, err := s.db.QueryContext(ctx, query)
	if err != nil {
		s.errLogger.Println(errLogMsg + err.Error())
		return nil, 0, internalErr
	}

	var res []entities.SupportedRegionSuggestion
	for rows.Next() {
		var x entities.SupportedRegionSuggestion
		if err := rows.Scan(
			&x.ID, &x.ProfileID, &x.Region, &x.Content, &x.Status,
			&x.CreatedBy, &x.ReviewedBy, &x.CreatedAt, &x.UpdatedAt); err != nil {
			s.errLogger.Println(errLogMsg + err.Error())
			return nil, 0, internalErr
		}

		res = append(res, x)
	}

	var totalRecords int
	s.db.QueryRowContext(ctx, generateCountTotalRecordsQuery(supported_region_proposal_table, queryCondition)).Scan(&totalRecords)

	return res, caculateTotalPages(totalRecords, req.PageSize), nil
}

// IsRegionRequested implements repository.ISupportedRegionSuggestionRepository.
func (s *supportedRegionProposalRepo) IsRegionRequested(region string, ctx context.Context) (bool, error) {
	var query string = "SELECT id FROM " + supported_region_proposal_table + " WHERE LOWER(region) = LOWER($1) ADN (status = 'Pending' OR status 'Approved') LIMIT 1"

	var id string
	if err := s.db.QueryRowContext(ctx, query, region).Scan(&id); err != nil {
		if err == sql.ErrNoRows {
			return false, nil
		}

		s.errLogger.Println(fmt.Sprintf(noti.REPO_ERR_MSG, shared.SUPPORTED_REGION_PROPOSAL) + "IsRegionRequested - " + err.Error())
		return false, errors.New(noti.INTERNALL_ERR_MSG)
	}

	return id != "", nil
}

// UpdateSupportedRegionSuggestion implements repository.ISupportedRegionSuggestionRepository.
func (s *supportedRegionProposalRepo) UpdateSupportedRegionSuggestion(proposal entities.SupportedRegionSuggestion, ctx context.Context) error {
	var errLogMsg string = fmt.Sprintf(noti.REPO_ERR_MSG, shared.SUPPORTED_REGION_PROPOSAL) + "UpdateSupportedRegionSuggestion - "
	var query string = "UPDATE " + volunteer_noti_table + " SET content = $1, status = $2, reviewed_by = $3 WHERE id = $4"

	res, err := s.db.ExecContext(ctx, query, proposal.Content, proposal.Status, proposal.ReviewedBy, proposal.ID)
	var internalErr error = errors.New(noti.INTERNALL_ERR_MSG)

	if err != nil {
		s.errLogger.Println(errLogMsg + err.Error())
		return internalErr
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		s.errLogger.Println(errLogMsg + err.Error())
		return internalErr
	}

	if rowsAffected == 0 {
		return errors.New(fmt.Sprintf(noti.UNDEFINED_OBJECT_WARN_MSG, supported_region_proposal_table))
	}

	return nil
}
