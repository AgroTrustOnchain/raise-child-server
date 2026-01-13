package apiroute

import (
	"raise-child/transport"
	"raise-child/util/middleware"

	"github.com/gin-gonic/gin"
)

func InitializePaymentsRoutes(server *gin.Engine) {
	var contextPath string = "payments"

	// Normal group
	var norGroup = server.Group(contextPath)
	norGroup.GET("/callback/:id", transport.CallbackTransaction)

	// Auth group
	var authGroup = server.Group(contextPath, middleware.Authorize)
	authGroup.POST("/donate", transport.Donate)
}
