package apiroute

import (
	"raise-child/transport"
	"raise-child/util/middleware"

	"github.com/gin-gonic/gin"
)

func InitializeCampaignRoutes(server *gin.Engine) {
	// var contextPath string = "pool-campaigns"

	// // Rate limits
	// var listLimit = middleware.InitializeRateLimiter(rate.Every(time.Second/2), 15)
	// var viewLimit = middleware.InitializeRateLimiter(rate.Every(time.Second/5), 20)
	// var postLimit = middleware.InitializeRateLimiter(rate.Every(time.Minute/5), 6)
	// var supportLimit = middleware.InitializeRateLimiter(rate.Every(time.Minute/6), 6)

	// // Nor group
	// var norGroup = server.Group(contextPath)
	// norGroup.GET("", middleware.RateLimitMiddleware(listLimit), transport.GetCampaigns)
	// norGroup.GET("/:id", middleware.RateLimitMiddleware(viewLimit), transport.GetCampaign)

	// // Auth group
	// var authGroup = server.Group(contextPath, middleware.Authorize)
	// authGroup.POST("/:id/approve", middleware.RateLimitMiddleware(supportLimit), transport.SupportCampaign)

	// // Manager group
	// var managerGroup = server.Group(contextPath, middleware.Authorize, middleware.ManagerRoleAuthorize)
	// managerGroup.POST("/withdraw-proposal", middleware.RateLimitMiddleware(postLimit), transport.CreateCampaignWithdrawProposal)

	var contextPath string = "pool-campaigns"

	// Nor group
	var norGroup = server.Group(contextPath)
	norGroup.GET("", transport.GetCampaigns)
	norGroup.GET("/:id", transport.GetCampaign)

	// Auth group
	var authGroup = server.Group(contextPath, middleware.Authorize)
	authGroup.POST("/:id/approve", transport.SupportCampaign)

	// Manager group
	var managerGroup = server.Group(contextPath, middleware.Authorize, middleware.ManagerRoleAuthorize)
	managerGroup.POST("/withdraw-proposal", transport.CreateCampaignWithdrawProposal)
}
