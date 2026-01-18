package entities

type MainPool struct {
	ID                ID       `json:"id"`
	Admin             string   `json:"admin"`
	LocalPools        []string `json:"local_pools"`
	WithDrawProposals []string `json:"withdraw_proposals"`
	Mods              []string `json:"mods"`
	TotalAmount       string   `json:"total_amount"`
}

type LocalPool struct {
	ID          ID       `json:"id"`
	Region      string   `json:"region"`
	Mods        []string `json:"mods"`
	TotalAmount string   `json:"total_amount"`
}
