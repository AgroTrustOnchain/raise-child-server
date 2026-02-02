package apiroute

import (
	"raise-child/transport"
	"raise-child/util/middleware"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
)

// func InitializeChildRoutes(server *gin.Engine) {
// 	var contextPath string = "children"

// 	// Normal group
// 	var norGroup = server.Group(contextPath)
// 	norGroup.GET("", transport.GetChildren)
// 	norGroup.GET("/:id", transport.GetChild)

// 	// Auth group
// 	var authGroup = server.Group(contextPath, middleware.Authorize)
// 	authGroup.PUT("/metadata/string/:id", transport.AddChildStringMetadata)
// 	authGroup.PUT("/metadata/number/:id", transport.AddChildNumberMetadata)
// 	authGroup.POST("", transport.UploadChild)
// }

func InitializeChildRoutes(server *gin.Engine) {
	var contextPath string = "children"

	// Rate limits
	var listLimit = middleware.InitalizeRateLimiter(rate.Every(time.Second/2), 10)
	var detailLimit = middleware.InitalizeRateLimiter(rate.Every(time.Second/5), 20)
	var uploadLimit = middleware.InitalizeRateLimiter(rate.Every(time.Minute/30), 30)
	var metadataLimit = middleware.InitalizeRateLimiter(rate.Every(time.Second/1), 10)

	// Normal group
	var norGroup = server.Group(contextPath)
	norGroup.GET("", middleware.RateLimitMiddleware(listLimit), transport.GetChildren)
	norGroup.GET("/:id", middleware.RateLimitMiddleware(detailLimit), transport.GetChild)

	// Auth group
	var authGroup = server.Group(contextPath, middleware.Authorize)
	authGroup.PUT("/metadata/string/:id", middleware.RateLimitMiddleware(metadataLimit), transport.AddChildStringMetadata)
	authGroup.PUT("/metadata/number/:id", middleware.RateLimitMiddleware(metadataLimit), transport.AddChildNumberMetadata)
	authGroup.POST("", middleware.RateLimitMiddleware(uploadLimit), transport.UploadChild)
}
