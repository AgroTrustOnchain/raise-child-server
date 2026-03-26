package request

type GetAdminsRequest struct {
	Keyword      string `form:"keyword"`
	Gender       string `form:"gender"`
	YearOfBirth  *int   `form:"year_of_birth"`
	SortCriteria string `form:"sort_criteria"`
	SortOrder    string `form:"sort_order"`
	PageSize     int    `form:"page_size"`
	Page         int    `form:"page"`
}
