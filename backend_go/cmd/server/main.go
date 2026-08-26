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
	"sync/atomic"
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

// AppHandlers holds references to all route handlers
type AppHandlers struct {
	GovJobHandler    *handlers.GovJobHandler
	PrivJobHandler   *handlers.PrivJobHandler
	CourseHandler    *handlers.CourseHandler
	VideoHandler     *handlers.VideoHandler
	SearchHandler    *handlers.SearchHandler
	FeatureHandler   *handlers.FeatureHandler
	CrawlerHandler   *handlers.CrawlerHandler
	LegalHandler     *handlers.LegalHandler
	WSHandler        *handlers.WSHandler
	InterviewHandler *handlers.InterviewHandler
	StatsHandler     *handlers.StatsHandler
	UploadHandler    *handlers.UploadHandler
}

// Global atomic variables for thread-safe state management
var (
	dbReady     atomic.Bool
	appHandlers atomic.Pointer[AppHandlers]
)

// dbReadinessMiddleware blocks route execution with 503 until DB connects
func dbReadinessMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Always allow health checks and metrics to pass through
		if c.Request.URL.Path == "/health" || c.Request.URL.Path == "/metrics" {
			c.Next()
			return
		}

		// If DB is not ready, return 503 Service Unavailable
		if !dbReady.Load() || appHandlers.Load() == nil {
			c.Header("Retry-After", "3")
			c.JSON(http.StatusServiceUnavailable, gin.H{
				"error":   "Database connection initializing",
				"message": "Service warming up, please retry shortly",
			})
			c.Abort()
			return
		}
		c.Next()
	}
}

