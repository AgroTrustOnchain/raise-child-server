package apiroute

import (
	"raise-child/transport"

	"github.com/gin-gonic/gin"
)

func InitializeChildNeedRoutes(server *gin.Engine) {
	// var contextPath string = "child-needs"

	// // Rate limits
	// var detailLimit = middleware.InitializeRateLimiter(rate.Every(time.Second/5), 20)

	// // Normal group
	// var norGroup = server.Group(contextPath)
	// norGroup.GET("/books-need/:id", middleware.RateLimitMiddleware(detailLimit), transport.GetChildBooksNeed)
	// norGroup.GET("/health-insurance-need/:id", middleware.RateLimitMiddleware(detailLimit), transport.GetChildHealthInsuranceNeed)
	// norGroup.GET("/meal-need/:id", middleware.RateLimitMiddleware(detailLimit), transport.GetChildMealNeed)

	var contextPath string = "child-needs"

	// Normal group
	var norGroup = server.Group(contextPath)
	norGroup.GET("/books-need/:id", transport.GetChildBooksNeed)
	norGroup.GET("/health-insurance-need/:id", transport.GetChildHealthInsuranceNeed)
	norGroup.GET("/meal-need/:id", transport.GetChildMealNeed)
}
