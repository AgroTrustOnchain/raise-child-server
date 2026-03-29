package request

type GetTransactionRecordsRequest struct {
	Keyword      string `form:"keyword"`
	PoolID       string `form:"pool_id"`
	Actor        string `form:"actor"`
	ActionType   string `form:"action_type"` // e.g. "Withdraw", "Donate"
	MinAmount    *int64 `form:"min_amount"`
	MaxAmount    *int64 `form:"max_amount"`
	SortCriteria string `form:"sort_criteria"`
	SortOrder    string `form:"sort_order"`
	PageSize     int    `form:"page_size"`
	Page         int    `form:"page"`
}
