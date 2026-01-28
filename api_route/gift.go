package apiroute

import (
	"raise-child/transport"
	"raise-child/util/middleware"

	"github.com/gin-gonic/gin"
)

func InitializeGiftRoute(server *gin.Engine) {
	var contextPath string = "gifts"

	// Normal group
	var norGroup = server.Group(contextPath)
	norGroup.GET("", transport.GetGift)
	norGroup.GET("/child/:id", transport.GetGiftsOfChild)

	// Auth group
	var authGroup = server.Group(contextPath, middleware.Authorize)
	authGroup.POST("", transport.CreateGift)
	authGroup.POST("/:id/confirm", transport.ConfirmReceiveGift)
}
