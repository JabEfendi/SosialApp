package models

import "time"

type UserReport struct {
	ID             uint      `gorm:"primaryKey"`
	ReporterID     uint      `gorm:"not null"`
	ReportedUserID uint      `gorm:"not null"`
	ReasonID       uint      `gorm:"not null"`
	Description    string    `gorm:"type:text"`
	Status         string    `gorm:"size:20;default:'pending'"`
	AdminNote      string    `gorm:"type:text"`
	CreatedAt      time.Time
	UpdatedAt      time.Time
	Reporter       User         `gorm:"foreignKey:ReporterID"`
	ReportedUser   User         `gorm:"foreignKey:ReportedUserID"`
	Reason         ReportReason `gorm:"foreignKey:ReasonID"`
}