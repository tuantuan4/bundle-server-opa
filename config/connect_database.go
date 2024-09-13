package config

import (
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"log"
	"opa-test/models"
)

func ConnectionDatabase() (db *gorm.DB) {
	db, err := gorm.Open(mysql.Open("jobber:api@tcp(localhost:3306)/jobber_auth?charset=utf8mb4&parseTime=True&loc=Local"), &gorm.Config{})
	if err != nil {
		log.Fatalln(err)
	} else {
		fmt.Println("connect DB done")
	}
	db.AutoMigrate(models.Role{}, models.Permission{}, models.RolePermission{})

	return db
}
