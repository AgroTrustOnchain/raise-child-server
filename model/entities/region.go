package entities

import "time"

type SupportedRegionSuggestion struct {
	ID           string    `json:"id"`
	ProfileID    string    `json:"profile_id"`
	Region       string    `json:"region"`
	Content      string    `json:"content"`
	Status       string    `json:"status"`
	RefuseReason *string   `json:"refuse_reason"`
	CreatedBy    string    `json:"created_by"`
	ReviewedBy   *string   `json:"reviewed_by"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}
