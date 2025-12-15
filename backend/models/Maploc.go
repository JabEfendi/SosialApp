package models

import "time"

type MapLoc struct {
    ID           uint      `gorm:"primaryKey" json:"id"`
    Name         string    `json:"name"`
    Address      string    `json:"address"`
    Description  string    `json:"description"`
    Latitude     float64   `json:"latitude"`
    Longitude    float64   `json:"longitude"`
    LocationType string    `json:"location_type"`
    IsActive     bool      `json:"is_active"`
    CreatedAt    time.Time `json:"created_at"`
    UpdatedAt    time.Time `json:"updated_at"`
}