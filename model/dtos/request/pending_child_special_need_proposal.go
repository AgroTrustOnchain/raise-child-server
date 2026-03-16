package request

type GetPendingChildSpecialNeedProposalsRequest struct {
	Keyword      string `json:"keyword"`
	Region       string `json:"region"`
	Status       string `json:"status"`
	Creator      string `json:"creator"`
	Reviewer     string `json:"reviewer"`
	MinAmount    *int64 `json:"min_amount"`
	MaxAmount    *int64 `json:"max_amount"`
	SortCriteria string `json:"sort_criteria"`
	SortOrder    string `json:"sort_order"`
	PageSize     int    `json:"page_size"`
	Page         int    `json:"page"`
}
