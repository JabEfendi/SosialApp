package controllers

import (
	"net/http"
	"time"

	"backend/db"
	"backend/models"
	"backend/helpers"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func CreateTopUp(c *gin.Context) {
	var input struct {
		PackageID     uint64 `json:"package_id" binding:"required"`
		PaymentMethod string `json:"payment_method" binding:"required"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}

	userID := c.GetUint("user_id") // dari middleware auth

	// Ambil paket token
	var tokenPackage models.TokenPackage
	if err := db.DB.
		Where("id = ? AND is_active = true", input.PackageID).
		First(&tokenPackage).Error; err != nil {

		c.JSON(http.StatusNotFound, gin.H{
			"message": "Token package not found",
		})
		return
	}

	// Buat transaksi top up
	topup := models.TopUpTransaction{
		InvoiceNumber: uuid.New().String(),
		UserID:        uint64(userID),
		TokenPackageID: &tokenPackage.ID,
		TokenAmount:   tokenPackage.TokenAmount,
		Price:         tokenPackage.Price,
		PaymentMethod: input.PaymentMethod,
		Status:        helpers.TopUpStatusPending,
	}

	if err := db.DB.Create(&topup).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Failed to create top up",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Top up created",
		"data":    topup,
	})
}

func PaymentCallback(c *gin.Context) {
	var input struct {
		InvoiceNumber    string `json:"invoice_number" binding:"required"`
		PaymentReference string `json:"payment_reference"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	tx := db.DB.Begin()

	var topup models.TopUpTransaction
	if err := tx.
		Where("invoice_number = ? AND status = ?", input.InvoiceNumber, helpers.TopUpStatusPending).
		First(&topup).Error; err != nil {

		tx.Rollback()
		c.JSON(http.StatusNotFound, gin.H{
			"message": "Top up not found or already processed",
		})
		return
	}

	// Ambil / buat wallet
	var wallet models.TokenWallet
	err := tx.Where("user_id = ?", topup.UserID).First(&wallet).Error
	if err != nil {
		wallet = models.TokenWallet{
			UserID:  topup.UserID,
			Balance: 0,
		}
		if err := tx.Create(&wallet).Error; err != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed create wallet"})
			return
		}
	}

	balanceBefore := wallet.Balance
	balanceAfter := balanceBefore + topup.TokenAmount

	// Update wallet
	if err := tx.Model(&wallet).
		Update("balance", balanceAfter).Error; err != nil {

		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed update wallet"})
		return
	}

	// Insert ledger
	ledger := models.TokenLedger{
		UserID:        topup.UserID,
		SourceType:    helpers.TokenSourceTopUp,
		SourceID:      &topup.ID,
		Amount:        topup.TokenAmount,
		BalanceBefore: balanceBefore,
		BalanceAfter:  balanceAfter,
		Description:   "Top up token",
	}

	if err := tx.Create(&ledger).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed insert ledger"})
		return
	}

	// Update top up status
	now := time.Now()
	if err := tx.Model(&topup).Updates(map[string]interface{}{
		"status":            helpers.TopUpStatusPaid,
		"payment_reference": input.PaymentReference,
		"paid_at":           &now,
	}).Error; err != nil {

		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed update top up"})
		return
	}

	tx.Commit()

	c.JSON(http.StatusOK, gin.H{
		"message": "Top up success",
	})
}

func GetTopUpHistory(c *gin.Context) {
	userID := c.GetUint("user_id")

	var history []models.TopUpTransaction
	if err := db.DB.
		Where("user_id = ?", userID).
		Order("created_at DESC").
		Find(&history).Error; err != nil {

		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Failed fetch history",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": history,
	})
}

