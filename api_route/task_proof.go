package apiroute

import (
	"raise-child/transport"
	"raise-child/util/middleware"

	"github.com/gin-gonic/gin"
)

func InitializeTaskProofRoutes(server *gin.Engine) {
	// var contextPath string = "task-proofs"

	// // Rate limits
	// var listLimit = middleware.InitializeRateLimiter(rate.Every(time.Second/2), 15)
	// var viewLimit = middleware.InitializeRateLimiter(rate.Every(time.Second/5), 20)
	// var submitLimit = middleware.InitializeRateLimiter(rate.Every(time.Minute/6), 6)
	// var reviewLimit = middleware.InitializeRateLimiter(rate.Every(time.Minute/10), 20)

	// // Normal group
	// var norGroup = server.Group(contextPath)
	// norGroup.GET("", middleware.RateLimitMiddleware(listLimit), transport.GetTaskProofs)
	// norGroup.GET("/:id", middleware.RateLimitMiddleware(viewLimit), transport.GetTaskProof)

	// // Staff group
	// var staffGroup = server.Group(contextPath, middleware.Authorize, middleware.StaffRoleAuthorize)
	// staffGroup.POST("/task/:id/submit", middleware.RateLimitMiddleware(submitLimit), transport.SubmitTaskProof)

	// // Leader group
	// var leaderGroup = server.Group(contextPath, middleware.Authorize, middleware.LeaderAuthorize)
	// leaderGroup.POST("/:id/approve", middleware.RateLimitMiddleware(reviewLimit), transport.ApproveTaskProof)
	// leaderGroup.POST("/:id/refuse", middleware.RateLimitMiddleware(reviewLimit), transport.RefuseTaskProof)

	var contextPath string = "task-proofs"

	// Normal group
	var norGroup = server.Group(contextPath)
	norGroup.GET("", transport.GetTaskProofs)
	norGroup.GET("/:id", transport.GetTaskProof)

	// Staff group
	var staffGroup = server.Group(contextPath, middleware.Authorize, middleware.StaffRoleAuthorize)
	staffGroup.POST("/task/:id/submit", transport.SubmitTaskProof)

	// Leader group
	var leaderGroup = server.Group(contextPath, middleware.Authorize, middleware.LeaderAuthorize)
	leaderGroup.POST("/:id/approve", transport.ApproveTaskProof)
	leaderGroup.POST("/:id/refuse", transport.RefuseTaskProof)
}
