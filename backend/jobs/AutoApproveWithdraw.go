package jobs

import (
	"time"

	"backend/db"
	"backend/models"
)

func AutoApproveWithdrawAccount() {
	var requests []models.WithdrawAccountRequest

	db.DB.
		Where("status = ? AND auto_approve_at <= ?", "pending", time.Now()).
		Find(&requests)

	for _, req := range requests {

		db.DB.Model(&models.UserWithdrawAccount{}).
			Where("id = ?", req.WithdrawAccountID).
			Updates(map[string]interface{}{
				"account_bank_id": req.AccountBankID,
				"account_number":  req.AccountNumber,
				"account_name":    req.AccountName,
				"status":          "active",
			})

		db.DB.Model(&models.WithdrawAccountRequest{}).
			Where("id = ?", req.ID).
			Updates(map[string]interface{}{
				"status":      "approved",
				"approved_at": time.Now(),
			})
	}
}