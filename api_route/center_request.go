package apiroute

import (
	"raise-child/transport"
	"raise-child/util/middleware"

	"github.com/gin-gonic/gin"
)

func InitializeCenterRequestRoute(server *gin.Engine) {
	var contextPath string = "center-reqs"

	// Normal group
	var norGroup = server.Group(contextPath)
	norGroup.GET("", transport.GetCenterRequests)
	norGroup.GET("/user/:id", transport.GetWalletCenterRequests)
	norGroup.GET("/:id", transport.GetCenterRequest)

	// Staff group
	var staffGroup = server.Group(contextPath, middleware.Authorize, middleware.StaffRoleAuthorize)
	staffGroup.POST("/:id/vote", transport.VoteCenterRequest)

	// Leader group
	var leaderGroup = server.Group(contextPath, middleware.Authorize, middleware.LeaderAuthorize)
	leaderGroup.POST("", transport.CreateCenterRequest)
	leaderGroup.POST("/:id/confirm", transport.ConfirmCenterRequest)
}
