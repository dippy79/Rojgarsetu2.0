package main

import (
	"context"
	"crypto/tls"
	"database/sql"
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
	serverPort int
	redisURL   string
)

func main() {
	// Initialize logger
	logger = zerolog.New(os.Stderr).Output(zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: time.RFC3339}).With().Timestamp().Caller().Logger()

	// Load config early
	cfg := config.Load()
	if cfg.JWT.Secret == "" {
		logger.Fatal().Msg("JWT_SECRET required")
		return
	}
	if len(cfg.JWT.Secret) < 32 {
		logger.Fatal().Msg("JWT_SECRET must be at least 32 characters")
		return
	}

	dbURL := os.Getenv("DATABASE_URL")
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
	if err := run(cfg, dbURL); err != nil {
		logger.Fatal().Err(err).Msg("Server failed")
	}
}

func run(cfg *config.Config, dbURL string) error {
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

	// Initialize database - FIXED
	sqlDB, err := sql.Open("postgres", dbURL)
	if err != nil {
		return fmt.Errorf("sql open: %w", err)
	}
	defer sqlDB.Close()
	if err = sqlDB.Ping(); err != nil {
		return fmt.Errorf("db ping: %w", err)
	}
	database := db.NewPostgresDB(sqlDB)
	logger.Info().Msg("Database connected")

	// Initialize services - FIXED to ContentService
	contentSvc := services.NewContentService(database)

	// Initialize handlers - FIXED to use ContentService
	govJobHandler := handlers.NewGovJobHandler(contentSvc)
	privJobHandler := handlers.NewPrivJobHandler(contentSvc)
	courseHandler := handlers.NewCourseHandler(contentSvc)
	videoHandler := handlers.NewVideoHandler(contentSvc)

	middleware.RegisterBuildInfo()

	// Setup Gin router
	gin.SetMode(gin.ReleaseMode)
	router := gin.New()
	router.Use(gin.Recovery())
	router.Use(middleware.LoggingMiddleware())
	router.Use(gin.Logger())

	// CORS
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

	// Health checks - FIXED DB access via sqlDB (keep alive)
	router.GET("/health", func(c *gin.Context) {
		requestID, _ := c.Get("requestID")
		status := "healthy"
		dbStatus := "ok"

		if err := sqlDB.PingContext(c.Request.Context()); err != nil {
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
		if err := sqlDB.PingContext(ctx); err != nil {
			status = "not_ready"
			msg = append(msg, "db: "+err.Error())
		}

		// Redis ping via TCP
		if redisConn, err := net.DialTimeout("tcp", "redis:6379", 2*time.Second); err != nil {
			status = "not_ready"
			msg = append(msg, "redis: "+err.Error())
		} else {
			redisConn.Close()
		}

		// Migrations dir check
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

	// API routes - FIXED to gin.HandlerFunc wrapper
	api := router.Group("/api/v1")
	{
		api.GET("/gov-jobs", gin.WrapF(govJobHandler.List))
		api.GET("/gov-jobs/:id", gin.WrapF(govJobHandler.GetByID))

		api.GET("/private-jobs", gin.WrapF(privJobHandler.List))
		api.GET("/private-jobs/:id", gin.WrapF(privJobHandler.GetByID))

		api.GET("/courses", gin.WrapF(courseHandler.List))
		api.GET("/courses/:id", gin.WrapF(courseHandler.GetByID))

		api.GET("/videos", gin.WrapF(videoHandler.List))
		api.GET("/videos/:id", gin.WrapF(videoHandler.GetByID))
	}

	// HTTP Server
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

	// HTTPS if enabled
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

	if cfg.TLS.Enabled && httpsSrv != nil {
		if err := httpsSrv.Shutdown(shutdownCtx); err != nil {
			logger.Error().Err(err).Msg("HTTPS server shutdown failed")
		}
	}

	logger.Info().Msg("Server shutdown complete")
	return nil
}
