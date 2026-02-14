package request

type GetGiftsRequest struct {
	Keyword   string `json:"keyword"`
	Status    string `json:"status"`
	Category  string `json:"category"`
	SortOrder string `json:"sort_order"`
	PageSize  int    `json:"page_size"`
	Page      int    `json:"page"`
}

type CreateGiftRequest struct {
	Recipient       string `json:"recipient" validate:"required"`
	TrackingCode    string `json:"tracking_code" validate:"required"`
	Carrier         string `json:"carrier" validate:"required"`
	GiftImageBlobID string `json:"gift_image_blob_id" validate:"required"`
	Category        string `json:"category" validate:"required"`
	GiftValue       int64  `json:"gift_value" validate:"required,min=2000"`
	Description     string `json:"description" validate:"required"`
	Message         string `json:"message" validate:"required"`
}

type ConfirmReceiveGiftRequest struct {
	DeliveredImageBlobID string `json:"delivered_image_blob_id" validate:"required"`
}

type CancelGiftRequest struct {
	CancelReason string `json:"cancel_reason" validate:"required"`
}
