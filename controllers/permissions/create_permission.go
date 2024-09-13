package permissions

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"opa-test/models"
)

func CreatePermission(db *gorm.DB) func(ctx *gin.Context) {
	return func(ctx *gin.Context) {
		var perm models.Permission

		err := ctx.BindJSON(&perm)
		if err != nil {
			ctx.JSON(400, gin.H{
				"error": "error",
			})
			return
		}

		result := db.Create(&perm)
		if result.Error != nil {
			ctx.JSON(400, gin.H{
				"error": "Failed to create permission",
			})
		}
		ctx.JSON(200, gin.H{
			"data": perm,
		})

	}
}
