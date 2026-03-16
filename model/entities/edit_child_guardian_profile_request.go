package entities

import "time"

type EditChildGuardianProfileRequest struct {
	ID                 string    `json:"id"`
	Purpose            string    `json:"purpose"`  // "Edit" || "Upload"
	Guardian           string    `json:"guardian"` // "First" || "Second"
	ChildID            string    `json:"child_id"`
	Region             string    `json:"region"`
	ActorProfileID     string    `json:"actor_profile_id"`
	ActorAddress       string    `json:"actor_address"`
	FullName           string    `json:"full_name"`
	PhoneNumber        string    `json:"phone_number"`
	Relation           string    `json:"relation"`
	IdentityCardBlobID string    `json:"identity_card_blob_id"`
	AIEvaluation       string    `json:"ai_evaluation"`
	ReviewStatus       string    `json:"review_status"`
	ReviewedBy         *string   `json:"reviewed_by"`
	CreatedAt          time.Time `json:"created_at"`
	UpdatedAt          time.Time `json:"updated_at"`
}

type EditChildGuardianProfilePurpose string

const (
	EditGuardianPurpose   EditChildGuardianProfilePurpose = "Edit"
	UploadGuardianPurpose EditChildGuardianProfilePurpose = "Upload"
)
