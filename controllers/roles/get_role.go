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
		var result []models.Role
		// Apply limit and offset for pagination
		if err := db.Limit(pageSizeInt).Offset(offset).Find(&result).Error; err != nil {
			ctx.JSON(400, gin.H{"error": err.Error()})
			return
		}
		// Query for the total count of roles
		var total int64
		if err := db.Model(&models.Role{}).Count(&total).Error; err != nil {
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
