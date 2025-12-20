package models

import "time"

type TokenPackage struct {
	ID          uint64    `gorm:"primaryKey;autoIncrement"`
	Name        string    `gorm:"type:varchar(100);not null"`
	TokenAmount int64     `gorm:"not null"`
	Price       int64     `gorm:"not null"`
	IsActive    bool      `gorm:"default:true"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

func (TokenPackage) TableName() string {
	return "token_packages"
}
