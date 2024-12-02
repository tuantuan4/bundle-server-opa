package routers

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"opa-test/controllers/RolePerm"
	"opa-test/controllers/UserRole"
	"opa-test/controllers/permissions"
	"opa-test/controllers/roles"
	"opa-test/controllers/users"
	"opa-test/middleware"
	"time"
)

func DefineRouter(r *gin.Engine, db *gorm.DB) {
	v := r.Group("/api/v1/admin")
	{
		//v.GET("/json", users.GetFileJson(db))
		//v.GET("/targz", controllers.GetFileTarGz(db)) //ver1
		v.GET("/targz", users.NotifyUpdate(db)) //ver2
		v.POST("/user", users.CreateUser(db))
		v.GET("/users", users.GetAllUser(db))

	}
	role := v.Group("/roles")
	{
		role.GET("", roles.GetAllRole(db))
		role.GET("/get", roles.GetRoleById(db))
		role.POST("", roles.CreateRole(db))
		role.POST("/count", roles.CreateRandomRoles(db))
	}
	perm := v.Group("/permissions")
	{
		perm.GET("", permissions.GetAllPerm(db))
		perm.GET("/:id_perm", permissions.GetPermById(db))
		perm.POST("", permissions.CreatePermission(db))
		perm.POST("/count", permissions.CreateRandomPermissions(db))
	}
	role_permission := v.Group("/rolePermission")
	{
		role_permission.POST("", RolePerm.CreateRolePerm(db))
		role_permission.Use(middleware.BasicAuth(middleware.ConvertBasicAuth())).POST("/checkPermission", RolePerm.GetRolePerm(db))
		role_permission.GET("/permissions/:id_role", RolePerm.GetListPermByRoleId(db))
		role_permission.GET("", RolePerm.GetAll(db))
		role_permission.GET("/exportJson", RolePerm.GetFileJsonData(db))
		role_permission.POST("/count", RolePerm.CreateRandomRolePerm(db))
	}
	user_role := v.Group("/userRole")
	{
		user_role.Use(middleware.BasicAuth(middleware.ConvertBasicAuth())).POST("", UserRole.CreateUserRole(db))
		user_role.GET("", UserRole.GetUserRole(db))
		user_role.POST("/count", UserRole.CreateUserRoleRandom(db))
	}
}

func Init(db *gorm.DB) {
	r := gin.Default()
	// Add CORS middleware before defining any routes
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:8889"},        // Allow requests from http://localhost:8889
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE"}, // HTTP methods to allow
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))
	DefineRouter(r, db)
	r.Run()
}
