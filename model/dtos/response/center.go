package response

import "time"

type CenterResponse struct {
	ID                string    `json:"id"`
	Region            string    `json:"region"`
	CenterAddress     string    `json:"center_address"`
	CenterPhoneNumber string    `json:"center_phone_number"`
	ImageBlobIDs      []string  `json:"image_blob_ids"`
	Gifts             []string  `json:"gifts"`
	UploadedAt        time.Time `json:"uploaded_at"`
	UpdatedAt         time.Time `json:"updated_at"`
}

type CenterCardMinimumResponse struct {
	ID                string    `json:"id"`
	Region            string    `json:"region"`
	CenterAddress     string    `json:"center_address"`
	CenterPhoneNumber string    `json:"center_phone_number"`
	UploadedAt        time.Time `json:"uploaded_at"`
	UpdatedAt         time.Time `json:"updated_at"`
}
