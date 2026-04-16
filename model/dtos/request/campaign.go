package request

type GetPendingCampaignsRequest struct {
	Keyword      string `form:"keyword"`
	PoolName     string `form:"pool_name"`
	Status       string `form:"status"`
	Creator      string `form:"creator"`
	Reviewer     string `form:"reviewer"`
	MinAmount    *int64 `form:"min_amount"`
	MaxAmount    *int64 `form:"max_amount"`
	SortCriteria string `form:"sort_criteria"`
	SortOrder    string `form:"sort_order"`
	PageSize     int    `form:"page_size"`
	Page         int    `form:"page"`
}

type CreatePendingCampaignRequest struct {
	PoolName    string  `json:"pool_name" validate:"required"`
	Target      int64   `json:"target" validate:"required,min=50000"`
	Description string  `json:"description" validate:"required"`
	ProofBlobID *string `json:"proof_blob_id"`
}

type GetCampaignsRequest struct {
	Keyword      string `form:"keyword"`
	PoolName     string `form:"pool_name"`
	Creator      string `form:"creator"`
	MinAmount    *int64 `form:"min_amount"`
	MaxAmount    *int64 `form:"max_amount"`
	SortCriteria string `form:"sort_criteria"`
	SortOrder    string `form:"sort_order"`
	PageSize     int    `form:"page_size"`
	Page         int    `form:"page"`
}

type CreateCampaignWithdrawProposalRequest struct {
	CampaignID  string  `json:"campaign_id" validate:"required"`
	Amount      int64   `json:"amount" validate:"required,min=2000"`
	Description string  `json:"description" validate:"required"`
	ProofBlobID *string `json:"proof_blob_id"`
}

type SupportCampaignRequest struct {
	Amount      int64  `json:"amount" validate:"required,min=2000"`
	Description string `json:"description" validate:"required"`
}
