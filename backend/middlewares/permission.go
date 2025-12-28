package middlewares

import (
	"net/http"
	"backend/db"
	"backend/models"

	"github.com/gin-gonic/gin"
)

func RequirePermission(permissionCode string) gin.HandlerFunc {
	return func(c *gin.Context) {

		adminIDValue, exists := c.Get("admin_id")
		if !exists {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "unauthorized",
			})
			return
		}

		adminID := adminIDValue.(uint)

		var admin models.Admin
		err := db.DB.
			Preload("Role.Permissions").
			First(&admin, adminID).
			Error

		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "admin not found",
			})
			return
		}

		for _, perm := range admin.Role.Permissions {
			if perm.Code == permissionCode {
				c.Next()
				return
			}
		}

		c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
			"error": "forbidden",
			"message": "you do not have permission to perform this action",
		})
	}
}

