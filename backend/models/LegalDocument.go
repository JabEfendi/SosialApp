package models

import "time"

type LegalDocument struct {
	ID        uint   `gorm:"primaryKey"`
	Type      string `gorm:"type:varchar(50);not null"`
	Version   string `gorm:"type:varchar(50);not null"`
	Content   string `gorm:"type:text;not null"`
	IsActive  bool   `gorm:"default:true"`
	CreatedAt time.Time
}