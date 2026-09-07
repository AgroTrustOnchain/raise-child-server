package entities

import (
	"raise-child/model/dtos/response"
	"raise-child/util"
	"strconv"
)

type Gift struct {
	ID                   ID     `json:"id"`
	Sender               string `json:"sender"`
	Recipient            string `json:"recipient"`
	IsForChild           bool   `json:"is_for_child"`
	TrackingCode         string `json:"tracking_code"`
	Carrier              string `json:"carrier"`
	GiftImageBlobID      string `json:"gift_image_blob_id"`
	Status               string `json:"status"`
	Category             string `json:"category"`
	Description          string `json:"description"`
	Message              string `json:"message"`
	CancelReason         string `json:"cancel_reason"`
	DeliveredImageBlobID string `json:"delivered_image_blob_id"`
	UploadedAt           string `json:"uploaded_at"`
	UpdatedAt            string `json:"updated_at"`
	DeliveredAt          string `json:"delivered_at"`
	ConfirmRecievedBy    string `json:"confirm_recieved_by"`
}

func (g Gift) ToGiftResponse() response.GiftResponse {
	if g.ID.ID == "" {
		return response.GiftResponse{}
	}

	uploadedAt, _ := strconv.ParseInt(g.UploadedAt, 10, 64)
	updatedAt, _ := strconv.ParseInt(g.UpdatedAt, 10, 64)
	deliveredAt, _ := strconv.ParseInt(g.DeliveredAt, 10, 64)

	return response.GiftResponse{
		ID:                   g.ID.ID,
		Sender:               g.Sender,
		Recipient:            g.Recipient,
		IsForChild:           g.IsForChild,
		TrackingCode:         g.TrackingCode,
		Carrier:              g.Carrier,
		GiftImageBlobID:      g.GiftImageBlobID,
		Status:               g.Status,
		Category:             g.Category,
		Description:          g.Description,
		Message:              g.Message,
		CancelReason:         g.CancelReason,
		DeliveredImageBlobID: g.DeliveredImageBlobID,
		UploadedAt:           util.MilliSecToTime(uploadedAt),
		UpdatedAt:            util.MilliSecToTime(updatedAt),
		DeliveredAt:          util.MilliSecToTime(deliveredAt),
		ConfirmRecievedBy:    g.ConfirmRecievedBy,
	}
}
