package permissions

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"net/http"
	"opa-test/models"
	"strconv"
)

func GetAllPerm(db *gorm.DB) func(ctx *gin.Context) {
	return func(ctx *gin.Context) {
		var result []models.Permission

		if err := db.Find(&result).Error; err != nil {
			ctx.JSONP(400, gin.H{
				"error": err.Error(),
			})
			return
		}

		ctx.JSONP(200, gin.H{
			"data": result,
		})
	}
}

func GetPermById(db *gorm.DB) func(ctx *gin.Context) {
	return func(ctx *gin.Context) {
		id, err := strconv.Atoi(ctx.Param("id_permission"))
		if err != nil {
			ctx.JSON(400, gin.H{
				"error": err.Error(),
			})
			return
		}

		var perm models.Permission
		if err := db.Where("id = ?", id).First(&perm).Error; err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error": "id permission not found",
			})
			return
		}

		ctx.JSON(http.StatusOK, gin.H{
			"data": perm,
		})
	}
}

func GetPermByRequest(db *gorm.DB) func(ctx *gin.Context) {
	return func(ctx *gin.Context) {
		id, err := strconv.Atoi(ctx.Query("url"))
		if err != nil {
			ctx.JSON(400, gin.H{
				"error": err.Error(),
			})
			return
		}

		var perm models.Permission
		if err := db.Where("id = ?", id).First(&perm).Error; err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error": "id permission not found",
			})
			return
		}

		ctx.JSON(http.StatusOK, gin.H{
			"data": perm,
		})
	}
}
