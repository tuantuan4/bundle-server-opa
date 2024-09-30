package routers

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"opa-test/controllers/RolePerm"
	"opa-test/controllers/UserRole"
	"opa-test/controllers/permissions"
	"opa-test/controllers/roles"
	"opa-test/controllers/users"
	"opa-test/middleware"
)

func DefineRouter(r *gin.Engine, db *gorm.DB) {
	v := r.Group("/api/v1/admin")
	{
		//v.GET("/json", users.GetFileJson(db))
		//v.GET("/targz", controllers.GetFileTarGz(db)) //ver1
		v.GET("/targz", users.NotifyUpdate(db)) //ver2
		v.POST("/user", users.CreateUser(db))
		v.GET("/user", users.GetAllUser(db))

	}
	role := v.Group("/roles")
	{
		role.GET("", roles.GetAllRole(db))
		role.GET("/get", roles.GetRoleById(db))
		role.POST("", roles.CreateRole(db))
	}
	perm := v.Group("/permissions")
	{
		perm.GET("", permissions.GetAllPerm(db))
		perm.GET("/:id_perm", permissions.GetPermById(db))
		perm.POST("", permissions.CreatePermission(db))
	}
	role_permission := v.Group("/rolePermission")
	{
		role_permission.POST("", RolePerm.CreateRolePerm(db))
		role_permission.Use(middleware.BasicAuth(middleware.ConvertBasicAuth())).POST("/checkPermission", RolePerm.GetRolePerm(db))
		role_permission.GET("/permissions/:id_role", RolePerm.GetListPermByRoleId(db))
		role_permission.GET("", RolePerm.GetAll(db))
	}
	user_role := v.Group("/userRole")
	{
		user_role.Use(middleware.BasicAuth(middleware.ConvertBasicAuth())).POST("/:role_id", UserRole.CreateUserRole(db))
		user_role.Use(middleware.BasicAuth(middleware.ConvertBasicAuth())).GET("", UserRole.GetUserRole(db))
	}
}

func Init(db *gorm.DB) {
	r := gin.Default()
	DefineRouter(r, db)
	r.Run()
}
