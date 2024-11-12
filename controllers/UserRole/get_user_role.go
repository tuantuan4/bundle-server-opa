package UserRole

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"opa-test/models"
	"strconv"
)

type UserRoleResponse struct {
	Username string `json:"username"`
	ListRole []uint `json:"list_role"`
}

func GetUserRole(db *gorm.DB) func(ctx *gin.Context) {
	return func(ctx *gin.Context) {
		pageSize := ctx.DefaultQuery("pageSize", "10")
		pageNum := ctx.DefaultQuery("pageNum", "1")
		email := ctx.Query("email")
		username := ctx.Query("username")
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

		//if email == "" || username == "" {
		//	ctx.JSON(400, gin.H{
		//		"error": "email and username are required",
		//	})
		//	return
		//}

		var listUserRole []models.UserRole
		query := db.Model(&models.UserRole{})
		if email != "" {
			query = query.Where("email = ?", email)
		}
		if username != "" {
			query = query.Where("username = ?", username)
		}
		if err := query.Limit(pageSizeInt).Offset(offset).Find(&listUserRole).Error; err != nil {
			ctx.JSONP(400, gin.H{
				"error": err.Error(),
			})
			return
		}

		userRoleMap := make(map[string][]uint)
		for _, rp := range listUserRole {
			userRoleMap[rp.Username] = append(userRoleMap[rp.Username], rp.RoleId)
		}
		// Tạo slice chứa các response với cấu trúc {username, list_role}
		var response []UserRoleResponse
		for username, roles := range userRoleMap {
			response = append(response, UserRoleResponse{
				Username: username,
				ListRole: roles,
			})
		}
		ctx.JSON(200, gin.H{
			"result": gin.H{
				"data":  response,
				"total": len(response),
			},
		})
	}
}
