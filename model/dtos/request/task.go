package request

type GetTasksRequest struct {
	Keyword       string `form:"keyword"`
	Region        string `form:"region"`
	Status        string `form:"status"`
	AssignedStaff string `form:"assgined_staff"`
	SortOrder     string `form:"sort_order"`
	PageSize      int    `form:"page_size"`
	Page          int    `form:"page"`
}

type GetTaskProofsRequest struct {
	Keyword      string `form:"keyword"`
	IsChildTask  *bool  `form:"is_child_task"`
	Status       string `form:"status"`
	Region       string `json:"region"`
	ActorAddress string `form:"actor_address"`
	ReviewedBy   string `form:"reviewed_by"`
	SortOrder    string `form:"sort_order"`
	PageSize     int    `form:"page_size"`
	Page         int    `form:"page"`
}

type CreateTaskRequest struct {
	Region      string  `json:"region" validate:"required"`
	IsChildTask *bool   `json:"is_child_task"`
	ChildID     *string `json:"child_id"`
	NeedID      *string `json:"need_id"`
	Description string  `json:"description" validate:"required"`
	StartPeriod string  `json:"start_period" validate:"required"`
	EndPeriod   string  `json:"end_period" validate:"required"`
}

type SubmitTaskProofRequest struct {
	ImageCloudinaryBlobID string `json:"image_cloudinary_blob_id" validate:"required"`
	ImageUrl              string `json:"image_url" validate:"required"`
}
