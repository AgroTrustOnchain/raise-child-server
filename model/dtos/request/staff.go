package request

type GetStaffsRequest struct {
	Keyword     string `json:"keyword"`
	Role        string `json:"role"`
	Region      string `json:"region"`
	Gender      string `json:"gender"`
	YearOfBirth *int   `json:"year_of_birth"`
	SortOrder   string `json:"sort_order"`
	PageSize    int    `json:"page_size"`
	Page        int    `json:"page"`
}
