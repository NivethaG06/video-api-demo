package main

import (
	"github.com/gin-contrib/gzip"
	"github.com/gin-gonic/gin"
	"interview-project/internal/app"
	"interview-project/internal/cache"
	"interview-project/internal/configloader"
	"interview-project/internal/routes"
	"interview-project/pkg/error"
	"interview-project/pkg/logger"
	"log"
	"os"
)

func main() {
	Fileloader()
	go Cacheloader()

	router := gin.Default()
	router.NoRoute(app.NoRouteHandler())
	router.Use(error.ErrorHandler())
	router.Use(gzip.Gzip(gzip.DefaultCompression))
	routes.SetupRoutes(router)

	router.Run(os.Getenv("LOCALSERVER_HOST") + os.Getenv("LOCALSERVER_PORT"))
}

func Fileloader() {
	if len(os.Args) < 2 {
		log.Fatal("Please provide the YAML file path as an argument.")
		//FatalErrorHandler()
	}
	loggerPath := os.Args[1]

	logger.Logloader(loggerPath)

	configPath := os.Args[2]

	config, err := configloader.LoadConfig(configPath)
	if err != nil {
		log.Fatalf("Error loading configloader: %v", err)
	}

	configloader.SetEnvVars(config, "")
}

func Cacheloader() {
	cache.Init()
}
