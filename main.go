package main

import (
	_ "github.com/binginx/bqd_chat_log/docs" // 导入swagger文档
	"github.com/binginx/bqd_chat_log/handlers"
	"github.com/binginx/bqd_chat_log/models"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @title           Chat Log API
// @version         1.0
// @description     记录聊天信息的API服务
// @termsOfService  http://swagger.io/terms/

// @contact.name   API Support
// @contact.url    http://www.example.com/support
// @contact.email  support@example.com

// @license.name  MIT
// @license.url   https://opensource.org/licenses/MIT

// @host      localhost:8080
// @BasePath  /api/v1

func main() {
	// 初始化存储系统
	models.InitStorage()

	// 创建Gin引擎
	r := gin.Default()

	// API路由组
	v1 := r.Group("/api/v1")
	{
		v1.POST("/logs", handlers.CreateLog)
		v1.GET("/logs", handlers.GetLogs)
	}

	// Swagger文档
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// 启动服务器
	r.Run(":8090")
}
