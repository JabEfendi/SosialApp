package models

import "time"

type MaintenanceSchedule struct {
	ID        uint      `gorm:"primaryKey"`
	Title     string    `gorm:"type:varchar(150);not null"`
	Message   string    `gorm:"type:text"`
	StartAt   time.Time `gorm:"not null"`
	EndAt     time.Time `gorm:"not null"`
	IsActive  bool      `gorm:"default:true"`
	CreatedBy *uint
	CreatedAt time.Time
}