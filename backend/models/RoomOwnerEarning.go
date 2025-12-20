package models

import "time"

type RoomOwnerEarning struct {
	ID        uint      `gorm:"primaryKey"`
	RoomID    uint
	OwnerID   uint
	OrderID   uint
	Amount    float64   `gorm:"type:decimal(10,2);not null"`
	CreatedAt time.Time
	Room      Room      `gorm:"foreignKey:RoomID"`
	Owner     User      `gorm:"foreignKey:OwnerID"`
	Order     RoomJoinOrder `gorm:"foreignKey:OrderID"`
}