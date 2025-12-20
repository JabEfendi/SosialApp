package models

import "time"

type RoomCommission struct {
	ID          uint      `gorm:"primaryKey"`
	RoomID      uint
	OrderID     uint
	PlatformFee float64   `gorm:"type:decimal(10,2);not null"`
	CreatedAt   time.Time
	Room        Room      `gorm:"foreignKey:RoomID"`
	Order       RoomJoinOrder `gorm:"foreignKey:OrderID"`
}