package utils

import "github.com/gin-gonic/gin"

func FormattedErrorGenerator(c *gin.Context, statusCode int, error string, message string) {
	c.AbortWithStatusJSON(statusCode, gin.H{
		"statusCode": statusCode,
		"error":      error,
		"message":    message,
	})
}
