package entities

type DaoStruct struct {
	ID              ID     `json:"id"`
	MinApprovedRate string `json:"min_approved_rate"`
	MinVoters       string `json:"min_voters"`
}
