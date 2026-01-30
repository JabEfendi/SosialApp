package controllers

import (
    "backend/db"
		"backend/models"
		"net/http"
		"encoding/json"

		"gorm.io/datatypes"
		"gorm.io/gorm"
		"github.com/gin-gonic/gin"
)

func SubmitOrUpdateKyc(c *gin.Context) {
    var input struct {
        UserID   uint                   `json:"user_id"`
        DataJSON map[string]interface{} `json:"data_json" binding:"required"`
    }

    if err := c.ShouldBindJSON(&input); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    jsonBytes, _ := json.Marshal(input.DataJSON)

    var kyc models.KycSession
    err := db.DB.Where("user_id = ?", input.UserID).First(&kyc).Error

    if err == gorm.ErrRecordNotFound {
        // CREATE KYC
        newKyc := models.KycSession{
            UserID:   input.UserID,
            DataJSON: datatypes.JSON(jsonBytes),
            Status:   "pending",
        }

        db.DB.Create(&newKyc)
        c.JSON(http.StatusOK, gin.H{
            "message": "KYC submitted and pending approval",
            "kyc":     newKyc,
        })
        return
    }

    // UPDATE KYC → reset ke pending
    kyc.DataJSON = datatypes.JSON(jsonBytes)
    kyc.Status = "pending"

    db.DB.Save(&kyc)

    c.JSON(http.StatusOK, gin.H{
        "message": "KYC updated and pending approval",
        "kyc":     kyc,
    })
}

func ApproveKyc(c *gin.Context) {
    var input struct {
        KycID uint `json:"kyc_id" binding:"required"`
    }

    if err := c.ShouldBindJSON(&input); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    var kyc models.KycSession
    if err := db.DB.First(&kyc, input.KycID).Error; err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "KYC not found"})
        return
    }

    kyc.Status = "approved"
    kyc.Used = true
    db.DB.Save(&kyc)

    c.JSON(http.StatusOK, gin.H{
        "message": "KYC approved",
        "kyc":     kyc,
    })
}

func RejectKyc(c *gin.Context) {
    var input struct {
        KycID uint `json:"kyc_id" binding:"required"`
    }

    if err := c.ShouldBindJSON(&input); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    var kyc models.KycSession
    if err := db.DB.First(&kyc, input.KycID).Error; err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "KYC not found"})
        return
    }

    kyc.Status = "rejected"
    kyc.Used = false
    db.DB.Save(&kyc)

    c.JSON(http.StatusOK, gin.H{
        "message": "KYC rejected",
        "kyc":     kyc,
    })
}
