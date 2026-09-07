package cmd

import (
	"log"
	"net/http"
	"os"
	"raise-child/constants/env"

	"strings"
	"time"

	"raise-child/docs"

	"github.com/fsnotify/fsnotify"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	swagger_files "github.com/swaggo/files"
	gin_swagger "github.com/swaggo/gin-swagger"
)

// Load .env file
func loadEnv() {
	godotenv.Load(".env")
}

// Enable CORS
func corsConfig(server *gin.Engine) {
	server.Use(cors.New(cors.Config{
		AllowAllOrigins: true,

		AllowMethods: []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS", "HEAD"},

		AllowHeaders: []string{
			"Origin",
			"Content-Type",
			"Accept",
			"Authorization",
			"Cache-Control",
			"X-Requested-With",
		},

		AllowCredentials: false,

		MaxAge: 12 * time.Hour,
	}))

	server.Use(func(ctx *gin.Context) {
		if ctx.Request.Method == http.MethodOptions {
			log.Printf("[PREFLIGHT] Origin=%s Path=%s RequestedMethod=%s RequestedHeaders=%s",
				ctx.GetHeader("Origin"),
				ctx.Request.URL.Path,
				ctx.GetHeader("Access-Control-Request-Method"),
				ctx.GetHeader("Access-Control-Request-Headers"),
			)
		}
		ctx.Next()
		if ctx.Request.Method == http.MethodOptions {
			log.Printf("[PREFLIGHT RESULT] Status=%d AllowOrigin=%s",
				ctx.Writer.Status(),
				ctx.Writer.Header().Get("Access-Control-Allow-Origin"),
			)
		}
	})
}

func setupSwagger(server *gin.Engine) {
	// Configure swagger info
	docs.SwaggerInfo.Title = "AgroTrust Server API"
	docs.SwaggerInfo.Version = "1.0"
	docs.SwaggerInfo.Schemes = []string{"https", "http"}
	docs.SwaggerInfo.Host = os.Getenv(env.SERVER_HOST)

	// Add swagger route
	server.GET("swagger/*any", gin_swagger.WrapHandler(swagger_files.Handler))
}

func watcherHttpOffConfig(logger *log.Logger) {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		logger.Println("ERROR creating watcher:", err)
		return
	}

	// get the current working directory and watch it
	currentDir, err := os.Getwd()
	if err != nil {
		logger.Println("ERROR getting current directory:", err)
		_ = watcher.Close()
		return
	}

	if err := watcher.Add(currentDir); err != nil {
		logger.Println("ERROR adding watcher to directory:", err)
		_ = watcher.Close()
		return
	}

	// watch for App_offline.htm and exit the program if present
	go func() {
		defer watcher.Close()
		for {
			select {
			case event := <-watcher.Events:
				if strings.HasSuffix(event.Name, "app_offline.htm") {
					logger.Println("Exiting due to app_offline.htm being present")
					os.Exit(0)
				}
			case err := <-watcher.Errors:
				logger.Println("Watcher error:", err)
			}
		}
	}()
}
