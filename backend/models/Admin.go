package models

import "time"

type Admin struct {
	ID           uint      `gorm:"primaryKey"`
	RoleID       uint      `gorm:"not null"`
	Role         AdminRole `gorm:"foreignKey:RoleID"`
	Name         string    `gorm:"type:varchar(100);not null"`
	Email        string    `gorm:"type:varchar(100);unique;not null"`
	Password     string    `gorm:"type:varchar(255);not null"`
	Status     	 string    `gorm:"default:'pending'"`
	LastLoginAt *time.Time
	ApprovedAt *time.Time
	ApprovedBy *uint
	CreatedAt time.Time
	UpdatedAt time.Time
}