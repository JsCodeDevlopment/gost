package guards

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func AuthGuard() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.GetHeader("Authorization")
		if token == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"statusCode": http.StatusUnauthorized,
				"error":      "Unauthorized",
				"message":    "Missing Authorization header",
			})
			return
		}

		c.Next()
	}
}
