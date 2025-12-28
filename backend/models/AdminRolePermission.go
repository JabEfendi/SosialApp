package models

type AdminRolePermission struct {
	ID           uint `gorm:"primaryKey"`
	RoleID       uint `gorm:"not null"`
	PermissionID uint `gorm:"not null"`
	Role       AdminRole       `gorm:"foreignKey:RoleID"`
	Permission AdminPermission `gorm:"foreignKey:PermissionID"`
}
