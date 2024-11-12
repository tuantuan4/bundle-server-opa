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
		pageSize := ctx.DefaultQuery("pageSize", "10")
		pageNum := ctx.DefaultQuery("pageNum", "1")
		url := ctx.Query("url")       // Retrieve the "url" parameter
		method := ctx.Query("method") // Retrieve the "method" parameter
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
		var result []models.Permission
		query := db.Model(&models.Permission{})
		if url != "" {
			query = query.Where("url LIKE ?", "%"+url+"%")
		}
		if method != "" {
			query = query.Where("method = ?", method)
		}
		if err := query.Limit(pageSizeInt).Offset(offset).Find(&result).Error; err != nil {
			ctx.JSONP(400, gin.H{
				"error": err.Error(),
			})
			return
		}

		var total int64
		if err := query.Model(&models.Permission{}).Count(&total).Error; err != nil {
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
