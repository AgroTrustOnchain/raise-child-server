package entities

import "time"

type OffChainDonation struct {
	ID          string
	Purpose     string
	Target      string
	StartPeriod string
	EndPeriod   string
	CreatedAt   time.Time
}
