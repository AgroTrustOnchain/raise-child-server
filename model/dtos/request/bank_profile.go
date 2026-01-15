package request

type CreateBankProfileRequest struct {
	BankOrg          string `json:"bank_org" validate:"required"`
	BankCode         string `json:"bank_code" validate:"required"`
	OwnerName        string `json:"owner_name" validate:"required"`
	PayosClientID    string `json:"payos_client_id"`
	PayosApiKey      string `json:"payos_api_key"`
	PayosCheckSumKey string `json:"payos_check_sum_key"`
}

type UpdateBankProfileRequest struct {
	BankOrg          string `json:"bank_org"`
	BankCode         string `json:"bank_code"`
	OwnerName        string `json:"owner_name"`
	PayosClientID    string `json:"payos_client_id"`
	PayosApiKey      string `json:"payos_api_key"`
	PayosCheckSumKey string `json:"payos_check_sum_key"`
}
