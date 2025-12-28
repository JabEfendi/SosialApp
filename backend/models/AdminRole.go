package models

import "time"

type AdminRole struct {
	ID          uint      `gorm:"primaryKey"`
	Name        string    `gorm:"type:varchar(50);unique;not null"`
	Description string    `gorm:"type:varchar(255)"`
	CreatedAt   time.Time
	UpdatedAt   time.Time

	Permissions []AdminPermission `gorm:"many2many:admin_role_permissions;"`
}