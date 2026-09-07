package apiroute

import (
	"raise-child/transport"
	"raise-child/util/middleware"

	"github.com/gin-gonic/gin"
)

// func InitializeWithdrawProposalRoute(server *gin.Engine) {
// 	var contextPath string = "withdraw-proposals"

// 	// Normal group
// 	var norGroup = server.Group(contextPath)
// 	norGroup.GET("", transport.GetWithdrawProposals)
// 	norGroup.GET("/:id", transport.GetWithdrawProposal)

// 	// Auth group
// 	var authGroup = server.Group(contextPath, middleware.Authorize)
// 	authGroup.POST("", transport.CreateWithdrawProposal)
// 	authGroup.POST("/:id/vote", transport.VoteWithdrawProposal)
// 	authGroup.POST("/:id/confirm", transport.ConfirmWithdrawProposal)
// 	authGroup.POST("/:id/main-pool-confirm", transport.ConfirmMainPoolWithdrawProposal)
// }

func InitializeWithdrawProposalRoute(server *gin.Engine) {
	// var contextPath string = "withdraw-proposals"

	// // Rate limits
	// var viewLimit = middleware.InitializeRateLimiter(rate.Every(time.Second/5), 20)
	// var createLimit = middleware.InitializeRateLimiter(rate.Every(time.Minute*2), 2)
	// var voteLimit = middleware.InitializeRateLimiter(rate.Every(time.Second/1), 5)
	// var confirmLimit = middleware.InitializeRateLimiter(rate.Every(time.Minute/4), 4)

	// // Normal group
	// var norGroup = server.Group(contextPath)
	// norGroup.GET("", middleware.RateLimitMiddleware(viewLimit), transport.GetWithdrawProposals)
	// norGroup.GET("/:id", middleware.RateLimitMiddleware(viewLimit), transport.GetWithdrawProposal)

	// // Auth group
	// var authGroup = server.Group(contextPath, middleware.Authorize)
	// authGroup.POST("", middleware.RateLimitMiddleware(createLimit), transport.CreateWithdrawProposal)
	// authGroup.POST("/:id/vote", middleware.RateLimitMiddleware(voteLimit), transport.VoteWithdrawProposal)
	// authGroup.POST("/:id/confirm", middleware.RateLimitMiddleware(confirmLimit), transport.ConfirmWithdrawProposal)
	// authGroup.POST("/:id/main-pool-confirm", middleware.RateLimitMiddleware(confirmLimit), transport.ConfirmMainPoolWithdrawProposal)

	var contextPath string = "withdraw-proposals"

	// Normal group
	var norGroup = server.Group(contextPath)
	norGroup.GET("", transport.GetWithdrawProposals)
	norGroup.GET("/pending", transport.GetOnChainPendingWithdrawProposals)
	norGroup.GET("/:id", transport.GetWithdrawProposal)

	// Auth group
	var authGroup = server.Group(contextPath, middleware.Authorize)
	authGroup.POST("/:id/vote", transport.VoteWithdrawProposal)
	authGroup.POST("/:id/main-pool-confirm", transport.ConfirmMainPoolWithdrawProposal)

	// Leader group
	var leaderGroup = server.Group(contextPath, middleware.Authorize, middleware.LeaderAuthorize)
	leaderGroup.POST("/children/propose", transport.ProposeChildrenWithdrawRequests)

	// Manage group
	var managerGroup = server.Group(contextPath, middleware.Authorize, middleware.ManagerRoleAuthorize)
	managerGroup.POST("", transport.CreateWithdrawProposal)

	// Admin Group
	var adminGroup = server.Group(contextPath, middleware.Authorize, middleware.AdminAuthorize)
	adminGroup.POST("/:id/confirm", transport.ConfirmWithdrawProposal)
}
