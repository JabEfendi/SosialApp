package models

import "time"

type CorporateUserLog struct {
	ID uint `gorm:"primaryKey"`
	UserID      uint `gorm:"not null"`
	CorporateID uint `gorm:"not null"`
	Action string `gorm:"type:varchar(50);not null"`
	CreatedAt time.Time
}