package error

import (
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"os"
)

type APIError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func SendSuccessBodyResponse(c *gin.Context, apierror *APIError, data interface{}) {
	c.IndentedJSON(apierror.Code, gin.H{
		"status": apierror.Message,
		"data":   data,
	})
}

func SendSuccessResponse(c *gin.Context, apierror *APIError) {
	c.IndentedJSON(apierror.Code, gin.H{
		"status":  apierror.Code,
		"message": apierror.Message,
	})
}

func SendErrorResponse(c *gin.Context, apierror *APIError) {
	c.IndentedJSON(apierror.Code, gin.H{
		"status":  "error",
		"message": apierror.Message,
	})
}

func ErrorHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()
		if len(c.Errors) > 0 {
			err := c.Errors.Last()
			if err.Type == http.StatusInternalServerError {
				SendErrorResponse(c, &APIError{Code: 500, Message: "Internal Server Error"})
			}
			if err.Type == http.StatusBadRequest {
				SendErrorResponse(c, &APIError{Code: 400, Message: "Bad Request"})
			}
		}
	}
}

func FatalErrorHandler(err error, message string) {
	if err != nil {
		log.Printf("%s: %v", message, err)
		// Exit the application if it's a critical error
		os.Exit(1)
	}
}
