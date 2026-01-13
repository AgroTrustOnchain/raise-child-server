package business

import (
	"context"
	"raise-child/model/dtos/request"
	"raise-child/model/dtos/response"
)

type IAuthService interface {
	GetNonce(address string, ctx context.Context) (response.GetNonceResponse, error)
	GetSalt(id string, ctx context.Context) (response.GetSaltResponse, error)
	Login(req request.LoginRequest, ctx context.Context) (response.LoginResponse, error)
	LoginV2(req request.LoginRequestV2, ctx context.Context) (response.LoginResponse, error)
	Logout(address string, ctx context.Context) error
	LogoutV2(ctx context.Context) error
}
