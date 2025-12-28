package models

import "time"

type Corporate struct {
	ID     uint   `gorm:"primaryKey"`
	Name   string `gorm:"type:varchar(150);not null"`
	Code   string `gorm:"type:varchar(50);uniqueIndex;not null"`
	Status string `gorm:"type:varchar(50);default:'active'"`
	CreatedBy *uint
	UpdatedBy *uint
	CreatedAt time.Time
	UpdatedAt time.Time
	Users []User `gorm:"foreignKey:CorporateID"`
}