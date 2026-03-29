package apiroute

import (
	"raise-child/transport"
	"raise-child/util/middleware"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
)

// func InitializeProfileRoutes(server *gin.Engine) {
// 	var contextPath string = "profiles"

// 	// Auth group
// 	var authGroup = server.Group(contextPath, middleware.Authorize)
// 	authGroup.POST("/:id", transport.UploadProfile)
// }

func InitializeProfileRoutes(server *gin.Engine) {
	var contextPath string = "profiles"

	// Rate limits
	var viewLimit = middleware.InitializeRateLimiter(rate.Every(time.Second/5), 20)

	// Normal group
	var norGroup = server.Group(contextPath)
	norGroup.GET("/personal-wallet-profile/:id", middleware.RateLimitMiddleware(viewLimit), transport.GetWalletPersonalProfile)

	// Auth group
	var authGroup = server.Group(contextPath, middleware.Authorize)
	authGroup.POST("/:id", middleware.RateLimitMiddleware(middleware.InitializeRateLimiter(rate.Every(time.Minute/2), 3)), transport.UploadProfile)
}
