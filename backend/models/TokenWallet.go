package models

import "time"

type TokenWallet struct {
	ID        uint64    `gorm:"primaryKey;autoIncrement"`
	UserID    uint64    `gorm:"not null;index"`
	Balance   int64     `gorm:"not null;default:0"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (TokenWallet) TableName() string {
	return "token_wallets"
}
