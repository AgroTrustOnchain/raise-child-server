package entities

import "time"

type LeaveRequest struct {
	ID           string    `json:"id"`
	ProfileID    string    `json:"profile_id"`
	ActorAddress string    `json:"actor_address"`
	Reason       string    `json:"reason"`
	StartDate    time.Time `json:"start_date"`
	EndDate      time.Time `json:"end_date"`
	Status       string    `json:"status"`
	ReviewedBy   *string   `json:"reviewed_by"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}
