package apiroute

import (
	"raise-child/transport"
	"raise-child/util/middleware"

	"github.com/gin-gonic/gin"
)

func InitializeWithdrawProposalRoute(server *gin.Engine) {
	var contextPath string = "withdraw-proposals"

	// Normal group
	var norGroup = server.Group(contextPath)
	norGroup.GET("", transport.GetWithdrawProposals)
	norGroup.GET("/:id", transport.GetWithdrawProposal)

	// Auth group
	var authGroup = server.Group(contextPath, middleware.Authorize)
	authGroup.POST("", transport.CreateWithdrawProposal)
	authGroup.POST("/:id/vote", transport.VoteWithdrawProposal)
	authGroup.POST("/:id/confirm", transport.ConfirmWithdrawProposal)
	authGroup.POST("/:id/main-pool-confirm", transport.ConfirmMainPoolWithdrawProposal)
}
