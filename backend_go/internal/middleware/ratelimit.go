package middleware

import (
	"sync"

	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
)

// IPRateLimiter stores per-IP limiters
type IPRateLimiter struct {
	limiters map[string]*rate.Limiter
	mu       sync.RWMutex
	rate     rate.Limit
	burst    int
}

// NewIPRateLimiter creates new rate limiter
func NewIPRateLimiter(rps int) *IPRateLimiter {
	limiterRate := rate.Limit(rps)
	return &IPRateLimiter{
		limiters: make(map[string]*rate.Limiter),
		rate:     limiterRate,
		burst:    rps * 60, // ~1 min burst
	}
}

// GetLimiter gets or creates limiter for IP
func (rl *IPRateLimiter) GetLimiter(ip string) *rate.Limiter {
	rl.mu.RLock()
	limiter, exists := rl.limiters[ip]
	rl.mu.RUnlock()

	if !exists {
		rl.mu.Lock()
		limiter, exists = rl.limiters[ip]
		if !exists {
			limiter = rate.NewLimiter(rl.rate, rl.burst)
			rl.limiters[ip] = limiter
		}
		rl.mu.Unlock()
	}

	return limiter
}

// RateLimitMiddleware middleware for Gin
func RateLimitMiddleware(rps int) gin.HandlerFunc {
	rl := NewIPRateLimiter(rps)
	return func(c *gin.Context) {
		clientIP := c.ClientIP()
		limiter := rl.GetLimiter(clientIP)

		if !limiter.Allow() {
			c.Header("Retry-After", "60")
			c.JSON(429, gin.H{
				"status": "error",
				"code":   429,
				"message": "Rate limit exceeded. Please try again later.",
			})
			c.Abort()
			return
		}

		c.Next()
	}
}

