package models

import "time"

type RoomRevenueHold struct {
	ID          uint64    `gorm:"primaryKey"`
	RoomID      uint      `gorm:"index;not null"`
	OwnerUserID uint      `gorm:"index;not null"`
	CorporateID *uint     `gorm:"index"`
	Amount      int64     `gorm:"not null"`
	Status      string    `gorm:"type:varchar(20);default:'pending'"`
	ReleaseAt time.Time
	CreatedAt time.Time
	UpdatedAt time.Time
}