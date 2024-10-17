package RolePerm

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"log"
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

// Hàm để chia dữ liệu thành các batch
func batchInsert(db *gorm.DB, data []models.RolePermission, batchSize int) error {
	for i := 0; i < len(data); i += batchSize {
		end := i + batchSize
		if end > len(data) {
			end = len(data)
		}
		// Chèn 1000 bản ghi mỗi lần
		if err := db.Create(data[i:end]).Error; err != nil {
			return err
		}
	}
	return nil
}

// Hàm tạo RolePerm ngẫu nhiên với role và permission đã tồn tại trong DB
func CreateRandomRolePerm(db *gorm.DB) func(ctx *gin.Context) {
	return func(ctx *gin.Context) {
		// Lấy tất cả role từ DB
		var roles []models.Role
		if err := db.Find(&roles).Error; err != nil || len(roles) == 0 {
			ctx.JSON(400, gin.H{"error": "No roles found"})
			return
		}

		// Lấy tất cả permission từ DB
		var permissions []models.Permission
		if err := db.Find(&permissions).Error; err != nil || len(permissions) == 0 {
			ctx.JSON(400, gin.H{"error": "No permissions found"})
			return
		}

		//// Bắt đầu một transaction
		//tx := db.Begin()
		//defer func() {
		//	if r := recover(); r != nil {
		//		tx.Rollback() // Rollback nếu có lỗi
		//		ctx.JSON(500, gin.H{"error": "Internal server error"})
		//	}
		//}()
		//
		//// Tạo từng tổ hợp role-permission và chèn vào DB
		//for _, role := range roles {
		//	for _, perm := range permissions {
		//		rolePerm := models.RolePermission{
		//			RoleId:       role.ID,
		//			PermissionId: perm.ID,
		//		}
		//		// Chèn từng bản ghi vào DB
		//		if err := tx.Create(&rolePerm).Error; err != nil {
		//			log.Printf("Failed to create role-permission combination for Role ID: %d and Permission ID: %d: %v", role.ID, perm.ID, err)
		//
		//			continue
		//		} else {
		//			log.Printf("Insert thanh cong to create role-permission combination for Role ID: %d and Permission ID: %d: ", role.ID, perm.ID)
		//		}
		//	}
		//}
		//
		//// Commit transaction nếu thành công
		//tx.Commit()
		//ctx.JSON(200, gin.H{"message": "All role-permission combinations created successfully"})
		// Tạo từng tổ hợp role-permission và chèn trực tiếp vào DB
		for _, role := range roles {
			for _, perm := range permissions {
				rolePerm := models.RolePermission{
					RoleId:       role.ID,
					PermissionId: perm.ID,
				}

				// Chèn trực tiếp từng bản ghi vào DB
				if err := db.Create(&rolePerm).Error; err != nil {
					// Bỏ qua lỗi và tiếp tục, có thể log lỗi
					log.Printf("Failed to create role-permission combination for Role ID: %d and Permission ID: %d: %v", role.ID, perm.ID, err)
					continue
				} else {
					log.Printf("Thanh cong to create role-permission combination for Role ID: %d and Permission ID: %d: ", role.ID, perm.ID)
				}
			}
		}

		// Phản hồi khi hoàn tất
		ctx.JSON(200, gin.H{"message": "All role-permission combinations created successfully"})
	}
}
