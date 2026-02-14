package entities

import "time"

type Payment struct {
	ID            string    `json:"id"`
	Actor         string    `json:"actor"`
	Sub           string    `json:"sub"`
	ProposalID    *string   `json:"proposal_id"`
	DonationID    *string   `json:"donation_id"`
	IsDonateTx    bool      `json:"is_donate_tx"`
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

type PaymentPurpose string

const (
	DONATE_PURPOSE       PaymentPurpose = "Donate"
	WITHDRAW_PURPOSE     PaymentPurpose = "Withdraw"
	BOOKS_NEED_PURPOSE   PaymentPurpose = "Child Books Need"
	MEAL_NEED_PURPOSE    PaymentPurpose = "Child Meal Need"
	SPECIAL_NEED_PURPOSE PaymentPurpose = "Child Special Need"
)




