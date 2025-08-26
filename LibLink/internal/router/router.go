package router

import (
	"liblink/internal/controllers/api"
	"liblink/internal/middleware"

	"github.com/gin-gonic/gin"
)

func Router() *gin.Engine {
	router := gin.Default()

	router.Use(middleware.CORS())
	router.POST("/register", api.Register)
	router.POST("/login", api.Login)
	router.GET("/ping_without_login", api.Ping)

	router.Static("/static", "./static")
	authRoutes := router.Group("/api")
	authRoutes.Use(middleware.JWTAuth())
	{
		authRoutes.GET("/ping", api.Ping)
		// 用户相关
		users := authRoutes.Group("/users")
		{
			users.GET("/summary", api.UsersSummary)
		}
		// 系统相关
		system := authRoutes.Group("/system")
		{
			notification := system.Group("/notifications")
			{
				notification.GET("/list", api.Notifications)
				notification.POST("/add", api.AddNotification)
			}
		}
		archives := authRoutes.Group("/archives")
		{
			archives.GET("/list", api.GetArchives)
			archives.POST("/add", api.AddArchive)
			archives.PATCH("/borrow", api.BorrowArchive)
			archives.PATCH("/return", api.ReturnArchive)
			archives.PUT("/update/:id", api.UpdateArchive)
			archives.POST("/batch_import", api.BatchImportArchives)
			archives.POST("/batch_operate", api.BatchOperateArchives)
		}
	}

	return router
}
