package models

import "time"

type AgreementRevision struct {
	ID uint `gorm:"primaryKey"`
	AgreementID uint   `gorm:"not null;index"`
	RequestedBy string `gorm:"type:varchar(20);not null"`
	Note string `gorm:"type:text;not null"`
	CreatedAt time.Time
}