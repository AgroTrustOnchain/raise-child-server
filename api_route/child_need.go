package apiroute

import (
	"raise-child/transport"
	"raise-child/util/middleware"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
)

func InitializeChildNeedRoutes(server *gin.Engine) {
	var contextPath string = "child-needs"

	// Rate limits
	var detailLimit = middleware.InitializeRateLimiter(rate.Every(time.Second/5), 20)

	// Normal group
	var norGroup = server.Group(contextPath)
	norGroup.GET("/books-need/:id", middleware.RateLimitMiddleware(detailLimit), transport.GetChildBooksNeed)
	norGroup.GET("/health-insurance-need/:id", middleware.RateLimitMiddleware(detailLimit), transport.GetChildHealthInsuranceNeed)
	norGroup.GET("/meal-need/:id", middleware.RateLimitMiddleware(detailLimit), transport.GetChildMealNeed)
}
