package apiroute

import (
	"raise-child/transport"
	"raise-child/util/middleware"

	"github.com/gin-gonic/gin"
)

func InitializeCenterRequestRoute(server *gin.Engine) {
	var contextPath string = "centers"

	// Normal group
	var norGroup = server.Group(contextPath)
	norGroup.GET("", transport.GetCenterRequests)
	norGroup.GET("/user/:id", transport.GetWalletCenterRequests)
	norGroup.GET("/:id", transport.GetCenterRequest)

	// Auth group
	var authGroup = server.Group(contextPath, middleware.Authorize)
	authGroup.POST("", transport.CreateCenterRequest)
	authGroup.POST("/:id/vote", transport.VoteCenterRequest)
	authGroup.POST("/:id/confirm", transport.ConfirmCenterRequest)
}
