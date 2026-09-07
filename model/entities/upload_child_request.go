package entities

import (
	"raise-child/model/dtos/response"
	"time"
)

type UploadChildRequest struct {
	ID                     string                `json:"id"`
	ProfileID              string                `json:"profile_id"`
	IdentityCode           string                `json:"identity_code"`
	BirthCertificateBlobID string                `json:"birth_certificate_blob_id"`
	AvatarBlobId           string                `json:"avatar_blob_id" validate:"required"`
	HomeBlobID             string                `json:"home_blob_id"`
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

type ChildGuardianProfile struct {
	FullName           string `json:"full_name"`
	PhoneNumber        string `json:"phone_number"`
	Relation           string `json:"relation"`
	IdentityCardBlobID string `json:"identity_card_blob_id"`
}

func (u *UploadChildRequest) ToUploadChildRequestResponse() response.UploadChildRequestResponse {
	if u == nil {
		return response.UploadChildRequestResponse{}
	}

	var secondGuardianProfile *response.ChildGuardianProfile
	if u.SecondGuardianProfile != nil {
		secondGuardianProfile = &response.ChildGuardianProfile{
			FullName:    u.SecondGuardianProfile.FullName,
			PhoneNumber: u.SecondGuardianProfile.PhoneNumber,
			Relation:    u.SecondGuardianProfile.Relation,
		}
	}

	return response.UploadChildRequestResponse{
		ID:                 u.ID,
		ProfileID:          u.ProfileID,
		IdentityCode:       u.IdentityCode,
		AvatarWalrusBlobID: u.AvatarBlobId,
		HomeWalrusBlobID:   u.HomeBlobID,
		Region:             u.Region,
		FirstName:          u.FirstName,
		LastName:           u.LastName,
		Gender:             u.Gender,
		DateOfBirth:        u.DateOfBirth,
		HomeAddress:        u.HomeAddress,
		FirstGuardianProfile: response.ChildGuardianProfile{
			FullName:    u.FirstGuardianProfile.FullName,
			PhoneNumber: u.FirstGuardianProfile.PhoneNumber,
			Relation:    u.FirstGuardianProfile.Relation,
		},
		SecondGuardianProfile: secondGuardianProfile,
		Status:                u.Status,
		IsConfirmUpload:       u.IsConfirmUpload,
		CreatedBy:             u.CreatedBy,
		ReviewedBy:            u.ReviewedBy,
		OnchainID:             u.OnchainID,
		CreatedAt:             u.CreatedAt,
		UpdatedAt:             u.UpdatedAt,
	}
}
