package models

import "time"

type UserFCMToken struct {
    ID        uint      `json:"id" gorm:"primaryKey"`
    UserID    uint      `json:"user_id"`
    FCMToken  string    `json:"fcm_token"`
    Device    string    `json:"device"`
    CreatedAt time.Time `json:"created_at"`
}
