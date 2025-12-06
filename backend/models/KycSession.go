package models

import (
	"gorm.io/datatypes"
)

type KycSession struct {
    ID        uint      `gorm:"primaryKey"`
    UserID    uint      `gorm:"not null"`
    DataJSON  datatypes.JSON `json:"data_json"`
    Used      bool      `gorm:"default:false"`
}
