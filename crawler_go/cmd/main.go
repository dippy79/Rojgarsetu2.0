package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/rojgarsetu/crawler/internal/browser"
	"github.com/rojgarsetu/crawler/internal/parser"
	"github.com/rojgarsetu/crawler/internal/proxy"
	"github.com/rojgarsetu/crawler/internal/sources"
	"github.com/rojgarsetu/crawler/internal/store"
	"github.com/rs/zerolog"

	"github.com/spf13/cobra"
)

var (
	version       = "2.0.0"
	logger        zerolog.Logger
	dbURL         string
	redisURL      string
	workers       int
	timeout       int
	retryCnt      int
	proxyEnabled  bool
	serverPort    int
	globalDB      *store.PostgresStore
	globalBrowser *browser.Pool
)

func main() {
	// Initialize logger
	output := zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: time.RFC3339}
	logger = logger.NewLogger()

	// Root command
	rootCmd := &cobra.Command{
		Use:   "crawler",
		Short: "RojgarSetu Job Crawler v2.0",
		Long:  "Scalable job crawler with browser automation and retry logic",
		Run:   runCrawler,
	}

	// Flags
	rootCmd.Flags().StringVar(&dbURL, "db-url", os.Getenv("DATABASE_URL"), "PostgreSQL connection URL")
	rootCmd.Flags().StringVar(&redisURL, "redis-url", os.Getenv("REDIS_URL"), "Redis connection URL")
	rootCmd.Flags().IntVarP(&workers, "workers", "w", 5, "Number of concurrent workers")
	rootCmd.Flags().IntVarP(&timeout, "timeout", "t", 30, "Request timeout in seconds")
	rootCmd.Flags().IntVarP(&retryCnt, "retry", "r", 3, "Number of retry attempts")
	rootCmd.Flags().BoolVar(&proxyEnabled, "proxy", false, "Enable proxy rotation")
	rootCmd.Flags().IntVarP(&serverPort, "port", "p", 8082, "HTTP server port")

	if err := rootCmd.Execute(); err != nil {
		logger.Fatal().Err(err).Msg("Failed to execute command")
	}
}

