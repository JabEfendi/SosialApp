package models

import "time"

type Corporate struct {
	ID     uint   `gorm:"primaryKey"`
	Name   string `gorm:"type:varchar(150);not null"`
	Logo   string `gorm:"type:varchar(255)"`
	Reffcorporate   string `gorm:"type:varchar(50);uniqueIndex;not null"`
	Email	 string `gorm:"type:varchar(100);uniqueIndex;not null"`
	Password string `gorm:"type:varchar(255);not null"`
	Description string `gorm:"type:text"`
	Status string `gorm:"type:varchar(50);default:'active'"`
	CreatedBy *uint
	UpdatedBy *uint
	CreatedAt time.Time
	UpdatedAt time.Time
}