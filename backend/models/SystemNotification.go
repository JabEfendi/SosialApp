package models

import "time"

type SystemNotification struct {
	ID        uint   `gorm:"primaryKey"`
	Type      string `gorm:"type:varchar(50);not null"`
	Title     string `gorm:"type:varchar(150);not null"`
	Message   string `gorm:"type:text"`
	IsActive  bool   `gorm:"default:true"`
	StartAt   *time.Time
	EndAt     *time.Time
	CreatedAt time.Time
}
