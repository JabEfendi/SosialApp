package models

import (
    "time"
    "gorm.io/datatypes"
)

type RoomLog struct {
    ID              uint    `gorm:"primaryKey"`
    RoomID          uint
    UserID          uint
    KycStatus       string
    Summary_json    datatypes.JSON `json:"summary_json"`
    CreatedAt       time.Time
}
