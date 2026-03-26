package request

type GetPaymentsRequest struct {
	Status     string `form:"status"`
	Method     string `form:"method"`
	Actor      string `form:"actor"`
	FilterProp string `form:"filter_prop"`
	SortOrder  string `form:"sort_order"`
	PageSize   int    `form:"page_size"`
	Page       int    `form:"page"`
}

type DonateRequest struct {
	PoolId  string `json:"pool_id" validate:"required"`
	Amount  int64  `json:"amount" validate:"required,min=2000"`
	Message string `json:"message"`
}