// safeHandler wraps handler execution to prevent nil pointer panics
func safeHandler(fn func(h *AppHandlers, c *gin.Context)) gin.HandlerFunc {
	return func(c *gin.Context) {
		h := appHandlers.Load()
		if h == nil {
			c.JSON(http.StatusServiceUnavailable, gin.H{"error": "Handlers unavailable"})
			return
		}
		fn(h, c)
	}
}

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

	// Track DB connection state
	var (
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

	// -------------------------------------------------------------
	// 1. CORS MIDDLEWARE (SABSE UPAR RAKHA HAI)
	// -------------------------------------------------------------
	origins := os.Getenv("CORS_ORIGINS")
	if origins == "" {
		origins = "http://localhost:8080,http://localhost:3000,http://localhost:3001,http://127.0.0.1:3000,http://127.0.0.1:8080"
	}

	corsMiddleware := cors.New(cors.Config{
		AllowOrigins:     strings.Split(origins, ","),
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization", "X-Requested-With"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	})
	router.Use(corsMiddleware)

	// Direct 204 response for all OPTIONS preflight requests
	router.OPTIONS("/*path", func(c *gin.Context) {
		c.Status(http.StatusNoContent)
	})

	// -------------------------------------------------------------
	// 2. DB READINESS MIDDLEWARE (Protects routes from DB race conditions)
	// -------------------------------------------------------------
	router.Use(dbReadinessMiddleware())

	// -------------------------------------------------------------
	// 3. SECURITY & RATE LIMITING MIDDLEWARES (CORS KE BAAD)
	// -------------------------------------------------------------
	router.Use(middleware.SecurityHeaders())
	router.Use(middleware.BodyLimit())
	router.Use(middleware.SanitizeInput())
	router.Use(middleware.RateLimitMiddleware(cfg.RateLimit))

	// Cache invalidation for write operations (if Redis is available)
	if redisClient != nil {
		router.Use(middleware.CacheInvalidationMiddleware(redisClient))
	}

	router.Use(middleware.PrometheusMiddleware())

	// Metrics endpoint
	router.GET("/metrics", middleware.MetricsHandler())

	// Swagger documentation
	router.GET("/docs/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Health check endpoint
	healthHandler := func(c *gin.Context) {
		status := "UP"
		if !dbReady.Load() {
			status = "INITIALIZING_DB"
		}
		c.JSON(http.StatusOK, gin.H{
			"status":    status,
			"service":   "backend-api",
			"version":   version,
			"db_ready":  dbReady.Load(),
			"timestamp": time.Now().Format(time.RFC3339),
		})
	}
	router.GET("/health", healthHandler)
	router.HEAD("/health", healthHandler)

	// -------------------------------------------------------------
	// 4. SYNCHRONOUS ROUTE REGISTRATION (Routes exist immediately on startup)
	// -------------------------------------------------------------
	api := router.Group("/api/v1")
	{
		// Apply caching with different TTLs based on endpoint
		if redisClient != nil {
			// 5 minute cache for jobs
			api.GET("/gov-jobs", middleware.CacheMiddleware(redisClient, 5*time.Minute),
				safeHandler(func(h *AppHandlers, c *gin.Context) { h.GovJobHandler.GetGovJobs(c) }))
			api.GET("/private-jobs", middleware.CacheMiddleware(redisClient, 5*time.Minute),
				safeHandler(func(h *AppHandlers, c *gin.Context) { h.PrivJobHandler.GetPrivJobs(c) }))
			api.GET("/priv-jobs", middleware.CacheMiddleware(redisClient, 5*time.Minute),
				safeHandler(func(h *AppHandlers, c *gin.Context) { h.PrivJobHandler.GetPrivJobs(c) }))

			// 10 minute cache for courses and videos
			api.GET("/courses", middleware.CacheMiddleware(redisClient, 10*time.Minute),
				safeHandler(func(h *AppHandlers, c *gin.Context) { h.CourseHandler.GetCourses(c) }))
			api.GET("/courses/providers", middleware.CacheMiddleware(redisClient, 10*time.Minute),
				safeHandler(func(h *AppHandlers, c *gin.Context) { h.CourseHandler.GetCourseProviders(c) }))
			api.GET("/videos", middleware.CacheMiddleware(redisClient, 10*time.Minute),
				safeHandler(func(h *AppHandlers, c *gin.Context) { h.VideoHandler.GetVideos(c) }))
			api.GET("/videos/channels", middleware.CacheMiddleware(redisClient, 10*time.Minute),
				safeHandler(func(h *AppHandlers, c *gin.Context) { h.VideoHandler.GetVideoChannels(c) }))
			api.GET("/videos/categories", middleware.CacheMiddleware(redisClient, 10*time.Minute),
				safeHandler(func(h *AppHandlers, c *gin.Context) { h.VideoHandler.GetVideoCategories(c) }))

			// 1 minute cache for stats
			api.GET("/crawler/stats", middleware.CacheMiddleware(redisClient, 1*time.Minute),
				safeHandler(func(h *AppHandlers, c *gin.Context) { h.CrawlerHandler.GetStats(c) }))
			api.GET("/forms", middleware.CacheMiddleware(redisClient, 10*time.Minute),
				safeHandler(func(h *AppHandlers, c *gin.Context) { h.LegalHandler.GetForms(c) }))
		} else {
			// No caching if Redis is unavailable
			api.GET("/gov-jobs", safeHandler(func(h *AppHandlers, c *gin.Context) { h.GovJobHandler.GetGovJobs(c) }))
			api.GET("/private-jobs", safeHandler(func(h *AppHandlers, c *gin.Context) { h.PrivJobHandler.GetPrivJobs(c) }))
			api.GET("/priv-jobs", safeHandler(func(h *AppHandlers, c *gin.Context) { h.PrivJobHandler.GetPrivJobs(c) }))
			api.GET("/courses", safeHandler(func(h *AppHandlers, c *gin.Context) { h.CourseHandler.GetCourses(c) }))
			api.GET("/courses/providers", safeHandler(func(h *AppHandlers, c *gin.Context) { h.CourseHandler.GetCourseProviders(c) }))
			api.GET("/videos", safeHandler(func(h *AppHandlers, c *gin.Context) { h.VideoHandler.GetVideos(c) }))
			api.GET("/videos/channels", safeHandler(func(h *AppHandlers, c *gin.Context) { h.VideoHandler.GetVideoChannels(c) }))
			api.GET("/videos/categories", safeHandler(func(h *AppHandlers, c *gin.Context) { h.VideoHandler.GetVideoCategories(c) }))
			api.GET("/crawler/stats", safeHandler(func(h *AppHandlers, c *gin.Context) { h.CrawlerHandler.GetStats(c) }))
		}

		// Uncached endpoints
		api.GET("/gov-jobs/:id", safeHandler(func(h *AppHandlers, c *gin.Context) { h.GovJobHandler.GetGovJobByID(c) }))
		api.GET("/private-jobs/:id", safeHandler(func(h *AppHandlers, c *gin.Context) { h.PrivJobHandler.GetPrivJobByID(c) }))
		api.GET("/priv-jobs/:id", safeHandler(func(h *AppHandlers, c *gin.Context) { h.PrivJobHandler.GetPrivJobByID(c) }))
		api.GET("/courses/:id", safeHandler(func(h *AppHandlers, c *gin.Context) { h.CourseHandler.GetCourseByID(c) }))
		api.GET("/videos/:id", safeHandler(func(h *AppHandlers, c *gin.Context) { h.VideoHandler.GetVideoByID(c) }))
		api.POST("/search", safeHandler(func(h *AppHandlers, c *gin.Context) { h.SearchHandler.Search(c) }))
		api.GET("/search", safeHandler(func(h *AppHandlers, c *gin.Context) { h.SearchHandler.SearchGET(c) }))

		// Public Stats Endpoint
		api.GET("/stats", safeHandler(func(h *AppHandlers, c *gin.Context) { h.StatsHandler.GetPlatformStats(c) }))

		// Crawler Endpoints
		api.POST("/crawler/crawl", safeHandler(func(h *AppHandlers, c *gin.Context) { h.CrawlerHandler.TriggerCrawl(c) }))
		api.GET("/crawler/health", safeHandler(func(h *AppHandlers, c *gin.Context) { h.CrawlerHandler.GetHealth(c) }))

		// Legal Endpoints
		api.GET("/legal/disclaimer", safeHandler(func(h *AppHandlers, c *gin.Context) { h.LegalHandler.GetDisclaimer(c) }))
		api.POST("/legal/takedown", safeHandler(func(h *AppHandlers, c *gin.Context) { h.LegalHandler.PostTakedown(c) }))
		api.GET("/crawler/forms", safeHandler(func(h *AppHandlers, c *gin.Context) { h.LegalHandler.GetForms(c) }))
		api.GET("/forms", safeHandler(func(h *AppHandlers, c *gin.Context) { h.LegalHandler.GetForms(c) }))

		// New Feature Endpoints
		api.POST("/company/reviews", safeHandler(func(h *AppHandlers, c *gin.Context) { h.FeatureHandler.CreateReviewHandler(c) }))
		api.POST("/jobs/report", safeHandler(func(h *AppHandlers, c *gin.Context) { h.FeatureHandler.ReportJobHandler(c) }))
		api.POST("/candidate/ratings", safeHandler(func(h *AppHandlers, c *gin.Context) { h.FeatureHandler.InternalRatingHandler(c) }))

		// WebSocket endpoint for real-time notifications
		api.GET("/ws", safeHandler(func(h *AppHandlers, c *gin.Context) { h.WSHandler.HandleWebSocket(c) }))

		// Interview Endpoints
		api.POST("/interviews", safeHandler(func(h *AppHandlers, c *gin.Context) { h.InterviewHandler.ScheduleInterview(c) }))
		api.GET("/interviews/:id", safeHandler(func(h *AppHandlers, c *gin.Context) { h.InterviewHandler.GetInterviewByID(c) }))

		// Upload Endpoints
		api.POST("/upload", safeHandler(func(h *AppHandlers, c *gin.Context) { h.UploadHandler.UploadFile(c) }))
	}

	// ASYNCHRONOUS DB CONNECTION & ATOMIC HANDLER SWAPPING
	go func() {
		logger.Info().Msg("[INIT] Connecting to Database in background...")

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
				sqlDB = sdb
				database = db.NewPostgresDB(sdb)
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

				// Initialize new feature repo
				featureRepo := db.NewFeatureRepository(sdb)

				logger.Info().Msg("[INIT] Database connected. Initializing route handlers...")

				// Atomically store initialized handlers
				initialized := &AppHandlers{
					GovJobHandler:    handlers.NewGovJobHandler(govJobService),
					PrivJobHandler:   handlers.NewPrivJobHandler(privJobService),
					CourseHandler:    handlers.NewCourseHandler(courseService),
					VideoHandler:     handlers.NewVideoHandler(videoService),
					SearchHandler:    handlers.NewSearchHandler(searchService),
					FeatureHandler:   handlers.NewFeatureHandler(featureRepo),
					CrawlerHandler:   handlers.NewCrawlerHandler(sdb),
					LegalHandler:     handlers.NewLegalHandler(sdb),
					WSHandler:        handlers.NewWSHandler(notificationService),
					InterviewHandler: handlers.NewInterviewHandler(interviewService),
					StatsHandler:     handlers.NewStatsHandler(database),
					UploadHandler:    handlers.NewUploadHandler(uploadService),
				}

				// Atomically swap the pointer and mark DB as ready
				appHandlers.Store(initialized)
				dbReady.Store(true)

				// Setup router with Analytics after DB ready
				router.Use(middleware.AnalyticsMiddleware(redisClient, database))

				logger.Info().Msg("[SUCCESS] All handlers initialized. Backend is fully operational.")

				// Start workers
				emailWorker := workers.NewEmailWorker(database)
				go emailWorker.Start(ctx)

				// Start Crawler Scheduler
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

	// Start HTTP server
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
	if crawlerScheduler != nil {
		logger.Info().Msg("Stopping crawler scheduler...")
		crawlerScheduler.Stop()
	}

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer shutdownCancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		return fmt.Errorf("server shutdown failed: %w", err)
	}

	// Close DB if connected
	if sqlDB != nil {
		sqlDB.Close()
	}

	// Close Redis connection if connected
	if redisClient != nil {
		redisClient.Close()
	}

	logger.Info().Msg("Server shutdown complete")
	return nil
}
