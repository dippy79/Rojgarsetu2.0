package middleware

import (
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

// LoggingMiddleware logs requests in structured JSON format with request ID
func LoggingMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Generate request ID
		requestID := c.GetHeader("X-Request-ID")
		if requestID == "" {
			requestID = fmt.Sprintf("req-%d", time.Now().UnixNano())
		}
		c.Set("requestID", requestID)

		// Log request start
		log.Info().
			Str("request_id", requestID).
			Str("method", c.Request.Method).
			Str("path", c.Request.URL.Path).
			Str("client_ip", c.ClientIP()).
			Msg("request started")

		// Record start time
		start := time.Now()

		// Process request
		c.Next()

		// Get userID from context if available
		userID, exists := c.Get("user_id")
		var userIDStr string
		if exists {
			userIDStr = userID.(string)
		}

		// Log request complete
		log.Info().
			Str("request_id", requestID).
			Str("user_id", userIDStr).
			Int("status", c.Writer.Status()).
			Dur("duration_ms", time.Since(start)).
			Msg("request completed")

		// Set response header
		c.Header("X-Request-ID", requestID)
	}
}

// JSONLoggerMiddleware replaces gin.Logger with structured JSON logs
func JSONLoggerMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		requestID := fmt.Sprintf("req-%d", time.Now().UnixNano())

		c.Next()

		latency := time.Since(start)

		log.Info().
			Str("request_id", requestID).
			Str("method", c.Request.Method).
			Str("path", path).
			Str("client_ip", c.ClientIP()).
			Int("status", c.Writer.Status()).
			Dur("duration", latency).
			Msg("request processed")

		c.Header("X-Request-ID", requestID)
	}
}
