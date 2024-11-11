package users

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"opa-test/models"
	"strconv"
)

func GetAllUser(db *gorm.DB) func(ctx *gin.Context) {
	return func(ctx *gin.Context) {
		pageSize := ctx.DefaultQuery("pageSize", "10")
		pageNum := ctx.DefaultQuery("pageNum", "1")
		// Convert parameters to integers
		pageSizeInt, err := strconv.Atoi(pageSize)
		if err != nil || pageSizeInt <= 0 {
			ctx.JSON(400, gin.H{"error": "Invalid pageSize"})
			return
		}

		pageNumInt, err := strconv.Atoi(pageNum)
		if err != nil || pageNumInt <= 0 {
			ctx.JSON(400, gin.H{"error": "Invalid pageNum"})
			return
		}
		offset := (pageNumInt - 1) * pageSizeInt
		var result []models.Auth
		if err := db.Limit(pageSizeInt).Offset(offset).Find(&result).Error; err != nil {
			ctx.JSONP(400, gin.H{
				"error": err.Error(),
			})
			return
		}
		var total int64
		if err := db.Model(&models.Auth{}).Count(&total).Error; err != nil {
			ctx.JSON(400, gin.H{"error": err.Error()})
			return
		}
		ctx.JSON(200, gin.H{
			"result": gin.H{
				"data":  result,
				"total": total,
			},
			"status": true,
		})
	}
}
