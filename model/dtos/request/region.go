package request

type GetSupportedRegionProposalsRequest struct {
	Keyword   string `json:"keyword"`
	CreatedBy string `json:"created_by"`
	SortOrder string `json:"sort_order"`
	PageSize  int    `json:"page_size"`
	Page      int    `json:"page"`
}

type CreateSupportedRegionProposalsRequest struct {
	Region  string `json:"region" validate:"required"`
	Content string `json:"content" validate:"required"`
}
