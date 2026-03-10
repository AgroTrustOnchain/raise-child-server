package entities

import "time"

type VolunteerTask struct {
	ID                string    `json:"id"`
	AssignedProfileID string    `json:"assigned_profile_id"`
	AssignedVolunteer string    `json:"assgined_volunteer"`
	ChildID           string    `json:"child_id"`
	Region            string    `json:"region"`
	Content           string    `json:"content"`
	StartPeriod       string    `json:"start_period"`
	EndPeriod         string    `json:"end_period"`
	IsEnd             bool      `json:"is_end"`
	CreatedAt         time.Time `json:"created_at"`
	UpdatedAt         time.Time `json:"updated_at"`
}
