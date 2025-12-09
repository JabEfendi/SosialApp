package models

import "time"

type ChatMessage struct {
	ID        int      `gorm:"primaryKey" json:"id"`
	RoomID    int      `json:"room_id"`
	UserID    int      `json:"user_id"`
	Message   string    `json:"message"`
	CreatedAt time.Time `json:"created_at"`
}
