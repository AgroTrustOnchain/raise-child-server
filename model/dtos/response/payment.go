package response

type PaymentUrlResponse struct {
	Url       string `json:"url"`
	PaymentID string `json:"payment_id"`
}
