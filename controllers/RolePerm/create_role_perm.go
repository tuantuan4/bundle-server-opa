package RolePerm

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"opa-test/models"
	"strconv"
)

func CreateRolePerm(db *gorm.DB) func(ctx *gin.Context) {
	return func(ctx *gin.Context) {

		role, err1 := strconv.Atoi(ctx.Query("role_id"))
		perm, err2 := strconv.Atoi(ctx.Query("permission_id"))
		if err1 != nil || err2 != nil {
			ctx.JSON(400, gin.H{"error": "invalid role or permission id"})
			return
		}
		var rolePerm models.RolePermission
		rolePerm.RoleId = uint(role)
		rolePerm.PermissionId = uint(perm)
		if err := db.Create(&rolePerm).Error; err != nil {
			ctx.JSON(400, gin.H{
				"error": "Failed to create permission",
			})
			return
		}

		ctx.JSON(200, gin.H{
			"data": rolePerm,
		})
	}
}
