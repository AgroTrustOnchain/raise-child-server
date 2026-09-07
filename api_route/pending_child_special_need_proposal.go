package apiroute

import (
	"raise-child/transport"
	"raise-child/util/middleware"

	"github.com/gin-gonic/gin"
)

func InitializeChildPendingSpecialNeedProposalRoutes(server *gin.Engine) {
	// var contextPath string = "pending-special-needs"

	// // Rate limits
	// var listLimit = middleware.InitializeRateLimiter(rate.Every(time.Second/2), 15)
	// var viewLimit = middleware.InitializeRateLimiter(rate.Every(time.Second/5), 20)
	// var postLimit = middleware.InitializeRateLimiter(rate.Every(time.Minute/4), 6)

	// // Manager group
	// var managerGroup = server.Group(contextPath, middleware.Authorize, middleware.ManagerRoleAuthorize)
	// managerGroup.GET("", middleware.RateLimitMiddleware(listLimit), transport.GetPendingChildSpecialNeedProposals)
	// managerGroup.GET("/:id", middleware.RateLimitMiddleware(viewLimit), transport.GetPendingChildSpecialNeedProposal)

	// // Admin group
	// var adminGroup = server.Group(contextPath, middleware.Authorize, middleware.AdminAuthorize)
	// adminGroup.POST("/:id/approve", middleware.RateLimitMiddleware(postLimit), transport.ApprovePendingChildSpecialNeedProposal)
	// adminGroup.POST("/:id/refuse", middleware.RateLimitMiddleware(postLimit), transport.RefusePendingChildSpecialNeedProposal)

	var contextPath string = "pending-special-needs"

	// Manager group
	var managerGroup = server.Group(contextPath, middleware.Authorize, middleware.ManagerRoleAuthorize)
	managerGroup.GET("", transport.GetPendingChildSpecialNeedProposals)
	managerGroup.GET("/:id", transport.GetPendingChildSpecialNeedProposal)

	// Admin group
	var adminGroup = server.Group(contextPath, middleware.Authorize, middleware.AdminAuthorize)
	adminGroup.POST("/:id/approve", transport.ApprovePendingChildSpecialNeedProposal)
	adminGroup.POST("/:id/refuse", transport.RefusePendingChildSpecialNeedProposal)
}
