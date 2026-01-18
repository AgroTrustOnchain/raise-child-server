package entities

import (
	"raise-child/model/dtos/response"
	"raise-child/util"
	"strconv"
	"time"
)

type WithDrawProposal struct {
	ID              ID       `json:"id"`
	PoolID          string   `json:"pool_id"`
	PoolName        string   `json:"pool_name"`
	Creator         string   `json:"creator"`
	WithdrawAmount  string   `json:"withdraw_amount"`
	Description     string   `json:"description"`
	Approvers       []string `json:"approvers"`
	Refusers        []string `json:"refusers"`
	ApproveWeight   string   `json:"approve_weight"`
	RefuseWeight    string   `json:"refuse_weight"`
	RefuseReasons   []string `json:"refuse_reasons"`
	IsExecuted      bool     `json:"is_executed"`
	IsFromLocalPool bool     `json:"is_from_local_pool"`
	AprrovedPeriods []string `json:"approved_periods"`
	RefusedPeriods  []string `json:"refused_periods"`
	CreatedAt       string   `json:"created_at"`
	UpdatedAt       string   `json:"updated_at"`
	ClosedAt        string   `json:"closed_at"`
}

func (w WithDrawProposal) ToWithDrawProposalResponse() response.WithDrawProposalResponse {
	if w.ID.ID == "" {
		return response.WithDrawProposalResponse{}
	}

	var loopLength int
	if len(w.AprrovedPeriods) >= len(w.RefusedPeriods) {
		loopLength = len(w.AprrovedPeriods)
	} else {
		loopLength = len(w.RefusedPeriods)
	}

	var aprrovedPeriods []time.Time
	var refusedPeriods []time.Time
	for i := 0; i < loopLength; i++ {
		if i < len(w.AprrovedPeriods) {
			period, _ := strconv.ParseInt(w.AprrovedPeriods[i], 10, 64)
			aprrovedPeriods = append(aprrovedPeriods, util.MilliSecToTime(period))
		}

		if i < len(w.RefusedPeriods) {
			period, _ := strconv.ParseInt(w.RefusedPeriods[i], 10, 64)
			refusedPeriods = append(refusedPeriods, util.MilliSecToTime(period))
		}
	}

	withdrawAmount, _ := strconv.ParseInt(w.WithdrawAmount, 10, 64)
	approveWeight, _ := strconv.ParseInt(w.ApproveWeight, 10, 64)
	refuseWeight, _ := strconv.ParseInt(w.RefuseWeight, 10, 64)
	createdAt, _ := strconv.ParseInt(w.CreatedAt, 10, 64)
	updatedAt, _ := strconv.ParseInt(w.UpdatedAt, 10, 64)
	closedAt, _ := strconv.ParseInt(w.ClosedAt, 10, 64)

	return response.WithDrawProposalResponse{
		ID:              w.ID.ID,
		PoolName:        w.PoolName,
		Creator:         w.Creator,
		WithdrawAmount:  withdrawAmount,
		Description:     w.Description,
		Approvers:       w.Approvers,
		Refusers:        w.Refusers,
		ApproveWeight:   approveWeight,
		RefuseWeight:    refuseWeight,
		RefuseReasons:   w.RefuseReasons,
		IsExecuted:      w.IsExecuted,
		IsFromLocalPool: w.IsFromLocalPool,
		AprrovedPeriods: aprrovedPeriods,
		RefusedPeriods:  refusedPeriods,
		CreatedAt:       util.MilliSecToTime(createdAt),
		UpdatedAt:       util.MilliSecToTime(updatedAt),
		ClosedAt:        util.MilliSecToTime(closedAt),
	}
}

func (w WithDrawProposal) ToMinimumWithDrawProposalResponse() response.WithDrawProposalResponse {
	if w.ID.ID == "" {
		return response.WithDrawProposalResponse{}
	}

	withdrawAmount, _ := strconv.ParseInt(w.WithdrawAmount, 10, 64)
	approveWeight, _ := strconv.ParseInt(w.ApproveWeight, 10, 64)
	refuseWeight, _ := strconv.ParseInt(w.RefuseWeight, 10, 64)
	createdAt, _ := strconv.ParseInt(w.CreatedAt, 10, 64)
	updatedAt, _ := strconv.ParseInt(w.UpdatedAt, 10, 64)
	closedAt, _ := strconv.ParseInt(w.ClosedAt, 10, 64)

	return response.WithDrawProposalResponse{
		ID:              w.ID.ID,
		PoolName:        w.PoolName,
		Creator:         w.Creator,
		WithdrawAmount:  withdrawAmount,
		Description:     w.Description,
		Approvers:       w.Approvers,
		Refusers:        w.Refusers,
		ApproveWeight:   approveWeight,
		RefuseWeight:    refuseWeight,
		RefuseReasons:   w.RefuseReasons,
		IsExecuted:      w.IsExecuted,
		IsFromLocalPool: w.IsFromLocalPool,
		CreatedAt:       util.MilliSecToTime(createdAt),
		UpdatedAt:       util.MilliSecToTime(updatedAt),
		ClosedAt:        util.MilliSecToTime(closedAt),
	}
}
