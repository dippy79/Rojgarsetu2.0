package middleware

import (
	"context"
	"log"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"github.com/rojgarsetu/backend/internal/db"
)

func AnalyticsMiddleware(redisClient *redis.Client, database *db.PostgresDB) gin.HandlerFunc {
	// Start background ticker to flush analytics to Postgres every 10 minutes
	go func() {
		ticker := time.NewTicker(10 * time.Minute)
		for range ticker.C {
			flushAnalytics(context.Background(), redisClient, database)
		}
	}()

	return func(c *gin.Context) {
		path := c.Request.URL.Path
		method := c.Request.Method

		if redisClient != nil {
			ctx := context.Background()
			// Visit count
			if strings.HasPrefix(path, "/api/v1") && method == "GET" {
				redisClient.Incr(ctx, "visits:daily")
			}

			// Application count
			if strings.HasSuffix(path, "/apply") && method == "POST" {
				redisClient.Incr(ctx, "applications:daily")
			}
		}

		c.Next()

		// Placement count (if status changed to selected)
		if strings.Contains(path, "/status") && method == "PATCH" && c.Writer.Status() == 200 {
			// Logic to check if status was set to 'selected'
			// This might require reading response body or checking context
		}
	}
}

func flushAnalytics(ctx context.Context, redisClient *redis.Client, database *db.PostgresDB) {
	if redisClient == nil {
		return
	}

	visits, _ := redisClient.Get(ctx, "visits:daily").Int64()
	apps, _ := redisClient.Get(ctx, "applications:daily").Int64()

	stats, err := database.Queries.GetPlatformStats(ctx)
	if err != nil {
		log.Printf("Failed to get platform stats: %v", err)
		return
	}

	_, err = database.Queries.UpdatePlatformStats(ctx, db.UpdatePlatformStatsParams{
		TotalJobs:         stats.TotalJobs,
		TotalCandidates:   stats.TotalCandidates,
		TotalCompanies:    stats.TotalCompanies,
		TotalPlacements:   stats.TotalPlacements,
		TotalApplications: stats.TotalApplications + apps,
		VisitsToday:       visits,
	})

	if err == nil {
		// Reset daily counters after flush if it's a new day
		// For simplicity, we just keep incrementing and reset manually if needed
	}
}
