package errors

import (
	"github.com/ellofae/payment-system-kafka/pkg/logger"
	"github.com/gin-gonic/gin"
)

type errorResponse struct {
	Message string `json:"message"`
}

func NewErrorResponse(c *gin.Context, statusCode int, message string) {
	c.AbortWithStatusJSON(statusCode, errorResponse{message})
}

func NewHTTPResponse(c *gin.Context, statusCode int, message string, err error) {
	log := logger.GetLogger()

	log.Error(message, "error", err.Error())
	c.JSON(statusCode, gin.H{
		"message": message,
	})
}

func NewHTTPErrorResposne(c *gin.Context, statusCode int, message string, err error) {
	log := logger.GetLogger()

	log.Error(message, "error", err.Error())
	c.JSON(statusCode, gin.H{
		"message":       message,
		"error_message": err.Error(),
	})
}
