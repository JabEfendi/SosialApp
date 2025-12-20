package controllers

import (
	"backend/db"
	"backend/models"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type JoinPaidRoomInput struct {
	RoomID uint `json:"room_id" binding:"required"`
	UserID uint `json:"user_id" binding:"required"`
}

const PlatformFeePercent = 0.10

func JoinPaidRoom(c *gin.Context) {
	var input JoinPaidRoomInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := db.DB.Transaction(func(tx *gorm.DB) error {
		var room models.Room
		if err := tx.First(&room, input.RoomID).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Room not found"})
			return err
		}

		if !room.IsPaid {
			c.JSON(http.StatusBadRequest, gin.H{"error": "This room is free. Use normal join endpoint"})
			return gorm.ErrInvalidTransaction
		}

		if room.UserID == input.UserID {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Owner cannot join own room"})
			return gorm.ErrInvalidTransaction
		}

		var count int64
		tx.Model(&models.RoomParticipant{}).Where("room_id = ?", room.ID).Count(&count)
		if room.Capacity > 0 && int(count) >= room.Capacity {
			c.JSON(http.StatusBadRequest, gin.H{"error": "The room is full"})
			return gorm.ErrInvalidTransaction
		}

		var existing models.RoomParticipant
		if err := tx.Where("room_id = ? AND user_id = ?", room.ID, input.UserID).First(&existing).Error; err == nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "User already joined this room"})
			return gorm.ErrInvalidTransaction
		}

		var wallet models.TokenWallet
		if err := tx.Where("user_id = ?", input.UserID).First(&wallet).Error; err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Wallet not found"})
			return err
		}

		fee := room.FeePerPerson
		platformFee := int64(float64(fee) * PlatformFeePercent)
		totalCharge := int64(fee) + platformFee

		if wallet.Balance < totalCharge {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Insufficient token balance"})
			return gorm.ErrInvalidTransaction
		}

		balanceBefore := wallet.Balance
		wallet.Balance -= totalCharge
		if err := tx.Save(&wallet).Error; err != nil {
			return err
		}

		var platformWallet models.TokenWallet
		if err := tx.Where("user_id = ?", 0).First(&platformWallet).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Platform wallet not found"})
			return err
		}
		platformBalanceBefore := platformWallet.Balance
		platformWallet.Balance += platformFee
		if err := tx.Save(&platformWallet).Error; err != nil {
			return err
		}

		now := time.Now()
		roomIDRef := uint64(room.ID)

		userLedger := models.TokenLedger{
			UserID:        uint64(input.UserID),
			SourceType:    "paid_room",
			SourceID:      &roomIDRef,
			Amount:        -totalCharge,
			BalanceBefore: balanceBefore,
			BalanceAfter:  wallet.Balance,
			Description:   "Join paid room",
			CreatedAt:     now,
		}
		if err := tx.Create(&userLedger).Error; err != nil {
			return err
		}

		platformLedger := models.TokenLedger{
			UserID:        0,
			SourceType:    "platform_fee",
			SourceID:      &roomIDRef,
			Amount:        platformFee,
			BalanceBefore: platformBalanceBefore,
			BalanceAfter:  platformWallet.Balance,
			Description:   "Platform fee from paid room",
			CreatedAt:     now,
		}
		if err := tx.Create(&platformLedger).Error; err != nil {
			return err
		}

		participant := models.RoomParticipant{
			RoomID:    room.ID,
			UserID:    input.UserID,
			CreatedAt: now,
		}
		if err := tx.Create(&participant).Error; err != nil {
			return err
		}

		c.JSON(http.StatusOK, gin.H{
			"message":       "Successfully joined paid room",
			"participant":   participant,
			"total_charge":  totalCharge,
			"platform_fee":  platformFee,
			"user_balance":  wallet.Balance,
		})

		return nil
	})

	if err != nil {
		return
	}
}