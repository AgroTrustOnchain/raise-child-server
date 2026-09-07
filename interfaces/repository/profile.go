package repository

import (
	"context"
	"raise-child/model/entities"
)

type IProfileRepository interface {
	GetFirstProfile(ctx context.Context) (*entities.Profile, error)
	GetProfileOfFirsts(position int, ctx context.Context) (*entities.Profile, error)
	GetProfile(id string, ctx context.Context) (*entities.Profile, error)
	GetProfileByWalletAddress(wallet string, ctx context.Context) (*entities.Profile, error)
	CreateProfile(pfl entities.Profile, ctx context.Context) error
	UploadProfile(pfl entities.Profile, ctx context.Context) error
	Login(id string, token string, walletAddress string, ctx context.Context) error
	Logout(id string, ctx context.Context) error
	IsPersonalInfoExist(identityCode, phoneNumber, email string, ctx context.Context) (bool, error)
	IsPhoneNumberRegistered(phoneNumber string, ctx context.Context) (bool, error)
	IsEmailRegistered(email string, ctx context.Context) (bool, error)
}
