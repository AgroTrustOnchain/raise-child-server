package entities

import "time"

type UploadChildRequest struct {
	ID              string    `json:"id"`
	ProfileID       string    `json:"profile_id"`
	IdentityCode    string    `json:"identity_code"`
	AvatarBlobId    string    `json:"avatar_blob_id" validate:"required"`
	Region          string    `json:"region" validate:"required"`
	FirstName       string    `json:"first_name" validate:"required"`
	LastName        string    `json:"last_name" validate:"required"`
	Gender          string    `json:"gender" validate:"required"`
	DateOfBirth     string    `json:"date_of_birth" validate:"required"`
	Approvers       []string  `json:"approvers"`
	Refusers        []string  `json:"refusers"`
	RefuseReasons   []string  `json:"refuse_reasons"`
	Status          string    `json:"status"`            // e.g. "Pending", "Approved", "Refused"
	IsConfirmUpload bool      `json:"is_confirm_upload"` // Default as false, when status aprroved, user clicks to update to true, call smart contract to register role
	CreatedBy       string    `json:"created_by"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
	ClosedAt        time.Time `json:"closed_at"`
}
