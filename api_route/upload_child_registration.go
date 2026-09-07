package apiroute

import (
	"raise-child/transport"
	"raise-child/util/middleware"

	"github.com/gin-gonic/gin"
)

// func InitializeUploadChildRequestRoute(server *gin.Engine) {
// 	var contextPath string = "child-upload-reqs"

// 	// Normal group
// 	var norGroup = server.Group(contextPath)
// 	norGroup.GET("", transport.GetUploadChildRequests)
// 	norGroup.GET("/user/:id", transport.GetWalletUploadChildRequests)
// 	norGroup.GET("/:id", transport.GetUploadChildRequest)

// 	// Auth group
// 	var authGroup = server.Group(contextPath, middleware.Authorize)
// 	authGroup.POST("", transport.CreateUploadChildRequest)
// 	authGroup.POST("/:id/vote", transport.VoteUploadChildRequest)
// 	authGroup.POST("/:id/confirm", transport.ConfirmUploadChildRequest)
// }

func InitializeUploadChildRequestRoute(server *gin.Engine) {
	var contextPath string = "child-upload-reqs"

	// Normal group
	var norGroup = server.Group(contextPath)
	norGroup.GET("", transport.GetUploadChildRequests)
	norGroup.GET("/user/:id", transport.GetWalletUploadChildRequests)
	norGroup.GET("/:id", transport.GetUploadChildRequest)

	// Auth group
	var authGroup = server.Group(contextPath, middleware.Authorize)
	authGroup.POST("", transport.CreateUploadChildRequest)

	// Admin group
	var leaderGroup = server.Group(contextPath, middleware.Authorize, middleware.LeaderAuthorize)
	leaderGroup.POST("/:id/review", transport.ReviewUploadChildRequest)
}
