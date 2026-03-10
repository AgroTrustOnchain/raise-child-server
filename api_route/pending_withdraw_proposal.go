package apiroute

import (
	"raise-child/transport"
	"raise-child/util/middleware"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
)

func InitializePendingWithdrawProposalRoute(server *gin.Engine) {
	var contextPath string = "pending-withdraw-proposals"

	// Rate limits
	var viewLimit = middleware.InitializeRateLimiter(rate.Every(time.Second/5), 20)
	var createLimit = middleware.InitializeRateLimiter(rate.Every(time.Minute), 2)
	var confirmLimit = middleware.InitializeRateLimiter(rate.Every(time.Minute/6), 10)

	// Auth group
	var managerGroup = server.Group(contextPath, middleware.Authorize, middleware.ManagerRoleAuthorize)
	managerGroup.GET("/:id", middleware.RateLimitMiddleware(viewLimit), transport.GetPendingWithdrawProposal)
	managerGroup.GET("", middleware.RateLimitMiddleware(viewLimit), transport.GetPendingWithdrawProposals)
	managerGroup.POST("", middleware.RateLimitMiddleware(createLimit), transport.CreatePendingWithdrawProposal)

	var adminGroup = server.Group(contextPath, middleware.Authorize, middleware.AdminAuthorize)
	adminGroup.POST("/:id/approve", middleware.RateLimitMiddleware(confirmLimit), transport.ApprovePendingWithdrawProposal)
	adminGroup.POST("/:id/refuse", middleware.RateLimitMiddleware(confirmLimit), transport.RefusePendingWithdrawProposal)
}
