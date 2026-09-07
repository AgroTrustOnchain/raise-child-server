package apiroute

import (
	"raise-child/transport"
	"raise-child/util/middleware"

	"github.com/gin-gonic/gin"
)

func InitializeBankProfileRoutes(server *gin.Engine) {
	// var contextPath string = "banks"

	// // Normal group
	// var norGroup = server.Group(contextPath)
	// norGroup.GET("/:id", transport.GetBankProfile)

	// // Auth group
	// var authGroup = server.Group(contextPath, middleware.Authorize)
	// var rateLimit = middleware.InitializeRateLimiter(rate.Every(time.Minute/5), 5)
	// authGroup.POST("", middleware.RateLimitMiddleware(rateLimit), transport.CreateBankProfile)
	// authGroup.PUT("/:id", middleware.RateLimitMiddleware(rateLimit), transport.UpdateBankProfile)

	var contextPath string = "banks"

	// Normal group
	var norGroup = server.Group(contextPath)
	norGroup.GET("/:id", transport.GetBankProfile)
	norGroup.GET("/user/:id", transport.GetBankProfileByOwner)

	// Manager group
	var authGroup = server.Group(contextPath, middleware.Authorize, middleware.ManagerRoleAuthorize)
	authGroup.POST("", transport.CreateBankProfile)
	authGroup.PUT("/:id", transport.UpdateBankProfile)
}
