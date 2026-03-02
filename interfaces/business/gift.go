package business

import (
	"context"
	"raise-child/model/dtos/request"
	"raise-child/model/dtos/response"
)

type IGiftService interface {
	GetGift(id string, ctx context.Context) (response.GiftResponse, error)
	GetGiftsOfChild(id string, req request.GetGiftsRequest, ctx context.Context) (response.PaginationDataResponse, error)
	GetGiftsOfRegion(region string, req request.GetGiftsRequest, ctx context.Context) (response.PaginationDataResponse, error)
	CreateGift(req request.CreateGiftRequest, ctx context.Context) (response.BuildTransactionResponse, error)
	ConfirmReceiveGift(id string, req request.ConfirmReceiveGiftRequest, ctx context.Context) (response.BuildTransactionResponse, error)
	CancelGift(id string, req request.CancelGiftRequest, ctx context.Context) (response.BuildTransactionResponse, error)
}

// Last uc: UC-26
