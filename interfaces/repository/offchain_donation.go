package repository

import (
	"context"
	"raise-child/model/entities"
)

type IOffChainDonationRepository interface {
	GetDonation(id string, ctx context.Context) (*entities.OffChainDonation, error)
	CreateDonation(donation entities.OffChainDonation, ctx context.Context) error
}
