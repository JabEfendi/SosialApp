package models

import "time"

type AdminPasswordReset struct {
	ID        uint      `gorm:"primaryKey"`
	AdminID   uint      `gorm:"not null"`
	Admin     Admin     `gorm:"foreignKey:AdminID"`
	Token     string    `gorm:"type:varchar(255);not null"`
	ExpiredAt time.Time
	CreatedAt time.Time
}