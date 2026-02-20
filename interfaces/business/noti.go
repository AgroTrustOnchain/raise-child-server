package business

import (
	"context"
	"raise-child/model/dtos/request"
	"raise-child/model/dtos/response"
)

type INotiService interface {
	GetCurrentWalletNotis(wallet string, req request.GetNotisRequest, ctx context.Context) (response.PaginationDataResponse, error)
}
