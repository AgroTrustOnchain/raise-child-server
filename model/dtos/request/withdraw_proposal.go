package request

type GetWithdrawProposalsRequest struct {
	Keyword      string `form:"keyword"` // Filter for pool ID, pool name, desc
	Creator      string `form:"creator"`
	MinAmount    *int64 `form:"min_amount"`
	MaxAmount    *int64 `form:"max_amount"`
	IsExecuted   *bool  `form:"is_executed"`
	IsClosed     *bool  `form:"is_closed"`
	SortCriteria string `form:"sort_criteria"`
	SortOrder    string `form:"sort_order"`
	PageSize     int    `form:"page_size"`
	Page         int    `form:"page"`
}

type GetOnchainPendingWithdrawProposalsRequest struct {
	Keyword      string `form:"keyword"` // Filter for pool ID, pool name, desc
	Creator      string `form:"creator"`
	MinAmount    *int64 `form:"min_amount"`
	MaxAmount    *int64 `form:"max_amount"`
	SortCriteria string `form:"sort_criteria"`
	SortOrder    string `form:"sort_order"`
	PageSize     int    `form:"page_size"`
	Page         int    `form:"page"`
}

type CreateWithdrawProposalRequest struct {
	PoolID         string  `json:"pool_id" validate:"required"`
	WithdrawAmount int64   `json:"withdraw_amount" validate:"required,min=10000,max=20000000"`
	Description    string  `json:"description" validate:"required"`
	ProofBlobID    *string `json:"proof_blob_id"`
}

type GetPendingWithdrawProposalsRequest struct {
	Keyword      string `form:"keyword"`
	Creator      string `form:"creator"`
	Reviewer     string `form:"reviewer"`
	MinAmount    *int64 `form:"min_amount"`
	MaxAmount    *int64 `form:"max_amount"`
	Status       string `form:"status"`
	SortCriteria string `form:"sort_criteria"`
	SortOrder    string `form:"sort_order"`
	PageSize     int    `form:"page_size"`
	Page         int    `form:"page"`
}

type CreatePendingWithdrawProposalRequest struct {
	PoolID         string  `json:"pool_id" validate:"required"`
	WithdrawAmount int64   `json:"withdraw_amount" validate:"required,min=10000,max=20000000"`
	Description    string  `json:"description" validate:"required"`
	ProofBlobID    *string `json:"proof_blob_id"`
}
