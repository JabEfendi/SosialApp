package middlewares

import (
	"os"
	"strings"
	
	"github.com/golang-jwt/jwt/v5"
	"github.com/gin-gonic/gin"
)

func TwoFAMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		auth := c.GetHeader("Authorization")
		if auth == "" {
			c.AbortWithStatusJSON(401, gin.H{"error": "missing token"})
			return
		}

		tokenString := strings.TrimPrefix(auth, "Bearer ")

		token, err := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
			return []byte(os.Getenv("JWT_SECRET")), nil
		})

		if err != nil || !token.Valid {
			c.AbortWithStatusJSON(401, gin.H{"error": "invalid token"})
			return
		}

		claims := token.Claims.(jwt.MapClaims)

		if claims["type"] != "2fa" {
			c.AbortWithStatusJSON(403, gin.H{"error": "invalid token type"})
			return
		}

		corporateID := uint(claims["corporate_id"].(float64))
		c.Set("corporate_id", corporateID)

		c.Next()
	}
}