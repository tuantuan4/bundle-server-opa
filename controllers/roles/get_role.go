package roles

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"net/http"
	"opa-test/models"
	"strconv"
)

func GetAllRole(db *gorm.DB) func(ctx *gin.Context) {
	return func(ctx *gin.Context) {
		var result []models.Role
		if err := db.Find(&result).Error; err != nil {
			ctx.JSON(400, gin.H{"error": err.Error()})
			return
		}
		ctx.JSON(200, gin.H{"data": result})
	}
}

func GetRoleById(db *gorm.DB) func(ctx *gin.Context) {
	return func(ctx *gin.Context) {
		//id, err := strconv.Atoi(ctx.Param("id_role"))
		id, err := strconv.Atoi(ctx.Query("id_role"))

		if err != nil {
			ctx.JSON(400, gin.H{
				"error": err,
			})
			return
		}

		var role models.Role
		if err := db.Where("id = ?", id).First(&role).Error; err != nil {
			ctx.JSON(http.StatusNotFound, gin.H{
				"error": "id not found",
			})
			return
		}
		ctx.JSON(200, gin.H{
			"data": role,
		})
	}
}
