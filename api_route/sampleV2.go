package apiroute

import (
	"raise-child/transport"

	"github.com/gin-gonic/gin"
)

func InitializeMintRoute(server *gin.Engine) {
	var contextPath string = "mints"
	var authGroup = server.Group(contextPath)
	authGroup.POST("", transport.MintCap)
	authGroup.POST("/caps", transport.MintCaps)
}
