package models

import "time"

type OTPVerification struct {
    ID        uint      `gorm:"primaryKey"`
    Email     string
    OTP       string
    ExpiredAt time.Time
    CreatedAt time.Time
}
