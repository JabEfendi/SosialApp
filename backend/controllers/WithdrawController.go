package controllers

import (
	"net/http"
	"backend/db"
	"backend/models"
	"backend/helpers"
	"time"

	"github.com/gin-gonic/gin"
)

func GetAccountBanks(c *gin.Context) {
	var banks []models.AccountBank

	if err := db.DB.
		Where("is_active = ?", true).
		Order("type ASC, name ASC").
		Find(&banks).Error; err != nil {

		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Failed to get banks",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Bank list",
		"data":    banks,
	})
}

func CreateWithdrawAccount(c *gin.Context) {
	userID := c.GetUint("user_id")

	var input struct {
		AccountBankID uint   `json:"account_bank_id" binding:"required"`
		AccountNumber string `json:"account_number" binding:"required"`
		AccountName   string `json:"account_name" binding:"required"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Invalid input",
		})
		return
	}

	var existing models.UserWithdrawAccount
	if err := db.DB.
		Where("user_id = ?", userID).
		First(&existing).Error; err == nil {

		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Withdraw account already exists",
		})
		return
	}

	account := models.UserWithdrawAccount{
		UserID:        userID,
		AccountBankID: input.AccountBankID,
		AccountNumber: input.AccountNumber,
		AccountName:   input.AccountName,
		Status:        "active",
	}

	if err := db.DB.Create(&account).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Failed to create withdraw account",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Withdraw account created",
		"data":    account,
	})
}

func RequestUpdateWithdrawAccount(c *gin.Context) {
	userID := c.GetUint("user_id")

	var input struct {
		AccountBankID uint   `json:"account_bank_id" binding:"required"`
		AccountNumber string `json:"account_number" binding:"required"`
		AccountName   string `json:"account_name" binding:"required"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Invalid input",
		})
		return
	}

	var account models.UserWithdrawAccount
	if err := db.DB.
		Where("user_id = ?", userID).
		First(&account).Error; err != nil {

		c.JSON(http.StatusNotFound, gin.H{
			"message": "Withdraw account not found",
		})
		return
	}

	var pending models.WithdrawAccountRequest
	if err := db.DB.
		Where("withdraw_account_id = ? AND status = ?", account.ID, "pending").
		First(&pending).Error; err == nil {

		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Update request still pending",
		})
		return
	}

	autoApproveAt := time.Now().Add(24 * time.Hour)

	req := models.WithdrawAccountRequest{
		UserID:            userID,
		WithdrawAccountID: account.ID,
		AccountBankID:     input.AccountBankID,
		AccountNumber:     input.AccountNumber,
		AccountName:       input.AccountName,
		Status:            "pending",
		AutoApproveAt:     autoApproveAt,
	}

	if err := db.DB.Create(&req).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Failed to request update",
		})
		return
	}

	db.DB.Model(&account).Update("status", "pending_update")

	c.JSON(http.StatusOK, gin.H{
		"message": "Update request submitted, auto approve in 24 hours",
		"data":    req,
	})
}

func GetMyWithdrawAccount(c *gin.Context) {
	userID := c.GetUint("user_id")

	var account models.UserWithdrawAccount
	if err := db.DB.
		Preload("AccountBank").
		Where("user_id = ?", userID).
		First(&account).Error; err != nil {

		c.JSON(http.StatusNotFound, gin.H{
			"message": "Withdraw account not found",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": account,
	})
}

func WithdrawCommission(c *gin.Context) {
	userID := c.GetUint("user_id")
	email := c.GetString("email")

	var input struct {
		Amount float64 `json:"amount" binding:"required,gt=0"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid input"})
		return
	}

	var account models.UserWithdrawAccount
	if err := db.DB.Where("user_id = ?", userID).First(&account).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Withdraw account not found"})
		return
	}

	if account.Status != "active" {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Withdraw account is not active"})
		return
	}

	_ = helpers.SendWithdrawSuccessEmail(email, input.Amount)

	c.JSON(http.StatusOK, gin.H{
		"message": "Withdraw success",
	})
}