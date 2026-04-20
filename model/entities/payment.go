package entities

import "time"

type Payment struct {
	ID            string     `json:"id"`
	Actor         string     `json:"actor"`
	ProfileID     string     `json:"profile_id"`
	ProposalID    *string    `json:"proposal_id"`
	DonationID    *string    `json:"donation_id"`
	IsDonateTx    bool       `json:"is_donate_tx"`
	TransactionId string     `json:"transaction_id"`
	Amount        int64      `json:"amount"`
	Currency      string     `json:"currency"` // e.g. "VND"
	Status        string     `json:"status"`   // e.g. "Pending", "Cancel", "Success"
	Method        string     `json:"method"`   // e.g. "VNPay", "Momo", "Payos"
	CancelReason  *string    `json:"cancel_reason"`
	Message       string     `json:"message"`
	ExpiredAt     time.Time  `json:"expired_at"`
	CreatedAt     time.Time  `json:"created_at"`
	UpdatedAt     time.Time  `json:"updated_at"`
	ProofBlobID   *string    `json:"proof_blob_id"`
	ReviewedBy    *string    `json:"reviewed_by"`
	ReviewStatus  string     `json:"review_status"`
	IsTransferred bool       `json:"is_transferred"`
	TransferredAt *time.Time `json:"transferred_at"`
}

type Purpose string

const (
	DONATE_PURPOSE                Purpose = "Donate"
	WITHDRAW_PURPOSE              Purpose = "Withdraw"
	BOOKS_NEED_PURPOSE            Purpose = "Child Books Need"
	MEAL_NEED_PURPOSE             Purpose = "Child Meal Need"
	SPECIAL_NEED_PURPOSE          Purpose = "Child Special Need"
	HEALTH_INSURANCE_NEED_PURPOSE Purpose = "Child Health Insurance Need"
	CAMPAIGN_PURPOSE              Purpose = "Pool Campaign"
)

type PaymentDescription string

const (
	MEAL_NEED_PAYMENT_DESCRIPTION             PaymentDescription = "Meal Need"
	BOOKS_NEED_PAYMENT_DESCRIPTION            PaymentDescription = "Books Need"
	HEALTH_INSRUANCE_PAYMENT_DESCRIPTION      PaymentDescription = "Health Insurance"
	SPECIAL_NEED_CAMPAIGN_PAYMENT_DESCRIPTION PaymentDescription = "Special Campaign"
	POOL_CAMPAIGN_PAYMENT_DESCRIPTION         PaymentDescription = "Pool Campaign"
)

func (p PaymentDescription) GenerateSupportPaymentDescription() string {
	return "Support " + string(p)
}

func (p PaymentDescription) GenerateWithdrawPaymentDescription() string {
	return "Withdraw " + string(p)
}
