package entities

import "time"

type OffChainWithdrawProposal struct {
	ID         string
	Purpose    string
	ProposalID string
	Target     string
	CreatedAt  time.Time
}

type PendingWithdrawProposal struct {
	ID             string
	ProfileID      string
	Creator        string
	PoolID         string
	PoolName       string
	WithdrawAmount int64
	ProofBlobID    *string
	Description    string
	Status         string
	AIEvaluation   string
	ReviewedBy     *string
	CreatedAt      time.Time
	UpdatedAt      time.Time
}
