package routes

import (
	"seldom-platform/config"
	"seldom-platform/handlers"
	"seldom-platform/middleware"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// Setup 创建并配置Gin引擎
func Setup(cfg *config.Config) *gin.Engine {
	r := gin.Default()
	SetupRoutes(r, cfg)
	return r
}

// SetupRoutes 设置路由
func SetupRoutes(r *gin.Engine, cfg *config.Config) {
	// 添加中间件
	r.Use(middleware.CORSMiddleware())

	// Swagger文档路由
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// 健康检查
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":  "ok",
			"message": "Seldom Platform API is running",
		})
	})

	// API路由组
	api := r.Group("/api")

	// 认证相关路由（不需要认证）
	authHandler := handlers.NewAuthHandler(cfg)
	auth := api.Group("/auth")
	{
		auth.POST("/login", authHandler.Login)
		auth.POST("/register", authHandler.Register)
	}

	// 需要认证的路由
	authenticated := api.Group("")
	authenticated.Use(middleware.AuthMiddleware(cfg))
	{
		// 用户信息路由
		authenticated.GET("/auth/profile", authHandler.GetProfile)
		authenticated.PUT("/auth/profile", authHandler.UpdateProfile)

		// 项目管理路由
		projectHandler := handlers.NewProjectHandler()
		projects := authenticated.Group("/projects")
		{
			projects.GET("", projectHandler.GetProjects)
			projects.POST("", projectHandler.CreateProject)
			projects.GET("/:id", projectHandler.GetProject)
			projects.PUT("/:id", projectHandler.UpdateProject)
			projects.DELETE("/:id", projectHandler.DeleteProject)
		}

		// 测试用例管理路由
		caseHandler := handlers.NewCaseHandler()
		cases := authenticated.Group("/cases")
		{
			cases.GET("", caseHandler.GetCases)
			cases.POST("", caseHandler.CreateCase)
			cases.GET("/:id", caseHandler.GetCase)
			cases.PUT("/:id", caseHandler.UpdateCase)
			cases.DELETE("/:id", caseHandler.DeleteCase)
			cases.POST("/:id/copy", caseHandler.CopyCase)
		}

		// 环境管理路由
		envHandler := handlers.NewEnvHandler()
		envs := authenticated.Group("/envs")
		{
			envs.GET("", envHandler.GetEnvs)
			envs.POST("", envHandler.CreateEnv)
			envs.GET("/:id", envHandler.GetEnv)
			envs.PUT("/:id", envHandler.UpdateEnv)
			envs.DELETE("/:id", envHandler.DeleteEnv)
		}

		// 任务管理路由
		taskHandler := handlers.NewTaskHandler()
		tasks := authenticated.Group("/tasks")
		{
			tasks.GET("", taskHandler.GetTasks)
			tasks.POST("", taskHandler.CreateTask)
			tasks.GET("/:id", taskHandler.GetTask)
			tasks.PUT("/:id", taskHandler.UpdateTask)
			tasks.DELETE("/:id", taskHandler.DeleteTask)
			tasks.POST("/:id/run", taskHandler.RunTask)
			tasks.GET("/:id/reports", taskHandler.GetTaskReports)
		}

		// 团队管理路由
		teamHandler := handlers.NewTeamHandler()
		teams := authenticated.Group("/teams")
		{
			teams.GET("", teamHandler.GetTeams)
			teams.POST("", teamHandler.CreateTeam)
			teams.GET("/:id", teamHandler.GetTeam)
			teams.PUT("/:id", teamHandler.UpdateTeam)
			teams.DELETE("/:id", teamHandler.DeleteTeam)
		}
	}
}