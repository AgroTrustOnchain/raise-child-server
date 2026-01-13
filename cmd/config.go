package cmd

import (
	"fmt"
	"log"
	"os"
	"raise-child/constants/noti"

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
func loadEnv(logger *log.Logger) {
	if err := godotenv.Load(); err != nil {
		logger.Println(fmt.Sprintf(noti.ENV_LOAD_ERR_MSG, "") + err.Error())
	}
}

// Enable CORS
func corsConfig(server *gin.Engine) {
	server.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"}, // Allow all origins, or specify ["http://example.com"]
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Content-Type", "Authorization"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))
}

func setupSwagger(server *gin.Engine, port string) {
	// Configure swagger info
	docs.SwaggerInfo.Title = "RaiseChild Server API"
	docs.SwaggerInfo.Version = "1.0"
	docs.SwaggerInfo.Schemes = []string{"http", "https"}
	docs.SwaggerInfo.Host = "localhost:" + port
	//docs.SwaggerInfo.Host = os.Getenv(env.SWAGGER_HOST)

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
