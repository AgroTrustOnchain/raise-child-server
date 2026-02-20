package request

type BuildTransactionRequest struct {
	Sender    string                 `json:"sender" validate:"required"`
	Function  string                 `json:"function" validate:"required"`
	Arguments map[string]interface{} `json:"arguments" validate:"required"`
}

type ExecuteTransactionRequest struct {
	TxBytes        string `json:"tx_bytes" validate:"required"`
	Signature      string `json:"signature" validate:"required"`
	ProposalID     string `json:"proposal_id"`
	CenterReq      string `json:"center_req"`
	UploadChildReq string `json:"upload_child_req"`
	RegistraionReq string `json:"registration_req"`
}

type MoneyActionRequest struct {
	Sender     string `json:"sender" validate:"required"`
	Amount     int    `json:"amoun" validate:"required,min=1"`
	CoinType   string `json:"coin_type"`                       // e.g. "SUI", "USDC"
	ActionType string `json:"action_type" validate:"required"` // "donate" || "withdraw"
	Message    string `json:"message"`
}
