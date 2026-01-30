package models

import "time"

type Agreement struct {
	ID uint `gorm:"primaryKey"`
	UserID      uint `gorm:"not null"`
	CorporateID uint `gorm:"not null"`
	AgreementNumber string `gorm:"type:varchar(100);uniqueIndex;not null"`
	StartDate time.Time
	EndDate   *time.Time
	RevenueUserPercent      float64 `gorm:"type:numeric(5,2);not null"`
	RevenueCorporatePercent float64 `gorm:"type:numeric(5,2);not null"`
	RevenueType   string `gorm:"type:varchar(20);not null"`
	PaymentPeriod string `gorm:"type:varchar(20);not null"`
	ScopeDescription string `gorm:"type:text"`
	Status string `gorm:"type:varchar(30);default:'requested'"`
	CorporateApprovedAt *time.Time
	UserApprovedAt      *time.Time
	TerminatedAt *time.Time
	TerminatedBy *string `gorm:"type:varchar(20)"`
	TerminationReason *string `gorm:"type:text"`
	CreatedAt time.Time
	UpdatedAt time.Time
}