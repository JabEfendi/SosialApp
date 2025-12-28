package helpers

import (
	"backend/db"
	"backend/models"
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
