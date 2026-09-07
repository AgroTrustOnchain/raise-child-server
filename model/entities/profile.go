package entities

import (
	"raise-child/model/dtos/response"
	"time"
)

type Profile struct {
	ID            string    `json:"id"`
	Salt          string    `json:"salt"`
	Status        string    `json:"status"`
	IdentityCode  *string   `json:"identity_code"`
	FirstName     *string   `json:"first_name"`
	LastName      *string   `json:"last_name"`
	Gender        *string   `json:"gender"`
	DateOfBirth   *string   `json:"date_of_birth"`
	PhoneNumber   *string   `json:"phone_number"`
	Email         *string   `json:"email"`
	Token         *string   `json:"token"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
	WalletAddress *string   `json:"wallet_address"`
}

func (p Profile) ToPersonalProfile() response.PersonalProfileResponse {
	if p.ID == "" {
		return response.PersonalProfileResponse{}
	}

	return response.PersonalProfileResponse{
		ID:           p.ID,
		IdentityCode: *p.IdentityCode,
		FirstName:    *p.FirstName,
		LastName:     *p.LastName,
		Gender:       *p.Gender,
		DateOfBirth:  *p.DateOfBirth,
		PhoneNumber:  *p.PhoneNumber,
		Email:        *p.Email,
	}
}
