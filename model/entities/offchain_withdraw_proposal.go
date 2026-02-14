package entities

import "time"

type OffChainWithdrawProposal struct {
	ID         string
	Purpose    string
	ProposalID string
	Target     string
	CreatedAt  time.Time
}
