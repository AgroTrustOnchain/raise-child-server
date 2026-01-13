package entities

type Manage struct {
	ID             string   `json:"id"`
	Admin          string   `json:"admin"`
	ChildIds       []string `json:"child_ids"`
	VolunteerIds   []string `json:"volunteer_ids"`
	LocalLeaderIds []string `json:"local_leader_ids"`
	SponsorIds     []string `json:"sponsor_ids"`
	SponsorNfts    []string `json:"sponsor_nfts"`
}
