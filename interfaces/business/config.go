package business

import (
	"context"
	"raise-child/model/dtos/request"
	"raise-child/model/dtos/response"
)

type IConfigService interface {
	UpdateChildEditMealNeedDates(req request.UpdateChildEditNeedDatesRequest, ctx context.Context) (response.BuildTransactionResponse, error)
	UpdateChildEditBooksNeedDates(req request.UpdateChildEditNeedDatesRequest, ctx context.Context) (response.BuildTransactionResponse, error)
	UpdateChildEditHealthInsuranceNeedDates(req request.UpdateChildEditNeedDatesRequest, ctx context.Context) (response.BuildTransactionResponse, error)
}
