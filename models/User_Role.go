package models

import "gorm.io/gorm"

type UserRole struct {
	gorm.Model
	Username string `json:"username" gorm:"type:varchar(355);uniqueIndex:idx_user_role"`
	Email    string `json:"email" gorm:"type:varchar(355);uniqueIndex:idx_user_role"`
	RoleId   uint   `json:"role_id" gorm:"uniqueIndex:idx_user_role"`
}
