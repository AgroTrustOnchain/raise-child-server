package request

type GetTasksRequest struct {
	Keyword       string `json:"keyword"`
	Region        string `json:"region"`
	Status        string `json:"status"`
	AssignedStaff string `json:"assgined_staff"`
	ReviewedBy    string `json:"reviewed_by"`
	SortOrder     string `json:"sort_order"`
	PageSize      int    `json:"page_size"`
	Page          int    `json:"page"`
}

type GetTaskProofsRequest struct {
	Keyword      string `json:"keyword"`
	Status       string `json:"status"`
	ActorAddress string `json:"actor_address"`
	ReviewedBy   string `json:"reviewed_by"`
	SortOrder    string `json:"sort_order"`
	PageSize     int    `json:"page_size"`
	Page         int    `json:"page"`
}

type CreateTaskRequest struct {
	Region      string `json:"region" validate:"required"`
	Description string `json:"description" validate:"required"`
	StartPeriod string `json:"start_period" validate:"required"`
	EndPeriod   string `json:"end_period" validate:"required"`
}

type SubmitTaskProofRequest struct {
	ImageBlobID string `json:"image_blob_id" validate:"required"`
}
