package apiroute

import (
	"raise-child/transport"
	"raise-child/util/middleware"

	"github.com/gin-gonic/gin"
)

func InitializeConfigRoutes(server *gin.Engine) {
	// var contextPath string = "configs"

	// // Rate limits
	// var editDatesLimit = middleware.InitializeRateLimiter(rate.Every(time.Minute*10), 3)

	// // Admin group
	// var adminGroup = server.Group(contextPath, middleware.Authorize, middleware.AdminAuthorize)
	// adminGroup.PUT("/books-need-edit-dates", middleware.RateLimitMiddleware(editDatesLimit), transport.UpdateChildEditBooksNeedDates)
	// adminGroup.PUT("/meal-need-edit-dates", middleware.RateLimitMiddleware(editDatesLimit), transport.UpdateChildEditMealNeedDates)
	// adminGroup.PUT("/health-insurance-need-edit-dates", middleware.RateLimitMiddleware(editDatesLimit), transport.UpdateChildEditHealthInsuranceNeedDates)

	var contextPath string = "configs"

	// Admin group
	var adminGroup = server.Group(contextPath, middleware.Authorize, middleware.AdminAuthorize)
	adminGroup.PUT("/books-need-edit-dates", transport.UpdateChildEditBooksNeedDates)
	adminGroup.PUT("/meal-need-edit-dates", transport.UpdateChildEditMealNeedDates)
	adminGroup.PUT("/health-insurance-need-edit-dates", transport.UpdateChildEditHealthInsuranceNeedDates)
	adminGroup.PUT("/child-special-need-dao", transport.EditSpecialNeedDao)
}
