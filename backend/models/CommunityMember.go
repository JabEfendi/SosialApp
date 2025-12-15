package models

import "time"

type CommunityMember struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	CommunityID uint      `json:"community_id"`
	UserID      uint      `json:"user_id"`
	Role        string    `json:"role"`
	Status      string    `json:"status"`
	JoinedAt    time.Time `json:"joined_at"`
}
