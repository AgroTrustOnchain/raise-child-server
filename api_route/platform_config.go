package apiroute

import (
	"raise-child/transport"
	"raise-child/util/middleware"

	"github.com/gin-gonic/gin"
)

func InitializePlatformConfigsRoutes(server *gin.Engine) {
	var contextPath string = "platform-configs"

	// Normal group
	var norGroup = server.Group(contextPath)
	norGroup.GET("/numeric", transport.GetNumericPlatformConfigs)
	norGroup.GET("/string", transport.GetStringPlatformConfigs)
	norGroup.GET("/numeric/:id", transport.GetNumericPlatformConfig)
	norGroup.GET("/string/:id", transport.GetStringPlatformConfig)

	// Manager group
	var adminGroup = server.Group(contextPath, middleware.Authorize, middleware.AdminAuthorize)
	adminGroup.PUT("/numeric/:id", transport.UpdateNumericPlatformConfig)
	adminGroup.PUT("/string/:id", transport.UpdateStringPlatformConfig)
}
