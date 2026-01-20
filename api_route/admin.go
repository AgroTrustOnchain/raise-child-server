package apiroute

import (
	"raise-child/transport"
	"raise-child/util/middleware"

	"github.com/gin-gonic/gin"
)

func InitializeAdminRoute(server *gin.Engine) {
	var contextPath string = "admins"

	var authGroup = server.Group(contextPath, middleware.Authorize)
	authGroup.POST("", transport.UpdatePublisherInfo)
}
