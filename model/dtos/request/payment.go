package request

type GetPaymentsRequest struct {
	Keyword          string `form:"keyword"`
	Status           string `form:"status"`
	Method           string `form:"method"`
	Actor            string `form:"actor"`
	MinAmount        *int64 `form:"min_amount"`
	MaxAmount        *int64 `form:"max_amount"`
	IsDonatePayment  *bool  `form:"is_donate_payment"`
	IsPaymentExpired *bool  `form:"is_payment_expired"`
	FilterProp       string `form:"filter_prop"`
	SortOrder        string `form:"sort_order"`
	PageSize         int    `form:"page_size"`
	Page             int    `form:"page"`
}

type DonateRequest struct {
	PoolId  string `json:"pool_id" validate:"required"`
	Amount  int64  `json:"amount" validate:"required,min=2000"`
	Message string `json:"message"`
}

type PaymentAuthCallbackRequest struct {
	ProofBlobID string `json:"proof_blob_id" validate:"required"`
}
