package request

type GetSponsorsRequest struct {
	Keyword string `json:"keyword"`
	Gender  string `json:"gender"`
	Page    int    `json:"page"`
}
