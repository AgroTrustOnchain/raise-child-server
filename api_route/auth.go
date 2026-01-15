package apiroute

import (
	"raise-child/transport"
	"raise-child/util/middleware"

	"github.com/gin-gonic/gin"
)

func InitializeAuthHandlerRoutes(server *gin.Engine) {
	var contextPath string = "auth"

	var norGroup = server.Group(contextPath)
	norGroup.POST("/login", transport.Login)
	norGroup.GET("/nonce/:address", transport.GetNonce)
	norGroup.GET("/salt/:id", transport.GetSalt)

	var authGroup = server.Group(contextPath, middleware.Authorize)
	authGroup.POST("/logout", transport.Logout)
}
