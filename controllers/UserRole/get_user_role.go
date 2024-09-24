package UserRole

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"net/http"
	"opa-test/models"
)

func GetUserRole(db *gorm.DB) func(ctx *gin.Context) {
	return func(ctx *gin.Context) {
		email := ctx.Query("email")
		username := ctx.Query("username")

		if email == "" || username == "" {
			ctx.JSON(400, gin.H{
				"error": "email and username are required",
			})
			return
		}

		var listUserRole []models.UserRole
		if err := db.Where("email = ? and username = ?", email, username).Find(&listUserRole).Error; err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error": "Query is failed",
			})
			return
		}

		userRoleMap := make(map[string][]uint)

		for _, rp := range listUserRole {
			userRoleMap[rp.Username] = append(userRoleMap[rp.Username], rp.RoleId)
		}

		ctx.JSON(200, gin.H{
			"data": userRoleMap[listUserRole[0].Username],
		})
	}
}
