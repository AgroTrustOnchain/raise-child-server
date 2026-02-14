package apiroute

import (
	"raise-child/transport"
	"raise-child/util/middleware"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
)

func InitializeBankProfileRoutes(server *gin.Engine) {
	var contextPath string = "banks"

	// Normal group
	var norGroup = server.Group(contextPath)
	norGroup.GET("/:id", transport.GetBankProfile)

	// Auth group
	var authGroup = server.Group(contextPath, middleware.Authorize)
	var rateLimit = middleware.InitializeRateLimiter(rate.Every(time.Minute/5), 5)
	authGroup.POST("", middleware.RateLimitMiddleware(rateLimit), transport.CreateBankProfile)
	authGroup.PUT("/:id", middleware.RateLimitMiddleware(rateLimit), transport.UpdateBankProfile)
}
