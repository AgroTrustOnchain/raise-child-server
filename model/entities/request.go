package entities

import "time"

type RequestType string
type RegisterRole string
type RequestStatus string

const (
	LOCAL_LEADER_REGISTRATION RequestType = "Local Leader Registration"
	VOLUNTEER_REGISTRATION    RequestType = "Volunteer Registration"
)

// Register Role
const (
	LOCAL_LEADER_ROLE     RegisterRole = "Local Leader"
	VOLUNTEER_LEADER_ROLE RegisterRole = "Volunteer"
)

// Request status
const (
	PENDING_REQUEST_STATUS  RequestStatus = "Pending"
	APPROVED_REQUEST_STATUS RequestStatus = "Approved"
	REFUSED_REQUEST_STATUS  RequestStatus = "Refused"
)

type Request struct {
	ID          string      `json:"id"`
	RequestType RequestType `json:"request_type"`
	Aprrovers   []string    `json:"approvers"`
	Refusers    []string    `json:"refusers"`
	CreatedBy   string      `json:"created_by"`
	CreatedAt   time.Time   `json:"created_at"`
	UpdatedAt   time.Time   `json:"updated_at"`
	ClosedAt    time.Time   `json:"closed_at"`
	IsClosed    bool        `json:"is_closed"`
}

type RegistrationRequest struct {
	ID                 string    `json:"id"`
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
	Aprrovers          []string  `json:"approvers"`
	Refusers           []string  `json:"refusers"`
	RefuseReasons      []string  `json:"refuse_reasons"`
	Status             string    `json:"status"`              // e.g. "Pending", "Approved", "Refused"
	IsConfirmRegister  bool      `json:"is_confirm_register"` // Default as false, when status aprroved, user clicks to update to true, call smart contract to register role
	CreatedBy          string    `json:"created_by"`
	CreatedAt          time.Time `json:"created_at"`
	UpdatedAt          time.Time `json:"updated_at"`
	ClosedAt           time.Time `json:"closed_at"`
}

type UploadChildRequest struct {
	ID              string    `json:"id"`
	IdentityCode    string    `json:"identity_code"`
	AvatarBlobId    string    `json:"avatar_blob_id" validate:"required"`
	Region          string    `json:"region" validate:"required"`
	FirstName       string    `json:"first_name" validate:"required"`
	LastName        string    `json:"last_name" validate:"required"`
	Gender          string    `json:"gender" validate:"required"`
	DateOfBirth     string    `json:"date_of_birth" validate:"required"`
	Aprrovers       []string  `json:"approvers"`
	Refusers        []string  `json:"refusers"`
	RefuseReasons   []string  `json:"refuse_reasons"`
	Status          string    `json:"status"`            // e.g. "Pending", "Approved", "Refused"
	IsConfirmUpload bool      `json:"is_confirm_upload"` // Default as false, when status aprroved, user clicks to update to true, call smart contract to register role
	CreatedBy       string    `json:"created_by"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
	ClosedAt        time.Time `json:"closed_at"`
}
