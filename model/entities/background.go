package entities

import "time"

type BackgroundRecord struct {
	ID        string
	Approvers []string
	Refusers  []string
	Role      string
	Sender    string
}

type BackgroundChildrenWithdrawProposalsRequest struct {
	ID              string
	ProfileID       string
	ActorAddress    string
	Region          string
	IsExecuted      bool
	RawProposedDate string
	CreatedAt       time.Time
	UpdatedAt       time.Time
}
