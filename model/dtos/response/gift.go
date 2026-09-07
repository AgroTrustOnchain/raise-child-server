package response

import "time"

type GiftResponse struct {
	ID                   string    `json:"id"`
	Sender               string    `json:"sender"`
	Recipient            string    `json:"recipient"`
	IsForChild           bool      `json:"is_for_child"`
	TrackingCode         string    `json:"tracking_code"`
	Carrier              string    `json:"carrier"`
	GiftImageBlobID      string    `json:"gift_image_blob_id"`
	Status               string    `json:"status"`
	Category             string    `json:"category"`
	Description          string    `json:"description"`
	Message              string    `json:"message"`
	CancelReason         string    `json:"cancel_reason"`
	DeliveredImageBlobID string    `json:"delivered_image_blob_id"`
	UploadedAt           time.Time `json:"uploaded_at"`
	UpdatedAt            time.Time `json:"updated_at"`
	DeliveredAt          time.Time `json:"delivered_at"`
	ConfirmRecievedBy    string    `json:"confirm_recieved_by"`
}
