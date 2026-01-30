package controllers

import (
	"backend/db"
	"backend/models"
	"net/http"
	"errors"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type JoinPaidRoomInput struct {
	RoomID uint `json:"room_id" binding:"required"`
	UserID uint `json:"user_id" binding:"required"`
}

const PlatformFeePercent = 0.10

// func JoinPaidRoom(c *gin.Context) {
// 	var input JoinPaidRoomInput
// 	if err := c.ShouldBindJSON(&input); err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
// 		return
// 	}

// 	err := db.DB.Transaction(func(tx *gorm.DB) error {
// 		var room models.Room
// 		if err := tx.First(&room, input.RoomID).Error; err != nil {
// 			c.JSON(http.StatusNotFound, gin.H{"error": "Room not found"})
// 			return err
// 		}

// 		if !room.IsPaid {
// 			c.JSON(http.StatusBadRequest, gin.H{"error": "This room is free. Use normal join endpoint"})
// 			return gorm.ErrInvalidTransaction
// 		}

// 		if room.UserID == input.UserID {
// 			c.JSON(http.StatusBadRequest, gin.H{"error": "Owner cannot join own room"})
// 			return gorm.ErrInvalidTransaction
// 		}

// 		var count int64
// 		tx.Model(&models.RoomParticipant{}).Where("room_id = ?", room.ID).Count(&count)
// 		if room.Capacity > 0 && int(count) >= room.Capacity {
// 			c.JSON(http.StatusBadRequest, gin.H{"error": "The room is full"})
// 			return gorm.ErrInvalidTransaction
// 		}

// 		var existing models.RoomParticipant
// 		if err := tx.Where("room_id = ? AND user_id = ?", room.ID, input.UserID).First(&existing).Error; err == nil {
// 			c.JSON(http.StatusBadRequest, gin.H{"error": "User already joined this room"})
// 			return gorm.ErrInvalidTransaction
// 		}

// 		var wallet models.TokenWallet
// 		if err := tx.Where("user_id = ?", input.UserID).First(&wallet).Error; err != nil {
// 			c.JSON(http.StatusBadRequest, gin.H{"error": "Wallet not found"})
// 			return err
// 		}

// 		fee := room.FeePerPerson
// 		platformFee := int64(float64(fee) * PlatformFeePercent)
// 		totalCharge := int64(fee) + platformFee

// 		if wallet.Balance < totalCharge {
// 			c.JSON(http.StatusBadRequest, gin.H{"error": "Insufficient token balance"})
// 			return gorm.ErrInvalidTransaction
// 		}

// 		balanceBefore := wallet.Balance
// 		wallet.Balance -= totalCharge
// 		if err := tx.Save(&wallet).Error; err != nil {
// 			return err
// 		}

// 		var platformWallet models.TokenWallet
// 		if err := tx.Where("user_id = ?", 0).First(&platformWallet).Error; err != nil {
// 			c.JSON(http.StatusInternalServerError, gin.H{"error": "Platform wallet not found"})
// 			return err
// 		}
// 		platformBalanceBefore := platformWallet.Balance
// 		platformWallet.Balance += platformFee
// 		if err := tx.Save(&platformWallet).Error; err != nil {
// 			return err
// 		}

// 		now := time.Now()
// 		roomIDRef := uint64(room.ID)

// 		userLedger := models.TokenLedger{
// 			UserID:        uint64(input.UserID),
// 			SourceType:    "paid_room",
// 			SourceID:      &roomIDRef,
// 			Amount:        -totalCharge,
// 			BalanceBefore: balanceBefore,
// 			BalanceAfter:  wallet.Balance,
// 			Description:   "Join paid room",
// 			CreatedAt:     now,
// 		}
// 		if err := tx.Create(&userLedger).Error; err != nil {
// 			return err
// 		}

// 		platformLedger := models.TokenLedger{
// 			UserID:        0,
// 			SourceType:    "platform_fee",
// 			SourceID:      &roomIDRef,
// 			Amount:        platformFee,
// 			BalanceBefore: platformBalanceBefore,
// 			BalanceAfter:  platformWallet.Balance,
// 			Description:   "Platform fee from paid room",
// 			CreatedAt:     now,
// 		}
// 		if err := tx.Create(&platformLedger).Error; err != nil {
// 			return err
// 		}

// 		participant := models.RoomParticipant{
// 			RoomID:    room.ID,
// 			UserID:    input.UserID,
// 			CreatedAt: now,
// 		}
// 		if err := tx.Create(&participant).Error; err != nil {
// 			return err
// 		}

// 		c.JSON(http.StatusOK, gin.H{
// 			"message":       "Successfully joined paid room",
// 			"participant":   participant,
// 			"total_charge":  totalCharge,
// 			"platform_fee":  platformFee,
// 			"user_balance":  wallet.Balance,
// 		})

// 		return nil
// 	})

// 	if err != nil {
// 		return
// 	}
// }

func JoinPaidRoom(c *gin.Context) {
	var input JoinPaidRoomInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := db.DB.Transaction(func(tx *gorm.DB) error {

		var room models.Room
		if err := tx.First(&room, input.RoomID).Error; err != nil {
			return err
		}

		if !room.IsPaid {
			return errors.New("room is free")
		}

		if room.UserID == input.UserID {
			return errors.New("owner cannot join own room")
		}

		var wallet models.TokenWallet
		if err := tx.Where("user_id = ?", input.UserID).First(&wallet).Error; err != nil {
			return err
		}

		fee := int64(room.FeePerPerson)

		platformFee := int64(float64(fee) * PlatformFeePercent)
		totalCharge := fee + platformFee

		if wallet.Balance < totalCharge {
			return errors.New("insufficient balance")
		}

		// 🔻 Deduct user wallet
		before := wallet.Balance
		wallet.Balance -= totalCharge
		if err := tx.Save(&wallet).Error; err != nil {
			return err
		}

		// 💰 Platform wallet
		var platformWallet models.TokenWallet
		if err := tx.Where("user_id = ?", 0).First(&platformWallet).Error; err != nil {
			return err
		}

		platformWallet.Balance += platformFee
		tx.Save(&platformWallet)

		now := time.Now()

		// 🧾 Ledger user
		tx.Create(&models.TokenLedger{
			UserID:        uint64(input.UserID),
			SourceType:    "paid_room",
			Amount:        -totalCharge,
			BalanceBefore: before,
			BalanceAfter:  wallet.Balance,
			Description:   "Join paid room",
			CreatedAt:     now,
		})

		// 👥 Participant
		tx.Create(&models.RoomParticipant{
			RoomID:    room.ID,
			UserID:    input.UserID,
			CreatedAt: now,
		})

		// 🔒 ESCROW
		var owner models.User
		if err := tx.First(&owner, room.UserID).Error; err != nil {
			return err
		}

		releaseAt := now.AddDate(0, 0, 3)
		if room.EndTime != nil {
			releaseAt = room.EndTime.AddDate(0, 0, 3)
		}

		hold := models.RoomRevenueHold{
			RoomID:      room.ID,
			OwnerUserID: room.UserID,
			Amount:      room.FeePerPerson,
			ReleaseAt:   releaseAt,
		}

		if owner.CorporateID != nil {
			hold.CorporateID = owner.CorporateID
		}

		if err := tx.Create(&hold).Error; err != nil {
			return err
		}

		c.JSON(http.StatusOK, gin.H{
			"message": "Joined paid room successfully",
		})

		return nil
	})

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}
}

