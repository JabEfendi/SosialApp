package models

import "time"

type TokenLedger struct {
	ID            uint64    `gorm:"primaryKey;autoIncrement"`
	UserID        uint64    `gorm:"not null;index"`
	SourceType    string    `gorm:"type:varchar(30);not null"`
	SourceID      *uint64
	Amount        int64     `gorm:"not null"`
	BalanceBefore int64     `gorm:"not null"`
	BalanceAfter  int64     `gorm:"not null"`
	Description   string    `gorm:"type:text"`
	CreatedAt     time.Time
}

func (TokenLedger) TableName() string {
	return "token_ledger"
}
