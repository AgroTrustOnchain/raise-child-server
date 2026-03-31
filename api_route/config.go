package apiroute

import (
	"raise-child/transport"
	"raise-child/util/middleware"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
)

func InitializeConfigRoutes(server *gin.Engine) {
	var contextPath string = "configs"

	// Rate limits
	var editDatesLimit = middleware.InitializeRateLimiter(rate.Every(time.Minute*10), 3)

	// Admin group
	var adminGroup = server.Group(contextPath, middleware.Authorize, middleware.AdminAuthorize)
	adminGroup.PUT("/books-need-edit-dates", middleware.RateLimitMiddleware(editDatesLimit), transport.UpdateChildEditBooksNeedDates)
	adminGroup.PUT("/meal-need-edit-dates", middleware.RateLimitMiddleware(editDatesLimit), transport.UpdateChildEditMealNeedDates)
	adminGroup.PUT("/health-insurance-need-edit-dates", middleware.RateLimitMiddleware(editDatesLimit), transport.UpdateChildEditHealthInsuranceNeedDates)
}