func ReleaseRoomRevenue() error {
	now := time.Now()

	var holds []models.RoomRevenueHold
	if err := db.DB.
		Where("status = ? AND release_at <= ?", "pending", now).
		Find(&holds).Error; err != nil {
		return err
	}

	for _, hold := range holds {
		err := db.DB.Transaction(func(tx *gorm.DB) error {

			targetID := hold.OwnerUserID
			if hold.CorporateID != nil {
				targetID = *hold.CorporateID
			}

			var wallet models.TokenWallet
			if err := tx.
				Where("user_id = ?", targetID).
				FirstOrCreate(&wallet, models.TokenWallet{
					UserID: uint64(targetID),
				}).Error; err != nil {
				return err
			}

			before := wallet.Balance
			wallet.Balance += hold.Amount
			tx.Save(&wallet)

			hold.Status = "released"
			tx.Save(&hold)

			tx.Create(&models.TokenLedger{
				UserID:        uint64(targetID),
				SourceType:    "room_revenue_release",
				Amount:        hold.Amount,
				BalanceBefore: before,
				BalanceAfter:  wallet.Balance,
				Description:   "Room revenue released",
				CreatedAt:     now,
			})

			return nil
		})

		if err != nil {
			return err
		}
	}

	return nil
}

func CorporateApprovePayout(c *gin.Context) {
	var req struct {
		RoomID uint `json:"room_id" binding:"required"` // room yang sudah selesai
		Amount int64 `json:"amount" binding:"required"`  // nominal revenue corporate mau bayar
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Ambil room
	var room models.Room
	if err := db.DB.First(&room, req.RoomID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "room not found"})
		return
	}

	// Pastikan room sudah selesai
	if room.EndTime == nil || time.Now().Before(*room.EndTime) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "room has not finished yet"})
		return
	}

	// Ambil host / owner
	var host models.User
	if err := db.DB.First(&host, room.UserID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "host user not found"})
		return
	}

	if host.CorporateID == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "host does not belong to a corporate"})
		return
	}

	// Ambil agreement host-corporate
	var agreement models.Agreement
	if err := db.DB.
		Where("user_id = ? AND corporate_id = ? AND status = ?", host.ID, *host.CorporateID, "active").
		First(&agreement).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "active agreement not found for host"})
		return
	}

	// Hitung share host sesuai agreement
	userShare := int64(float64(req.Amount) * agreement.RevenueUserPercent / 100)

	err := db.DB.Transaction(func(tx *gorm.DB) error {
		// Ambil corporate wallet
		var corpWallet models.TokenWallet
		if err := tx.Where("user_id = ?", *host.CorporateID).First(&corpWallet).Error; err != nil {
			return errors.New("corporate wallet not found")
		}

		if corpWallet.Balance < userShare {
			return errors.New("corporate balance insufficient")
		}

		// Kurangi corporate wallet
		beforeCorp := corpWallet.Balance
		corpWallet.Balance -= userShare
		if err := tx.Save(&corpWallet).Error; err != nil {
			return err
		}

		// Tambah host wallet
		var hostWallet models.TokenWallet
		if err := tx.Where("user_id = ?", host.ID).FirstOrCreate(&hostWallet, models.TokenWallet{
			UserID: uint64(host.ID),
		}).Error; err != nil {
			return err
		}

		beforeHost := hostWallet.Balance
		hostWallet.Balance += userShare
		if err := tx.Save(&hostWallet).Error; err != nil {
			return err
		}

		// Buat ledger
		now := time.Now()
		if err := tx.Create(&models.TokenLedger{
			UserID:        uint64(host.ID),
			SourceType:    "corporate_payout",
			Amount:        userShare,
			BalanceBefore: beforeHost,
			BalanceAfter:  hostWallet.Balance,
			Description:   "Revenue payout from corporate for room",
			CreatedAt:     now,
		}).Error; err != nil {
			return err
		}

		if err := tx.Create(&models.TokenLedger{
			UserID:        uint64(*host.CorporateID),
			SourceType:    "corporate_payout_deduction",
			Amount:        -userShare,
			BalanceBefore: beforeCorp,
			BalanceAfter:  corpWallet.Balance,
			Description:   "Corporate paid revenue to host",
			CreatedAt:     now,
		}).Error; err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":      "Payout to host completed",
		"host_user_id": host.ID,
		"amount":       userShare,
	})
}