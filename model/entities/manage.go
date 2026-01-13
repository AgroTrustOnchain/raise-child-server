package entities

type Manage struct {
	ID             string   `json:"id"`
	AdminIds       []string `json:"admin_ids"`
	ChildIds       []string `json:"child_ids"`
	VolunteerIds   []string `json:"volunteer_ids"`
	LocalLeaderIds []string `json:"local_leader_ids"`
	SponsorIds     []string `json:"sponsor_ids"`
	SponsorNfts    []string `json:"sponsor_nfts"`
}
