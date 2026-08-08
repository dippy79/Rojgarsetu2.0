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
	"github.com/rojgarsetu/backend/config"
	"github.com/rojgarsetu/backend/internal/db"
	"github.com/rojgarsetu/backend/internal/handlers"
	"github.com/rojgarsetu/backend/internal/middleware"
	"github.com/rojgarsetu/backend/internal/services"
	"github.com/rs/zerolog"
)

var (
	version    = "2.0.0"
	logger     zerolog.Logger
	dbURL      string
	serverPort int
	redisURL   string
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
		mu       sync.RWMutex
		dbReady  bool
		database *db.PostgresDB
		sqlDB    *sql.DB
	)

	middleware.RegisterBuildInfo()

	// Setup Gin router
	gin.SetMode(gin.ReleaseMode)
	router := gin.New()
	router.Use(gin.Recovery())
	router.Use(gin.Logger())

	// Security middleware first
	router.Use(middleware.SecurityHeaders())

	// Rate limiting
	router.Use(middleware.RateLimitMiddleware(cfg.RateLimit))

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

	api := router.Group("/api/v1")
	api.GET("/forms", func(c *gin.Context) {
		forms := []gin.H{
			{
				"id":         "stub-1",
				"title":      "SSC CGL 2025 Online Application",
				"department": "Staff Selection Commission",
				"last_date":  time.Now().Add(48 * time.Hour).Format(time.RFC3339),
				"apply_url":  "https://ssc.gov.in",
			},
			{
				"id":         "stub-2",
				"title":      "RRB NTPC 2025 Application Form",
				"department": "Railway Recruitment Board",
				"last_date":  time.Now().Add(6 * 24 * time.Hour).Format(time.RFC3339),
				"apply_url":  "https://rrbcdg.gov.in",
			},
			{
				"id":         "stub-3",
				"title":      "UPSC Civil Services Prelims 2025",
				"department": "Union Public Service Commission",
				"last_date":  time.Now().Add(10 * 24 * time.Hour).Format(time.RFC3339),
				"apply_url":  "https://upsc.gov.in",
			},
		}
		c.JSON(http.StatusOK, gin.H{"data": forms, "count": len(forms), "source": "stub"})
	})

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

				// Initialize handlers
				govJobHandler := handlers.NewGovJobHandler(govJobService)
				privJobHandler := handlers.NewPrivJobHandler(privJobService)
				courseHandler := handlers.NewCourseHandler(courseService)
				videoHandler := handlers.NewVideoHandler(videoService)
				searchHandler := handlers.NewSearchHandler(searchService)

				// Register DB-dependent routes
				api := router.Group("/api/v1")
				{
					api.GET("/gov-jobs", govJobHandler.GetGovJobs)
					api.GET("/gov-jobs/:id", govJobHandler.GetGovJobByID)
					api.GET("/private-jobs", privJobHandler.GetPrivJobs)
					api.GET("/private-jobs/:id", privJobHandler.GetPrivJobByID)
					api.GET("/courses", courseHandler.GetCourses)
					api.GET("/courses/providers", courseHandler.GetCourseProviders)
					api.GET("/courses/:id", courseHandler.GetCourseByID)
					api.GET("/videos", videoHandler.GetVideos)
					api.GET("/videos/channels", videoHandler.GetVideoChannels)
					api.GET("/videos/categories", videoHandler.GetVideoCategories)
					api.GET("/videos/:id", videoHandler.GetVideoByID)
					api.POST("/search", searchHandler.Search)
					api.GET("/search", searchHandler.SearchGET)
				}
				logger.Info().Msg("API routes registered")
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

	logger.Info().Msg("Server shutdown complete")
	return nil
}
