package models

import "time"

type RoomJoinOrder struct {
	ID           uint      `gorm:"primaryKey"`
	OrderCode    string    `gorm:"uniqueIndex;size:50"`
	RoomID       uint
	UserID       uint
	BasePrice    float64   `gorm:"type:decimal(10,2);not null"`
	PlatformFee  float64   `gorm:"type:decimal(10,2);not null"`
	TotalPrice   float64   `gorm:"type:decimal(10,2);not null"`
	Status       string    `gorm:"size:20;default:'pending'"`
	PaidAt       *time.Time
	CreatedAt    time.Time
	Room         Room      `gorm:"foreignKey:RoomID"`
	User         User      `gorm:"foreignKey:UserID"`
}