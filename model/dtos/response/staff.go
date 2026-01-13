package response

import (
	"time"
)

type StaffResponse struct {
	ID           string             `json:"id"`
	User         string             `json:"user"`
	Role         string             `json:"string"`
	IdentityCode string             `json:"identity_code"`
	AvatarBlobID string             `json:"avatar_blob_id"`
	Region       string             `json:"region"`
	FirstName    string             `json:"first_name"`
	LastName     string             `json:"last_name"`
	Gender       string             `json:"gender"`
	DateOfBirth  string             `json:"date_of_birth"`
	PhoneNumber  string             `json:"phone_number"`
	Email        string             `json:"email"`
	Nfts         []StaffNftResponse `json:"nfts"`
	UploadedAt   time.Time          `json:"uploaded_at"`
}
