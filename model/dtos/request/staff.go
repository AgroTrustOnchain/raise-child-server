package request

type GetStaffsRequest struct {
	Keyword     string `form:"keyword"`
	Role        string `form:"role"`
	Region      string `form:"region"`
	Gender      string `form:"gender"`
	YearOfBirth *int   `form:"year_of_birth"`
	SortOrder   string `form:"sort_order"`
	PageSize    int    `form:"page_size"`
	Page        int    `form:"page"`
}
