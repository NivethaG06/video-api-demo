package app

import (
	"github.com/gin-gonic/gin"
	"interview-project/internal/services"
	"interview-project/pkg/models"
	"log"
	"strconv"
)
import "interview-project/pkg/error"

func NoRouteHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		error.SendErrorResponse(c, &error.APIError{Code: 404, Message: "Page Not Found"})
	}
}

func CreateVideo(c *gin.Context) {
	var video models.VideoMetadata
	if err := c.ShouldBindJSON(&video); err != nil {
		error.SendErrorResponse(c, &error.APIError{Code: 400, Message: "Bad Request"})
		return
	}

	if err := services.CreateVideo(&video); err != nil {
		error.SendErrorResponse(c, &error.APIError{Code: 2, Message: "Failed to save video"})
		return
	}

	error.SendSuccessResponse(c, &error.APIError{Code: 200, Message: "Data inserted"})
}

func GetVideos(c *gin.Context) {
	offsetStr := c.DefaultQuery("offset", "0")
	limitStr := c.DefaultQuery("limit", "10")

	offset, err := strconv.Atoi(offsetStr)
	if err != nil {
		log.Println("Error converting offset to int:", err)
		c.JSON(400, gin.H{"error": "Invalid offset value"})
		return
	}

	limit, err := strconv.Atoi(limitStr)
	if err != nil {
		log.Println("Error converting limit to int:", err)
		c.JSON(400, gin.H{"error": "Invalid limit value"})
		return
	}

	videos, err := services.GetVideos(c, &models.LimitOffset{Offset: offset, Limit: limit})
	if err != nil {
		error.SendErrorResponse(c, &error.APIError{Code: 2, Message: "Failed to fetch videos"})
		return
	}

	error.SendSuccessBodyResponse(c, &error.APIError{Code: 200, Message: "Fetched Data succesfully"}, videos)
}
