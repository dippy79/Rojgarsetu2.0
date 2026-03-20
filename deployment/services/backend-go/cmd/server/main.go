package main

import (
	"context"
	"crypto/tls"
	"fmt"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
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

	// Root command
	dbURL = os.Getenv("DATABASE_URL")
	if dbURL == "" {
		logger.Fatal().Msg("DATABASE_URL environment variable required")
		return
	}

	redisURL = os.Getenv("REDIS_URL")
	if redisURL == "" {
		redisURL = "redis://localhost:6379"
	}

	serverPort = 8083
	if port := os.Getenv("PORT"); port != "" {
		fmt.Sscanf(port, "%d", &serverPort)
	}

	// Start server
	if err := run(); err != nil {
		logger.Fatal().Err(err).Msg("Server failed")
	}
}

func run() error {
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

	// Initialize database
	database, err := db.NewPostgresDB(dbURL)
	if err != nil {
		return fmt.Errorf("failed to connect to database: %w", err)
	}
	defer database.Close()
	logger.Info().Msg("Database connected")

	// Initialize services
	govJobService := services.NewGovJobService(database)
	privJobService := services.NewPrivJobService(database)
	courseService := services.NewCourseService(database)
	videoService := services.NewVideoService(database)

	// Initialize handlers
	govJobHandler := handlers.NewGovJobHandler(govJobService)
	privJobHandler := handlers.NewPrivJobHandler(privJobService)
	courseHandler := handlers.NewCourseHandler(courseService)
	videoHandler := handlers.NewVideoHandler(videoService)

	middleware.RegisterBuildInfo()

	// Setup Gin router
	gin.SetMode(gin.ReleaseMode)
	router := gin.New()
	router.Use(gin.Recovery())
	router.Use(middleware.LoggingMiddleware())
	router.Use(gin.Logger())

	// Explicit CORS - allows React localhost:3000, all methods, common headers
	corsMiddleware := cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000", "http://127.0.0.1:3000"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	})
	router.Use(corsMiddleware)
	router.Use(middleware.PrometheusMiddleware())

	// Metrics endpoint
	router.GET("/metrics", middleware.MetricsHandler())

	// Health checks
	router.GET("/health", func(c *gin.Context) {
		requestID, _ := c.Get("requestID")
		status := "healthy"
		dbStatus := "ok"

		if err := database.DB.PingContext(c.Request.Context()); err != nil {
			dbStatus = "error: " + err.Error()
			status = "degraded"
		}

		c.JSON(http.StatusOK, gin.H{
			"status":     status,
			"service":    "backend-api",
			"version":    version,
			"db_status":  dbStatus,
			"timestamp":  time.Now().Format(time.RFC3339),
			"request_id": requestID,
		})
	})

	router.GET("/live", func(c *gin.Context) {
		requestID, _ := c.Get("requestID")
		c.JSON(http.StatusOK, gin.H{
			"status":     "live",
			"service":    "backend-api",
			"timestamp":  time.Now().Format(time.RFC3339),
			"request_id": requestID,
		})
	})

	router.GET("/ready", func(c *gin.Context) {
		requestID, _ := c.Get("requestID")
		status := "ready"
		msg := []string{}

		ctx := c.Request.Context()
		if err := database.DB.PingContext(ctx); err != nil {
			status = "not_ready"
			msg = append(msg, "db: "+err.Error())
		}

		// Redis ping via TCP (no client dep)
		if redisConn, err := net.DialTimeout("tcp", "redis:6379", 2*time.Second); err != nil {
			status = "not_ready"
			msg = append(msg, "redis: "+err.Error())
		} else {
			redisConn.Close()
		}

		// Goose status mock (check migration dir exists)
		if _, err := os.Stat("migrations"); os.IsNotExist(err) {
			msg = append(msg, "migrations: missing")
		}

		c.JSON(http.StatusOK, gin.H{
			"status":     status,
			"service":    "backend-api",
			"checks":     msg,
			"timestamp":  time.Now().Format(time.RFC3339),
			"request_id": requestID,
		})
	})

	// API routes
	api := router.Group("/api/v1")
	{
		// Government Jobs
		api.GET("/gov-jobs", govJobHandler.GetGovJobs)
		api.GET("/gov-jobs/:id", govJobHandler.GetGovJobByID)

		// Private Jobs
		api.GET("/private-jobs", privJobHandler.GetPrivJobs)
		api.GET("/private-jobs/:id", privJobHandler.GetPrivJobByID)

		// Courses
		api.GET("/courses", courseHandler.GetCourses)
		api.GET("/courses/:id", courseHandler.GetCourseByID)

		// Videos
		api.GET("/videos", videoHandler.GetVideos)
		api.GET("/videos/:id", videoHandler.GetVideoByID)
	}

	cfg := config.Load()

	// HTTP Server (:8083 for metrics/internal)
	httpSrv := &http.Server{
		Addr:    fmt.Sprintf(":%d", serverPort),
		Handler: router,
	}

	go func() {
		logger.Info().Int("port", serverPort).Msg("Starting HTTP server")
		if err := httpSrv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Error().Err(err).Msg("HTTP server error")
		}
	}()

	var httpsSrv *http.Server

	// HTTPS Server (:8443) if enabled
	if cfg.TLS.Enabled {
		httpsSrv = &http.Server{
			Addr:      ":8443",
			Handler:   router,
			TLSConfig: &tls.Config{Certificates: []tls.Certificate{cfg.LoadTLSCert()}},
		}
		go func() {
			logger.Info().Msg("Starting HTTPS server on :8443")
			if err := httpsSrv.ListenAndServeTLS("", ""); err != nil && err != http.ErrServerClosed {
				logger.Error().Err(err).Msg("HTTPS server error")
			}
		}()
		logger.Info().Msg("TLS enabled, certs loaded from " + cfg.TLS.Cert)
	}

	logger.Info().Msg("Backend API service ready")

	// Wait for shutdown
	<-ctx.Done()
	logger.Info().Msg("Shutting down servers...")

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer shutdownCancel()

	if err := httpSrv.Shutdown(shutdownCtx); err != nil {
		logger.Error().Err(err).Msg("HTTP server shutdown failed")
	}

	if cfg.TLS.Enabled {
		if err := httpsSrv.Shutdown(shutdownCtx); err != nil {
			logger.Error().Err(err).Msg("HTTPS server shutdown failed")
		}
	}

	logger.Info().Msg("Server shutdown complete")
	return nil
}
