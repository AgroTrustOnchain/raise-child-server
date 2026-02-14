package request

type GetDonorsRequest struct {
	Keyword  string `json:"keyword"`
	Gender   string `json:"gender"`
	PageSize int    `json:"page_size"`
	Page     int    `json:"page"`
}
