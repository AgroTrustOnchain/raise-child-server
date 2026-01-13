package apiroute

import (
	"raise-child/transport"
	"raise-child/util/middleware"

	"github.com/gin-gonic/gin"
)

func InitializeOnChainRoutes(server *gin.Engine) {
	var contextPath string = "tx"
	var authGroup = server.Group(contextPath, middleware.Authorize)

	authGroup.POST("/execute", transport.ExecuteTransaction)
	authGroup.POST("/build/money", transport.BuildMoneyTransaction)
}
