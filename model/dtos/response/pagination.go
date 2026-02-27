package response

type PaginationDataResponse struct {
	Data       interface{} `json:"data"`
	Amount     int         `json:"amount"`
	Page       int         `json:"page"`
	TotalPages int         `json:"total_pages"`
}
