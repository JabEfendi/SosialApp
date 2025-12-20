package models

import "time"

type RoomAccess struct {
	ID        uint      `gorm:"primaryKey"`
	RoomID    uint
	UserID    uint
	OrderID   uint
	GrantedAt time.Time `gorm:"autoCreateTime"`
	Room      Room      `gorm:"foreignKey:RoomID"`
	User      User      `gorm:"foreignKey:UserID"`
	Order     RoomJoinOrder `gorm:"foreignKey:OrderID"`
}