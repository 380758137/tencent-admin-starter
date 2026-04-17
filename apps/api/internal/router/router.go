package router

import (
	"net/http"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/liusheng/tencent-admin-starter/apps/api/internal/config"
	"github.com/liusheng/tencent-admin-starter/apps/api/internal/handlers"
	"github.com/liusheng/tencent-admin-starter/apps/api/internal/middleware"
	"github.com/liusheng/tencent-admin-starter/apps/api/internal/response"
	"gorm.io/gorm"
)

func New(cfg config.Config, db *gorm.DB) *gin.Engine {
	r := gin.New()
	r.Use(gin.Logger(), gin.Recovery())
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Authorization", "Content-Type"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: false,
	}))

	r.GET("/healthz", func(c *gin.Context) {
		c.JSON(http.StatusOK, response.Success(gin.H{"status": "ok"}))
	})

	authHandler := handlers.NewAuthHandler(cfg, db)
	userHandler := handlers.NewUserHandler(db)
	departmentHandler := handlers.NewDepartmentHandler(db)
	roleHandler := handlers.NewRoleHandler(db)
	menuHandler := handlers.NewMenuHandler(db)
	dictionaryHandler := handlers.NewDictionaryHandler(db)
	systemParamHandler := handlers.NewSystemParamHandler(db)
	logHandler := handlers.NewLogHandler(db)
	monitorHandler := handlers.NewMonitorHandler(db)
	positionHandler := handlers.NewPositionHandler(db)
	noticeHandler := handlers.NewNoticeHandler(db)
	onlineUserHandler := handlers.NewOnlineUserHandler(db)
	scheduledJobHandler := handlers.NewScheduledJobHandler(db)

	api := r.Group("/api")
	{
		api.Use(middleware.OperationLog(db))
		api.POST("/auth/login", authHandler.Login)

		secured := api.Group("")
		secured.Use(middleware.Auth(cfg, db))
		{
			secured.GET("/auth/me", authHandler.Me)

			secured.GET("/users", middleware.RequirePerms(db, "user:list"), userHandler.List)
			secured.POST("/users", middleware.RequirePerms(db, "user:create"), userHandler.Create)
			secured.PUT("/users/:id", middleware.RequirePerms(db, "user:update"), userHandler.Update)
			secured.DELETE("/users/:id", middleware.RequirePerms(db, "user:delete"), userHandler.Delete)
			secured.GET("/users/export", middleware.RequirePerms(db, "user:export"), userHandler.ExportCSV)
			secured.POST("/users/import", middleware.RequirePerms(db, "user:import"), userHandler.ImportCSV)

			secured.GET("/departments", middleware.RequirePerms(db, "department:list"), departmentHandler.List)
			secured.POST("/departments", middleware.RequirePerms(db, "department:create"), departmentHandler.Create)
			secured.PUT("/departments/:id", middleware.RequirePerms(db, "department:update"), departmentHandler.Update)
			secured.DELETE("/departments/:id", middleware.RequirePerms(db, "department:delete"), departmentHandler.Delete)

			secured.GET("/roles", middleware.RequirePerms(db, "role:list"), roleHandler.List)
			secured.POST("/roles", middleware.RequirePerms(db, "role:create"), roleHandler.Create)
			secured.PUT("/roles/:id", middleware.RequirePerms(db, "role:update"), roleHandler.Update)
			secured.DELETE("/roles/:id", middleware.RequirePerms(db, "role:delete"), roleHandler.Delete)

			secured.GET("/menus", middleware.RequirePerms(db, "menu:list"), menuHandler.List)
			secured.POST("/menus", middleware.RequirePerms(db, "menu:create"), menuHandler.Create)
			secured.PUT("/menus/:id", middleware.RequirePerms(db, "menu:update"), menuHandler.Update)
			secured.DELETE("/menus/:id", middleware.RequirePerms(db, "menu:delete"), menuHandler.Delete)

			secured.GET("/dictionary-items", middleware.RequirePerms(db, "dictionary:list"), dictionaryHandler.List)
			secured.POST("/dictionary-items", middleware.RequirePerms(db, "dictionary:create"), dictionaryHandler.Create)
			secured.PUT("/dictionary-items/:id", middleware.RequirePerms(db, "dictionary:update"), dictionaryHandler.Update)
			secured.DELETE("/dictionary-items/:id", middleware.RequirePerms(db, "dictionary:delete"), dictionaryHandler.Delete)

			secured.GET("/system-params", middleware.RequirePerms(db, "param:list"), systemParamHandler.List)
			secured.POST("/system-params", middleware.RequirePerms(db, "param:create"), systemParamHandler.Create)
			secured.PUT("/system-params/:id", middleware.RequirePerms(db, "param:update"), systemParamHandler.Update)
			secured.DELETE("/system-params/:id", middleware.RequirePerms(db, "param:delete"), systemParamHandler.Delete)

			secured.GET("/operation-logs", middleware.RequirePerms(db, "log:operation:list"), logHandler.ListOperations)
			secured.GET("/login-logs", middleware.RequirePerms(db, "log:login:list"), logHandler.ListLogins)
			secured.GET("/monitor/overview", middleware.RequirePerms(db, "monitor:view"), monitorHandler.Overview)

			secured.GET("/positions", middleware.RequirePerms(db, "position:list"), positionHandler.List)
			secured.POST("/positions", middleware.RequirePerms(db, "position:create"), positionHandler.Create)
			secured.PUT("/positions/:id", middleware.RequirePerms(db, "position:update"), positionHandler.Update)
			secured.DELETE("/positions/:id", middleware.RequirePerms(db, "position:delete"), positionHandler.Delete)

			secured.GET("/notices", middleware.RequirePerms(db, "notice:list"), noticeHandler.List)
			secured.POST("/notices", middleware.RequirePerms(db, "notice:create"), noticeHandler.Create)
			secured.PUT("/notices/:id", middleware.RequirePerms(db, "notice:update"), noticeHandler.Update)
			secured.DELETE("/notices/:id", middleware.RequirePerms(db, "notice:delete"), noticeHandler.Delete)

			secured.GET("/online-users", middleware.RequirePerms(db, "online-user:list"), onlineUserHandler.List)

			secured.GET("/scheduled-jobs", middleware.RequirePerms(db, "job:list"), scheduledJobHandler.List)
			secured.POST("/scheduled-jobs", middleware.RequirePerms(db, "job:create"), scheduledJobHandler.Create)
			secured.PUT("/scheduled-jobs/:id", middleware.RequirePerms(db, "job:update"), scheduledJobHandler.Update)
			secured.DELETE("/scheduled-jobs/:id", middleware.RequirePerms(db, "job:delete"), scheduledJobHandler.Delete)
			secured.POST("/scheduled-jobs/:id/run", middleware.RequirePerms(db, "job:run"), scheduledJobHandler.RunOnce)
			secured.GET("/scheduled-job-logs", middleware.RequirePerms(db, "job-log:list"), scheduledJobHandler.ListLogs)
		}
	}

	return r
}
