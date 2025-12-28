package controllers

import (
	"net/http"
	"backend/db"
	"backend/helpers"
	"backend/models"

	"github.com/gin-gonic/gin"
)

func GetCorporates(c *gin.Context) {
	var corporates []models.Corporate

	if err := db.DB.Order("id desc").Find(&corporates).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": corporates,
	})
}

func CreateCorporate(c *gin.Context) {
	adminID := c.MustGet("admin_id").(uint)

	var input struct {
		Name string `json:"name" binding:"required"`
		Code string `json:"code" binding:"required"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	corporate := models.Corporate{
		Name:      input.Name,
		Code:      input.Code,
		Status:    "active",
		CreatedBy: &adminID,
	}

	if err := db.DB.Create(&corporate).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	// AUDIT LOG
	_ = helpers.CreateAdminAuditLog(helpers.AuditPayload{
		AdminID:    adminID,
		Action:     "CREATE",
		TargetType: "corporate",
		TargetID:   &corporate.ID,
		After:      corporate,
		Context:    c,
	})

	c.JSON(http.StatusCreated, gin.H{
		"message": "corporate created",
		"data":    corporate,
	})
}

func UpdateCorporate(c *gin.Context) {
	adminID := c.MustGet("admin_id").(uint)
	id := c.Param("id")

	var corporate models.Corporate
	if err := db.DB.First(&corporate, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "corporate not found",
		})
		return
	}

	before := corporate

	var input struct {
		Name   string `json:"name"`
		Code   string `json:"code"`
		Status string `json:"status"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	if input.Name != "" {
		corporate.Name = input.Name
	}
	if input.Code != "" {
		corporate.Code = input.Code
	}
	if input.Status != "" {
		corporate.Status = input.Status
	}

	corporate.UpdatedBy = &adminID

	if err := db.DB.Save(&corporate).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	// AUDIT LOG
	_ = helpers.CreateAdminAuditLog(helpers.AuditPayload{
		AdminID:    adminID,
		Action:     "UPDATE",
		TargetType: "corporate",
		TargetID:   &corporate.ID,
		Before:     before,
		After:      corporate,
		Context:    c,
	})

	c.JSON(http.StatusOK, gin.H{
		"message": "corporate updated",
	})
}

func ChangeCorporateStatus(c *gin.Context) {
	adminID := c.MustGet("admin_id").(uint)
	id := c.Param("id")

	var corporate models.Corporate
	if err := db.DB.First(&corporate, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "corporate not found",
		})
		return
	}

	before := corporate

	var input struct {
		Status string `json:"status" binding:"required"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	corporate.Status = input.Status
	corporate.UpdatedBy = &adminID

	if err := db.DB.Save(&corporate).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	// AUDIT LOG
	_ = helpers.CreateAdminAuditLog(helpers.AuditPayload{
		AdminID:    adminID,
		Action:     "UPDATE_STATUS",
		TargetType: "corporate",
		TargetID:   &corporate.ID,
		Before:     before,
		After:      corporate,
		Context:    c,
	})

	c.JSON(http.StatusOK, gin.H{
		"message": "status updated",
	})
}

