package request

type EditDaoRequest struct {
	MinVoters *int `json:"min_voters"`
	MinRate   *int `json:"min_rate"`
}
