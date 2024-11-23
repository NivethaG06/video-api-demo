package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"interview-project/pkg/models"
	"log"
	"net/http"
)

func ValidateRequestMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		var videoMetadata *models.VideoMetadata

		if err := c.ShouldBindJSON(&videoMetadata); err != nil {
			log.Printf("Request body error: %v", err)
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			c.Abort()
			return
		}
		validate := validator.New()
		err := validate.Struct(videoMetadata)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"validation_error": err.Error()})
			c.Abort()
			return
		}
		c.Next()
	}
}
