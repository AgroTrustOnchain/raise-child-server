package response

type BankProfileResponse struct {
	ID        string `json:"id"`
	Owner     string `json:"owner"`
	BankOrg   string `json:"bank_org"`
	BankCode  string `json:"bank_code"`
	OwnerName string `json:"owner_name"`
}
