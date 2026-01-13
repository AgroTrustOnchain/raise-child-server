package cmd

import (
	"os"
	"raise-child/constants/env"
	"raise-child/constants/shared"
	"raise-child/util"

	"github.com/gin-gonic/gin"
)

func Execute() {
	var errLogger = util.GetLogConfig(shared.ERROR_LEVEL)

	// Load env
	loadEnv(errLogger)

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

	// Run server
	if err := server.Run(":" + apiPort); err != nil {
		errLogger.Fatalln("Error run server - " + err.Error())
	}

	var infoLogger = util.GetLogConfig(shared.INFO_LEVEL)
	infoLogger.Println("Server starts on port ", apiPort)
}
