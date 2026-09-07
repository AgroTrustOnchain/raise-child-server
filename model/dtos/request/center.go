package request

type GetCentersRequest struct {
	Keyword   string `form:"keyword"`
	SortOrder string `form:"sort_order"`
	PageSize  int    `form:"page_size"`
	Page      int    `form:"page"`
}
