package middlewares

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {

		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"message": "Authorization header required",
			})
			c.Abort()
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"message": "Invalid authorization format",
			})
			c.Abort()
			return
		}
		token := parts[1]

		// ini dummy, nanti bisa diganti JWT validation
		if token == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"message": "Invalid token",
			})
			c.Abort()
			return
		}

		c.Set("user_id", uint(1))
		c.Next()
	}
}