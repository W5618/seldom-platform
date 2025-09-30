package middleware

import (
	"fmt"
	"net/http"
	"runtime/debug"

	"github.com/gin-gonic/gin"
	"seldom-platform/utils"
)

// RecoveryMiddleware 错误恢复中间件
func RecoveryMiddleware() gin.HandlerFunc {
	return gin.CustomRecovery(func(c *gin.Context, recovered interface{}) {
		// 记录panic信息
		stack := debug.Stack()
		errorMsg := fmt.Sprintf("Panic recovered: %v\nStack trace:\n%s", recovered, string(stack))
		
		// 记录到日志
		if logger := utils.GetLogger(); logger != nil {
			logger.LogError("PANIC", errorMsg, map[string]interface{}{
				"method": c.Request.Method,
				"path":   c.Request.URL.Path,
				"ip":     c.ClientIP(),
			})
		}

		// 返回500错误
		utils.InternalServerError(c, "服务器内部错误")
		c.Abort()
	})
}

// CustomRecoveryWithWriter 自定义恢复中间件（带写入器）
func CustomRecoveryWithWriter() gin.HandlerFunc {
	return gin.CustomRecoveryWithWriter(nil, func(c *gin.Context, recovered interface{}) {
		// 获取错误信息
		var errorMsg string
		switch err := recovered.(type) {
		case string:
			errorMsg = err
		case error:
			errorMsg = err.Error()
		default:
			errorMsg = fmt.Sprintf("%v", err)
		}

		// 记录详细的错误信息
		stack := debug.Stack()
		fullErrorMsg := fmt.Sprintf("Panic: %s\nStack: %s", errorMsg, string(stack))
		
		// 记录到日志
		if logger := utils.GetLogger(); logger != nil {
			logger.LogError("PANIC_RECOVERY", fullErrorMsg, map[string]interface{}{
				"method":     c.Request.Method,
				"path":       c.Request.URL.Path,
				"ip":         c.ClientIP(),
				"user_agent": c.Request.UserAgent(),
				"headers":    c.Request.Header,
			})
		}

		// 根据Accept头返回不同格式的错误
		accept := c.GetHeader("Accept")
		if gin.IsDebugging() {
			// 开发模式下返回详细错误信息
			c.JSON(http.StatusInternalServerError, gin.H{
				"error":   "Internal Server Error",
				"message": errorMsg,
				"stack":   string(stack),
			})
		} else {
			// 生产模式下返回简单错误信息
			if accept == "application/json" {
				utils.InternalServerError(c, "服务器内部错误")
			} else {
				c.String(http.StatusInternalServerError, "Internal Server Error")
			}
		}
		
		c.Abort()
	})
}

// PanicHandler 处理panic的函数
func PanicHandler(c *gin.Context, err interface{}) {
	// 获取堆栈信息
	stack := debug.Stack()
	
	// 构造错误信息
	errorInfo := map[string]interface{}{
		"error":      fmt.Sprintf("%v", err),
		"method":     c.Request.Method,
		"path":       c.Request.URL.Path,
		"ip":         c.ClientIP(),
		"user_agent": c.Request.UserAgent(),
		"stack":      string(stack),
	}

	// 记录错误日志
	if logger := utils.GetLogger(); logger != nil {
		logger.LogError("PANIC", fmt.Sprintf("Panic recovered: %v", err), errorInfo)
	}

	// 返回错误响应
	utils.InternalServerError(c, "服务器发生了意外错误")
}

// SafeHandler 安全处理器包装函数
func SafeHandler(handler gin.HandlerFunc) gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				PanicHandler(c, err)
			}
		}()
		handler(c)
	}
}

// SafeAsyncHandler 安全异步处理器
func SafeAsyncHandler(handler func()) {
	go func() {
		defer func() {
			if err := recover(); err != nil {
				// 记录异步任务中的panic
				stack := debug.Stack()
				errorMsg := fmt.Sprintf("Async panic recovered: %v\nStack: %s", err, string(stack))
				
				if logger := utils.GetLogger(); logger != nil {
					logger.LogError("ASYNC_PANIC", errorMsg, map[string]interface{}{
						"type": "async_task",
					})
				}
			}
		}()
		handler()
	}()
}