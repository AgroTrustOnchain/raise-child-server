package apiroute

import (
	"raise-child/transport"
	"raise-child/util/middleware"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
)

func InitializeTransactionRecordRoutes(server *gin.Engine) {
	var contextPath string = "profiles"

	// Rate limits
	var viewLimit = middleware.InitializeRateLimiter(rate.Every(time.Second/5), 20)
	var listLimit = middleware.InitializeRateLimiter(rate.Every(time.Second/2), 10)

	// Normal group
	var norGroup = server.Group(contextPath)
	norGroup.GET("", middleware.RateLimitMiddleware(listLimit), transport.GetTransactionRecords)
	norGroup.GET("/:id", middleware.RateLimitMiddleware(viewLimit), transport.GetTransactionRecord)
}
