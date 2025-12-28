package models

type AdminPermission struct {
	ID          uint   `gorm:"primaryKey"`
	Code        string `gorm:"type:varchar(100);unique;not null"`
	Description string `gorm:"type:varchar(255)"`
}
