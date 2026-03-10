package repository

import (
	"raise-child/model/entities"

	"golang.org/x/net/context"
)

type INftProfileRepository interface {
	CreateNftProfile(nft entities.NftProfile, ctx context.Context) error
	GetUserNfts(id string, ctx context.Context) ([]entities.NftProfile, error)
	GetUserNftByRole(id, role string, ctx context.Context) (*entities.NftProfile, error)
	GetNftProfile(id string, ctx context.Context) (*entities.NftProfile, error)
}
