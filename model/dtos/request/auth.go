package request

type GetNonceRequest struct {
	Address string `json:"address" form:"address" validate:"required"`
}

type LoginRequest struct {
	Address   string `json:"address" validate:"required"`
	PublicKey string `json:"public_key" validate:"required"`
	Message   string `json:"message" validate:"required"`
	Signature string `json:"signature" validate:"required"`
}

type LoginRequestV2 struct {
	Address string `json:"address" validate:"required"`
	Sub     string `json:"sub" validate:"required"`
}

type LogoutRequest struct {
	Address string `json:"address" validate:"required"`
}
