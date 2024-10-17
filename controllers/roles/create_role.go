package roles

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"math/rand"
	"opa-test/models"
	"strconv"
	"time"
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

// Hàm tạo chuỗi ngẫu nhiên
func RandStringBytes(n int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	rand.Seed(time.Now().UnixNano())
	b := make([]byte, n)
	for i := range b {
		b[i] = charset[rand.Intn(len(charset))]
	}
	return string(b)
}

// API tạo role ngẫu nhiên với số lượng role được truyền vào
func CreateRandomRoles(db *gorm.DB) func(ctx *gin.Context) {
	return func(ctx *gin.Context) {

		//count, err := strconv.Atoi(ctx.Param("count"))
		//if err != nil {
		//	ctx.JSON(400, gin.H{
		//		"error": err.Error(),
		//	})
		//	return
		//}
		var roles []models.Role

		for i := 0; i < 20000; i++ {
			role := models.Role{
				Name:        "Role_" + RandStringBytes(5),
				Description: "Description_" + strconv.Itoa(rand.Intn(100)),
			}

			result := db.Create(&role)
			if result.Error != nil {
				ctx.JSON(400, gin.H{"error": "Failed to create role"})
				return
			}

			roles = append(roles, role)
		}

		ctx.JSON(201, gin.H{"data": roles})
	}
}
