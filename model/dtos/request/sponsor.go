package request

type GetDonorsRequest struct {
	Keyword  string `form:"keyword"`
	Gender   string `form:"gender"`
	PageSize int    `form:"page_size"`
	Page     int    `form:"page"`
}
