package models

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Name     string `json:"name"`
	EndPoint string `json:"endPoint"`
	Role     string `json:"role"`
}
