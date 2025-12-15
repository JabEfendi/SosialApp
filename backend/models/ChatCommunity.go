package models

import "time"

type CommunityChatMessage struct {
    ID          uint      `json:"id" gorm:"primaryKey"`
    CommunityID uint      `json:"community_id"`
    UserID      uint      `json:"user_id"`
    Message     string    `json:"message"`
    CreatedAt   time.Time `json:"created_at"`
}
