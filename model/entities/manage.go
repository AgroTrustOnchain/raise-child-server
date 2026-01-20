package entities

type Manage struct {
	ID              ID       `json:"id"`
	AdminIds        []string `json:"admin_ids"`
	AdminNfts       []string `json:"admin_nfts"`
	ChildIds        []string `json:"child_ids"`
	VolunteerIds    []string `json:"volunteer_ids"`
	VolunteerNfts   []string `json:"volunteer_nfts"`
	LocalLeaderNfts []string `json:"local_leader_nfts"`
	LocalLeaderIds  []string `json:"local_leader_ids"`
	LocalRegions    []string `json:"local_regions"`
	SponsorIds      []string `json:"sponsor_ids"`
	SponsorNfts     []string `json:"sponsor_nfts"`
}
