package apiroute

import (
	"raise-child/transport"
	"raise-child/util/middleware"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
)

// func InitializeAdminRoute(server *gin.Engine) {
// 	var contextPath string = "admins"

// 	// Auth group
// 	var authGroup = server.Group(contextPath, middleware.Authorize)
// 	authGroup.POST("", transport.UpdatePublisherInfo)

// 	// Normal group
// 	var norGroup = server.Group(contextPath)
// 	norGroup.GET("", transport.GetAdmins)
// }

func InitializeAdminRoute(server *gin.Engine) {
	var contextPath string = "admins"

	// Auth group
	var authGroup = server.Group(contextPath, middleware.Authorize)
	authGroup.POST("", transport.UpdatePublisherInfo)

	// Normal group
	var norGroup = server.Group(contextPath)
	var listLimit = middleware.InitializeRateLimiter(rate.Every(time.Second/2), 20)
	norGroup.GET("", middleware.RateLimitMiddleware(listLimit), transport.GetAdmins)
}
