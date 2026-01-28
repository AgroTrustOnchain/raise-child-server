package cmd

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"raise-child/constants/env"
	"raise-child/constants/shared"
	"raise-child/util"
	"sync"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
)

func Execute() {
	var errLogger = util.GetLogConfig(shared.ERROR_LEVEL)

	// Load env
	loadEnv(errLogger)

	// Initialize context for backgroun goroutines management
	ctx, cancel := context.WithCancel(context.Background())
	var wg sync.WaitGroup

	// Run background services
	setupBackgroundService(ctx, &wg)

	// Initialize gin server for API
	var server = gin.Default()

	// Config CORS for requests
	corsConfig(server)

	// Get API port
	var apiPort string = os.Getenv("HTTP_PLATFORM_PORT")
	if apiPort == "" {
		apiPort = os.Getenv(env.API_PORT)
	}

	// Set up API routes
	setupApiRoutes(server)

	// Set up swagger
	setupSwagger(server, apiPort)

	// Watcher http offline
	watcherHttpOffConfig(errLogger)

	// Setup payments
	setupPayments(errLogger)

	// Convert gin server to HTTP server
	var httpServer = &http.Server{
		Addr:    ":" + apiPort,
		Handler: server,
	}

	// Execute gin server in another goroutine
	var infoLogger = util.GetLogConfig(shared.INFO_LEVEL)
	go func() {
		infoLogger.Println("Server starts on port ", apiPort)
		if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			errLogger.Fatalln("Error run server - " + err.Error())
		}
	}()

	// Listen for gracefull shutdown
	var quit = make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	infoLogger.Println("Server shutting down...")

	// Flag signal for all goroutines
	cancel()

	// Shutdwon gin server
	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer shutdownCancel()

	if err := httpServer.Shutdown(shutdownCtx); err != nil {
		errLogger.Println("Error while shutting down gin server: " + err.Error())
	}

	// Wait for all goroutines to finish
	wg.Wait()

	infoLogger.Println("Server shutdown")
}
