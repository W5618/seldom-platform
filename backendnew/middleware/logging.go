package middleware

import (
	"seldom-platform/utils"
	"time"

	"github.com/gin-gonic/gin"
)

// LoggingMiddleware 日志记录中间件
func LoggingMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 记录开始时间
		startTime := time.Now()

		// 处理请求
		c.Next()

		// 计算处理时间
		duration := time.Since(startTime)

		// 记录请求日志
		utils.LogRequest(
			c.Request.Method,
			c.Request.URL.Path,
			c.ClientIP(),
			c.Writer.Status(),
			duration,
		)
	}
}

// ErrorLoggingMiddleware 错误日志记录中间件
func ErrorLoggingMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		// 检查是否有错误
		if len(c.Errors) > 0 {
			for _, err := range c.Errors {
				utils.LogError("Request error: %s - Path: %s, Method: %s, IP: %s", 
					err.Error(), c.Request.URL.Path, c.Request.Method, c.ClientIP())
			}
		}
	}
}