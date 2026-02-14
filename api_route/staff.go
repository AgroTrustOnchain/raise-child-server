package apiroute

import (
	"raise-child/transport"
	"raise-child/util/middleware"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
)

// func InitializeStaffRoute(server *gin.Engine) {
// 	var contextPath string = "staffs"

// 	// Normal group
// 	var norGroup = server.Group(contextPath)
// 	norGroup.GET("/:id", transport.GetStaff)
// 	norGroup.GET("", transport.GetStaffs)
// }

func InitializeStaffRoute(server *gin.Engine) {
	var contextPath string = "staffs"

	// Rate limits
	var listLimit = middleware.InitializeRateLimiter(rate.Every(time.Second/2), 20)
	var detailLimit = middleware.InitializeRateLimiter(rate.Every(time.Second/5), 30)

	// Normal group
	var norGroup = server.Group(contextPath)
	norGroup.GET("", middleware.RateLimitMiddleware(listLimit), transport.GetStaffs)
	norGroup.GET("/:id", middleware.RateLimitMiddleware(detailLimit), transport.GetStaff)
}
