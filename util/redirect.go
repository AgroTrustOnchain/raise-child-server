package util

import (
	"net/url"
	"os"
	"raise-child/constants/env/payment"
)

type GenerateRedirectParamRequest struct {
	OrderCode string
	Status    string
	Message   string
	TxBytes   string
}

func GeneratePaymentRedirectUrl(req GenerateRedirectParamRequest) string {
	res, _ := url.Parse(os.Getenv(payment.PAYMENT_REDIRECT_URL))

	var query = res.Query()
	query.Set("transaction", req.OrderCode)
	query.Set("status", req.Status)
	query.Set("message", req.Message)
	if req.TxBytes != "" {
		query.Set("tx_bytes", req.TxBytes)
	}

	res.RawQuery = query.Encode()
	return res.String()
}
