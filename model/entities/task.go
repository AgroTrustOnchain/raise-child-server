package entities

import "time"

type Task struct {
	ID                string    `json:"id"`
	IsChildTask       bool      `json:"is_child_task"`
	ChildTaskDetailID *string   `json:"child_task_detail_id"`
	CreatedBy         string    `json:"created_by"`
	AssignedProfileID *string   `json:"assigned_profile_id"`
	AssignedStaff     *string   `json:"assgined_staff"`
	Region            string    `json:"region"`
	Description       string    `json:"description"`
	StartPeriod       time.Time `json:"start_period"`
	EndPeriod         time.Time `json:"end_period"`
	CreatedAt         time.Time `json:"created_at"`
	UpdatedAt         time.Time `json:"updated_at"`
}

type TaskV2 struct {
	ID                  string    `json:"id"`
	IsChildTask         bool      `json:"is_child_task"`
	ChildTaskDetailID   *string   `json:"child_task_detail_id"`
	CreatedBy           string    `json:"created_by"`
	AssignedProfileID   *string   `json:"assigned_profile_id"`
	AssignedStaff       *string   `json:"assgined_staff"`
	ReviewProfileStatus string    `json:"review_profile_status"`
	ReviewedBy          *string   `json:"reviewed_by"`
	Region              string    `json:"region"`
	Description         string    `json:"description"`
	StartPeriod         time.Time `json:"start_period"`
	EndPeriod           time.Time `json:"end_period"`
	CreatedAt           time.Time `json:"created_at"`
	UpdatedAt           time.Time `json:"updated_at"`
	IsSubmitted         bool      `json:"is_submitted"`
}

type ChildTaskDetail struct {
	ID        string    `json:"id"`
	ChildID   string    `json:"child_id"`
	Purpose   string    `json:"purpose"`
	Target    string    `json:"target"`
	CreatedAt time.Time `json:"created_at"`
}

// type TaskProof struct {
// 	ID             string    `json:"id"`
// 	TaskID         string    `json:"task_id"`
// 	Description    string    `json:"description"`
// 	ActorProfileID string    `json:"actor_profile_id"`
// 	ActorAddress   string    `json:"actor_address"`
// 	ImageBlobID    string    `json:"image_blob_id"`
// 	ReviewedBy     *string   `json:"reviewed_by"`
// 	AIEvaluation   string    `json:"ai_evaluation"`
// 	AIReason       string    `json:"ai_reason"`
// 	ReviewStatus   string    `json:"review_status"`
// 	RawSubmitDate  string    `json:"raw_submit_date"`
// 	CreatedAt      time.Time `json:"created_at"`
// 	UpdatedAt      time.Time `json:"updated_at"`
// }

type TaskProof struct {
	ID                    string    `json:"id"`
	TaskID                string    `json:"task_id"`
	Description           string    `json:"description"`
	ActorProfileID        string    `json:"actor_profile_id"`
	ActorAddress          string    `json:"actor_address"`
	ImageWalrusBlobID     string    `json:"image_blob_id"`
	ImageCloudinaryBlobID string    `json:"image_cloudinary_blob_id"`
	ReviewedBy            *string   `json:"reviewed_by"`
	AIEvaluation          string    `json:"ai_evaluation"`
	AIReason              string    `json:"ai_reason"`
	ReviewStatus          string    `json:"review_status"`
	RawSubmitDate         string    `json:"raw_submit_date"`
	CreatedAt             time.Time `json:"created_at"`
	UpdatedAt             time.Time `json:"updated_at"`
}

type TaskProofWithIsChildTask struct {
	ID                    string    `json:"id"`
	TaskID                string    `json:"task_id"`
	Description           string    `json:"description"`
	ActorProfileID        string    `json:"actor_profile_id"`
	ActorAddress          string    `json:"actor_address"`
	ImageBlobID           string    `json:"image_blob_id"`
	ImageCloudinaryBlobID string    `json:"image_cloudinary_blob_id"`
	ReviewedBy            *string   `json:"reviewed_by"`
	AIEvaluation          string    `json:"ai_evaluation"`
	AIReason              string    `json:"ai_reason"`
	ReviewStatus          string    `json:"review_status"`
	RawSubmitDate         string    `json:"raw_submit_date"`
	CreatedAt             time.Time `json:"created_at"`
	UpdatedAt             time.Time `json:"updated_at"`
	IsChildTask           bool      `json:"is_child_task"`
}
