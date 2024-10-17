package permissions

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"math/rand"
	"opa-test/models"
	"time"
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

func RandStringBytes(n int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	rand.Seed(time.Now().UnixNano())
	b := make([]byte, n)
	for i := range b {
		b[i] = charset[rand.Intn(len(charset))]
	}
	return string(b)
}

// API tạo permission ngẫu nhiên với số lượng được truyền vào
func CreateRandomPermissions(db *gorm.DB) func(ctx *gin.Context) {
	return func(ctx *gin.Context) {

		var permissions []models.Permission

		// Các method hợp lệ
		methods := []string{"GET", "POST", "PUT", "DELETE"}

		for i := 0; i < 2000; i++ {
			perm := models.Permission{
				Name:   "Permission_" + RandStringBytes(5),       // Tên ngẫu nhiên
				Url:    "/api/v1/resource/" + RandStringBytes(3), // URL ngẫu nhiên
				Method: methods[rand.Intn(len(methods))],         // Method ngẫu nhiên
			}

			result := db.Create(&perm)
			if result.Error != nil {
				ctx.JSON(400, gin.H{"error": "Failed to create permission"})
				return
			}

			permissions = append(permissions, perm)
		}

		ctx.JSON(201, gin.H{"data": permissions})
	}
}
