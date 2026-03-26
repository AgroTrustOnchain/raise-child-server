package response

import "time"

type StaffNftResponse struct {
	ID                 string    `json:"id"`
	Owner              string    `json:"owner"`
	Role               string    `json:"role"`
	IdentityCode       string    `json:"identity_code"`
	IdentityCardBlobID string    `json:"identity_card_blob_id"`
	AvatarBlobID       string    `json:"avatar_blob_id"`
	Region             string    `json:"region"`
	FirstName          string    `json:"first_name"`
	LastName           string    `json:"last_name"`
	Gender             string    `json:"gender"`
	DateOfBirth        time.Time `json:"date_of_birth"`
	PhoneNumber        string    `json:"phone_number"`
	Email              string    `json:"email"`
	UploadedAt         time.Time `json:"uploaded_at"`
	Name               string    `json:"name"`
	Url                string    `json:"url"`
}

type AdminNftResponse struct {
	ID                 string    `json:"id"`
	IdentityCode       string    `json:"identity_code"`
	IdentityCardBlobID string    `json:"identity_card_blob_id"`
	AvatarBlobID       string    `json:"avatar_blob_id"`
	FirstName          string    `json:"first_name"`
	LastName           string    `json:"last_name"`
	Gender             string    `json:"gender"`
	DateOfBirth        time.Time `json:"date_of_birth"`
	PhoneNumber        string    `json:"phone_number"`
	Email              string    `json:"email"`
	UploadedAt         time.Time `json:"uploaded_at"`
	Name               string    `json:"name"`
	Url                string    `json:"url"`
}
