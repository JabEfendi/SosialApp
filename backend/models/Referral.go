package models

import "time"

type Referral struct {
	ID uint `gorm:"primaryKey"`
	ReferrerID uint `gorm:"not null"`
	ReferredID uint `gorm:"not null;unique"`
	Status string `gorm:"size:20;default:pending"`
	RewardAmount int
	RewardGivenAt *time.Time
	Referrer User `gorm:"foreignKey:ReferrerID"`
	Referred User `gorm:"foreignKey:ReferredID"`
	Rewards []ReferralReward `gorm:"foreignKey:ReferralID"`
	CreatedAt time.Time
}