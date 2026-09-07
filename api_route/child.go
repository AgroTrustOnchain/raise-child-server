package apiroute

import (
	"raise-child/transport"
	"raise-child/util/middleware"

	"github.com/gin-gonic/gin"
)

func InitializeChildRoutes(server *gin.Engine) {
	var contextPath string = "children"

	// Normal group
	var norGroup = server.Group(contextPath)
	norGroup.GET("", transport.GetChildren)
	norGroup.GET("/user/:id/supported", transport.GetUserSupportedChildren)
	norGroup.GET("/:id", transport.GetChild)

	// Auth group
	var authGroup = server.Group(contextPath, middleware.Authorize)
	authGroup.PUT("/Metadata/string/:id", transport.AddChildStringMetadata)
	authGroup.PUT("/Metadata/number/:id", transport.AddChildNumberMetadata)
	authGroup.POST("", transport.UploadChild)
	authGroup.POST("/books-need/:id/support", transport.SupportBooksNeed)
	authGroup.POST("/health-insurance-need/:id/support", transport.SupportHealthInsuranceNeed)
	authGroup.POST("/meal-need/:id/support", transport.SupportMealNeed)
	authGroup.POST("/special-need/:id/support", transport.SupportSpecialNeed)
	authGroup.POST("/special-need/proposal/:id/vote", transport.VoteSpecialNeedProposal)

	// Manager group
	var managerGroup = server.Group(contextPath, middleware.Authorize, middleware.ManagerRoleAuthorize)
	managerGroup.POST("/books-need/withdraw-proposal", transport.CreateBooksNeedWithdrawProposal)
	managerGroup.POST("/health-insurance-need/withdraw-proposal", transport.CreateHealthInsuranceNeedWithdrawProposal)
	managerGroup.POST("/meal-need/withdraw-proposal", transport.CreateMealNeedWithdrawProposal)
	managerGroup.POST("/special-need/withdraw-proposal", transport.CreateSpecialNeedWithdrawProposal)
	managerGroup.POST("/special-need/proposal", transport.CreateSpecialNeedProposal)
	managerGroup.POST("/special-need/proposal/:id/confirm", transport.ConfirmSpecialNeedProposal)

	// Staff group
	var staffGroup = server.Group(contextPath, middleware.Authorize, middleware.StaffRoleAuthorize)
	staffGroup.POST("/:id/provide-meal/confirm", transport.ConfirmProvideMealForChild)

	// Leader group
	var leaderGroup = server.Group(contextPath, middleware.Authorize, middleware.LeaderAuthorize)
	leaderGroup.PUT("/books-need", transport.UpdateBooksNeed)
	leaderGroup.PUT("/meal-need", transport.UpdateMealNeed)
	leaderGroup.PUT("/health-insurance-need", transport.UpdateHealthInsuranceNeed)
}
