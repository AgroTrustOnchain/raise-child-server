package apiroute

import (
	"raise-child/transport"

	"github.com/gin-gonic/gin"
)

func InitializeNotiRoute(server *gin.Engine) {
	// var contextPath string = "notis"

	// // Rate limits
	// var listLimit = middleware.InitializeRateLimiter(rate.Every(time.Second/2), 20)

	// // Normal group
	// var norGroup = server.Group(contextPath)
	// norGroup.GET("/user/:id", middleware.RateLimitMiddleware(listLimit), transport.GetWalletNotifications)

	var contextPath string = "notis"

	// Normal group
	var norGroup = server.Group(contextPath)
	norGroup.GET("/user/:id", transport.GetWalletNotifications)
}
