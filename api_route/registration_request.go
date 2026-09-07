package apiroute

import (
	"raise-child/transport"
	"raise-child/util/middleware"

	"github.com/gin-gonic/gin"
)

func InitializeRegistrationRequestRoute(server *gin.Engine) {
	// var contextPath string = "registrations"

	// // Rate limits
	// var listLimit = middleware.InitializeRateLimiter(rate.Every(time.Second/2), 15)
	// var viewLimit = middleware.InitializeRateLimiter(rate.Every(time.Second/5), 20)
	// var createLimit = middleware.InitializeRateLimiter(rate.Every(time.Minute), 2)
	// var voteLimit = middleware.InitializeRateLimiter(rate.Every(time.Second/2), 10)
	// var confirmLimit = middleware.InitializeRateLimiter(rate.Every(time.Minute/2), 2)

	// // Normal group
	// var norGroup = server.Group(contextPath)
	// norGroup.GET("", middleware.RateLimitMiddleware(listLimit), transport.GetRegistrationRequests)
	// norGroup.GET("/user/:id", middleware.RateLimitMiddleware(listLimit), transport.GetWalletRegistrationRequests)
	// norGroup.GET("/:id", middleware.RateLimitMiddleware(viewLimit), transport.GetRegistrationRequest)

	// // Auth group
	// var authGroup = server.Group(contextPath, middleware.Authorize)
	// authGroup.POST("", middleware.RateLimitMiddleware(createLimit), transport.CreateRegistrationRequest)
	// authGroup.POST("/:id/vote", middleware.RateLimitMiddleware(voteLimit), transport.VoteRegistrationRequest)
	// authGroup.POST("/:id/confirm", middleware.RateLimitMiddleware(confirmLimit), transport.ConfirmRegistrationRequest)

	var contextPath string = "registrations"

	// Normal group
	var norGroup = server.Group(contextPath)
	norGroup.GET("", transport.GetRegistrationRequests)
	norGroup.GET("/user/:id", transport.GetWalletRegistrationRequests)
	norGroup.GET("/:id", transport.GetRegistrationRequest)

	// Auth group
	var authGroup = server.Group(contextPath, middleware.Authorize)
	authGroup.POST("", transport.CreateRegistrationRequest)
	authGroup.POST("/:id/vote", transport.VoteRegistrationRequest)
}
