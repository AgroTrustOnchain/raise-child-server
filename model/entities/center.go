package entities

import (
	"raise-child/model/dtos/response"
	"raise-child/util"
	"strconv"
	"time"
)

type CenterRequest struct {
	ID                   string    `json:"id"`
	ProfileID            string    `json:"profile_id"`
	Region               string    `json:"region"`
	Address              string    `json:"address"`
	PhoneNumber          string    `json:"phone_number"`
	ImageBlobID          string    `json:"image_blob_id"`
	Approvers            []string  `json:"approvers"`
	Refusers             []string  `json:"refusers"`
	RefuseReasons        []string  `json:"refuse_reasons"`
	Status               string    `json:"status"` // e.g. "Pending", "Approved", "Refused"
	IsAvailableToConfirm bool      `jsosn:"is_available_to_confirm"`
	IsConfirmRegister    bool      `json:"is_confirm_register"` // Default as false, when status aprroved, user clicks to update to true, call smart contract to register role
	CreatedBy            string    `json:"created_by"`
	CreatedAt            time.Time `json:"created_at"`
	UpdatedAt            time.Time `json:"updated_at"`
	ClosedAt             time.Time `json:"closed_at"`
}

type Center struct {
	ID                ID       `json:"id"`
	Region            string   `json:"region"`
	CenterAddress     string   `json:"center_address"`
	CenterPhoneNumber string   `json:"center_phone_number"`
	ImageBlobIDs      []string `json:"image_blob_ids"`
	ChildIDs          []string `json:"child_ids"`
	Gifts             []string `json:"gifts"`
	AllGifts          []string `json:"all_gifts"`
	TaskProofs        []string `json:"task_proofs"`
	UploadedAt        string   `json:"uploaded_at"`
	UpdatedAt         string   `json:"updated_at"`
}

func (c Center) ToCenterResponse() response.CenterResponse {
	if c.ID.ID == "" {
		return response.CenterResponse{}
	}

	uploadedAt, _ := strconv.ParseInt(c.UploadedAt, 10, 64)
	updatedAt, _ := strconv.ParseInt(c.UpdatedAt, 10, 64)

	return response.CenterResponse{
		ID:                c.ID.ID,
		Region:            c.Region,
		CenterAddress:     c.CenterAddress,
		CenterPhoneNumber: c.CenterPhoneNumber,
		ImageBlobIDs:      c.ImageBlobIDs,
		Gifts:             c.Gifts,
		UploadedAt:        util.MilliSecToTime(uploadedAt),
		UpdatedAt:         util.MilliSecToTime(updatedAt),
	}
}
