package cmd

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	api_route "raise-child/api_route"
	"raise-child/business"
	"raise-child/constants/env/payment"
	"raise-child/constants/noti"
	"sync"
	"time"

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

	// Donor API endpoints
	api_route.InitializeDonorRoutes(server)

	// Region API endpoints
	api_route.InitializeRegionRoutes(server)

	// Upload Child Request API endpoints
	api_route.InitializeUploadChildRequestRoute(server)

	// Bank Profile API endpoints
	api_route.InitializeBankProfileRoutes(server)

	// Withdraw Proposal API endpoints
	api_route.InitializeWithdrawProposalRoute(server)

	// Admin API endpoints
	api_route.InitializeAdminRoute(server)

	// Center Request API endpoints
	api_route.InitializeCenterRequestRoute(server)

	// Gift API endpoints
	api_route.InitializeGiftRoute(server)

	// Noti API endpoints
	api_route.InitializeNotiRoute(server)

	// Pending Withdraw Proposal API endpoints
	api_route.InitializePendingWithdrawProposalRoute(server)

	// Child Pending Special Need Proposal APi endpoints
	api_route.InitializeChildPendingSpecialNeedProposalRoutes(server)

	// Task API endpoints
	api_route.InitializeTaskRoutes(server)

	// Task Proof API endpoints
	api_route.InitializeTaskProofRoutes(server)

	// Transaction Record API endpoints
	api_route.InitializeTransactionRecordRoutes(server)

	// Config API endpoints
	api_route.InitializeConfigRoutes(server)

	// Child Need API endpoints
	api_route.InitializeChildNeedRoutes(server)

	// Pending Pool Campaign API endpoints
	api_route.InitializePendingCampaignRoutes(server)

	// Pool Campaign API endpoints
	api_route.InitializeCampaignRoutes(server)

	// Center API enpoints
	api_route.InitializeCenterRoute(server)

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

func setupBackgroundService(ctx context.Context, wg *sync.WaitGroup) {
	var services = []func(context.Context, *sync.WaitGroup, time.Duration){
		processCenterBackgroundService,
		processRegistrationBackgroundService,
	}

	wg.Add(len(services))
	var duration time.Duration = time.Minute
	for _, service := range services {
		go service(ctx, wg, duration)
	}
}

func processCenterBackgroundService(ctx context.Context, wg *sync.WaitGroup, duration time.Duration) {
	defer wg.Done()
	var ticker = time.NewTicker(duration)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			if service, err := business.GenerateBackgroundService(); err == nil {
				service.ProcessBackgroundCenterRequests(ctx)
			}
		}
	}
}

func processRegistrationBackgroundService(ctx context.Context, wg *sync.WaitGroup, duration time.Duration) {
	defer wg.Done()
	var ticker = time.NewTicker(duration)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			if service, err := business.GenerateBackgroundService(); err == nil {
				service.ProcessBackgroundRegistrationRequests(ctx)
			}
		}
	}
}

// func processUploadChildBackgroundService(ctx context.Context, wg *sync.WaitGroup, duration time.Duration) {
// 	defer wg.Done()
// 	var ticker = time.NewTicker(duration)
// 	defer ticker.Stop()
// 	for {
// 		select {
// 		case <-ctx.Done():
// 			return
// 		case <-ticker.C:
// 			if service, err := business.GenerateBackgroundService(); err == nil {
// 				service.ProcessBackgroundUploadChildRequests(ctx)
// 			}
// 		}
// 	}
// }
