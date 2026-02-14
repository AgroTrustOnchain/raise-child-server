package apiroute

import (
	"raise-child/transport"
	"raise-child/util/middleware"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
)

func InitializeGiftRoute(server *gin.Engine) {
	var contextPath string = "gifts"

	// Rate limits
	var viewLimit = middleware.InitializeRateLimiter(rate.Every(time.Second/5), 20)
	var createLimit = middleware.InitializeRateLimiter(rate.Every(time.Minute/1), 3)
	var confirmLimit = middleware.InitializeRateLimiter(rate.Every(time.Minute/5), 2)

	// Normal group
	var norGroup = server.Group(contextPath)
	norGroup.GET("/:id", middleware.RateLimitMiddleware(viewLimit), transport.GetGift)
	norGroup.GET("/child/:id", middleware.RateLimitMiddleware(viewLimit), transport.GetGiftsOfChild)

	// Auth group
	var authGroup = server.Group(contextPath, middleware.Authorize)
	authGroup.POST("", middleware.RateLimitMiddleware(createLimit), transport.CreateGift)
	authGroup.POST("/:id/confirm", middleware.StaffRoleAuthorize, middleware.RateLimitMiddleware(confirmLimit), transport.ConfirmReceiveGift)
}
