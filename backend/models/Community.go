package models

import "time"

type Community struct {
	ID               uint      `gorm:"primaryKey" json:"id"`
	CreatorID        uint      `json:"creator_id"`
	Name             string    `json:"name"`
	Description      string    `json:"description"`
	CountryRegion    string    `json:"country_region"`
	Interests        string    `json:"interests"`
	Type             string    `json:"type"`
	AutoApprove      bool      `json:"auto_approve"`
	ChatNotifications bool     `json:"chat_notifications"`
	InviteCode       string    `json:"invite_code"`
	CoverImage       string    `json:"cover_image"`
	Avatar     			 string 	 `gorm:"type:text"`
	CreatedAt        time.Time `json:"created_at"`
	UpdatedAt        time.Time `json:"updated_at"`
}