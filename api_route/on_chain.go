package apiroute

import (
	"raise-child/transport"
	"raise-child/util/middleware"

	"github.com/gin-gonic/gin"
)

// func InitializeOnChainRoutes(server *gin.Engine) {
// 	var contextPath string = "tx"
// 	var authGroup = server.Group(contextPath, middleware.Authorize)

// 	authGroup.POST("/execute", transport.ExecuteTransaction)
// 	authGroup.POST("/build/money", transport.BuildMoneyTransaction)
// }

func InitializeOnChainRoutes(server *gin.Engine) {
	// var contextPath string = "tx"
	// var authGroup = server.Group(contextPath, middleware.Authorize)
	// authGroup.POST("/execute", middleware.RateLimitMiddleware(middleware.InitializeRateLimiter(rate.Every(time.Second/5), 20)), transport.ExecuteTransaction)

	var contextPath string = "tx"
	var authGroup = server.Group(contextPath, middleware.Authorize)
	authGroup.POST("/execute", transport.ExecuteTransaction)
}
