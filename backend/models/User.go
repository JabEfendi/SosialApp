package models

import "time"

type User struct {
	ID uint `gorm:"primaryKey"`
	Name     string `gorm:"size:255;not null"`
	Username string `gorm:"size:50;unique"`
	Email    string `gorm:"size:255;unique;not null"`
	Password string `gorm:"size:255;not null"`
	Gender   string `gorm:"size:20;not null"`
	Birthdate *time.Time
	Phone     string `gorm:"size:20"`
	Bio       string `gorm:"size:255"`
	Country   string `gorm:"size:100"`
	Address   string `gorm:"type:text"`
	Provider   string `gorm:"size:50;default:local"`
	ProviderID string `gorm:"size:255"`
	Avatar     string `gorm:"type:text"`
	// CoinBalance int64 `gorm:"default:0"`
	ReferralCode string `gorm:"size:20;unique"`
	ReferredBy   *uint
	Referrer *User `gorm:"foreignKey:ReferredBy"`
	Referrals []Referral `gorm:"foreignKey:ReferrerID"`
	CreatedAt time.Time
	UpdatedAt time.Time
	CorporateID       *uint
	JoinedCorporateAt *time.Time
	Corporate *Corporate `gorm:"foreignKey:CorporateID"`
	IsReported  bool   `gorm:"default:false"`
	ReportCount int    `gorm:"default:0"`
	Status      string `gorm:"size:20;default:'active'"`
}