package apiroute

import (
	"raise-child/transport"
	"raise-child/util/middleware"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
)

// func InitializeChildRoutes(server *gin.Engine) {
// 	var contextPath string = "children"

// 	// Normal group
// 	var norGroup = server.Group(contextPath)
// 	norGroup.GET("", transport.GetChildren)
// 	norGroup.GET("/:id", transport.GetChild)

// 	// Auth group
// 	var authGroup = server.Group(contextPath, middleware.Authorize)
// 	authGroup.PUT("/metadata/string/:id", transport.AddChildStringMetadata)
// 	authGroup.PUT("/metadata/number/:id", transport.AddChildNumberMetadata)
// 	authGroup.POST("", transport.UploadChild)
// }

func InitializeChildRoutes(server *gin.Engine) {
	var contextPath string = "children"

	// Rate limits
	var listLimit = middleware.InitializeRateLimiter(rate.Every(time.Second/2), 10)
	var detailLimit = middleware.InitializeRateLimiter(rate.Every(time.Second/5), 20)
	var uploadLimit = middleware.InitializeRateLimiter(rate.Every(time.Minute/30), 30)
	var metadataLimit = middleware.InitializeRateLimiter(rate.Every(time.Second/1), 10)
	var donateLimit = middleware.InitializeRateLimiter(rate.Every(time.Minute/4), 7)
	var proposalLimit = middleware.InitializeRateLimiter(rate.Every(time.Minute/3), 5)
	var voteLimit = middleware.InitializeRateLimiter(rate.Every(time.Second/1), 5)
	var confirmLimit = middleware.InitializeRateLimiter(rate.Every(time.Minute/2), 5)

	// Normal group
	var norGroup = server.Group(contextPath)
	norGroup.GET("", middleware.RateLimitMiddleware(listLimit), transport.GetChildren)
	norGroup.GET("/:id", middleware.RateLimitMiddleware(detailLimit), transport.GetChild)

	// Auth group
	var authGroup = server.Group(contextPath, middleware.Authorize)
	authGroup.PUT("/metadata/string/:id", middleware.RateLimitMiddleware(metadataLimit), transport.AddChildStringMetadata)
	authGroup.PUT("/metadata/number/:id", middleware.RateLimitMiddleware(metadataLimit), transport.AddChildNumberMetadata)
	authGroup.POST("", middleware.RateLimitMiddleware(uploadLimit), transport.UploadChild)
	authGroup.POST("/books-need/:id/support", middleware.RateLimitMiddleware(donateLimit), transport.SupportBooksNeed)
	authGroup.POST("/health-insurance-need/:id/support", middleware.RateLimitMiddleware(donateLimit), transport.SupportHealthInsuranceNeed)
	authGroup.POST("/meal-need/:id/support", middleware.RateLimitMiddleware(donateLimit), transport.SupportMealNeed)
	authGroup.POST("/special-need/:id/support", middleware.RateLimitMiddleware(donateLimit), transport.SupportSpecialNeed)
	authGroup.POST("/special-need/proposal/:id/vote", middleware.RateLimitMiddleware(voteLimit), transport.VoteSpecialNeedProposal)

	// Manager group
	var managerGroup = server.Group(contextPath, middleware.Authorize, middleware.ManagerRoleAuthorize)
	managerGroup.POST("/books-need/withdraw-proposal", middleware.RateLimitMiddleware(proposalLimit), transport.CreateBooksNeedWithdrawProposal)
	managerGroup.POST("/health-insurance-need/withdraw-proposal", middleware.RateLimitMiddleware(proposalLimit), transport.CreateHealthInsuranceNeedWithdrawProposal)
	managerGroup.POST("/meal-need/withdraw-proposal", middleware.RateLimitMiddleware(proposalLimit), transport.CreateMealNeedWithdrawProposal)
	managerGroup.POST("/special-need/withdraw-proposal", middleware.RateLimitMiddleware(proposalLimit), transport.CreateSpecialNeedWithdrawProposal)
	managerGroup.POST("/special-need/proposal", middleware.RateLimitMiddleware(proposalLimit), transport.CreateSpecialNeedProposal)
	managerGroup.POST("/special-need/proposal/:id/confirm", middleware.RateLimitMiddleware(proposalLimit), transport.ConfirmSpecialNeedProposal)

	// Staff group
	var staffGroup = server.Group(contextPath, middleware.Authorize, middleware.StaffRoleAuthorize)
	staffGroup.POST("/:id/provide-meal/confirm", middleware.RateLimitMiddleware(confirmLimit), transport.ConfirmProvideMealForChild)
}
