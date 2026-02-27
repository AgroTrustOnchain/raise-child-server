package apiroute

import (
	"raise-child/transport"
	"raise-child/util/middleware"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
)

func InitializeRegionRoutes(server *gin.Engine) {
	var contextPath string = "regions"

	// Rate limits
	var listLimit = middleware.InitializeRateLimiter(rate.Every(time.Second/2), 15)
	var viewLimit = middleware.InitializeRateLimiter(rate.Every(time.Second/5), 20)
	var createLimit = middleware.InitializeRateLimiter(rate.Every(time.Minute), 2)

	// Normal group
	var norGroup = server.Group(contextPath)
	norGroup.GET("", middleware.RateLimitMiddleware(listLimit), transport.GetRegions)
	norGroup.GET("/supported-proposal", middleware.RateLimitMiddleware(listLimit), transport.GetSupportedRegionProposals)
	norGroup.GET("/supported-proposal/user/:id", middleware.RateLimitMiddleware(listLimit), transport.GetUserSupportedRegionProposals)
	norGroup.GET("/supported-proposal/:id", middleware.RateLimitMiddleware(viewLimit), transport.GetSupportedRegionProposal)

	var authGroup = server.Group(contextPath, middleware.Authorize, middleware.ManagerRoleAuthorize)
	authGroup.POST("/supported-proposal", middleware.RateLimitMiddleware(createLimit), transport.CreateSupportedRegionProposal)
}
