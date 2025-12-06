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

    var user models.User
    if err := db.DB.First(&user, input.UserID).Error; err != nil {
        c.JSON(http.StatusBadRequest, gin.H{
            "error": "KYC ditolak, Mohon re-login",
        })
        return
    }

    jsonBytes, _ := json.Marshal(input.DataJSON)

    var kyc models.KycSession
    err := db.DB.Where("user_id = ? AND used = false", input.UserID).First(&kyc).Error

    if err == gorm.ErrRecordNotFound {
        newKyc := models.KycSession{
            UserID:   input.UserID,
            DataJSON: datatypes.JSON(jsonBytes),
            Used:     false,
        }

        db.DB.Create(&newKyc)
        c.JSON(http.StatusOK, gin.H{
            "message": "KYC created",
            "kyc":     newKyc,
        })
        return
    }

    kyc.DataJSON = datatypes.JSON(jsonBytes)
    db.DB.Save(&kyc)

    c.JSON(http.StatusOK, gin.H{
        "message": "KYC updated",
        "kyc":     kyc,
    })
}



