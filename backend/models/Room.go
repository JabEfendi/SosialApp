package models

import "time"

type Room struct {
    ID               uint      `gorm:"primaryKey"`
    UserID           uint
    Name             string
    Description      string
    StartTime        time.Time
    EndTime          time.Time
    Duration         time.Duration
    MapLocID         *uint      `json:"map_loc_id"`
    MapLoc           MapLoc    `gorm:"foreignKey:MapLocID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL" json:"map_loc"`
    Capacity         int
    FeePerPerson     float64
    IsPaid           bool `json:"is_paid"`
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
