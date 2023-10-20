package config

import (
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"log"
	"opa-test/models"
)

func ConnectionDatabase() (db *gorm.DB) {
	db, err := gorm.Open(mysql.Open("tuan:Tuanstudent_123@tcp(127.0.0.1:3306)/data_opa?charset=utf8mb4&parseTime=True&loc=Local"), &gorm.Config{})
	if err != nil {
		log.Fatalln(err)
	}
	db.AutoMigrate(models.User{})

	return db
}
