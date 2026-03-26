package request

type GetPendingChildSpecialNeedProposalsRequest struct {
	Keyword      string `form:"keyword"`
	Region       string `form:"region"`
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
