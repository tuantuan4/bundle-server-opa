package RolePerm

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"net/http"
	"opa-test/models"
	"strconv"
)

type RolePermRequest struct {
	RoleId []uint `json:"role_id"`
	Url    string `json:"url"`
	Method string `json:"method"`
}

func GetRolePerm(db *gorm.DB) func(ctx *gin.Context) {
	return func(ctx *gin.Context) {
		var request RolePermRequest

		err := ctx.BindJSON(&request)
		if err != nil {
			ctx.JSON(400, gin.H{
				"error": "error",
			})
			return
		}
		var method, url, id_role = request.Method, request.Url, request.RoleId
		if method == "" || url == "" || len(id_role) == 0 {
			ctx.JSON(400, gin.H{
				"message": "Method and URL and ID_ROLE is required",
			})
		}
		var perm models.Permission

		if err := db.Where("method = ? AND url = ?", method, url).First(&perm).Error; err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error": "permission is not found",
			})
			return
		}

		var rolePerm models.RolePermission

		if err := db.Where("role_id IN ? AND permission_id = ?", id_role, perm.ID).First(&rolePerm).Error; err != nil {
			ctx.JSON(http.StatusOK, gin.H{
				"result": "false",
			})
			return
		}
		ctx.JSON(http.StatusOK, gin.H{
			"result": "true",
		})

	}
}

func GetListPermByRoleId(db *gorm.DB) func(ctx *gin.Context) {
	return func(ctx *gin.Context) {

		id, err := strconv.Atoi(ctx.Param("id_role"))
		if err != nil {
			ctx.JSON(400, gin.H{
				"message": "error",
			})
		}

		var role models.Role
		if err := db.Where("id = ?", id).First(&role).Error; err != nil {
			ctx.JSON(http.StatusNotFound, gin.H{
				"error": "id role not found",
			})
			return
		}
		var listPerm []uint
		var rolePerm []models.RolePermission
		if err := db.Where("role_id = ?", uint(id)).Find(&rolePerm).Error; err != nil {
			ctx.JSON(200, gin.H{
				"error": "Get list error",
			})
		}
		for _, obj := range rolePerm {
			listPerm = append(listPerm, obj.PermissionId)
		}
		ctx.JSON(200, gin.H{
			"data": listPerm,
		})
	}
}

func GetAll(db *gorm.DB) func(ctx *gin.Context) {
	return func(ctx *gin.Context) {
		var result []models.RolePermission
		if err := db.Find(&result).Error; err != nil {
			ctx.JSON(400, gin.H{"error": err.Error()})
			return
		}
		rolePermissionMap := make(map[uint][]uint)
		for _, rp := range result {
			rolePermissionMap[rp.RoleId] = append(rolePermissionMap[rp.RoleId], rp.PermissionId)
		}
		ctx.JSON(200, gin.H{"data": rolePermissionMap})
	}
}
