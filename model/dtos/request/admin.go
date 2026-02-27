package request

type GetAdminsRequest struct {
	Keyword      string `json:"keyword"`
	Gender       string `json:"gender"`
	YearOfBirth  *int   `json:"year_of_birth"`
	SortCriteria string `json:"sort_criteria"`
	SortOrder    string `json:"sort_order"`
	PageSize     int    `json:"page_size"`
	Page         int    `json:"page"`
}
