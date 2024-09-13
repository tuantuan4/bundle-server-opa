package models

import (
	"gorm.io/gorm"
)

type Permission struct {
	gorm.Model
	Name            string           `json:"name"`
	Url             string           `json:"url"`
	Method          string           `json:"method"`
	RolePermissions []RolePermission `gorm:"many2many:role_permissions;"`
}
