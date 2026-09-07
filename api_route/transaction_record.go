package apiroute

import (
	"raise-child/transport"

	"github.com/gin-gonic/gin"
)

func InitializeTransactionRecordRoutes(server *gin.Engine) {
	// var contextPath string = "tx-records"

	// // Rate limits
	// var viewLimit = middleware.InitializeRateLimiter(rate.Every(time.Second/5), 20)
	// var listLimit = middleware.InitializeRateLimiter(rate.Every(time.Second/2), 10)

	// // Normal group
	// var norGroup = server.Group(contextPath)
	// norGroup.GET("", middleware.RateLimitMiddleware(listLimit), transport.GetTransactionRecords)
	// norGroup.GET("/:id", middleware.RateLimitMiddleware(viewLimit), transport.GetTransactionRecord)

	var contextPath string = "tx-records"

	// Normal group
	var norGroup = server.Group(contextPath)
	norGroup.GET("", transport.GetTransactionRecords)
	norGroup.GET("/:id", transport.GetTransactionRecord)
}
