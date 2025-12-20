package models

import "time"

type TopUpTransaction struct {
	ID              uint64     `gorm:"primaryKey;autoIncrement"`
	InvoiceNumber   string     `gorm:"type:uuid;uniqueIndex;not null"`
	UserID          uint64     `gorm:"not null;index"`
	TokenPackageID  *uint64
	TokenAmount     int64      `gorm:"not null"`
	Price           int64      `gorm:"not null"`
	PaymentMethod   string     `gorm:"type:varchar(50);not null"`
	PaymentReference *string
	Status          string     `gorm:"type:varchar(20);default:PENDING"`
	PaidAt          *time.Time
	CreatedAt       time.Time
	UpdatedAt       time.Time
	TokenPackage *TokenPackage `gorm:"foreignKey:TokenPackageID"`
}

func (TopUpTransaction) TableName() string {
	return "topup_transactions"
}
