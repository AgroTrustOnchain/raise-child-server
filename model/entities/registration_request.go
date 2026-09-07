package entities

import (
	"raise-child/model/dtos/response"
	"time"
)

type AdminRegistrationRequest struct {
	ID                   string    `json:"id"`
	IdentityCode         string    `json:"identity_code"`
	IdentityCardBlobID   string    `json:"identity_card_blob_id"`
	AvatarBlobID         string    `json:"avatar_blob_id"`
	FirstName            string    `json:"first_name"`
	LastName             string    `json:"last_name"`
	Gender               string    `json:"gender"`
	DateOfBirth          string    `json:"date_of_birth"`
	PhoneNumber          string    `json:"phone_number"`
	Email                string    `json:"email"`
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

type VolunteerRegistrationRequest struct {
	Region string `json:"region"`
	AdminRegistrationRequest
}

type LocalLeaderRegistrationRequest struct {
	CenterAddress     string `json:"center_address"`
	CenterPhoneNumber string `json:"center_phone_number"`
	CenterImageBlobID string `json:"center_image_blob_id"`
	VolunteerRegistrationRequest
}

type RegistrationRequest struct {
	ID                 string    `json:"id"`
	ProfileID          string    `json:"profile_id"`
	RegisterRole       string    `json:"register_role"`
	IdentityCode       string    `json:"identity_code"`
	IdentityCardBlobID string    `json:"identity_card_blob_id"`
	AvatarBlobID       string    `json:"avatar_blob_id"`
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
	OnchainID          *string   `json:"on_chain_id"`
	CreatedAt          time.Time `json:"created_at"`
	UpdatedAt          time.Time `json:"updated_at"`
	ClosedAt           time.Time `json:"closed_at"`
}

type GetRegistrationRequestWithDetail struct {
	RegisterRole      string
	IdentityCode      string
	AvatarBlobID      string
	Region            string
	FirstName         string
	LastName          string
	Gender            string
	DateOfBirth       string
	PhoneNumber       string
	Email             string
	Status            string
	IsConfirmRegister bool
	CreatedBy         string
}

type RegistrationForm struct {
	ID                   string    `json:"id"`
	Sub                  string    `json:"sub"`
	RegisterRole         string    `json:"register_role"`
	IdentityCode         string    `json:"identity_code"`
	IdentityCardBlobID   string    `json:"identity_card_blob_id"`
	AvatarBlobID         string    `json:"avatar_blob_id"`
	Region               string    `json:"region"`
	FirstName            string    `json:"first_name"`
	LastName             string    `json:"last_name"`
	Gender               string    `json:"gender"`
	DateOfBirth          string    `json:"date_of_birth"`
	PhoneNumber          string    `json:"phone_number"`
	Email                string    `json:"email"`
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

func (r *RegistrationRequest) ToRegistrationRequestResponse() response.RegistrationRequestResponse {
	if r == nil {
		return response.RegistrationRequestResponse{}
	}

	return response.RegistrationRequestResponse{
		ID:                 r.ID,
		ProfileID:          r.ProfileID,
		RegisterRole:       r.RegisterRole,
		IdentityCode:       r.IdentityCode,
		AvatarWalrusBlobID: r.AvatarBlobID,
		Region:             r.Region,
		FirstName:          r.FirstName,
		LastName:           r.LastName,
		Gender:             r.Gender,
		DateOfBirth:        r.DateOfBirth,
		PhoneNumber:        r.PhoneNumber,
		Email:              r.Email,
		Approvers:          r.Approvers,
		Refusers:           r.Refusers,
		RefuseReasons:      r.RefuseReasons,
		Status:             r.Status,            // e.g. "Pending", "Approved", "Refused"
		IsConfirmRegister:  r.IsConfirmRegister, // Default as false, when status aprroved, user clicks to update to true, call smart contract to register role
		CreatedBy:          r.CreatedBy,
		CreatedAt:          r.CreatedAt,
		UpdatedAt:          r.UpdatedAt,
		ClosedAt:           r.ClosedAt,
	}
}
