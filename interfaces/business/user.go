package business

// import (
// 	"context"

// 	"github.com/blockuniorg/blockuni/model/dtos/request"
// 	"github.com/blockuniorg/blockuni/model/dtos/response"
// 	"github.com/blockuniorg/blockuni/model/entities"
// )

// type IUserService interface {
// 	// Auth
// 	RegisterMember(req request.RegisterRequest, ctx context.Context) (string, error)
// 	Login(req request.LoginRequest, ctx context.Context) (response.LoginResponse, error)
// 	Logout(id string, ctx context.Context) error
// 	VerifyAction(rawToken string, ctx context.Context) (string, error)
// 	ResetPassword(input string, ctx context.Context) (string, error)
// 	VerifyResetPassword(req request.VerifyResetPasswordRequest, ctx context.Context) error

// 	// User
// 	GetUsers(req request.GetUsersRequest, ctx context.Context) (*[]entities.User, error)
// 	GetUserById(id string, ctx context.Context) (*entities.User, error)
// 	UpdateUser(req request.UpdateUserRequest, ctx context.Context) (string, error)
// 	RemoveAccount(req request.RemoveAccountRequest, ctx context.Context) error
// }

// GetAllUsers(pageNumber int, ctx context.Context) (response.PaginationDataResponse, error)
// 	GetUsersByRole(pageNumber int, role string, ctx context.Context) (response.PaginationDataResponse, error)
// 	GetUsersByStatus(pageNumber int, status *bool, ctx context.Context) (response.PaginationDataResponse, error)
// 	GetUsersByKeyword(pageNumber int, keyword string, ctx context.Context) (response.PaginationDataResponse, error)
// 	GetUser(id string, ctx context.Context) (*entity.User, error)
// 	CreateUser(req request.CreateUserRequest, ctx context.Context) (string, error)
// 	CreateVipCode(req request.CreateVipCodeRequest, ctx context.Context) error
// 	UpdateUser(req request.UpdateUserRequest, ctx context.Context) (string, error)
// 	ChangeUserStatus(req request.ChangeUserStatusRequest, ctx context.Context) (string, error)
// 	Login(req request.LoginRequest, ctx context.Context) (string, string, error)
// 	Logout(id string, ctx context.Context) error

// 	VerifyAction(rawToken string, ctx context.Context) (string, error)
// 	ResetPassword(newPass, re_newPass, token string, ctx context.Context) (string, error)
// 	RefreshToken(req request.RefreshTokenRequest, ctx context.Context) (string, error)
