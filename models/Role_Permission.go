package models

import (
	"gorm.io/gorm"
)

type RolePermission struct {
	gorm.Model
	RoleId       uint `json:"role_id" gorm:"uniqueIndex:idx_role_permission"`
	PermissionId uint `json:"permission_id" gorm:"uniqueIndex:idx_role_permission"`
}
