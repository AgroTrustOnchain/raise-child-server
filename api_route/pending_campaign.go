package apiroute

import (
	"raise-child/transport"
	"raise-child/util/middleware"

	"github.com/gin-gonic/gin"
)

func InitializePendingCampaignRoutes(server *gin.Engine) {
	// var contextPath string = "pending-pool-campaigns"

	// // Rate limits
	// var listLimit = middleware.InitializeRateLimiter(rate.Every(time.Second/2), 15)
	// var viewLimit = middleware.InitializeRateLimiter(rate.Every(time.Second/5), 20)
	// var postLimit = middleware.InitializeRateLimiter(rate.Every(time.Minute/4), 6)

	// // Nor group
	// var norGroup = server.Group(contextPath)
	// norGroup.GET("", middleware.RateLimitMiddleware(listLimit), transport.GetPendingCampaigns)
	// norGroup.GET("/:id", middleware.RateLimitMiddleware(viewLimit), transport.GetPendingCampaign)

	// // Manager group
	// var managerGroup = server.Group(contextPath, middleware.Authorize, middleware.ManagerRoleAuthorize)
	// managerGroup.POST("", middleware.RateLimitMiddleware(postLimit), transport.CreatePendingCampaign)

	// // Admin group
	// var adminGroup = server.Group(contextPath, middleware.Authorize, middleware.AdminAuthorize)
	// adminGroup.POST("/:id/approve", middleware.RateLimitMiddleware(postLimit), transport.ApprovePendingCampaign)
	// adminGroup.POST("/:id/refuse", middleware.RateLimitMiddleware(postLimit), transport.RefusePendingCampaign)

	var contextPath string = "pending-pool-campaigns"

	// Nor group
	var norGroup = server.Group(contextPath)
	norGroup.GET("", transport.GetPendingCampaigns)
	norGroup.GET("/:id", transport.GetPendingCampaign)

	// Manager group
	var managerGroup = server.Group(contextPath, middleware.Authorize, middleware.ManagerRoleAuthorize)
	managerGroup.POST("", transport.CreatePendingCampaign)

	// Admin group
	var adminGroup = server.Group(contextPath, middleware.Authorize, middleware.AdminAuthorize)
	adminGroup.POST("/:id/approve", transport.ApprovePendingCampaign)
	adminGroup.POST("/:id/refuse", transport.RefusePendingCampaign)
}
