package apiroute

import (
	"raise-child/transport"

	"github.com/gin-gonic/gin"
)

func InitializeRegionRoutes(server *gin.Engine) {
	var contextPath string = "regions"

	// Normal group
	var norGroup = server.Group(contextPath)
	norGroup.GET("", transport.GetRegions)
}
