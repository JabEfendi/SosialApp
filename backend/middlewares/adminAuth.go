package middlewares

import (
	"net/http"
	"strings"
	// "log"

	"backend/helpers"

	"github.com/gin-gonic/gin"
)

func AdminAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		auth := c.GetHeader("Authorization")
		if auth == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "missing token",
			})
			return
		}

		tokenString := strings.TrimPrefix(auth, "Bearer ")

		claims, err := helpers.ValidateAdminToken(tokenString)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "invalid token",
			})
			return
		}

		// log.Println("JWT ROLE:", claims.Role)
		
		c.Set("admin_id", claims.AdminID)
		c.Set("admin_role", claims.Role)

		// log.Println("ADMIN ID FROM JWT:", claims.AdminID)
		// log.Println("ROLE FROM JWT:", claims.Role)
		c.Next()
	}
}