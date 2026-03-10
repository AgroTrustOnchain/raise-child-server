package request

type GetWithdrawProposalsRequest struct {
	Keyword      string `json:"keyword"` // Filter for pool ID, pool name, desc
	Creator      string `json:"creator"`
	MinAmount    *int64 `json:"min_amount"`
	MaxAmount    *int64 `json:"max_amount"`
	IsExecuted   *bool  `json:"is_executed"`
	IsClosed     *bool  `json:"is_closed"`
	SortCriteria string `json:"sort_criteria"`
	SortOrder    string `json:"sort_order"`
	PageSize     int    `json:"page_size"`
	Page         int    `json:"page"`
}

type CreateWithdrawProposalRequest struct {
	PoolID         string `json:"pool_id" validate:"required"`
	WithdrawAmount int64  `json:"withdraw_amount" validate:"required,min=10000,max=20000000"`
	Description    string `json:"description" validate:"required"`
	ProofBlobID    string `json:"proof_blob_id"`
}

type GetPendingWithdrawProposalsRequest struct {
	Keyword      string `json:"keyword"`
	Creator      string `json:"creator"`
	Reviewer     string `json:"reviewer"`
	MinAmount    *int64 `json:"min_amount"`
	MaxAmount    *int64 `json:"max_amount"`
	Status       string `json:"status"`
	SortCriteria string `json:"sort_criteria"`
	SortOrder    string `json:"sort_order"`
	PageSize     int    `json:"page_size"`
	Page         int    `json:"page"`
}

type CreatePendingWithdrawProposalRequest struct {
	PoolID         string  `json:"pool_id" validate:"required"`
	WithdrawAmount int64   `json:"withdraw_amount" validate:"required,min=10000,max=20000000"`
	Description    string  `json:"description" validate:"required"`
	ProofBlobID    *string `json:"proof_blob_id"`
}
