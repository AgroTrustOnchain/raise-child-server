package apiroute

import (
	"raise-child/transport"
	"raise-child/util/middleware"

	"github.com/gin-gonic/gin"
)

func InitializeAdminRoute(server *gin.Engine) {
	var contextPath string = "admins"

	// Auth group
	var authGroup = server.Group(contextPath, middleware.Authorize)
	authGroup.POST("", transport.UpdatePublisherInfo)

	// Normal group
	var norGroup = server.Group(contextPath)
	norGroup.GET("", transport.GetAdmins)
}
