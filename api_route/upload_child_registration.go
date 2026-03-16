package apiroute

import (
	"raise-child/transport"
	"raise-child/util/middleware"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
)

// func InitializeUploadChildRequestRoute(server *gin.Engine) {
// 	var contextPath string = "child-upload-reqs"

// 	// Normal group
// 	var norGroup = server.Group(contextPath)
// 	norGroup.GET("", transport.GetUploadChildRequests)
// 	norGroup.GET("/user/:id", transport.GetWalletUploadChildRequests)
// 	norGroup.GET("/:id", transport.GetUploadChildRequest)

// 	// Auth group
// 	var authGroup = server.Group(contextPath, middleware.Authorize)
// 	authGroup.POST("", transport.CreateUploadChildRequest)
// 	authGroup.POST("/:id/vote", transport.VoteUploadChildRequest)
// 	authGroup.POST("/:id/confirm", transport.ConfirmUploadChildRequest)
// }

func InitializeUploadChildRequestRoute(server *gin.Engine) {
	var contextPath string = "child-upload-reqs"

	// Rate limits
	var listLimit = middleware.InitializeRateLimiter(rate.Every(time.Second/2), 15)
	var viewLimit = middleware.InitializeRateLimiter(rate.Every(time.Second/5), 20)
	var createLimit = middleware.InitializeRateLimiter(rate.Every(time.Minute), 2)
	var voteLimit = middleware.InitializeRateLimiter(rate.Every(time.Second/2), 10)
	var confirmLimit = middleware.InitializeRateLimiter(rate.Every(time.Minute/2), 2)

	// Normal group
	var norGroup = server.Group(contextPath)
	norGroup.GET("", middleware.RateLimitMiddleware(listLimit), transport.GetUploadChildRequests)
	norGroup.GET("/user/:id", middleware.RateLimitMiddleware(listLimit), transport.GetWalletUploadChildRequests)
	norGroup.GET("/:id", middleware.RateLimitMiddleware(viewLimit), transport.GetUploadChildRequest)

	// Auth group
	var authGroup = server.Group(contextPath, middleware.Authorize)
	authGroup.POST("", middleware.RateLimitMiddleware(createLimit), transport.CreateUploadChildRequest)
	authGroup.POST("/:id/vote", middleware.RateLimitMiddleware(voteLimit), transport.VoteUploadChildRequest)
	authGroup.POST("/:id/confirm", middleware.RateLimitMiddleware(confirmLimit), transport.ConfirmUploadChildRequest)

	// Admin group
	var adminGroup = server.Group(contextPath, middleware.Authorize, middleware.AdminAuthorize)
	adminGroup.POST("/:id/review", middleware.RateLimitMiddleware(voteLimit), transport.ReviewUploadChildRequest)
}
