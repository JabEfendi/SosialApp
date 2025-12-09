package models

import "time"

type ChatDirectMessage struct {
    ID         uint      `gorm:"primaryKey"`
    ThreadID   uint      `json:"thread_id"`
    SenderID   uint      `json:"sender_id"`
    ReceiverID uint      `json:"receiver_id"`
    Message    string    `json:"message"`
    Status     string    `json:"status"
    CreatedAt  time.Time `json:"created_at"`
}
