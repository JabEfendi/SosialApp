package models

import "time"

type RoomParticipant struct {
    ID        uint      `gorm:"primaryKey" json:"id"`
    RoomID    uint      `json:"room_id"`
    UserID    uint      `json:"user_id"`
    CreatedAt time.Time `json:"created_at"`
}
