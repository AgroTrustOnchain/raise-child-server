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

	// Auth group
	var authGroup = server.Group(contextPath, middleware.Authorize)
	authGroup.POST("/:id", middleware.RateLimitMiddleware(middleware.InitializeRateLimiter(rate.Every(time.Minute/2), 3)), transport.UploadProfile)
}
