package apiroute

import (
	"raise-child/transport"
	"raise-child/util/middleware"

	"github.com/gin-gonic/gin"
)

func InitializeTaskRoutes(server *gin.Engine) {
	// var contextPath string = "tasks"

	// // Rate limits
	// var listLimit = middleware.InitializeRateLimiter(rate.Every(time.Second/2), 15)
	// var viewLimit = middleware.InitializeRateLimiter(rate.Every(time.Second/5), 20)
	// var postLimit = middleware.InitializeRateLimiter(rate.Every(time.Minute/6), 6)
	// var reviewLimit = middleware.InitializeRateLimiter(rate.Every(time.Minute/10), 20)

	// // Normal group
	// var norGroup = server.Group(contextPath)
	// norGroup.GET("", middleware.RateLimitMiddleware(listLimit), transport.GetTasks)
	// norGroup.GET("/:id", middleware.RateLimitMiddleware(viewLimit), transport.GetTask)

	// // Staff group
	// var staffGroup = server.Group(contextPath, middleware.Authorize, middleware.StaffRoleAuthorize)
	// staffGroup.POST("/:id/claim", middleware.RateLimitMiddleware(postLimit), transport.ClaimTask)

	// // Manager group
	// var managerGroup = server.Group(contextPath, middleware.Authorize, middleware.ManagerRoleAuthorize)
	// managerGroup.POST("", middleware.RateLimitMiddleware(postLimit), transport.CreateTask)
	// managerGroup.POST("/:id/review", middleware.RateLimitMiddleware(reviewLimit), transport.ReviewAssignedProfileOfTask)

	var contextPath string = "tasks"

	// Normal group
	var norGroup = server.Group(contextPath)
	norGroup.GET("", transport.GetTasks)
	norGroup.GET("/:id", transport.GetTask)

	// Staff group
	var staffGroup = server.Group(contextPath, middleware.Authorize, middleware.StaffRoleAuthorize)
	staffGroup.POST("/:id/claim", transport.ClaimTask)
	staffGroup.GET("/region/staff/:id", transport.GetTasksOfRegionOnUser)
	staffGroup.GET("/staff/:id", transport.GetUserTasks)

	// Manager group
	var managerGroup = server.Group(contextPath, middleware.Authorize, middleware.ManagerRoleAuthorize)
	managerGroup.POST("", transport.CreateTask)
}
