package apiroute

import (
	"raise-child/transport"
	"raise-child/util/middleware"

	"github.com/gin-gonic/gin"
)

func InitializeORCRoute(server *gin.Engine) {
	var contextPath string = "ocr"

	// Auth group
	var authGroup = server.Group(contextPath, middleware.Authorize)
	authGroup.POST("/extract-child-info", transport.GetExtractChildInfo)
}
