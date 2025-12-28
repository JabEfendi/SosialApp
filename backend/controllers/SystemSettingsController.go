package controllers

import (
	"net/http"
	"time"
	"backend/db"
	"backend/models"

	"github.com/gin-gonic/gin"
)

func CreatePermission(c *gin.Context) {
	var req struct {
		Code        string `json:"code" binding:"required"`
		Description string `json:"description"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	data := models.AdminPermission{
		Code:        req.Code,
		Description: req.Description,
	}

	if err := db.DB.Create(&data).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, data)
}


func UpdateRolePermissions(c *gin.Context) {
	roleID := c.Param("role_id")

	var role models.AdminRole
	if err := db.DB.First(&role, roleID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "role not found"})
		return
	}

	var req struct {
		PermissionIDs []uint `json:"permission_ids" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var permissions []models.AdminPermission
	db.DB.Where("id IN ?", req.PermissionIDs).Find(&permissions)

	if err := db.DB.Model(&role).
		Association("Permissions").
		Replace(&permissions); err != nil {

		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "role permissions updated"})
}


func GetNotificationSettings(c *gin.Context) {
	var data []models.NotificationSetting
	db.DB.Find(&data)
	c.JSON(http.StatusOK, data)
}

func UpdateNotificationSetting(c *gin.Context) {
	id := c.Param("id")

	var req struct {
		IsEnabled bool `json:"is_enabled"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	db.DB.Model(&models.NotificationSetting{}).
		Where("id = ?", id).
		Update("is_enabled", req.IsEnabled)

	c.JSON(http.StatusOK, gin.H{"message": "notification updated"})
}


func CreateLegalDocument(c *gin.Context) {
	var req struct {
		Type    string `json:"type" binding:"required"`
		Version string `json:"version" binding:"required"`
		Content string `json:"content" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	doc := models.LegalDocument{
		Type:    req.Type,
		Version: req.Version,
		Content: req.Content,
	}

	db.DB.Create(&doc)
	c.JSON(http.StatusCreated, doc)
}


func CreateEmailCampaign(c *gin.Context) {
	adminID := c.MustGet("admin_id").(uint)

	var req struct {
		Title       string     `json:"title" binding:"required"`
		Subject     string     `json:"subject" binding:"required"`
		Content     string     `json:"content" binding:"required"`
		ScheduledAt *time.Time `json:"scheduled_at"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	data := models.EmailCampaign{
		Title:       req.Title,
		Subject:     req.Subject,
		Content:     req.Content,
		ScheduledAt: req.ScheduledAt,
		CreatedBy:   &adminID,
	}

	db.DB.Create(&data)
	c.JSON(http.StatusCreated, data)
}


func CreateMaintenance(c *gin.Context) {
	adminID := c.MustGet("admin_id").(uint)

	var req struct {
		Title   string    `json:"title" binding:"required"`
		Message string    `json:"message"`
		StartAt time.Time `json:"start_at" binding:"required"`
		EndAt   time.Time `json:"end_at" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	data := models.MaintenanceSchedule{
		Title:     req.Title,
		Message:   req.Message,
		StartAt:   req.StartAt,
		EndAt:     req.EndAt,
		CreatedBy: &adminID,
	}

	db.DB.Create(&data)
	c.JSON(http.StatusCreated, data)
}