package main

import (
	"context"
	"crypto/tls"
	"database/sql"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"sync"
	"syscall"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
	"github.com/redis/go-redis/v9"
	"github.com/rojgarsetu/backend/config"
	_ "github.com/rojgarsetu/backend/docs" // swagger docs
	"github.com/rojgarsetu/backend/internal/crawler"
	"github.com/rojgarsetu/backend/internal/db"
	"github.com/rojgarsetu/backend/internal/handlers"
	"github.com/rojgarsetu/backend/internal/middleware"
	"github.com/rojgarsetu/backend/internal/services"
	"github.com/rojgarsetu/backend/internal/workers"
	"github.com/rs/zerolog"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

var (
	version     = "2.0.0"
	logger      zerolog.Logger
	dbURL       string
	serverPort  int
	redisURL    string
	redisClient *redis.Client
)

func main() {
	// Initialize logger
	logger = zerolog.New(os.Stderr).Output(zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: time.RFC3339}).With().Timestamp().Caller().Logger()

	// Load config
	cfg := config.Load()

	// Root command
	dbURL = os.Getenv("DATABASE_URL")
	if dbURL == "" {
		logger.Fatal().Msg("DATABASE_URL environment variable required")
		return
	}

	if err := runMigrations(dbURL); err != nil {
		logger.Fatal().Err(err).Msg("Database migration failed")
		return
	}

	redisURL = os.Getenv("REDIS_URL")
	if redisURL == "" {
		redisURL = "redis://localhost:6379"
	}

	// Initialize Redis client
	opt, err := redis.ParseURL(redisURL)
	if err != nil {
		logger.Fatal().Err(err).Msg("Failed to parse Redis URL")
		return
	}
	redisClient = redis.NewClient(opt)

	// Test Redis connection
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := redisClient.Ping(ctx).Err(); err != nil {
		logger.Warn().Err(err).Msg("Redis connection failed, caching will be disabled")
		redisClient = nil
	} else {
		logger.Info().Msg("Redis connected successfully")
	}

	serverPort = 8083
	if port := os.Getenv("PORT"); port != "" {
		if p, err := strconv.Atoi(port); err == nil {
			serverPort = p
		}
	}
	// Start server
	if err := run(cfg); err != nil {
		logger.Fatal().Err(err).Msg("Server failed")
	}
}

func runMigrations(dbURL string) error {
	migrationsPath := os.Getenv("MIGRATIONS_PATH")
	if migrationsPath == "" {
		migrationsPath = "file://./migrations"
	}
	logger.Info().Str("migrations_path", migrationsPath).Msg("Running database migrations")

	m, err := migrate.New(migrationsPath, dbURL)
	if err != nil {
		return fmt.Errorf("failed to initialize migrations: %w", err)
	}
	defer func() {
		if _, err := m.Close(); err != nil {
			logger.Warn().Err(err).Msg("failed to close migration instance")
		}
	}()

	if err := m.Up(); err != nil {
		if err == migrate.ErrNoChange {
			logger.Info().Msg("Database migrations already up to date")
			return nil
		}
		return fmt.Errorf("migration up failed: %w", err)
	}

	logger.Info().Msg("Database migrations applied successfully")
	return nil
}

