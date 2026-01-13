package response

type LoginResponse struct {
	Token string `json:"token"`
}

type GetNonceResponse struct {
	Nonce string `json:"nonce"`
}

type GetSaltResponse struct {
	Salt string `json:"salt"`
}
