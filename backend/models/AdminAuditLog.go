package models

import (
	"time"

	"gorm.io/datatypes"
)

type AdminAuditLog struct {
	ID         uint      `gorm:"primaryKey"`
	AdminID    uint      `gorm:"not null"`
	Admin      Admin     `gorm:"foreignKey:AdminID"`
	Action     string    `gorm:"type:varchar(50);not null"`
	TargetType string    `gorm:"type:varchar(50);not null"`
	TargetID   *uint
	BeforeData datatypes.JSON
	AfterData  datatypes.JSON
	IPAddress  string    `gorm:"type:varchar(45)"`
	UserAgent  string    `gorm:"type:text"`
	CreatedAt time.Time
}