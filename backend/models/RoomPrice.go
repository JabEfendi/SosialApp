package models

import "time"

type RoomPrice struct {
	ID                 uint      `gorm:"primaryKey"`
	RoomID             uint
	BasePrice          float64   `gorm:"type:decimal(10,2);not null"`
	PlatformFeePercent float64   `gorm:"type:decimal(5,2);default:10.00"`
	IsActive           bool      `gorm:"default:true"`
	CreatedAt          time.Time
}