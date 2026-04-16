package business

import (
	"context"
	"raise-child/model/dtos/response"
)

type IChildNeedService interface {
	GetChildBooksNeed(id string, ctx context.Context) (response.BooksNeedResponse, error)
	GetChildHealthInsuranceNeed(id string, ctx context.Context) (response.HealthInsuranceNeedResponse, error)
	GetChildMealNeed(id string, ctx context.Context) (response.MealNeedResponse, error)
}
