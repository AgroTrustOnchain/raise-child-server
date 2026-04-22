package request

type GetSupportedRegionSuggestionsRequest struct {
	Keyword   string `form:"keyword"`
	CreatedBy string `form:"created_by"`
	SortOrder string `form:"sort_order"`
	PageSize  int    `form:"page_size"`
	Page      int    `form:"page"`
}

type CreateSupportedRegionSuggestionsRequest struct {
	Region  string `json:"region" validate:"required"`
	Content string `json:"content" validate:"required"`
}

type GetChildrenFromRegionDetailRequest struct {
	Keyword   string `form:"keyword"`
	SortOrder string `form:"sort_order"`
	PageSize  int    `form:"page_size"`
	Page      int    `form:"page"`
}
