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
		if err := db.DB.
			Select("id", "role_id").
			First(&admin, adminID).
			Error; err != nil {

			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "admin not found",
			})
			return
		}

		var rolePermissions []models.AdminRolePermission
		if err := db.DB.
			Where("role_id = ?", admin.RoleID).
			Preload("Permission").
			Find(&rolePermissions).
			Error; err != nil {

			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
				"error": "failed to load permissions",
			})
			return
		}

		for _, rp := range rolePermissions {
			if rp.Permission.Code == permissionCode {
				c.Next()
				return
			}
		}

		c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
			"error":   "forbidden",
			"message": "you do not have permission to perform this action",
		})
	}
}