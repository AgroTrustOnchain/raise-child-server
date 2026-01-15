package cmd

import (
	"fmt"
	"log"
	"net/http"
	"os"
	api_route "raise-child/api_route"
	"raise-child/constants/env/payment"
	"raise-child/constants/noti"

	"github.com/gin-gonic/gin"
	"github.com/payOSHQ/payos-lib-golang"
)

func setupApiRoutes(server *gin.Engine) {
	// Auth API endpoints
	api_route.InitializeAuthHandlerRoutes(server)

	// On-chain API endpoints
	api_route.InitializeOnChainRoutes(server)

	// Child API endpoints
	api_route.InitializeChildRoutes(server)

	// Registraion Request API endpoints
	api_route.InitializeRegistrationRequestRoute(server)

	// Profile API endpoints
	api_route.InitializeProfileRoutes(server)

	// Payment API endpoints
	api_route.InitializePaymentsRoutes(server)

	// Sponsor API endpoints
	api_route.InitializeSponsorRoutes(server)

	// Region API endpoints
	api_route.InitializeRegionRoutes(server)

	// Upload Child Request API endpoints
	api_route.InitializeUploadChildRequestRoute(server)

	// Bank Profile API endpoints
	api_route.InitializeBankProfileRoutes(server)

	// Default route to Swagger documentation
	server.GET("/", func(ctx *gin.Context) {
		ctx.Redirect(http.StatusMovedPermanently, "/swagger/index.html#")
	})
}

func setupPayments(errLogger *log.Logger) {
	// Payos
	if err := payos.Key(os.Getenv(payment.PAYOS_CLIENT_ID), os.Getenv(payment.PAYOS_API_KEY), os.Getenv(payment.PAYOS_CHECKSUM_KEY)); err != nil {
		errLogger.Println(fmt.Sprintf(noti.PAYMENT_INIT_ENV_ERR_MSG, "payos") + err.Error())
	}
}
