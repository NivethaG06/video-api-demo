package routes

import (
	"github.com/gin-gonic/gin"
	"interview-project/internal/app"
	"interview-project/pkg/middleware"
)

func SetupRoutes(router *gin.Engine) {
	router.POST("/video", middleware.ValidateRequestMiddleware(), app.CreateVideo)
	router.GET("/videos", app.GetVideos)
}
