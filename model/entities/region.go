package entities

import "time"

type SupportedRegionProposal struct {
	ID        string
	Sub       string
	Region    string
	Content   string
	CreatedBy string
	CreatedAt time.Time
	UpdatedAt time.Time
}
