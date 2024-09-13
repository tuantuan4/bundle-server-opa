package users

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"opa-test/models"
)

func GetAllUser(db *gorm.DB) func(ctx *gin.Context) {
	return func(ctx *gin.Context) {
		var result []models.User
		if err := db.Find(&result).Error; err != nil {
			ctx.JSONP(400, gin.H{
				"error": err.Error(),
			})
			return
		}
		ctx.JSON(200, gin.H{
			"bundle": result,
		})
	}
}
