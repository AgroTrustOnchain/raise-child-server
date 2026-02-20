package request

type GetCenterRequests struct {
	Keyword              string `json:"keyword"`
	Region               string `json:"region"`
	Status               string `json:"status"`
	IsClosed             *bool  `json:"is_closed"`
	IsAvailableToConfirm *bool  `json:"is_available_to_confirm"`
	SortOrder            string `json:"sort_order"`
	PageSize             int    `json:"page_size"`
	Page                 int    `json:"page"`
}

type CreateCenterRequest struct {
	Region      string `json:"region" validate:"required"`
	Address     string `json:"address" validate:"required"`
	PhoneNumber string `json:"phone_number" validate:"required"`
	ImageBlobID string `json:"image_blob_id" validate:"required"`
}

type EditStaffNumbersToCenterRequest struct {
	MinStaffNumber *int `json:"min_staff_number"`
}