func run(cfg *config.Config) error {
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

	// Track DB readiness with thread-safe flag
	var (
		mu               sync.RWMutex
		dbReady          bool
		database         *db.PostgresDB
		sqlDB            *sql.DB
		crawlerScheduler *crawler.Scheduler
	)

	middleware.RegisterBuildInfo()

	// Setup Gin router
	gin.SetMode(gin.ReleaseMode)
	router := gin.New()
	router.Use(gin.Recovery())
	router.Use(gin.Logger())

	// Security middleware first
	router.Use(middleware.SecurityHeaders())
	router.Use(middleware.BodyLimit())
	router.Use(middleware.SanitizeInput())

	// Rate limiting
	router.Use(middleware.RateLimitMiddleware(cfg.RateLimit))

	// Cache invalidation for write operations (if Redis is available)
	if redisClient != nil {
		router.Use(middleware.CacheInvalidationMiddleware(redisClient))
	}

	// CORS Configuration (Dynamic via Environment Variables)
	origins := os.Getenv("CORS_ORIGINS")
	if origins == "" {
		origins = "http://localhost:3000,http://localhost:3001,http://127.0.0.1:3000"
	}

	corsMiddleware := cors.New(cors.Config{
		AllowOrigins:     strings.Split(origins, ","),
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	})
	router.Use(corsMiddleware)
	router.Use(middleware.PrometheusMiddleware())

	// Metrics endpoint
	router.GET("/metrics", middleware.MetricsHandler())

	// Swagger documentation
	router.GET("/docs/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Health check - responds immediately with 200 even before DB is ready
	// Handles both GET and HEAD (the latter is used by Docker wget --spider healthcheck)
	healthHandler := func(c *gin.Context) {
		mu.RLock()
		ready := dbReady
		mu.RUnlock()
		status := "healthy"
		if !ready {
			status = "starting"
		}
		c.JSON(http.StatusOK, gin.H{
			"status":    status,
			"service":   "backend-api",
			"version":   version,
			"db_ready":  ready,
			"timestamp": time.Now().Format(time.RFC3339),
		})
	}
	router.GET("/health", healthHandler)
	router.HEAD("/health", healthHandler)

	// Attempt database connection in background, then register DB routes
	go func() {
		maxRetries := 30
		for i := 0; i < maxRetries; i++ {
			select {
			case <-ctx.Done():
				return
			default:
			}

			sdb, err := sql.Open("postgres", dbURL)
			if err != nil {
				logger.Warn().Err(err).Msgf("Failed to open database (attempt %d/%d)", i+1, maxRetries)
				time.Sleep(2 * time.Second)
				continue
			}

			// Configure connection pool
			sdb.SetMaxOpenConns(cfg.Database.MaxOpenConns)
			sdb.SetMaxIdleConns(cfg.Database.MaxIdleConns)
			sdb.SetConnMaxLifetime(cfg.Database.ConnMaxLifetime)
			sdb.SetConnMaxIdleTime(cfg.Database.ConnMaxIdleTime)

			logger.Info().
				Int("max_open_conns", cfg.Database.MaxOpenConns).
				Int("max_idle_conns", cfg.Database.MaxIdleConns).
				Dur("conn_max_lifetime", cfg.Database.ConnMaxLifetime).
				Dur("conn_max_idle_time", cfg.Database.ConnMaxIdleTime).
				Msg("Database connection pool configured")

			pingCtx, pingCancel := context.WithTimeout(context.Background(), 5*time.Second)
			err = sdb.PingContext(pingCtx)
			pingCancel()
			if err == nil {
				mu.Lock()
				sqlDB = sdb
				database = db.NewPostgresDB(sdb)
				dbReady = true
				mu.Unlock()
				logger.Info().Msg("Database connected")

				// Initialize services
				govJobService := services.NewGovJobService(database)
				privJobService := services.NewPrivJobService(database)
				courseService := services.NewCourseService(database)
				videoService := services.NewVideoService(database)
				searchService := services.NewSearchService(database)
				notificationService := services.NewNotificationService(database)
				interviewService := services.NewInterviewService(database)
				uploadService := services.NewUploadService(database)

				// Initialize new feature repo and handlers
				featureRepo := db.NewFeatureRepository(sdb)

				// Initialize handlers
				govJobHandler := handlers.NewGovJobHandler(govJobService)
				privJobHandler := handlers.NewPrivJobHandler(privJobService)
				courseHandler := handlers.NewCourseHandler(courseService)
				videoHandler := handlers.NewVideoHandler(videoService)
				searchHandler := handlers.NewSearchHandler(searchService)
				featureHandler := handlers.NewFeatureHandler(featureRepo)
				crawlerHandler := handlers.NewCrawlerHandler(sdb)
				legalHandler := handlers.NewLegalHandler(sdb)
				wsHandler := handlers.NewWSHandler(notificationService)
				interviewHandler := handlers.NewInterviewHandler(interviewService)
				statsHandler := handlers.NewStatsHandler(database)
				uploadHandler := handlers.NewUploadHandler(uploadService)

				// Start workers
				emailWorker := workers.NewEmailWorker(database)
				go emailWorker.Start(ctx)

				// Setup router with Analytics
				router.Use(middleware.AnalyticsMiddleware(redisClient, database))

				// Register DB-dependent routes
				api := router.Group("/api/v1")
				{
					// Apply caching with different TTLs based on endpoint
					if redisClient != nil {
						// 5 minute cache for jobs
						api.GET("/gov-jobs", middleware.CacheMiddleware(redisClient, 5*time.Minute), govJobHandler.GetGovJobs)
						api.GET("/private-jobs", middleware.CacheMiddleware(redisClient, 5*time.Minute), privJobHandler.GetPrivJobs)

						// 10 minute cache for courses and videos
						api.GET("/courses", middleware.CacheMiddleware(redisClient, 10*time.Minute), courseHandler.GetCourses)
						api.GET("/courses/providers", middleware.CacheMiddleware(redisClient, 10*time.Minute), courseHandler.GetCourseProviders)
						api.GET("/videos", middleware.CacheMiddleware(redisClient, 10*time.Minute), videoHandler.GetVideos)
						api.GET("/videos/channels", middleware.CacheMiddleware(redisClient, 10*time.Minute), videoHandler.GetVideoChannels)
						api.GET("/videos/categories", middleware.CacheMiddleware(redisClient, 10*time.Minute), videoHandler.GetVideoCategories)

						// 1 minute cache for stats
						api.GET("/crawler/stats", middleware.CacheMiddleware(redisClient, 1*time.Minute), crawlerHandler.GetStats)
					} else {
						// No caching if Redis is unavailable
						api.GET("/gov-jobs", govJobHandler.GetGovJobs)
						api.GET("/private-jobs", privJobHandler.GetPrivJobs)
						api.GET("/courses", courseHandler.GetCourses)
						api.GET("/courses/providers", courseHandler.GetCourseProviders)
						api.GET("/videos", videoHandler.GetVideos)
						api.GET("/videos/channels", videoHandler.GetVideoChannels)
						api.GET("/videos/categories", videoHandler.GetVideoCategories)
						api.GET("/crawler/stats", crawlerHandler.GetStats)
					}

					// Uncached endpoints
					api.GET("/gov-jobs/:id", govJobHandler.GetGovJobByID)
					api.GET("/private-jobs/:id", privJobHandler.GetPrivJobByID)
					api.GET("/courses/:id", courseHandler.GetCourseByID)
					api.GET("/videos/:id", videoHandler.GetVideoByID)
					api.POST("/search", searchHandler.Search)
					api.GET("/search", searchHandler.SearchGET)

					// Crawler Endpoints
					api.POST("/crawler/crawl", crawlerHandler.TriggerCrawl)
					api.GET("/crawler/health", crawlerHandler.GetHealth)
					// Legal Endpoints
					api.GET("/legal/disclaimer", legalHandler.GetDisclaimer)
					api.POST("/legal/takedown", legalHandler.PostTakedown)
					api.GET("/crawler/forms", legalHandler.GetForms)
					// New Feature Endpoints
					api.POST("/company/reviews", featureHandler.CreateReviewHandler)
					api.POST("/jobs/report", featureHandler.ReportJobHandler)
					api.POST("/candidate/ratings", featureHandler.InternalRatingHandler)

					// WebSocket endpoint for real-time notifications
					api.GET("/ws", wsHandler.HandleWebSocket)

					// Interview Endpoints
					api.POST("/interviews", interviewHandler.ScheduleInterview)
					api.GET("/interviews/:id", interviewHandler.GetInterviewByID)

					// Public Stats
					router.GET("/api/v1/stats", statsHandler.GetPlatformStats)

					// Upload Endpoints
					api.POST("/upload", uploadHandler.UploadFile)
				}
				logger.Info().Msg("API routes registered")

				// Start Crawler Scheduler (Runs every 6 hours in background)
				crawlerScheduler = crawler.NewScheduler(sdb, 6*time.Hour)
				crawlerScheduler.Start()
				logger.Info().Msg("Crawler scheduler started with 6-hour interval")
				return
			}
			sdb.Close()
			logger.Warn().Err(err).Msgf("Database ping failed (attempt %d/%d) - retrying in 2s", i+1, maxRetries)
			time.Sleep(2 * time.Second)
		}
		logger.Error().Msg("Failed to connect to database after all retries")
	}()

	// Start HTTP server (independent of DB readiness)
	srv := &http.Server{
		Addr:         fmt.Sprintf(":%d", serverPort),
		Handler:      router,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 30 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	if cfg.TLS.Enabled {
		cert, err := tls.LoadX509KeyPair(cfg.TLS.Cert, cfg.TLS.Key)
		if err != nil {
			return fmt.Errorf("failed to load TLS cert: %w", err)
		}
		srv.TLSConfig = &tls.Config{
			Certificates: []tls.Certificate{cert},
		}
		logger.Info().Msgf("Starting HTTPS server on :8443 with TLS cert %s", cfg.TLS.Cert)
		go func() {
			if err := srv.ListenAndServeTLS("", ""); err != nil && err != http.ErrServerClosed {
				logger.Error().Err(err).Msg("HTTPS server error")
			}
		}()
	} else {
		logger.Info().Int("port", serverPort).Msg("Starting HTTP server")
		go func() {
			if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
				logger.Error().Err(err).Msg("HTTP server error")
			}
		}()
	}

	logger.Info().Msg("Backend API service ready")

	// Wait for shutdown
	<-ctx.Done()
	logger.Info().Msg("Shutting down server...")

	// Stop crawler scheduler if running
	mu.RLock()
	if crawlerScheduler != nil {
		logger.Info().Msg("Stopping crawler scheduler...")
		crawlerScheduler.Stop()
	}
	mu.RUnlock()

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer shutdownCancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		return fmt.Errorf("server shutdown failed: %w", err)
	}

	// Close DB if connected
	mu.RLock()
	if sqlDB != nil {
		sqlDB.Close()
	}
	mu.RUnlock()

	// Close Redis connection if connected
	if redisClient != nil {
		redisClient.Close()
	}

	logger.Info().Msg("Server shutdown complete")
	return nil
}
