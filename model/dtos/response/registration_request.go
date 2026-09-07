package response

import "time"

type RegistrationRequestResponse struct {
	ID                 string    `json:"id"`
	ProfileID          string    `json:"profile_id"`
	RegisterRole       string    `json:"register_role"`
	IdentityCode       string    `json:"identity_code"`
	IdentityCardImgUrl string    `json:"identity_card_img_url"`
	AvatarWalrusBlobID string    `json:"avatar_blob_id"`
	Region             string    `json:"region"`
	FirstName          string    `json:"first_name"`
	LastName           string    `json:"last_name"`
	Gender             string    `json:"gender"`
	DateOfBirth        string    `json:"date_of_birth"`
	PhoneNumber        string    `json:"phone_number"`
	Email              string    `json:"email"`
	Approvers          []string  `json:"approvers"`
	Refusers           []string  `json:"refusers"`
	RefuseReasons      []string  `json:"refuse_reasons"`
	Status             string    `json:"status"`              // e.g. "Pending", "Approved", "Refused"
	IsConfirmRegister  bool      `json:"is_confirm_register"` // Default as false, when status aprroved, user clicks to update to true, call smart contract to register role
	CreatedBy          string    `json:"created_by"`
	CreatedAt          time.Time `json:"created_at"`
	UpdatedAt          time.Time `json:"updated_at"`
	ClosedAt           time.Time `json:"closed_at"`
}
