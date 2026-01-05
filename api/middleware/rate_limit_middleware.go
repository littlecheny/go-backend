package middleware

import (
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/littlecheny/go-backend/domain"
)

// RateLimiter 简单的限流器
type RateLimiter struct {
	visitors map[string]*Visitor
	mu       sync.RWMutex
	rate     int           // 每分钟允许的请求数
	duration time.Duration // 时间窗口
}

type Visitor struct {
	lastSeen time.Time
	count    int
}

// NewRateLimiter 创建新的限流器
func NewRateLimiter(rate int, duration time.Duration) *RateLimiter {
	rl := &RateLimiter{
		visitors: make(map[string]*Visitor),
		rate:     rate,
		duration: duration,
	}
	go rl.cleanupVisitors()
	return rl
}

// RateLimitMiddleware 限流中间件
func (rl *RateLimiter) RateLimitMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		ip := c.ClientIP()
		
		rl.mu.Lock()
		visitor, exists := rl.visitors[ip]
		if !exists {
			rl.visitors[ip] = &Visitor{
				lastSeen: time.Now(),
				count:    1,
			}
			rl.mu.Unlock()
			c.Next()
			return
		}

		// 检查时间窗口
		if time.Since(visitor.lastSeen) > rl.duration {
			visitor.count = 1
			visitor.lastSeen = time.Now()
			rl.mu.Unlock()
			c.Next()
			return
		}

		// 检查请求次数
		if visitor.count >= rl.rate {
			rl.mu.Unlock()
			c.JSON(http.StatusTooManyRequests, domain.ErrorResponse{
				Message: "Rate limit exceeded. Please try again later.",
			})
			c.Abort()
			return
		}

		visitor.count++
		rl.mu.Unlock()
		c.Next()
	}
}

// 定期清理过期的访问者记录
func (rl *RateLimiter) cleanupVisitors() {
	for {
		time.Sleep(time.Minute)
		rl.mu.Lock()
		for ip, visitor := range rl.visitors {
			if time.Since(visitor.lastSeen) > rl.duration {
				delete(rl.visitors, ip)
			}
		}
		rl.mu.Unlock()
	}
}

