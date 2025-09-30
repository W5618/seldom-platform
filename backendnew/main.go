package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"seldom-platform/config"
	"seldom-platform/database"
	"seldom-platform/routes"
	"seldom-platform/services"
	"seldom-platform/utils"
)

// @title Seldom Platform API
// @version 1.0
// @description 测试平台后端API接口文档
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host localhost:8080
// @BasePath /api

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
func main() {
	// 加载配置
	cfg := config.Load()

	// 初始化日志记录器
	if err := utils.InitLogger(); err != nil {
		log.Fatal("Failed to initialize logger:", err)
	}

	// 初始化数据库
	db, err := database.Init(cfg.Database)
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	defer database.Close(db)

	// 设置Gin模式
	if cfg.Server.Mode == "release" {
		gin.SetMode(gin.ReleaseMode)
	}

	// 初始化并启动调度服务
	if err := services.InitGlobalScheduler(); err != nil {
		log.Fatal("Failed to start scheduler service:", err)
	}
	defer services.StopGlobalScheduler()

	// 初始化路由
	r := routes.Setup(cfg)

	// 启动服务器
	log.Printf("Server starting on port %s", cfg.Server.Port)
	if err := r.Run(":" + cfg.Server.Port); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}