package apiroute

import (
	"raise-child/transport"
	"raise-child/util/middleware"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
)

// func InitializeSponsorRoutes(server *gin.Engine) {
// 	var contextPath string = "sponsors"

// 	// Normal group
// 	var norGroup = server.Group(contextPath)
// 	norGroup.GET("", transport.GetSponsors)
// 	norGroup.GET("/:id", transport.GetSponsor)
// }

func InitializeSponsorRoutes(server *gin.Engine) {
	var contextPath string = "sponsors"

	// Rate limits
	var listLimit = middleware.InitalizeRateLimiter(rate.Every(time.Second/2), 20)
	var detailLimit = middleware.InitalizeRateLimiter(rate.Every(time.Second/5), 30)

	// Normal group
	var norGroup = server.Group(contextPath)
	norGroup.GET("", middleware.RateLimitMiddleware(listLimit), transport.GetSponsors)
	norGroup.GET("/:id", middleware.RateLimitMiddleware(detailLimit), transport.GetSponsor)
}
