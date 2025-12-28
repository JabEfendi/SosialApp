package models

import "time"

type AdminNotification struct {
	ID        uint      `gorm:"primaryKey"`
	AdminID   uint      `gorm:"not null"`
	Admin     Admin     `gorm:"foreignKey:AdminID"`
	Type      string    `gorm:"type:varchar(50);not null"`
	Title     string    `gorm:"type:varchar(100);not null"`
	Message   string    `gorm:"type:text"`
	IsRead    bool      `gorm:"default:false"`
	CreatedAt time.Time
}