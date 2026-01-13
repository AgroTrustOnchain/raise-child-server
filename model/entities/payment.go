package entities

import "time"

type Payment struct {
	ID            string    `json:"id"`
	Actor         string    `json:"actor"`
	TransactionId string    `json:"transaction_id"`
	Amount        int64     `json:"amount"`
	Currency      string    `json:"currency"` // e.g. "VND"
	Status        string    `json:"status"`   // e.g. "Pending", "Cancel", "Success"
	Method        string    `json:"method"`   // e.g. "VNPay", "Momo", "Payos"
	CancelReason  string    `json:"cancel_reason"`
	Message       string    `json:"message"`
	ExpiredAt     time.Time `json:"expired_at"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}
