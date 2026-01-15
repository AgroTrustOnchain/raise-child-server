package apiroute

import (
	"raise-child/transport"
	"raise-child/util/middleware"

	"github.com/gin-gonic/gin"
)

func InitializeBankProfileRoutes(server *gin.Engine) {
	var contextPath string = "banks"

	// Normal group
	var norGroup = server.Group(contextPath)
	norGroup.GET("/:id", transport.GetBankProfile)

	// Auth group
	var authGroup = server.Group(contextPath, middleware.Authorize)
	authGroup.POST("", transport.CreateBankProfile)
	authGroup.PUT("/:id", transport.UpdateBankProfile)
}
