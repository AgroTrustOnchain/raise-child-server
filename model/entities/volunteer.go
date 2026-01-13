package entities

type Volunteer struct {
	ID             string `json:"id"`
	IdentityCode   string `json:"identity_code"`
	FirstName      string `json:"first_name"`
	LastName       string `json:"last_name"`
	Gender         string `json:"gender"`
	AvatarBlobId   string `json:"avatar_blob_id"`
	VolunteerNftId string `json:"volunteer_nft_id"`
	UploadedAt     int64  `json:"uploaded_at"`
	UpdatedAt      int64  `json:"updated_at"`
}
