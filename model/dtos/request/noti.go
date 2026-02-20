package request

type GetNotisRequest struct {
	PageSize int `json:"page_size"`
	Page     int `json:"page"`
}
