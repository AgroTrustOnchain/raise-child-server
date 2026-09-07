package apiroute

import (
	"raise-child/transport"
	"raise-child/util/middleware"

	"github.com/gin-gonic/gin"
)

func InitializePoolRoutes(server *gin.Engine) {
	var contextPath string = "pools"
	var leaderGroup = server.Group(contextPath, middleware.Authorize, middleware.LeaderAuthorize)
	leaderGroup.GET("/leader/:id", transport.GetLeaderPool)
}
