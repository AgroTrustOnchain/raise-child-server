package business

import (
	"context"
	"raise-child/model/dtos/request"
	"raise-child/model/dtos/response"
	"raise-child/model/entities"
)

type IBankProfileService interface {
	GetBankProfile(id string, ctx context.Context) (response.BankProfileResponse, error)
	GetBankProfileByOwner(id string, ctx context.Context) (response.BankProfileResponse, error)
	CreateBankProfile(req request.CreateBankProfileRequest, ctx context.Context) (*entities.BankProfile, error)
	UpdateBankProfile(id string, req request.UpdateBankProfileRequest, ctx context.Context) (*entities.BankProfile, error)
}
