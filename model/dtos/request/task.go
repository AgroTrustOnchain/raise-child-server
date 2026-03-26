package request

type GetTasksRequest struct {
	Keyword       string `form:"keyword"`
	Region        string `form:"region"`
	Status        string `form:"status"`
	AssignedStaff string `form:"assgined_staff"`
	ReviewedBy    string `form:"reviewed_by"`
	SortOrder     string `form:"sort_order"`
	PageSize      int    `form:"page_size"`
	Page          int    `form:"page"`
}

type GetTaskProofsRequest struct {
	Keyword      string `form:"keyword"`
	Status       string `form:"status"`
	ActorAddress string `form:"actor_address"`
	ReviewedBy   string `form:"reviewed_by"`
	SortOrder    string `form:"sort_order"`
	PageSize     int    `form:"page_size"`
	Page         int    `form:"page"`
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
