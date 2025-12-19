package models

import "time"

type ReferralReward struct {
	ID uint `gorm:"primaryKey"`
	UserID     uint `gorm:"not null"`
	ReferralID uint `gorm:"not null"`
	Amount     int    `gorm:"not null"`
	RewardType string `gorm:"size:50"`
	Note       string `gorm:"type:text"`
	User     User     `gorm:"foreignKey:UserID"`
	Referral Referral `gorm:"foreignKey:ReferralID"`
	CreatedAt time.Time
}