package RolePerm

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"io/ioutil"
	"net/http"
	"opa-test/dto"
	"opa-test/models"
	"os"
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
		rolePermissionMap := make(map[uint][]models.Permission)
		for _, rp := range result {
			var perm models.Permission
			if err := db.Where("id = ?", rp.PermissionId).First(&perm).Error; err != nil {
				continue
			}
			rolePermissionMap[rp.RoleId] = append(rolePermissionMap[rp.RoleId], perm)
		}
		ctx.JSON(200, gin.H{"data": rolePermissionMap})
	}
}

func GetFileJsonData(db *gorm.DB) func(ctx *gin.Context) {
	return func(ctx *gin.Context) {
		var userRole []models.UserRole
		if err := db.Find(&userRole).Error; err != nil {
			ctx.JSON(400, gin.H{"error": err.Error()})
		}
		userRoleMap := make(map[string][]uint)
		for _, user := range userRole {
			userRoleMap[user.Username] = append(userRoleMap[user.Username], user.RoleId)
		}

		var result []models.RolePermission
		if err := db.Find(&result).Error; err != nil {
			ctx.JSON(400, gin.H{"error": err.Error()})
			return
		}
		rolePermissionMap := make(map[uint][]dto.PermissionDto)
		for _, rp := range result {
			var perm models.Permission
			if err := db.Where("id = ?", rp.PermissionId).First(&perm).Error; err != nil {
				continue
			}
			rolePermissionMap[rp.RoleId] = append(rolePermissionMap[rp.RoleId], dto.ToPermissionDto(perm))
		}

		response := ResponseData{
			userRoleMap, rolePermissionMap,
		}
		// Chuyển đổi dữ liệu thành JSON
		jsonData, err := json.MarshalIndent(response, "", "  ")
		if err != nil {
			ctx.JSON(500, gin.H{"error": "Unable to generate JSON"})
			return
		}

		// Xác định thư mục để lưu file (nếu chưa có, tạo thư mục)
		outputDir := "./bundle/data"
		if _, err := os.Stat(outputDir); os.IsNotExist(err) {
			err := os.Mkdir(outputDir, 0755) // Tạo thư mục với quyền truy cập đọc/ghi cho người sở hữu
			if err != nil {
				ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to create directory"})
				return
			}
		}

		filePath := fmt.Sprintf("%s/data.json", outputDir)
		// Ghi file JSON vào thư mục
		err = ioutil.WriteFile(filePath, jsonData, 0644) // Quyền 0644 cho phép đọc/ghi đối với chủ sở hữu, chỉ đọc cho những người khác
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to write file"})
			return
		}
		// Trả về phản hồi thành công
		ctx.JSON(http.StatusOK, gin.H{
			"message": "File exported successfully",
			"path":    filePath,
		})
		//ctx.JSON(200, gin.H{
		//	"user_roles":  userRoleMap,
		//	"role_grants": rolePermissionMap,
		//})
	}
}

type ResponseData struct {
	UseRoles   map[string][]uint            `json:"user_roles"`
	RoleGrants map[uint][]dto.PermissionDto `json:"role_grants"`
}
