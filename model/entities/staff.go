package entities

import (
	"raise-child/model/dtos/response"
	"raise-child/util"
	"strconv"
)

type Staff struct {
	ID           ID     `json:"id"`
	User         string `json:"user"`
	Role         string `json:"role"`
	IdentityCode string `json:"identity_code"`
	AvatarBlobID string `json:"avatar_blob_id"`
	Region       string `json:"region"`
	FirstName    string `json:"first_name"`
	LastName     string `json:"last_name"`
	Gender       string `json:"gender"`
	DateOfBirth  string `json:"date_of_birth"`
	PhoneNumber  string `json:"phone_number"`
	Email        string `json:"email"`
	UploadedAt   string `json:"uploaded_at"`
}

func (s Staff) ToStaffResponse() response.StaffResponse {
	if s.ID.ID == "" {
		return response.StaffResponse{}
	}

	uploadedAt, _ := strconv.ParseInt(s.UploadedAt, 10, 64)

	return response.StaffResponse{
		ID:           s.ID.ID,
		User:         s.User,
		Role:         s.Role,
		IdentityCode: s.IdentityCode,
		AvatarBlobID: s.AvatarBlobID,
		Region:       s.Region,
		FirstName:    s.FirstName,
		LastName:     s.LastName,
		Gender:       s.Gender,
		DateOfBirth:  s.DateOfBirth,
		PhoneNumber:  s.PhoneNumber,
		Email:        s.Email,
		UploadedAt:   util.MilliSecToTime(uploadedAt),
	}
}
