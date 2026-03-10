package pipes

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func ValidateBody[T any](c *gin.Context) (*T, error) {
	var dto T
	if err := c.ShouldBindJSON(&dto); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"statusCode": http.StatusBadRequest,
			"error":      "Bad Request",
			"message":    err.Error(),
		})
		return nil, err
	}
	return &dto, nil
}
