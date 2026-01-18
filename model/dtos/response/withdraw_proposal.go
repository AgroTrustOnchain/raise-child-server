package response

import "time"

type WithDrawProposalResponse struct {
	ID              string      `json:"id"`
	PoolName        string      `json:"pool_name"`
	Creator         string      `json:"creator"`
	WithdrawAmount  int64       `json:"withdraw_amount"`
	Description     string      `json:"description"`
	Approvers       []string    `json:"approvers"`
	Refusers        []string    `json:"refusers"`
	ApproveWeight   int64       `json:"approve_weight"`
	RefuseWeight    int64       `json:"refuse_weight"`
	RefuseReasons   []string    `json:"refuse_reasons"`
	IsExecuted      bool        `json:"is_executed"`
	IsFromLocalPool bool        `json:"is_from_local_pool"`
	AprrovedPeriods []time.Time `json:"approved_periods"`
	RefusedPeriods  []time.Time `json:"refused_periods"`
	CreatedAt       time.Time   `json:"created_at"`
	UpdatedAt       time.Time   `json:"updated_at"`
	ClosedAt        time.Time   `json:"closed_at"`
}
