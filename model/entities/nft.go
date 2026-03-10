package entities

import (
	"raise-child/model/dtos/response"
	"raise-child/util"
	"strconv"
	"time"
)

type StaffNft struct {
	ID                 ID     `json:"id"`
	Owner              string `json:"owner"`
	Role               string `json:"role"`
	IdentityCode       string `json:"identity_code"`
	IdentityCardBlobID string `json:"identity_card_blob_id"`
	AvatarBlobID       string `json:"avatar_blob_id"`
	Region             string `json:"region"`
	FirstName          string `json:"first_name"`
	LastName           string `json:"last_name"`
	Gender             string `json:"gender"`
	DateOfBirth        string `json:"date_of_birth"`
	PhoneNumber        string `json:"phone_number"`
	Email              string `json:"email"`
	UploadedAt         string `json:"uploaded_at"`
	Name               string `json:"name"`
	Url                string `json:"url"`
}

type AdminNft struct {
	ID                 ID     `json:"id"`
	Owner              string `json:"owner"`
	IdentityCode       string `json:"identity_code"`
	IdentityCardBlobID string `json:"identity_card_blob_id"`
	AvatarBlobID       string `json:"avatar_blob_id"`
	FirstName          string `json:"first_name"`
	LastName           string `json:"last_name"`
	Gender             string `json:"gender"`
	DateOfBirth        string `json:"date_of_birth"`
	PhoneNumber        string `json:"phone_number"`
	Email              string `json:"email"`
	UploadedAt         string `json:"uploaded_at"`
	Name               string `json:"name"`
	Url                string `json:"url"`
}

type NftProfile struct {
	ID        string
	ProfileID string
	NftID     *string
	Role      string
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (s StaffNft) ToStaffNftResponse() response.StaffNftResponse {
	if s.ID.ID == "" {
		return response.StaffNftResponse{}
	}

	uploadedAt, _ := strconv.ParseInt(s.UploadedAt, 10, 64)
	return response.StaffNftResponse{
		ID:                 s.Owner,
		Role:               s.Role,
		IdentityCode:       s.IdentityCode,
		IdentityCardBlobID: s.IdentityCardBlobID,
		AvatarBlobID:       s.AvatarBlobID,
		Region:             s.Region,
		FirstName:          s.FirstName,
		LastName:           s.LastName,
		Gender:             s.Gender,
		DateOfBirth:        util.RawDateToTime(s.DateOfBirth),
		PhoneNumber:        s.PhoneNumber,
		Email:              s.Email,
		UploadedAt:         util.MilliSecToTime(uploadedAt),
		Name:               s.Name,
		Url:                s.Url,
	}
}

func (a AdminNft) ToAdminNftResponse() response.AdminNftResponse {
	if a.ID.ID == "" {
		return response.AdminNftResponse{}
	}

	uploadedAt, _ := strconv.ParseInt(a.UploadedAt, 10, 64)
	return response.AdminNftResponse{
		ID:                 a.Owner,
		IdentityCode:       a.IdentityCode,
		IdentityCardBlobID: a.IdentityCardBlobID,
		AvatarBlobID:       a.AvatarBlobID,
		FirstName:          a.FirstName,
		LastName:           a.LastName,
		Gender:             a.Gender,
		DateOfBirth:        util.RawDateToTime(a.DateOfBirth),
		PhoneNumber:        a.PhoneNumber,
		Email:              a.Email,
		UploadedAt:         util.MilliSecToTime(uploadedAt),
		Name:               a.Name,
		Url:                a.Url,
	}
}
