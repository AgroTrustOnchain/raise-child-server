package request

type GetPaymentsRequest struct {
	Status     string `json:"status"`
	Method     string `json:"method"`
	Actor      string `json:"actor"`
	FilterProp string `json:"filter_prop"`
	SortOrder  string `json:"sort_order"`
	PageSize   int    `json:"page_size"`
	Page       int    `json:"page"`
}

type DonateRequest struct {
	PoolId  string `json:"pool_id" validate:"required"`
	Amount  int64  `json:"amount" validate:"required,min=2000"`
	Message string `json:"message"`
}
