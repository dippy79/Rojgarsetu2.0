package middleware

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

// CacheMiddleware creates a Redis caching middleware with configurable TTL
func CacheMiddleware(redisClient *redis.Client, ttl time.Duration) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Skip caching for non-GET requests
		if c.Request.Method != "GET" {
			c.Next()
			return
		}

		// Generate cache key
		cacheKey := generateCacheKey(c)

		// Try to get from cache
		ctx := context.Background()
		cachedData, err := redisClient.Get(ctx, cacheKey).Bytes()
		if err == nil {
			// Cache hit
			var response interface{}
			if err := json.Unmarshal(cachedData, &response); err == nil {
				c.Header("X-Cache", "HIT")
				c.JSON(200, response)
				c.Abort()
				return
			}
		}

		// Cache miss - proceed with request
		c.Header("X-Cache", "MISS")

		// Store original writer
		writer := c.Writer

		// Create custom response writer to capture response
		blw := &bodyLogWriter{body: []byte{}, ResponseWriter: writer}
		c.Writer = blw

		c.Next()

		// Cache successful responses
		if c.Writer.Status() == 200 && len(blw.body) > 0 {
			// Only cache if response is valid JSON
			var js interface{}
			if json.Unmarshal(blw.body, &js) == nil {
				if err := redisClient.Set(ctx, cacheKey, blw.body, ttl).Err(); err == nil {
					c.Header("X-Cache-Status", "CACHED")
				}
			}
		}
	}
}

// generateCacheKey creates a unique cache key based on request method, path, and query
func generateCacheKey(c *gin.Context) string {
	return fmt.Sprintf("cache:%s:%s:%s", c.Request.Method, c.Request.URL.Path, c.Request.URL.RawQuery)
}

// InvalidateCache invalidates cache entries matching a pattern
func InvalidateCache(redisClient *redis.Client, pattern string) error {
	ctx := context.Background()
	iter := redisClient.Scan(ctx, 0, pattern, 0).Iterator()
	for iter.Next(ctx) {
		if err := redisClient.Del(ctx, iter.Val()).Err(); err != nil {
			return err
		}
	}
	return iter.Err()
}

// InvalidateCacheByPattern invalidates cache entries by path pattern
func InvalidateCacheByPath(redisClient *redis.Client, method, path string) error {
	pattern := fmt.Sprintf("cache:%s:%s:*", method, path)
	return InvalidateCache(redisClient, pattern)
}

// CacheInvalidationMiddleware invalidates cache on write operations
func CacheInvalidationMiddleware(redisClient *redis.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Only invalidate for write operations
		if c.Request.Method == "POST" || c.Request.Method == "PUT" || c.Request.Method == "DELETE" {
			c.Next()

			// Invalidate cache based on the path pattern after successful write
			if c.Writer.Status() >= 200 && c.Writer.Status() < 300 {
				path := c.Request.URL.Path

				// Invalidate related cache entries
				switch {
				case strings.Contains(path, "/gov-jobs"):
					InvalidateCacheByPath(redisClient, "GET", "/api/v1/gov-jobs")
				case strings.Contains(path, "/private-jobs"):
					InvalidateCacheByPath(redisClient, "GET", "/api/v1/private-jobs")
				case strings.Contains(path, "/courses"):
					InvalidateCacheByPath(redisClient, "GET", "/api/v1/courses")
				case strings.Contains(path, "/videos"):
					InvalidateCacheByPath(redisClient, "GET", "/api/v1/videos")
				}
			}
		} else {
			c.Next()
		}
	}
}

// bodyLogWriter is a custom response writer to capture response body
type bodyLogWriter struct {
	gin.ResponseWriter
	body []byte
}

func (w *bodyLogWriter) Write(b []byte) (int, error) {
	w.body = append(w.body, b...)
	return w.ResponseWriter.Write(b)
}
