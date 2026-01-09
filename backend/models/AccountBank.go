package models

import "time"

type AccountBank struct {
	ID        uint      `gorm:"primaryKey"`
	Code      string    `gorm:"size:50;unique;not null"`
	Name      string    `gorm:"size:100;not null"`
	Type      string    `gorm:"type:enum('bank','ewallet');not null"`
	Logo      *string
	IsActive  bool      `gorm:"default:true"`
	CreatedAt time.Time
	UpdatedAt time.Time
}