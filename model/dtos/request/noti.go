package request

type GetNotisRequest struct {
	PageSize int `form:"page_size"`
	Page     int `form:"page"`
}
