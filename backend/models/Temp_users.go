package models

import "time"

type TempUser struct {
	ID uint `gorm:"primaryKey"`
	Email    string `gorm:"size:255;unique;not null"`
	Name     string `gorm:"size:255"`
	Username string `gorm:"size:255"`
	Password string `gorm:"size:255"`
	Gender   string `gorm:"size:50"`
	Birthdate *time.Time
	Phone     string `gorm:"size:50"`
	Bio       string `gorm:"type:text"`
	Country   string `gorm:"size:100"`
	Address   string `gorm:"type:text"`
	ReferralCode string `gorm:"size:20"`
	CreatedAt time.Time
}