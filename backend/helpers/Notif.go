package helpers

import (
	"fmt"

	"backend/db"
	"backend/models"
	"backend/services"
)

type AdminNotificationPayload struct {
	AdminID uint
	Type    string
	Title   string
	Message string
}

func CreateAdminNotification(p AdminNotificationPayload) error {
	notif := models.AdminNotification{
		AdminID: p.AdminID,
		Type:    p.Type,
		Title:   p.Title,
		Message: p.Message,
	}

	return db.DB.Create(&notif).Error
}

func SendWithdrawSuccessEmail(email string, amount float64) error {
	subject := "Withdraw Berhasil"
	body := fmt.Sprintf(
		"Halo,\n\nWithdraw sebesar Rp %.0f berhasil diproses.\n\nTerima kasih.\n\n— SosialApp",
		amount,
	)

	return services.SendEmail(email, subject, body)
}