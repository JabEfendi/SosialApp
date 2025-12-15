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
        DataJSON map[string]interface{} `json:"data_json"`
    }

    if err := c.ShouldBindJSON(&input); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    jsonBytes, _ := json.Marshal(input.DataJSON)

    var kyc models.KycSession
    err := db.DB.Where("user_id = ?", input.UserID).First(&kyc).Error

    if err == gorm.ErrRecordNotFound {
        newKyc := models.KycSession{
            UserID:   input.UserID,
            DataJSON: datatypes.JSON(jsonBytes),
            Used:     false,
            Status:   "pending",
        }

        db.DB.Create(&newKyc)
        c.JSON(http.StatusOK, gin.H{"message": "KYC created and pending", "kyc": newKyc})
        return
    }

    kyc.DataJSON = datatypes.JSON(jsonBytes)
    kyc.Status = "pending"
    db.DB.Save(&kyc)

    c.JSON(http.StatusOK, gin.H{"message": "KYC updated and pending", "kyc": kyc})
}

func ApproveKyc(c *gin.Context) {
    var input struct {
        KycID uint `json:"kyc_id"`
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

    kyc.Status = "active"
    kyc.Used = false
    db.DB.Save(&kyc)

    c.JSON(http.StatusOK, gin.H{"message": "KYC approved", "kyc": kyc})
}

func RejectKyc(c *gin.Context) {
    var input struct {
        KycID uint `json:"kyc_id"`
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

    db.DB.Delete(&kyc)
    c.JSON(http.StatusOK, gin.H{"message": "KYC rejected and deleted"})
}
