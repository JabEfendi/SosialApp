package middlewares

import (
	"net/http"
	// "log"

	"github.com/gin-gonic/gin"
)

func SuperAdminOnly() gin.HandlerFunc {
	return func(c *gin.Context) {
		role := c.GetString("admin_role")
		// log.Println("ROLE IN CONTEXT:", role)

		if role != "superadmin" {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
				"error": "superadmin access only",
			})
			return
		}

		c.Next()
	}
}