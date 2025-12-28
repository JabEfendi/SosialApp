package models

import "time"

type EmailCampaign struct {
	ID          uint      `gorm:"primaryKey"`
	Title       string    `gorm:"type:varchar(150);not null"`
	Subject     string    `gorm:"type:varchar(150);not null"`
	Content     string    `gorm:"type:text;not null"`
	Status      string    `gorm:"type:varchar(50);default:'draft'"`
	ScheduledAt *time.Time
	CreatedBy   *uint
	CreatedAt   time.Time
}