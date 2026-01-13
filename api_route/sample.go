package apiroute

import (
	"raise-child/transport"

	"github.com/gin-gonic/gin"
)

func InitSampleRoutes(server *gin.Engine) {
	var contextPath string = "tx"
	var txGroup = server.Group(contextPath)
	txGroup.POST("/build", transport.BuildCampaignTransaction)
	//txGroup.POST("/execute", transport.ExecuteCampaignTransaction)
}
