package repository

// import (
// 	"context"

// 	"github.com/blockuniorg/blockuni/model/entities"
// )

// type IUserRepository interface {
// 	GetAllUsers(pageNumber int, cursorID string, ctx context.Context) (*[]entities.User, int, error)
// 	GetUsersByRole(pageNumber int, role string, cursorID string, ctx context.Context) (*[]entities.User, int, error)
// 	GetUsersByStatus(pageNumber int, status bool, cursorID string, ctx context.Context) (*[]entities.User, int, error)
// 	GetUsersByKeyword(pageNumber int, keyword string, cursorID string, ctx context.Context) (*[]entities.User, int, error)
// 	GetUserById(id string, ctx context.Context) (*entities.User, error)
// 	GetUserByLogin(input string, ctx context.Context) (*entities.User, error)
// 	GetUserByUsername(username string, ctx context.Context) (*entities.User, error)
// 	GetUserByEmail(email string, ctx context.Context) (*entities.User, error)
// 	CreateUser(user entities.User, ctx context.Context) error
// 	UpdateUser(user entities.User, ctx context.Context) error
// }
