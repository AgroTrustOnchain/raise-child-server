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
	NeedID                  string    `json:"need_id"`
	NeedType                string    `json:"need_tpe"`
	ChildID                 string    `json:"child_id"`
	Region                  string    `json:"region"`
	AssignedLeaders         []string  `json:"assgined_leaders"`
	ExpectedWithdrawPeriods []string  `json:"expected_withdraw_periods"`
	GeneralContent          string    `json:"general_content"`
	Contents                []string  `json:"contents"`
	CreatedAt               time.Time `json:"created_at"`
	UpdatedAt               time.Time `json:"updated_at"`
}
