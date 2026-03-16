package entities

import (
	"raise-child/model/dtos/response"
	"raise-child/util"
	"strconv"
	"time"
)

type WithdrawProposal struct {
	ID                  ID       `json:"id"`
	PoolID              string   `json:"pool_id"`
	PoolName            string   `json:"pool_name"`
	Creator             string   `json:"creator"`
	WithdrawAmount      string   `json:"withdraw_amount"`
	Description         string   `json:"description"`
	ProofBlobID         string   `json:"proof_blob_id"`
	Approvers           []string `json:"approvers"`
	Refusers            []string `json:"refusers"`
	ApproveWeight       string   `json:"approve_weight"`
	RefuseWeight        string   `json:"refuse_weight"`
	RefuseReasons       []string `json:"refuse_reasons"`
	IsExecuted          bool     `json:"is_executed"`
	IsFromLocalPool     bool     `json:"is_from_local_pool"`
	AprrovedPeriods     []string `json:"approved_periods"`
	RefusedPeriods      []string `json:"refused_periods"`
	TransactionRecordID *string  `json:"transaction_record_id"`
	CreatedAt           string   `json:"created_at"`
	UpdatedAt           string   `json:"updated_at"`
	ClosedAt            string   `json:"closed_at"`
}

type SpecialNeedProposal struct {
	ID              ID       `json:"id"`
	ChildID         string   `json:"child"`
	Creator         string   `json:"creator"`
	Target          string   `json:"target"`
	Description     string   `json:"description"`
	ProofBlobID     string   `json:"proof_blob_id"`
	Approvers       []string `json:"approvers"`
	Refusers        []string `json:"refusers"`
	ApproveWeight   string   `json:"approve_weight"`
	RefuseWeight    string   `json:"refuse_weight"`
	RefuseReasons   []string `json:"refuse_reasons"`
	AprrovedPeriods []string `json:"approved_periods"`
	RefusedPeriods  []string `json:"refused_periods"`
	IsConfirm       bool     `json:"is_confirm"`
	CreatedAt       string   `json:"created_at"`
	UpdatedAt       string   `json:"updated_at"`
	ClosedAt        string   `json:"closed_at"`
}

func (w WithdrawProposal) ToMinimumWithdrawProposalResponse() response.WithdrawProposalResponse {
	if w.ID.ID == "" {
		return response.WithdrawProposalResponse{}
	}

	withdrawAmount, _ := strconv.ParseInt(w.WithdrawAmount, 10, 64)
	approveWeight, _ := strconv.ParseInt(w.ApproveWeight, 10, 64)
	refuseWeight, _ := strconv.ParseInt(w.RefuseWeight, 10, 64)
	createdAt, _ := strconv.ParseInt(w.CreatedAt, 10, 64)
	updatedAt, _ := strconv.ParseInt(w.UpdatedAt, 10, 64)
	closedAt, _ := strconv.ParseInt(w.ClosedAt, 10, 64)

	return response.WithdrawProposalResponse{
		ID:             w.ID.ID,
		PoolName:       w.PoolName,
		Creator:        w.Creator,
		WithdrawAmount: withdrawAmount,
		Description:    w.Description,

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

func (w WithdrawProposal) ToWithdrawProposalResponse() response.WithdrawProposalResponse {
	if w.ID.ID == "" {
		return response.WithdrawProposalResponse{}
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

	return response.WithdrawProposalResponse{
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

func (s SpecialNeedProposal) ToMinimumSpecialNeedProposalResponse() response.SpecialNeedProposalResponse {
	if s.ID.ID == "" {
		return response.SpecialNeedProposalResponse{}
	}

	target, _ := strconv.ParseInt(s.Target, 10, 64)
	approveWeight, _ := strconv.ParseInt(s.ApproveWeight, 10, 64)
	refuseWeight, _ := strconv.ParseInt(s.RefuseWeight, 10, 64)
	createdAt, _ := strconv.ParseInt(s.CreatedAt, 10, 64)
	updatedAt, _ := strconv.ParseInt(s.UpdatedAt, 10, 64)
	closedAt, _ := strconv.ParseInt(s.ClosedAt, 10, 64)

	return response.SpecialNeedProposalResponse{
		ID:            s.ID.ID,
		ChildID:       s.ChildID,
		Creator:       s.Creator,
		Target:        target,
		Description:   s.Description,
		Approvers:     s.Approvers,
		Refusers:      s.Refusers,
		ApproveWeight: approveWeight,
		RefuseWeight:  refuseWeight,
		RefuseReasons: s.RefuseReasons,
		IsConfirm:     s.IsConfirm,
		CreatedAt:     util.MilliSecToTime(createdAt),
		UpdatedAt:     util.MilliSecToTime(updatedAt),
		ClosedAt:      util.MilliSecToTime(closedAt),
	}
}

func (s SpecialNeedProposal) ToSpecialNeedProposalResponse() response.SpecialNeedProposalResponse {
	if s.ID.ID == "" {
		return response.SpecialNeedProposalResponse{}
	}

	var loopLength int
	if len(s.AprrovedPeriods) >= len(s.RefusedPeriods) {
		loopLength = len(s.AprrovedPeriods)
	} else {
		loopLength = len(s.RefusedPeriods)
	}

	var aprrovedPeriods []time.Time
	var refusedPeriods []time.Time
	for i := 0; i < loopLength; i++ {
		if i < len(s.AprrovedPeriods) {
			period, _ := strconv.ParseInt(s.AprrovedPeriods[i], 10, 64)
			aprrovedPeriods = append(aprrovedPeriods, util.MilliSecToTime(period))
		}

		if i < len(s.RefusedPeriods) {
			period, _ := strconv.ParseInt(s.RefusedPeriods[i], 10, 64)
			refusedPeriods = append(refusedPeriods, util.MilliSecToTime(period))
		}
	}

	target, _ := strconv.ParseInt(s.Target, 10, 64)
	approveWeight, _ := strconv.ParseInt(s.ApproveWeight, 10, 64)
	refuseWeight, _ := strconv.ParseInt(s.RefuseWeight, 10, 64)
	createdAt, _ := strconv.ParseInt(s.CreatedAt, 10, 64)
	updatedAt, _ := strconv.ParseInt(s.UpdatedAt, 10, 64)
	closedAt, _ := strconv.ParseInt(s.ClosedAt, 10, 64)

	return response.SpecialNeedProposalResponse{
		ID:              s.ID.ID,
		ChildID:         s.ChildID,
		Creator:         s.Creator,
		Target:          target,
		Description:     s.Description,
		Approvers:       s.Approvers,
		Refusers:        s.Refusers,
		ApproveWeight:   approveWeight,
		RefuseWeight:    refuseWeight,
		RefuseReasons:   s.RefuseReasons,
		AprrovedPeriods: aprrovedPeriods,
		RefusedPeriods:  refusedPeriods,
		IsConfirm:       s.IsConfirm,
		CreatedAt:       util.MilliSecToTime(createdAt),
		UpdatedAt:       util.MilliSecToTime(updatedAt),
		ClosedAt:        util.MilliSecToTime(closedAt),
	}
}
