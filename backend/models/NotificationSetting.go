package models

import "time"

type NotificationSetting struct {
	ID        uint   `gorm:"primaryKey"`
	Channel   string `gorm:"type:varchar(50);not null"`
	IsEnabled bool   `gorm:"default:true"`
	UpdatedBy *uint
	UpdatedAt time.Time
}
