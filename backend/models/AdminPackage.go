package models

import "time"

type AdminPackage struct {
	ID         uint   `gorm:"primaryKey"`
	Name       string `gorm:"type:varchar(100);not null"`
	CoinAmount int     `gorm:"not null"`
	Price      float64 `gorm:"type:numeric(15,2);not null"`
	BonusCoin  int     `gorm:"default:0"`
	Status     string  `gorm:"default:'active'"`
	CreatedBy  *uint
	UpdatedBy  *uint
	CreatedAt time.Time
	UpdatedAt time.Time
}