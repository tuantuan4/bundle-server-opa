package roles

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"opa-test/models"
)

func CreateRole(db *gorm.DB) func(ctx *gin.Context) {
	return func(ctx *gin.Context) {
		var role models.Role
		err := ctx.BindJSON(&role)
		if err != nil {
			ctx.JSON(400, gin.H{"error": err.Error()})
			return
		}

		result := db.Create(&role)
		if result.Error != nil {
			ctx.JSON(400, gin.H{"error": "Failed to create role"})
		}
		ctx.JSON(201, gin.H{"data": role})
	}
}
