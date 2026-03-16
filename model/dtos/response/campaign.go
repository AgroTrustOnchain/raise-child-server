package response

import "time"

type SpecialNeedCampaignResponse struct {
	ID                string    `json:"id"`
	ChildID           string    `json:"child"`
	Creator           string    `json:"creator"`
	Target            int64     `json:"target"`
	Description       string    `json:"description"`
	ProofBlobID       string    `json:"proof_blob_id"`
	TotalDonated      int64     `json:"total_donated"`
	WithdrawAmount    int64     `json:"withdraw_amount"`
	Donations         []string  `json:"donations"`
	Withdraws         []string  `json:"withdraws"`
	WithdrawProposals []string  `json:"withdraw_proposals"`
	CreatedAt         time.Time `json:"created_at"`
	UpdatedAt         time.Time `json:"updated_at"`
}
