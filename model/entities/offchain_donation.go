package entities

import "time"

type OffChainDonation struct {
	ID             string
	Purpose        string
	Target         string
	MealDurationID *string
	CreatedAt      time.Time
}
