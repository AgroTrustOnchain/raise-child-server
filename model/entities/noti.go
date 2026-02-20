package entities

import "time"

type VolunteerNoti struct {
	ID                 string    `json:"id"`
	ChildID            string    `json:"child_id"`
	Region             string    `json:"region"`
	AssginedVolunteers []string  `json:"assgined_volunteers"`
	Content            string    `json:"content"`
	StartPeriod        time.Time `json:"start_period"`
	EndPeriod          time.Time `json:"end_period"`
	CreatedAt          time.Time `json:"created_at"`
	UpdatedAt          time.Time `json:"updated_at"`
}

type LeaderNoti struct {
	ID                      string    `json:"id"`
	MealNeedID              string    `json:"meal_need_id"`
	ChildID                 string    `json:"child_id"`
	Region                  string    `json:"region"`
	AssignedLeaders         []string  `json:"assgined_leaders"`
	ExpectedWithdrawPeriods []string  `json:"expected_withdraw_periods"`
	Content                 string    `json:"content"`
	CreatedAt               time.Time `json:"created_at"`
	UpdatedAt               time.Time `json:"updated_at"`
}
