package UserRole

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"log"
	"opa-test/models"
	"strconv"
	"strings"
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
func CreateUserRoleRandom(db *gorm.DB) func(ctx *gin.Context) {
	return func(ctx *gin.Context) {

		user := ctx.Query("username")
		email := ctx.Query("email")
		roleIdList := ctx.Query("role")
		stringList := strings.Split(roleIdList, ",")
		var intList []uint
		for _, str := range stringList {
			num, err := strconv.ParseUint(str, 10, 0)
			if err != nil {
				// Xử lý lỗi nếu không thể chuyển đổi
				fmt.Printf("Error converting '%s' to int: %v\n", str, err)
				continue // Bỏ qua phần tử này
			}
			intList = append(intList, uint(num))
		}

		if user == "" || email == "" {
			ctx.JSON(400, gin.H{"error": "username and email are required"})
			return
		}
		// Chèn các bản ghi UserRole cho từng roleId
		for _, roleId := range intList {
			var userRole models.UserRole
			userRole.Username = user
			userRole.Email = email
			userRole.RoleId = roleId // Gán RoleId từ danh sách

			if err := db.Create(&userRole).Error; err != nil {
				// Xử lý lỗi nếu không thể chèn bản ghi
				fmt.Printf("Failed to create user role for roleId %d: %v\n", roleId, err)
				continue // Bỏ qua lỗi và tiếp tục với roleId tiếp theo
			} else {
				log.Printf("Thanh cong to create user-role combination for username: %d and role ID: %d: ", user, roleId)
			}
		}

		// Phản hồi thành công
		ctx.JSON(200, gin.H{
			"message": "User roles created successfully",
		})
	}
}
