package repository

import (
	"context"
	"raise-child/model/entities"
)

type IBankProfileRepository interface {
	GetBankProfileByOwner(owner string, ctx context.Context) (*entities.BankProfile, error)
	GetBankProfileById(id string, ctx context.Context) (*entities.BankProfile, error)
	CreateBankProfile(bp entities.BankProfile, ctx context.Context) error
	UpdateBankProfile(bp entities.BankProfile, ctx context.Context) error
}
