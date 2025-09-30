package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"seldom-platform/utils"
)

// SecurityHeaders 安全头中间件
func SecurityHeaders() gin.HandlerFunc {
	return func(c *gin.Context) {
		// X-Content-Type-Options
		c.Header("X-Content-Type-Options", "nosniff")
		
		// X-Frame-Options
		c.Header("X-Frame-Options", "DENY")
		
		// X-XSS-Protection
		c.Header("X-XSS-Protection", "1; mode=block")
		
		// Referrer-Policy
		c.Header("Referrer-Policy", "strict-origin-when-cross-origin")
		
		// Content-Security-Policy
		c.Header("Content-Security-Policy", "default-src 'self'; script-src 'self' 'unsafe-inline'; style-src 'self' 'unsafe-inline'; img-src 'self' data: https:; font-src 'self' https:; connect-src 'self' https:; media-src 'self'; object-src 'none'; child-src 'none'; worker-src 'none'; frame-ancestors 'none'; form-action 'self'; base-uri 'self'; manifest-src 'self'")
		
		// Strict-Transport-Security (仅在HTTPS下)
		if c.Request.TLS != nil {
			c.Header("Strict-Transport-Security", "max-age=31536000; includeSubDomains; preload")
		}
		
		// Permissions-Policy
		c.Header("Permissions-Policy", "geolocation=(), microphone=(), camera=()")
		
		c.Next()
	}
}

// IPWhitelistMiddleware IP白名单中间件
func IPWhitelistMiddleware(allowedIPs []string) gin.HandlerFunc {
	return func(c *gin.Context) {
		clientIP := utils.GetClientIP(c.Request)
		
		// 如果没有配置白名单，则允许所有IP
		if len(allowedIPs) == 0 {
			c.Next()
			return
		}
		
		// 检查IP是否在白名单中
		allowed := false
		for _, ip := range allowedIPs {
			if ip == clientIP || ip == "*" {
				allowed = true
				break
			}
		}
		
		if !allowed {
			utils.Forbidden(c, "IP地址不在允许范围内")
			c.Abort()
			return
		}
		
		c.Next()
	}
}

// UserAgentFilterMiddleware 用户代理过滤中间件
func UserAgentFilterMiddleware(blockedUserAgents []string) gin.HandlerFunc {
	return func(c *gin.Context) {
		userAgent := c.GetHeader("User-Agent")
		
		// 检查是否为被阻止的用户代理
		for _, blocked := range blockedUserAgents {
			if strings.Contains(strings.ToLower(userAgent), strings.ToLower(blocked)) {
				utils.Forbidden(c, "请求被拒绝")
				c.Abort()
				return
			}
		}
		
		c.Next()
	}
}

// RequestSizeLimit 请求大小限制中间件
func RequestSizeLimit(maxSize int64) gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.Request.ContentLength > maxSize {
			utils.BadRequest(c, "请求体过大")
			c.Abort()
			return
		}
		
		c.Next()
	}
}

// HTTPSRedirect HTTPS重定向中间件
func HTTPSRedirect() gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.Request.Header.Get("X-Forwarded-Proto") == "http" {
			httpsURL := "https://" + c.Request.Host + c.Request.RequestURI
			c.Redirect(http.StatusMovedPermanently, httpsURL)
			c.Abort()
			return
		}
		c.Next()
	}
}

// APIKeyMiddleware API密钥验证中间件
func APIKeyMiddleware(validAPIKeys []string) gin.HandlerFunc {
	return func(c *gin.Context) {
		apiKey := c.GetHeader("X-API-Key")
		if apiKey == "" {
			apiKey = c.Query("api_key")
		}
		
		if apiKey == "" {
			utils.Unauthorized(c, "缺少API密钥")
			c.Abort()
			return
		}
		
		// 验证API密钥
		valid := false
		for _, key := range validAPIKeys {
			if key == apiKey {
				valid = true
				break
			}
		}
		
		if !valid {
			utils.Unauthorized(c, "无效的API密钥")
			c.Abort()
			return
		}
		
		c.Next()
	}
}

// MethodOverrideMiddleware HTTP方法覆盖中间件
func MethodOverrideMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 检查_method参数或X-HTTP-Method-Override头
		method := c.PostForm("_method")
		if method == "" {
			method = c.GetHeader("X-HTTP-Method-Override")
		}
		
		if method != "" {
			method = strings.ToUpper(method)
			if method == "PUT" || method == "PATCH" || method == "DELETE" {
				c.Request.Method = method
			}
		}
		
		c.Next()
	}
}

// NoCache 禁用缓存中间件
func NoCache() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Cache-Control", "no-cache, no-store, must-revalidate")
		c.Header("Pragma", "no-cache")
		c.Header("Expires", "0")
		c.Next()
	}
}

// RequestIDMiddleware 请求ID中间件
func RequestIDMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		requestID := c.GetHeader("X-Request-ID")
		if requestID == "" {
			// 生成新的请求ID
			if id, err := utils.GenerateRandomString(16); err == nil {
				requestID = id
			} else {
				requestID = "unknown"
			}
		}
		
		c.Header("X-Request-ID", requestID)
		c.Set("request_id", requestID)
		c.Next()
	}
}

// ContentTypeValidation 内容类型验证中间件
func ContentTypeValidation(allowedTypes []string) gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.Request.Method == "POST" || c.Request.Method == "PUT" || c.Request.Method == "PATCH" {
			contentType := c.GetHeader("Content-Type")
			
			// 检查内容类型是否被允许
			allowed := false
			for _, allowedType := range allowedTypes {
				if strings.Contains(contentType, allowedType) {
					allowed = true
					break
				}
			}
			
			if !allowed {
				utils.BadRequest(c, "不支持的内容类型")
				c.Abort()
				return
			}
		}
		
		c.Next()
	}
}