package entities

import (
	"raise-child/model/dtos/response"
	"time"
)

type BankProfile struct {
	ID               string    `json:"id"`
	ProfileID        string    `json:"profile_id"`
	Owner            string    `json:"owner"`
	BankOrg          string    `json:"bank_org"`
	BankCode         string    `json:"bank_code"`
	OwnerName        string    `json:"owner_name"`
	PayosClientID    string    `json:"payos_client_id"`
	PayosApiKey      string    `json:"payos_api_key"`
	PayosCheckSumKey string    `json:"payos_check_sum_key"`
	CreatedAt        time.Time `json:"created_at"`
	UpdatedAt        time.Time `json:"updated_at"`
}

func (b BankProfile) ToBankProfileResponse() response.BankProfileResponse {
	if b.ID == "" {
		return response.BankProfileResponse{}
	}

	var isUploadedInfo bool = b.PayosApiKey != ""

	return response.BankProfileResponse{
		ID:             b.ID,
		Owner:          b.Owner,
		BankOrg:        b.BankOrg,
		BankCode:       b.BankCode,
		OwnerName:      b.OwnerName,
		IsUploadedInfo: isUploadedInfo,
	}
}
