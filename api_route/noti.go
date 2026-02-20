package apiroute

import (
	"raise-child/transport"
	"raise-child/util/middleware"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
)

func InitializeNotiRoute(server *gin.Engine) {
	var contextPath string = "notis"

	// Rate limits
	var listLimit = middleware.InitializeRateLimiter(rate.Every(time.Second/2), 20)

	// Normal group
	var norGroup = server.Group(contextPath)
	norGroup.GET("/user/:id", middleware.RateLimitMiddleware(listLimit), transport.GetWalletNotifications)
}
