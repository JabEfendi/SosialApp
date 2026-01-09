package models

import "time"

type UserWithdrawAccount struct {
	ID            uint      `gorm:"primaryKey"`
	UserID        uint      `gorm:"uniqueIndex;not null"`
	AccountBankID uint      `gorm:"not null"`
	AccountBank   AccountBank
	AccountNumber string    `gorm:"size:50;not null"`
	AccountName   string    `gorm:"size:100;not null"`
	Status        string    `gorm:"type:enum('active','pending_update');default:'active'"`
	CreatedAt     time.Time
	UpdatedAt     time.Time
}