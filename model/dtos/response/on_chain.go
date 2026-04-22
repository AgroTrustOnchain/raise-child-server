package response

type Contents struct {
	Json map[string]interface{} `json:"json"`
}

type AsMoveObject struct {
	Contents Contents `json:"contents"`
}

type Node struct {
	Address      string       `json:"address"`
	AsMoveObject AsMoveObject `json:"asMoveObject"`
}

type Objects struct {
	Nodes []Node `json:"nodes"`
}

type SuiGraphQlObjectResponse struct {
	Objects Objects `json:"objects"`
}

type BuildTransactionResponse struct {
	TxBytes         string `json:"tx_bytes"`
	ProposalId      string `json:"proposal_id"`
	CenterReq       string `json:"center_req"`
	UploadChildReq  string `json:"upload_child_req"`
	RegistrationReq string `json:"registration_req"`
}
