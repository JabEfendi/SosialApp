package models

import "time"

type Room struct {
    ID               uint      `gorm:"primaryKey"`
    UserID           uint
    Name             string
    Description      string
    StartTime        time.Time
    EndTime          time.Time
    Duration         string
    Location         string
    Capacity         int
    FeePerPerson     float64
    Gender           string
    AgeGroup         string
    IsRegular        bool
    AutoApprove      bool
    SendNotification bool
    ImageURL         string
    Status           string
    CreatedAt        time.Time
    UpdatedAt        time.Time
}
