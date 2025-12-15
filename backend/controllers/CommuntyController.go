package controllers

import (
	"backend/db"
	"backend/models"
	"math/rand"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func generateInviteCode() string {
	const charset = "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	rand.Seed(time.Now().UnixNano())

	b := make([]byte, 6)
	for i := range b {
		b[i] = charset[rand.Intn(len(charset))]
	}
	return string(b)
}

type CreateCommunityInput struct {
	UserID            uint   `json:"user_id" binding:"required"`
	Name             string `json:"name" binding:"required"`
	Description      string `json:"description"`
	CountryRegion    string `json:"country_region"`
	Interests        string `json:"interests"`
	Type             string `json:"type" binding:"required"`
	AutoApprove      bool   `json:"auto_approve"`
	ChatNotifications bool  `json:"chat_notifications"`
	CoverImage       string `json:"cover_image"`
}

func CreateCommunity(c *gin.Context) {

	var input CreateCommunityInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// userID := c.GetUint("user_id")

	invite := generateInviteCode()

	community := models.Community{
		CreatorID:        input.UserID,
		Name:             input.Name,
		Description:      input.Description,
		CountryRegion:    input.CountryRegion,
		Interests:        input.Interests,
		Type:             input.Type,
		AutoApprove:      input.AutoApprove,
		ChatNotifications: input.ChatNotifications,
		InviteCode:       invite,
		CoverImage:       input.CoverImage,
	}

	if err := db.DB.Create(&community).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to create a community"})
		return
	}

	db.DB.Create(&models.CommunityMember{
		CommunityID: community.ID,
		// UserID:      userID,
		UserID:      input.UserID,
		Role:        "admin",
		Status:      "approved",
		JoinedAt:    time.Now(),
	})

	c.JSON(http.StatusOK, gin.H{
		"message": "Community created",
		"data":    community,
	})
}

type JoinInput struct {
	UserID      uint   `json:"user_id" binding:"required"`
	InviteCode  string `json:"invite_code"`
	CommunityID uint   `json:"community_id"`
}

func JoinCommunity(c *gin.Context) {
	var input JoinInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	// userID := c.GetUint("user_id")
	userID := input.UserID

	var community models.Community

	if input.InviteCode != "" {
		if err := db.DB.Where("invite_code = ?", input.InviteCode).First(&community).Error; err != nil {
			c.JSON(400, gin.H{"error": "Invalid invitation code"})
			return
		}
	} else {
		if err := db.DB.First(&community, input.CommunityID).Error; err != nil {
			c.JSON(400, gin.H{"error": "Community not found"})
			return
		}
	}

	var existing models.CommunityMember
	err := db.DB.Where("community_id = ? AND user_id = ?", community.ID, userID).First(&existing).Error
	if err == nil {
		c.JSON(400, gin.H{"error": "Already a member"})
		return
	}

	status := "pending"
	if community.Type == "public" || community.AutoApprove {
		status = "approved"
	}

	db.DB.Create(&models.CommunityMember{
		CommunityID: community.ID,
		UserID:      userID,
		Role:        "member",
		Status:      status,
		JoinedAt:    time.Now(),
	})

	c.JSON(200, gin.H{
		"message": "Successfully joined",
		"status":  status,
	})
}

func GetCommunityDetail(c *gin.Context) {
	id := c.Param("id")

	var community models.Community
	if err := db.DB.First(&community, id).Error; err != nil {
		c.JSON(404, gin.H{"error": "Community not found"})
		return
	}

	c.JSON(200, gin.H{
		"data": community,
	})
}

func GetCommunityMembers(c *gin.Context) {
	id := c.Param("id")

	var members []models.CommunityMember
	db.DB.Where("community_id = ? AND status = 'approved'", id).Find(&members)

	c.JSON(200, gin.H{
		"members": members,
	})
}
