package apiroute

import (
	"raise-child/transport"
	"raise-child/util/middleware"

	"github.com/gin-gonic/gin"
)

func InitializeProfileRoutes(server *gin.Engine) {
	var contextPath string = "profiles"

	// Auth group
	var authGroup = server.Group(contextPath, middleware.Authorize)
	authGroup.POST("/:id", transport.UploadProfile)
}
