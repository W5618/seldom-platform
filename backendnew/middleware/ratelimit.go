package middleware

import (
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

// RateLimiter 限流器结构
type RateLimiter struct {
	visitors map[string]*Visitor
	mu       sync.RWMutex
	rate     int           // 每分钟允许的请求数
	window   time.Duration // 时间窗口
}

// Visitor 访问者信息
type Visitor struct {
	requests []time.Time
	mu       sync.Mutex
}

// NewRateLimiter 创建新的限流器
func NewRateLimiter(rate int, window time.Duration) *RateLimiter {
	rl := &RateLimiter{
		visitors: make(map[string]*Visitor),
		rate:     rate,
		window:   window,
	}

	// 启动清理goroutine
	go rl.cleanup()

	return rl
}

// Allow 检查是否允许请求
func (rl *RateLimiter) Allow(ip string) bool {
	rl.mu.RLock()
	visitor, exists := rl.visitors[ip]
	rl.mu.RUnlock()

	if !exists {
		rl.mu.Lock()
		visitor = &Visitor{
			requests: make([]time.Time, 0),
		}
		rl.visitors[ip] = visitor
		rl.mu.Unlock()
	}

	visitor.mu.Lock()
	defer visitor.mu.Unlock()

	now := time.Now()
	
	// 清理过期的请求记录
	cutoff := now.Add(-rl.window)
	validRequests := make([]time.Time, 0)
	for _, reqTime := range visitor.requests {
		if reqTime.After(cutoff) {
			validRequests = append(validRequests, reqTime)
		}
	}
	visitor.requests = validRequests

	// 检查是否超过限制
	if len(visitor.requests) >= rl.rate {
		return false
	}

	// 添加当前请求
	visitor.requests = append(visitor.requests, now)
	return true
}

// cleanup 清理过期的访问者记录
func (rl *RateLimiter) cleanup() {
	ticker := time.NewTicker(time.Minute)
	defer ticker.Stop()

	for range ticker.C {
		rl.mu.Lock()
		now := time.Now()
		cutoff := now.Add(-rl.window * 2) // 保留更长时间以避免频繁创建

		for ip, visitor := range rl.visitors {
			visitor.mu.Lock()
			if len(visitor.requests) == 0 || (len(visitor.requests) > 0 && visitor.requests[len(visitor.requests)-1].Before(cutoff)) {
				delete(rl.visitors, ip)
			}
			visitor.mu.Unlock()
		}
		rl.mu.Unlock()
	}
}

// RateLimitMiddleware 限流中间件
func RateLimitMiddleware(rate int, window time.Duration) gin.HandlerFunc {
	limiter := NewRateLimiter(rate, window)

	return func(c *gin.Context) {
		ip := c.ClientIP()
		
		if !limiter.Allow(ip) {
			c.JSON(http.StatusTooManyRequests, gin.H{
				"error": "Too many requests",
				"code":  429,
			})
			c.Abort()
			return
		}

		c.Next()
	}
}

// APIRateLimitMiddleware API限流中间件（更严格的限制）
func APIRateLimitMiddleware() gin.HandlerFunc {
	return RateLimitMiddleware(100, time.Minute) // 每分钟100个请求
}

// AuthRateLimitMiddleware 认证接口限流中间件（防止暴力破解）
func AuthRateLimitMiddleware() gin.HandlerFunc {
	return RateLimitMiddleware(10, time.Minute) // 每分钟10个请求
}