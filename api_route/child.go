package apiroute

import (
	"raise-child/transport"
	"raise-child/util/middleware"

	"github.com/gin-gonic/gin"
)

func InitializeChildRoutes(server *gin.Engine) {
	var contextPath string = "children"

	// Normal group
	var norGroup = server.Group(contextPath)
	norGroup.GET("", transport.GetChildren)
	norGroup.GET("/:id", transport.GetChild)

	// Auth group
	var authGroup = server.Group(contextPath, middleware.Authorize)
	authGroup.PUT("/metadata/string/:id", transport.AddChildStringMetadata)
	authGroup.PUT("/metadata/number/:id", transport.AddChildNumberMetadata)
	authGroup.POST("", transport.UploadChild)
}
