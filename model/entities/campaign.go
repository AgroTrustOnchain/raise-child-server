package entities

import (
	"raise-child/model/dtos/response"
	"raise-child/util"
	"strconv"
	"time"
)

type SpecialNeedCampaign struct {
	ID                      ID       `json:"id"`
	ChildID                 string   `json:"child"`
	Creator                 string   `json:"creator"`
	Target                  string   `json:"target"`
	Description             string   `json:"description"`
	ProofBlobID             string   `json:"proof_blob_id"`
	TotalDonated            string   `json:"total_donated"`
	WithdrawAmount          string   `json:"withdraw_amount"`
	Donors                  []string `json:"donors"`
	DonorTotalContributions []string `json:"donor_total_contributions"`
	RemainVotePowers        []string `json:"remain_vote_powers"`
	Donations               []string `json:"donations"`
	Withdraws               []string `json:"withdraws"`
	WithdrawProposals       []string `json:"withdraw_proposals"`
	CreatedAt               string   `json:"created_at"`
	UpdatedAt               string   `json:"updated_at"`
}

type OnChainCampaign struct {
	ID                      ID       `json:"id"`
	PoolID                  string   `json:"pool_id"`
	PoolName                string   `json:"pool_name"`
	Creator                 string   `json:"creator"`
	Target                  string   `json:"target"`
	Description             string   `json:"description"`
	ProofBlobID             string   `json:"proof_blob_id"`
	TotalDonated            string   `json:"total_donated"`
	WithdrawAmount          string   `json:"withdraw_amount"`
	Donors                  []string `json:"donors"`
	DonorTotalContributions []string `json:"donor_total_contributions"`
	RemainVotePowers        []string `json:"remain_vote_powers"`
	Donations               []string `json:"donations"`
	Withdraws               []string `json:"withdraws"`
	WithdrawProposals       []string `json:"withdraw_proposals"`
	CreatedAt               string   `json:"created_at"`
	UpdatedAt               string   `json:"updated_at"`
}

type PendingCampaign struct {
	ID             string    `json:"id"`
	ActorProfileID string    `json:"actor_profile_id"`
	ActorAddress   string    `json:"actor_address"`
	PoolID         string    `json:"pool_id"`
	PoolName       string    `json:"pool_name"`
	Target         int64     `json:"target"`
	Description    string    `json:"description"`
	ProofBlobID    *string   `json:"proof_blob_id"`
	ReviewStatus   string    `json:"review_status"`
	AIEvaluation   string    `json:"ai_evaluation"`
	ReviewedBy     *string   `json:"reviewed_by"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

func (s SpecialNeedCampaign) ToSpecialNeedCampaignResponse() response.SpecialNeedCampaignResponse {
	if s.ID.ID == "" {
		return response.SpecialNeedCampaignResponse{}
	}

	target, _ := strconv.ParseInt(s.Target, 10, 64)
	totalDonated, _ := strconv.ParseInt(s.TotalDonated, 10, 64)
	withdrawAmount, _ := strconv.ParseInt(s.WithdrawAmount, 10, 64)
	createdAt, _ := strconv.ParseInt(s.CreatedAt, 10, 64)
	updatedAt, _ := strconv.ParseInt(s.UpdatedAt, 10, 64)

	return response.SpecialNeedCampaignResponse{
		ID:                s.ID.ID,
		ChildID:           s.ChildID,
		Creator:           s.Creator,
		Target:            target,
		Description:       s.Description,
		ProofBlobID:       s.ProofBlobID,
		TotalDonated:      totalDonated,
		WithdrawAmount:    withdrawAmount,
		Donors:            s.Donors,
		Donations:         s.Donations,
		Withdraws:         s.Withdraws,
		WithdrawProposals: s.WithdrawProposals,
		CreatedAt:         util.MilliSecToTime(createdAt),
		UpdatedAt:         util.MilliSecToTime(updatedAt),
	}
}

func (s OnChainCampaign) ToOnChainCampaignResponse() response.OnChainCampaignResponse {
	if s.ID.ID == "" {
		return response.OnChainCampaignResponse{}
	}

	target, _ := strconv.ParseInt(s.Target, 10, 64)
	totalDonated, _ := strconv.ParseInt(s.TotalDonated, 10, 64)
	withdrawAmount, _ := strconv.ParseInt(s.WithdrawAmount, 10, 64)
	createdAt, _ := strconv.ParseInt(s.CreatedAt, 10, 64)
	updatedAt, _ := strconv.ParseInt(s.UpdatedAt, 10, 64)

	return response.OnChainCampaignResponse{
		ID:                s.ID.ID,
		PoolID:            s.PoolID,
		PoolName:          s.PoolName,
		Creator:           s.Creator,
		Target:            target,
		Description:       s.Description,
		ProofBlobID:       s.ProofBlobID,
		TotalDonated:      totalDonated,
		WithdrawAmount:    withdrawAmount,
		Donors:            s.Donors,
		Donations:         s.Donations,
		Withdraws:         s.Withdraws,
		WithdrawProposals: s.WithdrawProposals,
		CreatedAt:         util.MilliSecToTime(createdAt),
		UpdatedAt:         util.MilliSecToTime(updatedAt),
	}
}
