package apiroute

import (
	"raise-child/transport"

	"github.com/gin-gonic/gin"
)

// func InitializeDonorRoutes(server *gin.Engine) {
// 	var contextPath string = "donors"

// 	// Normal group
// 	var norGroup = server.Group(contextPath)
// 	norGroup.GET("", transport.GetDonors)
// 	norGroup.GET("/:id", transport.GetDonor)
// }

func InitializeDonorRoutes(server *gin.Engine) {
	// var contextPath string = "donors"

	// // Rate limits
	// var listLimit = middleware.InitializeRateLimiter(rate.Every(time.Second/2), 20)
	// var detailLimit = middleware.InitializeRateLimiter(rate.Every(time.Second/5), 30)

	// // Normal group
	// var norGroup = server.Group(contextPath)
	// norGroup.GET("", middleware.RateLimitMiddleware(listLimit), transport.GetDonors)
	// norGroup.GET("/:id", middleware.RateLimitMiddleware(detailLimit), transport.GetDonor)

	var contextPath string = "donors"

	// Normal group
	var norGroup = server.Group(contextPath)
	norGroup.GET("", transport.GetDonors)
	norGroup.GET("/:id", transport.GetDonor)
}
