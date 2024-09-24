package UserRole

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"opa-test/models"
	"strconv"
)

func CreateUserRole(db *gorm.DB) func(ctx *gin.Context) {
	return func(ctx *gin.Context) {

		user := ctx.Query("username")
		email := ctx.Query("email")
		role, err2 := strconv.Atoi(ctx.Param("role_id"))
		if err2 != nil {
			ctx.JSON(400, gin.H{"error": "invalid role id"})
			return
		}

		if user == "" || email == "" {
			ctx.JSON(400, gin.H{"error": "username and email are required"})
			return
		}
		var userRole models.UserRole
		userRole.Username = user
		userRole.Email = email
		userRole.RoleId = uint(role)
		if err := db.Create(&userRole).Error; err != nil {
			ctx.JSON(400, gin.H{
				"error": "Failed to create user role",
			})
			return
		}

		ctx.JSON(200, gin.H{
			"data": userRole,
		})
	}
}
