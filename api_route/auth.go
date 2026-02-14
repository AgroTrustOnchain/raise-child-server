package apiroute

import (
	"raise-child/transport"
	"raise-child/util/middleware"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
)

func InitializeAuthHandlerRoutes(server *gin.Engine) {
	var contextPath string = "auth"

	var norGroup = server.Group(contextPath)
	var loginLimit = middleware.InitializeRateLimiter(rate.Every(time.Minute/10), 10)
	var commonLimit = middleware.InitializeRateLimiter(rate.Every(time.Second/2), 10)
	norGroup.POST("/login", middleware.RateLimitMiddleware(loginLimit), transport.Login)
	norGroup.GET("/salt/:id", middleware.RateLimitMiddleware(commonLimit), transport.GetSalt)

	var authGroup = server.Group(contextPath, middleware.Authorize)
	authGroup.POST("/logout", middleware.RateLimitMiddleware(commonLimit), transport.Logout)
}
