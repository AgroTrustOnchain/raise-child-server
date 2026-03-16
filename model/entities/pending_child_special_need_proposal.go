package entities

import "time"

type PendingChildSpecialNeedProposal struct {
	ID             string    `json:"id"`
	ChildID        string    `json:"child_id"`
	Region         string    `json:"region"`
	ActorProfileID string    `json:"actor_profile_id"`
	ActorAddress   string    `json:"actor_address"`
	Target         int64     `json:"target"`
	Description    string    `json:"description"`
	ProofBlobID    *string   `json:"proof_blob_id"`
	AIEvaluation   string    `json:"ai_evaluation"`
	ReviewStatus   string    `json:"review_status"`
	ReviewedBy     *string   `json:"reviewed_by"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}
