package entities

import "time"

type SupportedRegionSuggestion struct {
	ID         string
	ProfileID  string
	Region     string
	Content    string
	Status     string
	CreatedBy  string
	ReviewedBy *string
	CreatedAt  time.Time
	UpdatedAt  time.Time
}
