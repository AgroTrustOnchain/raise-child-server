package apiroute

import (
	"raise-child/transport"
	"raise-child/util/middleware"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
)

// func InitializePaymentsRoutes(server *gin.Engine) {
// 	var contextPath string = "payments"

// 	// Normal group
// 	var norGroup = server.Group(contextPath)
// 	norGroup.GET("/callback/:id", transport.CallbackTransaction)

// 	// Auth group
// 	var authGroup = server.Group(contextPath, middleware.Authorize)
// 	authGroup.POST("/donate", transport.Donate)
// }

func InitializePaymentsRoutes(server *gin.Engine) {
	var contextPath string = "payments"

	// Rate limits
	var callbackLimit = middleware.InitializeRateLimiter(rate.Every(time.Second/5), 30)
	var donateLimit = middleware.InitializeRateLimiter(rate.Every(time.Minute/2), 5)

	// Normal group
	var norGroup = server.Group(contextPath)
	norGroup.GET("/callback/:id", middleware.RateLimitMiddleware(callbackLimit), transport.CallbackTransaction)

	// Auth group
	var authGroup = server.Group(contextPath, middleware.Authorize)
	authGroup.POST("/donate", middleware.RateLimitMiddleware(donateLimit), transport.Donate)
	authGroup.GET("/auth-callback/:id", middleware.AdminAuthorize, middleware.RateLimitMiddleware(donateLimit), transport.CallbackWithAuthTransaction)
}
