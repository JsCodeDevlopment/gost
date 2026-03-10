package app

import (
	"log"
	"os"

	"gost/src/common/filters"
	"gost/src/common/interceptors"
	"gost/src/config"
	"gost/src/modules/users"

	"github.com/gin-gonic/gin"
)

func SetupApp() *gin.Engine {
	config.LoadEnv()

	config.ConnectDatabase()
	config.ConnectRedis()

	router := gin.Default()

	router.Use(config.SetupCors())

	router.Use(interceptors.LoggerInterceptor())
	router.Use(filters.ErrorHandler())

	api := router.Group("/api/v1")

	users.InitModule(api)

	return router
}

func Bootstrap() {
	router := SetupApp()

	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}

	log.Printf("Application is starting on port %s...", port)
	if err := router.Run(":" + port); err != nil {
		log.Fatalf("Could not start server: %v", err)
	}
}
