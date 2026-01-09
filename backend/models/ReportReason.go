package models

type ReportReason struct {
	ID        uint      `gorm:"primaryKey"`
	Name      string    `gorm:"size:100;not null"`
	IsActive  bool      `gorm:"default:true"`
}
