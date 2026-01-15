package controllers

import (
	"backend/db"
	"backend/models"
	"math/rand"
	"net/http"
	"time"
	"fmt"
	"os"
	"strings"
	"path/filepath"

	"github.com/google/uuid"
	"github.com/spf13/cast"
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

func UploadCommunityAvatar(c *gin.Context) {
    communityID := c.PostForm("community_id")
    if communityID == "" {
        c.JSON(400, gin.H{"error": "community_id is required"})
        return
    }

    form, err := c.MultipartForm()
    if err != nil {
        c.JSON(400, gin.H{"error": "Failed to read photos"})
        return
    }

    files := form.File["photos"]
    if len(files) == 0 {
        c.JSON(400, gin.H{"error": "No photos uploaded"})
        return
    }

    uploadPath := "uploads/communities"
    baseURL := "https://testtestdomaingweh.com/"
    os.MkdirAll(uploadPath, os.ModePerm)

    uploadedURLs := []string{}
    var community models.Community

    if err := db.DB.First(&community, communityID).Error; err != nil {
        c.JSON(404, gin.H{"error": "Community not found"})
        return
    }

    for _, file := range files {

        if !strings.Contains(file.Header.Get("Content-Type"), "image") {
            continue
        }

        filename := fmt.Sprintf("%s-%d-%s%s",
            communityID,
            time.Now().Unix(),
            uuid.New().String(),
            filepath.Ext(file.Filename),
        )

        fullPath := uploadPath + "/" + filename
        photoURL := baseURL + "/" + fullPath

        if err := c.SaveUploadedFile(file, fullPath); err != nil {
            continue
        }

        db.DB.Create(&models.CommunityPhoto{
            CommunityID: cast.ToUint(communityID),
            Photo:  photoURL,
        })

        uploadedURLs = append(uploadedURLs, photoURL)

        if community.Avatar == "" {
						community.Avatar = photoURL
						db.DB.Model(&community).Update("avatar", photoURL)
				}
    }

    c.JSON(200, gin.H{
        "message": "Photos uploaded successfully",
        "photos":  uploadedURLs,
        "avatar":  community.Avatar,
    })
}

func GetPhotosCommunity(c *gin.Context) {
    communityID := c.Param("id")

    var photos []models.CommunityPhoto
    db.DB.Where("community_id = ?", communityID).Find(&photos)

    c.JSON(200, gin.H{
        "photos": photos,
    })
}

func SetProfilePhotoCommunity(c *gin.Context) {
    communityID := c.PostForm("community_id")
    photoID := c.PostForm("photo_id")

    if communityID == "" || photoID == "" {
        c.JSON(400, gin.H{"error": "community_id & photo_id required"})
        return
    }

    var photo models.CommunityPhoto
    if err := db.DB.Where("id = ? AND community_id = ?", photoID, communityID).First(&photo).Error; err != nil {
        c.JSON(404, gin.H{"error": "Photo not found"})
        return
    }

    db.DB.Model(&models.Community{}).
        Where("id = ?", communityID).
        Update("avatar", photo.Photo)

    c.JSON(200, gin.H{
        "message": "Avatar updated",
        "avatar":  photo.Photo,
    })
}

func DeletePhotoCommunity(c *gin.Context) {
    photoID := c.Param("id")

    var photo models.CommunityPhoto
    if err := db.DB.First(&photo, photoID).Error; err != nil {
        c.JSON(404, gin.H{"error": "Photo not found"})
        return
    }

    var community models.Community
    db.DB.First(&community, photo.CommunityID)

    parts := strings.Split(photo.Photo, "/")
    filename := parts[len(parts)-1]
    os.Remove("uploads/communities/" + filename)

    db.DB.Delete(&photo)

    if community.Avatar == photo.Photo {
        defaultPhoto := "https://testtestdomaingweh.com/default-avatar.png"
        db.DB.Model(&community).Update("avatar", defaultPhoto)
    }

    c.JSON(200, gin.H{
        "message": "Photo deleted",
    })
}

func GetAllCommunity(c *gin.Context) {
	var communities []models.Community

	countryRegion := c.Query("country_region")
	dbQuery := db.DB

	if countryRegion != "" {
		dbQuery = dbQuery.Where("country_region = ?", countryRegion)
	}

	interests := c.Query("interests")
	if interests != "" {
		dbQuery = dbQuery.Where("interests ILIKE ?", "%"+interests+"%")
	}

	if err := dbQuery.Find(&communities).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Failed to retrieve communities",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Community data retrieved",
		"total":   len(communities),
		"data":    communities,
	})
}