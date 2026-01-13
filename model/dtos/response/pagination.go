package response

type PaginationDataResponse struct {
	Data       interface{} `json:"data"`
	Page       int         `json:"page"`
	TotalPages int         `json:"total_pages"`
}
