package routers

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"opa-test/controllers"
)

func DefineRouter(r *gin.Engine, db *gorm.DB) {
	v := r.Group("/api/v1")
	{
		v.GET("/json", controllers.GetFileJson(db))
		//v.GET("/targz", controllers.GetFileTarGz(db)) //ver1
		v.GET("/targz", controllers.NotifyUpdate(db)) //ver2
		v.POST("/user", controllers.CreateUser(db))
		v.GET("/user", controllers.GetAllUser(db))

	}
}

func Init(db *gorm.DB) {
	r := gin.Default()
	DefineRouter(r, db)
	r.Run()
}
