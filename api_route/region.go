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
	var postLimit = middleware.InitializeRateLimiter(rate.Every(time.Minute/2), 2)

	// Normal group
	var norGroup = server.Group(contextPath)
	norGroup.GET("", middleware.RateLimitMiddleware(listLimit), transport.GetRegions)
	norGroup.GET("/established", middleware.RateLimitMiddleware(listLimit), transport.GetEstablishedRegions)
	norGroup.GET("/established/:region", middleware.RateLimitMiddleware(viewLimit), transport.GetRegionDetail)
	norGroup.GET("/supported-suggestions", middleware.RateLimitMiddleware(listLimit), transport.GetSupportedRegionSuggestions)
	norGroup.GET("/supported-suggestions/:id", middleware.RateLimitMiddleware(viewLimit), transport.GetSupportedRegionSuggestion)

	// Auth group
	var authGroup = server.Group(contextPath, middleware.Authorize)
	norGroup.GET("/user/:id/supported-suggestions", middleware.RateLimitMiddleware(listLimit), transport.GetWalletSupportedRegionSuggestions)
	authGroup.POST("/supported-suggestions", middleware.RateLimitMiddleware(postLimit), transport.CreateSupportedRegionSuggestion)

	// Admin group
	var adminGroup = server.Group(contextPath, middleware.Authorize, middleware.AdminAuthorize)
	adminGroup.GET("/admin/supported-suggestions", middleware.RateLimitMiddleware(viewLimit), transport.AdminGetSupportedRegionSuggestions)
	adminGroup.POST("/supported-suggestions/:id/review", middleware.RateLimitMiddleware(postLimit), transport.ReviewRegionSuggestion)
}
