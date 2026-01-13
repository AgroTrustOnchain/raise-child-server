package apiroute

import (
	"raise-child/transport"

	"github.com/gin-gonic/gin"
)

func InitializeSponsorRoutes(server *gin.Engine) {
	var contextPath string = "sponsors"

	// Normal group
	var norGroup = server.Group(contextPath)
	norGroup.GET("", transport.GetSponsors)
	norGroup.GET("/:id", transport.GetSponsor)
}
