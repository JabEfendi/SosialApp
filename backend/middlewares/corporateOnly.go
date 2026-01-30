package middlewares

import (
	"net/http"
	"strings"

	"backend/helpers"
	"github.com/gin-gonic/gin"
)

func CorporateAuth() gin.HandlerFunc {
	return func(c *gin.Context) {

		auth := c.GetHeader("Authorization")
		if auth == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "missing token",
			})
			return
		}

		tokenString := strings.TrimPrefix(auth, "Bearer ")

		claims, err := helpers.ValidateCorporateToken(tokenString)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "invalid token",
			})
			return
		}

		c.Set("corporate_id", claims.CorporateID)

		c.Next()
	}
}