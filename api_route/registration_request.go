package apiroute

import (
	"raise-child/transport"
	"raise-child/util/middleware"

	"github.com/gin-gonic/gin"
)

func InitializeRegistraionRequestRoute(server *gin.Engine) {
	var contextPath string = "registrations"

	// Normal group
	var norGroup = server.Group(contextPath)
	norGroup.GET("", transport.GetRegistrationRequests)
	norGroup.GET("/user/:id", transport.GetWalletRegistrationRequests)
	norGroup.GET("/:id", transport.GetRegistrationRequest)

	// Auth group
	var authGroup = server.Group(contextPath, middleware.Authorize)
	authGroup.POST("", transport.CreateRegistrationRequest)
	authGroup.POST("/:id/vote", transport.VoteRegistrationRequest)
	authGroup.POST("/:id/confirm", transport.ConfirmRegistrationRequest)
}
