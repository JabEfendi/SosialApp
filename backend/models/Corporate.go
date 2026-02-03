package models

import "time"

type Corporate struct {
	ID uint `gorm:"primaryKey"`

	// corporate identity
	Name          string `gorm:"type:varchar(150);not null"`
	Logo          string `gorm:"type:varchar(255)"`
	Reffcorporate string `gorm:"type:varchar(50);uniqueIndex;not null"`
	Description   string `gorm:"type:text"`
	Status        string `gorm:"type:varchar(50);default:'active'"`

	// corporate contact (NON LOGIN)
	EmailCorporate string `gorm:"type:varchar(100)"`
	Phone          string `gorm:"type:varchar(20)"`
	Address        string `gorm:"type:text"`
	City           string `gorm:"type:varchar(100)"`
	State          string `gorm:"type:varchar(100)"`
	Country        string `gorm:"type:varchar(100)"`
	ZipCode        string `gorm:"type:varchar(20)"`

	// PIC (LOGIN ACCOUNT)
	Email    string `gorm:"type:varchar(100);uniqueIndex;not null"` // PIC email
	Password string `gorm:"type:varchar(255);not null"`
	NamePIC  string `gorm:"type:varchar(150)"`
	PhonePIC string `gorm:"type:varchar(20)"`
	AgePIC   int

	TwoFAEnabled bool   `gorm:"column:two_fa_enabled;default:false"`
	TwoFASecret  string `gorm:"column:two_fa_secret;type:varchar(100)"`

	// audit
	CreatedBy *uint
	UpdatedBy *uint
	CreatedAt time.Time
	UpdatedAt time.Time
}