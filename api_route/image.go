package apiroute

import (
	"raise-child/transport"
	"raise-child/util/middleware"

	"github.com/gin-gonic/gin"
)

func InitializeImageRoute(server *gin.Engine) {
	var contextPath string = "images"

	// Auth group
	var authGroup = server.Group(contextPath, middleware.Authorize)
	authGroup.GET("/presigned-url", transport.GetPresignedUrl)
}
