package models

import "time"

type Notification struct {
    ID        uint      `json:"id" gorm:"primaryKey"`
    UserID    uint      `json:"user_id"`
    Title     string    `json:"title"`
    Message   string    `json:"message"`
    Type      string    `json:"type"`
    IsRead    bool      `json:"is_read"`
    CreatedAt time.Time `json:"created_at"`
}
