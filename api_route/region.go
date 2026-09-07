package apiroute

import (
	"raise-child/transport"
	"raise-child/util/middleware"

	"github.com/gin-gonic/gin"
)

func InitializeRegionRoutes(server *gin.Engine) {
	var contextPath string = "regions"

	// Normal group
	var norGroup = server.Group(contextPath)
	norGroup.GET("", transport.GetRegions)
	norGroup.GET("/established", transport.GetEstablishedRegions)
	norGroup.GET("/established/:region", transport.GetRegionDetail)
	norGroup.GET("/supported-suggestions", transport.GetSupportedRegionSuggestions)
	norGroup.GET("/supported-suggestions/:id", transport.GetSupportedRegionSuggestion)

	// Optional Auth group
	var optionalAuthGroup = server.Group(contextPath, middleware.OptinalAuthorize)
	optionalAuthGroup.GET("/supported-suggestionsv2", transport.GetSupportedRegionSuggestionsV2)

	// Auth group
	var authGroup = server.Group(contextPath, middleware.Authorize)
	authGroup.GET("/user/:id/supported-suggestions", transport.GetWalletSupportedRegionSuggestions)
	authGroup.POST("/supported-suggestions", transport.CreateSupportedRegionSuggestion)

	// Admin group
	var adminGroup = server.Group(contextPath, middleware.Authorize, middleware.AdminAuthorize)
	adminGroup.GET("/admin/supported-suggestions", transport.AdminGetSupportedRegionSuggestions)
	adminGroup.POST("/supported-suggestions/:id/review", transport.ReviewRegionSuggestion)
}
