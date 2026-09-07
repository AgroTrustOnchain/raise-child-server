package request

type GetCenterRequests struct {
	Keyword   string `form:"keyword"`
	Region    string `form:"region"`
	Status    string `form:"status"`
	IsClosed  *bool  `form:"is_closed"`
	SortOrder string `form:"sort_order"`
	PageSize  int    `form:"page_size"`
	Page      int    `form:"page"`
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
