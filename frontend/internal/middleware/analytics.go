package middleware

import (
    "github.com/gin-gonic/gin"
    "github.com/go-redis/redis/v8"
)

func AnalyticsMiddleware(rdb *redis.Client) gin.HandlerFunc {
    return func(c *gin.Context) {
        ctx := c.Request.Context()
        if c.Request.Method == "GET" {
            rdb.Incr(ctx, "visits:daily")
        } else if c.Request.Method == "POST" && c.FullPath() == "/api/v1/jobs/:id/apply" {
            rdb.Incr(ctx, "applications:daily")
        }
        c.Next()
    }
}
