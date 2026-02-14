package entities

type MainPool struct {
	ID                ID       `json:"id"`
	LocalPools        []string `json:"local_pools"`
	WithdrawProposals []string `json:"withdraw_proposals"`
	TotalAmount       string   `json:"total_amount"`
}

type LocalPool struct {
	ID          ID       `json:"id"`
	Region      string   `json:"region"`
	Mods        []string `json:"mods"`
	TotalAmount string   `json:"total_amount"`
}
