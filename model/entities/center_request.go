package entities

import "time"

type CenterRequest struct {
	ID                   string    `json:"id"`
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
