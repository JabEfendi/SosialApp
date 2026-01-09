package models

import "time"

type WithdrawAccountRequest struct {
	ID                uint      `gorm:"primaryKey"`
	UserID            uint      `gorm:"not null"`
	WithdrawAccountID uint      `gorm:"not null"`
	AccountBankID     uint      `gorm:"not null"`
	AccountBank       AccountBank
	AccountNumber     string    `gorm:"size:50;not null"`
	AccountName       string    `gorm:"size:100;not null"`
	Status            string     `gorm:"type:enum('pending','approved','rejected');default:'pending'"`
	AutoApproveAt     time.Time
	ApprovedAt        *time.Time
	RejectedReason    *string
	CreatedAt         time.Time
	UpdatedAt         time.Time
}