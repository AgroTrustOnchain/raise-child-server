package response

type PersonalProfileResponse struct {
	ID           string `json:"id"`
	IdentityCode string `json:"identity_code"`
	FirstName    string `json:"first_name"`
	LastName     string `json:"last_name"`
	Gender       string `json:"gender"`
	DateOfBirth  string `json:"date_of_birth"`
	PhoneNumber  string `json:"phone_number"`
	Email        string `json:"email"`
}

type PersonalWalletProfileResponse struct {
	WalletAddress   string                `json:"wallet_address"`
	FirstName       string                `json:"first_name"`
	LastName        string                `json:"last_name"`
	TotalDonation   int64                 `json:"total_donation"`
	SupportedChilds []string              `json:"supported_childs"`
	TxRecords       []TransactionResponse `json:"transaction_records"`
	RecordAmount    int                   `json:"record_amount"`
	Page            int                   `json:"page"`
	TotalPages      int                   `json:"total_pages"`
}
