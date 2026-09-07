package apiroute

import (
	"raise-child/transport"
	"raise-child/util/middleware"

	"github.com/gin-gonic/gin"
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
	// var contextPath string = "payments"

	// // Rate limits
	// var callbackLimit = middleware.InitializeRateLimiter(rate.Every(time.Second/5), 30)
	// var donateLimit = middleware.InitializeRateLimiter(rate.Every(time.Minute/2), 5)
	// var listLimit = middleware.InitializeRateLimiter(rate.Every(time.Second/2), 15)
	// var postLimit = middleware.InitializeRateLimiter(rate.Every(time.Minute/4), 6)

	// // Normal group
	// var norGroup = server.Group(contextPath)
	// norGroup.GET("/callback/:id", middleware.RateLimitMiddleware(callbackLimit), transport.CallbackTransaction)
	// norGroup.GET("", middleware.RateLimitMiddleware(listLimit), transport.GetPayments)

	// // Auth group
	// var authGroup = server.Group(contextPath, middleware.Authorize)
	// authGroup.POST("/donate", middleware.RateLimitMiddleware(donateLimit), transport.Donate)
	// authGroup.GET("/auth-callback/:id", middleware.AdminAuthorize, middleware.RateLimitMiddleware(donateLimit), transport.CallbackWithAuthTransaction)

	// // Admin group
	// var adminGroup = server.Group(contextPath, middleware.Authorize, middleware.AdminAuthorize)
	// authGroup.POST("/:id/approve", middleware.RateLimitMiddleware(postLimit), transport.ApprovePayment)
	// adminGroup.POST("/:id/refuse", middleware.RateLimitMiddleware(postLimit), transport.RefusePayment)

	var contextPath string = "payments"

	// Normal group
	var norGroup = server.Group(contextPath)
	norGroup.GET("/callback/:id", transport.CallbackTransaction)
	norGroup.GET("", transport.GetPayments)
	norGroup.GET("/:id", transport.GetPayment)

	// Auth group
	var authGroup = server.Group(contextPath, middleware.Authorize)
	authGroup.POST("/donate", transport.Donate)

	// Admin group
	var adminGroup = server.Group(contextPath, middleware.Authorize, middleware.AdminAuthorize)
	authGroup.POST("/:id/approve", transport.ApprovePayment)
	adminGroup.POST("/:id/refuse", transport.RefusePayment)
	adminGroup.POST("/auth-callback/:id", transport.CallbackWithAuthTransaction)
}
