package models

import "time"

type UserLegalAcceptance struct {
	ID              uint      `gorm:"primaryKey"`
	UserID          uint      `gorm:"not null"`
	LegalDocumentID uint      `gorm:"not null"`
	AcceptedAt      time.Time
}