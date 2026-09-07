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

type platformConfigRepo struct {
	db        *sql.DB
	errLogger *log.Logger
}

func InitializePlatformConfigRepository(db *sql.DB, errLogger *log.Logger) repository.IPlatformConfigRepository {
	return &platformConfigRepo{
		db:        db,
		errLogger: errLogger,
	}
}

// CreateConfig implements repository.IPlatformConfigRepository.
func (p *platformConfigRepo) CreateConfig(config entities.PlatformConfig, table string, ctx context.Context) error {
	panic("unimplemented")
}

// GetConfig implements repository.IPlatformConfigRepository.
func (p *platformConfigRepo) GetConfig(id string, table string, ctx context.Context) (*entities.PlatformConfig, error) {
	var query string = "SELECT * FROM " + table + " WHERE id = $1"
	var errLogMsg string = fmt.Sprintf(noti.REPO_ERR_MSG, shared.PLATFORM_CONFIG_REPOSITORY) + "GetConfig - "

	var res entities.PlatformConfig
	if err := p.db.QueryRowContext(ctx, query, id).Scan(
		&res.ID, &res.Key, &res.Value, &res.Description, &res.ActorProfileID, &res.ActorAddress,
		&res.CreatedAt, &res.UpdatedAt); err != nil {

		if err == sql.ErrNoRows {
			return nil, nil
		}

		p.errLogger.Println(errLogMsg + err.Error())
		return nil, errors.New(noti.INTERNALL_ERR_MSG)
	}

	return &res, nil
}

// GetConfigByKey implements repository.IPlatformConfigRepository.
func (p *platformConfigRepo) GetConfigByKey(table string, key string, ctx context.Context) (*entities.PlatformConfig, error) {
	var query string = "SELECT * FROM " + table + " WHERE key = $1"
	var errLogMsg string = fmt.Sprintf(noti.REPO_ERR_MSG, shared.PLATFORM_CONFIG_REPOSITORY) + "GetConfigByKey - "

	var res entities.PlatformConfig
	if err := p.db.QueryRowContext(ctx, query, key).Scan(
		&res.ID, &res.Key, &res.Value, &res.Description, &res.ActorProfileID, &res.ActorAddress,
		&res.CreatedAt, &res.UpdatedAt); err != nil {

		if err == sql.ErrNoRows {
			return nil, nil
		}

		p.errLogger.Println(errLogMsg + err.Error())
		return nil, errors.New(noti.INTERNALL_ERR_MSG)
	}

	return &res, nil
}

// GetConfigs implements repository.IPlatformConfigRepository.
func (p *platformConfigRepo) GetConfigs(req request.GetPlatformConfigsRequest, table string, ctx context.Context) ([]entities.PlatformConfig, int, error) {
	var errLogMsg string = fmt.Sprintf(noti.REPO_ERR_MSG, shared.PLATFORM_CONFIG_REPOSITORY) + "GetConfigs - "
	var internalErr error = errors.New(noti.INTERNALL_ERR_MSG)

	var queryCondition string
	var isHavePreviosCondition bool = false
	if req.Keyword != "" {
		queryCondition += fmt.Sprintf("(LOWER(key) LIKE LOWER('%s') OR LOWER(description) LIKE LOWER('%s'))", req.Keyword, req.Keyword)
		isHavePreviosCondition = true
	}

	if req.ActorAddress != "" {
		if isHavePreviosCondition {
			queryCondition += " AND "
		}

		queryCondition += fmt.Sprintf("actor_address = '%s'", req.ActorAddress)
		isHavePreviosCondition = true
	}

	var order string = "DESC"
	if req.SortOrder != "" {
		order = req.SortOrder
	}

	var query string = generateRetrieveQuery(generateRetrieveQueryRequest{
		table:       table,
		limitAmount: req.PageSize,
		condition:   queryCondition,
		order:       " ORDER BY created_at " + order,
		page:        req.Page,
		isGetCount:  false,
	})

	rows, err := p.db.QueryContext(ctx, query)
	if err != nil {
		p.errLogger.Println(errLogMsg + err.Error())
		return nil, 0, internalErr
	}
	defer rows.Close()

	var res []entities.PlatformConfig
	for rows.Next() {
		var x entities.PlatformConfig
		if err := rows.Scan(
			&x.ID, &x.Key, &x.Value, &x.Description, &x.ActorProfileID, &x.ActorAddress,
			&x.CreatedAt, &x.UpdatedAt); err != nil {
			return nil, 0, internalErr
		}

		res = append(res, x)
	}

	var totalRecords int
	p.db.QueryRowContext(ctx, generateCountTotalRecordsQuery(table, queryCondition)).Scan(&totalRecords)

	return res, caculateTotalPages(totalRecords, req.PageSize), nil
}

// UpdateConfig implements repository.IPlatformConfigRepository.
func (p *platformConfigRepo) UpdateConfig(config entities.PlatformConfig, table string, ctx context.Context) error {
	var errLogMsg string = fmt.Sprintf(noti.REPO_ERR_MSG, shared.PLATFORM_CONFIG_REPOSITORY) + "UpdateConfig - "
	var query string = "UPDATE " + table + " SET value = $1, description = $2, actor_profile_id = $3, actor_address = $4, updated_at = $5 WHERE id = $6"

	res, err := p.db.ExecContext(ctx, query, config.Value, config.Description, config.ActorProfileID, config.ActorAddress, config.UpdatedAt, config.ID)
	var internalErr error = errors.New(noti.INTERNALL_ERR_MSG)

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
		return fmt.Errorf(noti.UNDEFINED_OBJECT_WARN_MSG, table)
	}

	return nil
}
