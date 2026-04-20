package apiroute

import (
	"raise-child/transport"
	"raise-child/util/middleware"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
)

func InitializeCenterRequestRoute(server *gin.Engine) {
	var contextPath string = "center-reqs"

	// Rate limits
	var listLimit = middleware.InitializeRateLimiter(rate.Every(time.Second/2), 15)
	var detailLimit = middleware.InitializeRateLimiter(rate.Every(time.Second/5), 20)
	var createLimit = middleware.InitializeRateLimiter(rate.Every(time.Minute/5), 5)
	var voteLimit = middleware.InitializeRateLimiter(rate.Every(time.Second/1), 5)

	// Normal group
	var norGroup = server.Group(contextPath)
	norGroup.GET("", middleware.RateLimitMiddleware(listLimit), transport.GetCenterRequests)
	norGroup.GET("/user/:id", middleware.RateLimitMiddleware(listLimit), transport.GetWalletCenterRequests)
	norGroup.GET("/:id", middleware.RateLimitMiddleware(detailLimit), transport.GetCenterRequest)

	// Auth group
	var authGroup = server.Group(contextPath, middleware.Authorize)
	authGroup.POST("", middleware.RateLimitMiddleware(createLimit), transport.CreateCenterRequest)
	authGroup.POST("/:id/vote", middleware.RateLimitMiddleware(voteLimit), transport.VoteCenterRequest)
	authGroup.POST("/:id/confirm", middleware.RateLimitMiddleware(createLimit), transport.ConfirmCenterRequest)
}
