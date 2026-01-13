package apiroute

import (
	"raise-child/transport"

	"github.com/gin-gonic/gin"
)

func InitializeStaffRoute(server *gin.Engine) {
	var contextPath string = "staffs"

	// Normal group
	var norGroup = server.Group(contextPath)
	norGroup.GET("/:id", transport.GetStaff)
	norGroup.GET("", transport.GetStaffs)
}
