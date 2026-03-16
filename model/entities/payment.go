package entities

import "time"

type Payment struct {
	ID            string    `json:"id"`
	Actor         string    `json:"actor"`
	ProfileID     string    `json:"profile_id"`
	ProposalID    *string   `json:"proposal_id"`
	DonationID    *string   `json:"donation_id"`
	IsDonateTx    bool      `json:"is_donate_tx"`
	TransactionId string    `json:"transaction_id"`
	Amount        int64     `json:"amount"`
	Currency      string    `json:"currency"` // e.g. "VND"
	Status        string    `json:"status"`   // e.g. "Pending", "Cancel", "Success"
	Method        string    `json:"method"`   // e.g. "VNPay", "Momo", "Payos"
	CancelReason  *string   `json:"cancel_reason"`
	Message       string    `json:"message"`
	ExpiredAt     time.Time `json:"expired_at"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

type Purpose string

const (
	DONATE_PURPOSE                Purpose = "Donate"
	WITHDRAW_PURPOSE              Purpose = "Withdraw"
	BOOKS_NEED_PURPOSE            Purpose = "Child Books Need"
	MEAL_NEED_PURPOSE             Purpose = "Child Meal Need"
	SPECIAL_NEED_PURPOSE          Purpose = "Child Special Need"
	HEALTH_INSURANCE_NEED_PURPOSE Purpose = "Child Health Insurance Need"
)