// Health check handler
func healthHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"status":"healthy","service":"crawler","version":"2.0.0"}`))
}

// Metrics handler (simple implementation)
func metricsHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("# HELP crawler_jobs_processed Total jobs processed\n"))
	w.Write([]byte("# TYPE crawler_jobs_processed counter\n"))
	w.Write([]byte("crawler_jobs_processed 0\n"))
}

// Stats handler - returns job count
func statsHandler(w http.ResponseWriter, r *http.Request) {
	if globalDB == nil {
		http.Error(w, "Database not initialized", http.StatusInternalServerError)
		return
	}

	count, err := globalDB.GetJobCount()
	if err != nil {
		logger.Error().Err(err).Msg("Failed to get job count")
		http.Error(w, fmt.Sprintf("Failed to get job count: %v", err), http.StatusInternalServerError)
		return
	}

	logger.Info().Int("totalJobsInDatabase", count).Msg("Total jobs currently in database")

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status":    "success",
		"totalJobs": count,
		"timestamp": time.Now().Format(time.RFC3339),
	})
}

// Crawl handler - manual trigger
func crawlHandler(w http.ResponseWriter, r *http.Request) {
	logger.Info().Msg("Manual crawl triggered via HTTP")

	ctx := context.Background()

	if globalDB == nil {
		http.Error(w, "Database not initialized", http.StatusInternalServerError)
		return
	}

	// Initialize proxy rotator
	var proxyRotator *proxy.Rotator

	// Create sources
	naukriSource := sources.NewNaukriSource(globalBrowser, proxyRotator)

	// First run browser test
	if globalBrowser != nil {
		logger.Info().Msg("Running browser test...")
		if err := globalBrowser.RunBrowserTest(); err != nil {
			logger.Warn().Err(err).Msg("Browser test had issues but continuing with crawl")
		}
	}

	// Fetch jobs
	jobs, err := naukriSource.Fetch(ctx)
	if err != nil {
		logger.Error().Err(err).Msg("Failed to fetch jobs")
		http.Error(w, fmt.Sprintf("Failed to fetch jobs: %v", err), http.StatusInternalServerError)
		return
	}

	logger.Info().Int("jobsFound", len(jobs)).Msg("Jobs fetched from source")
	logger.Info().Int("totalJobsExtracted", len(jobs)).Msg("Total jobs extracted")

	// Save jobs to database
	logger.Info().Msg("Inserting jobs into database")
	successCount := 0
	for _, job := range jobs {
		parsedJob := parser.ParseJob(&job)
		if err := parser.ValidateJob(parsedJob); err != nil {
			logger.Warn().Err(err).Msg("Invalid job data, skipping")
			continue
		}
		if err := globalDB.SaveJob(parsedJob); err != nil {
			logger.Error().Err(err).Str("job", job.Title).Msg("Failed to save job")
		} else {
			successCount++
		}
	}
	logger.Info().Int("saved", successCount).Msg("Jobs inserted successfully")

	// Get job count after insertion
	jobCountAfter, err := globalDB.GetJobCount()
	if err != nil {
		logger.Warn().Err(err).Msg("Failed to get job count after insertion")
	} else {
		logger.Info().Int("totalJobsCurrentlyInDatabase", jobCountAfter).Msg("Total jobs currently in database")
	}

	// Return JSON response
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status":          "success",
		"jobsFound":       len(jobs),
		"jobsSaved":       successCount,
		"jobsInDatabase":  jobCountAfter,
		"jobsInsertedNow": successCount,
		"source":          "naukri",
		"timestamp":       time.Now().Format(time.RFC3339),
	})
}

func runCrawler(cmd *cobra.Command, args []string) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Setup signal handling
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-sigCh
		logger.Info().Msg("Received shutdown signal")
		cancel()
	}()

	// Start HTTP server for health checks and crawl endpoint
	go func() {
		http.HandleFunc("/health", healthHandler)
		http.HandleFunc("/metrics", metricsHandler)
		http.HandleFunc("/crawl", crawlHandler)
		http.HandleFunc("/stats", statsHandler)
		logger.Info().Int("port", serverPort).Msg("Starting HTTP server")
		if err := http.ListenAndServe(fmt.Sprintf(":%d", serverPort), nil); err != nil && err != http.ErrServerClosed {
			logger.Error().Err(err).Msg("HTTP server error")
		}
	}()

	logger.Info().Int("workers", workers).Msg("Starting RojgarSetu Crawler v2.0")

	// Initialize database
	db, err := store.NewPostgresStore(dbURL)
	if err != nil {
		logger.Fatal().Err(err).Msg("Failed to connect to database")
	}
	globalDB = db
	defer db.Close()
	logger.Info().Msg("Database connected")

	// Initialize Redis
	redisOpts, err := redis.ParseURL(redisURL)
	if err != nil {
		logger.Fatal().Err(err).Msg("Failed to parse Redis URL")
	}
	redisClient := redis.NewClient(redisOpts)
	defer redisClient.Close()

	if err := redisClient.Ping(ctx).Err(); err != nil {
		logger.Fatal().Err(err).Msg("Failed to connect to Redis")
	}
	logger.Info().Msg("Redis connected")

	// Initialize browser pool (if needed)
	var browserPool *browser.Pool
	if browserPool, err = browser.NewPool(workers); err != nil {
		logger.Warn().Err(err).Msg("Failed to initialize browser pool, using HTTP only")
	} else {
		logger.Info().Msg("Browser pool initialized")
	}
	globalBrowser = browserPool

	// Note: Job crawling is triggered via the /crawl HTTP endpoint
	// The crawler service is now ready to accept crawl requests

	logger.Info().Msg("Crawler service ready - waiting for crawl requests via /crawl endpoint")

	// Wait for shutdown signal
	<-ctx.Done()
	logger.Info().Msg("Crawler shutdown complete")
}
