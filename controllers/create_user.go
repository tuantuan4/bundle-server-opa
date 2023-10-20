package controllers

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"opa-test/models"
)

func CreateUser(db *gorm.DB) func(ctx *gin.Context) {
	return func(ctx *gin.Context) {
		var user models.User
		err := ctx.BindJSON(&user)
		if err != nil {
			ctx.JSON(400, gin.H{
				"error": "error",
			})
			return
		}
		result := db.Create(&user)
		if result.Error != nil {
			ctx.JSON(400, gin.H{
				"error": "Failed to create user",
			})
		}
		ctx.JSON(200, gin.H{
			"bundle": "OK",
			"user":   user,
		})
	}
}
