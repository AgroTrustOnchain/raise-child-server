package business

import (
	"context"
	"raise-child/model/dtos/request"
	"raise-child/model/dtos/response"
	"raise-child/model/entities"
)

// UpdateChildEditMealNeedDates(req request.UpdateChildEditNeedDatesRequest, ctx context.Context) (response.BuildTransactionResponse, error)
// UpdateChildEditBooksNeedDates(req request.UpdateChildEditNeedDatesRequest, ctx context.Context) (response.BuildTransactionResponse, error)
// UpdateChildEditHealthInsuranceNeedDates(req request.UpdateChildEditNeedDatesRequest, ctx context.Context) (response.BuildTransactionResponse, error)

type IConfigService interface {
	UpdateChildEditMealNeedDates(req request.UpdateChildEditNeedDatesRequest, ctx context.Context) error
	UpdateChildEditBooksNeedDates(req request.UpdateChildEditNeedDatesRequest, ctx context.Context) error
	UpdateChildEditHealthInsuranceNeedDates(req request.UpdateChildEditNeedDatesRequest, ctx context.Context) error
	EditSpecialNeedDao(req request.EditDaoRequest, ctx context.Context) error
	EditChildMaxAgeInSupport(ctx context.Context) error
}

type IPlatformConfigService interface {
	GetConfigs(req request.GetPlatformConfigsRequest, isNumericConfig bool, ctx context.Context) (response.PaginationDataResponse, error)
	GetConfig(id string, isNumericConfig bool, ctx context.Context) (*entities.PlatformConfig, error)
	UpdateConfig(id string, req request.UpdatePlatformConfigRequest, isNumericConfig bool, ctx context.Context) error
}
