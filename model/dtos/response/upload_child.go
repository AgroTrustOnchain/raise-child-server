package response

import "time"

type ExtractChildUploadInfoResponse struct {
	IdentityCode           string `json:"identity_code" `
	FirstName              string `json:"first_name"`
	LastName               string `json:"last_name"`
	Gender                 string `json:"gender"`
	DateOfBirth            string `json:"date_of_birth"`
	FirstGuardianFullName  string `json:"first_guardian_full_name"`
	SecondGuardianFullName string `json:"second_guardian_full_name"`
}

type UploadChildRequestResponse struct {
	ID                     string                `json:"id"`
	ProfileID              string                `json:"profile_id"`
	IdentityCode           string                `json:"identity_code"`
	BirthCertificateImgUrl string                `json:"birth_certificate_img_url"`
	AvatarWalrusBlobID     string                `json:"avatar_blob_id"`
	HomeWalrusBlobID       string                `json:"home_blob_id"`
	Region                 string                `json:"region" validate:"required"`
	FirstName              string                `json:"first_name" validate:"required"`
	LastName               string                `json:"last_name" validate:"required"`
	Gender                 string                `json:"gender" validate:"required"`
	DateOfBirth            string                `json:"date_of_birth" validate:"required"`
	HomeAddress            string                `json:"home_address"`
	FirstGuardianProfile   ChildGuardianProfile  `json:"first_guardian_profile"`
	SecondGuardianProfile  *ChildGuardianProfile `json:"second_guardian_profile"`
	Status                 string                `json:"status"`
	IsConfirmUpload        bool                  `json:"is_confirm_upload"`
	CreatedBy              string                `json:"created_by"`
	ReviewedBy             *string               `json:"reviewed_by"`
	OnchainID              *string               `json:"on_chain_id"`
	CreatedAt              time.Time             `json:"created_at"`
	UpdatedAt              time.Time             `json:"updated_at"`
}
