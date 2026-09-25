package tests

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
	"github.com/rojgarsetu/backend/config"
	"github.com/rojgarsetu/backend/internal/db"
	"github.com/rojgarsetu/backend/internal/handlers"
	"github.com/rojgarsetu/backend/internal/middleware"
	"github.com/rojgarsetu/backend/internal/services"
)

func TestSessionPersistenceAndRefresh(t *testing.T) {
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		dbURL = "postgres://postgres:postgres@127.0.0.1:5435/rojgarsetu2?sslmode=disable"
	}

	sqlDB, err := sql.Open("postgres", dbURL)
	if err != nil {
		t.Fatalf("Failed to open DB: %v", err)
	}
	defer sqlDB.Close()

	if err := sqlDB.Ping(); err != nil {
		t.Skipf("Skipping session test: db ping failed: %v", err)
		return
	}

	cfg := config.Load()
	if cfg.JWT.Secret == "" {
		cfg.JWT.Secret = "super-secret-jwt-key-minimum-32-characters-long"
	}

	postgresDB := db.NewPostgresDB(sqlDB)
	userSvc := services.NewUserService(postgresDB)
	tokenSvc := services.NewTokenService(postgresDB)
	authSvc := services.NewAuthService(userSvc, tokenSvc, cfg)
	authHandler := handlers.NewAuthHandler(authSvc)

	gin.SetMode(gin.TestMode)
	router := gin.New()

	api := router.Group("/api/v1")
	authRoutes := api.Group("/auth")
	authRoutes.POST("/login", authHandler.Login)
	authRoutes.POST("/refresh", authHandler.Refresh)
	authRoutes.GET("/me", middleware.AuthMiddleware(cfg), authHandler.Me)

	// 1. Create User
	email := fmt.Sprintf("sess_user_%d@example.com", time.Now().UnixNano())
	pass := "UserPass123456!"
	_, err = userSvc.CreateUser(context.Background(), db.RegisterRequest{
		Name:     "Session User",
		Email:    email,
		Password: pass,
		Role:     "candidate",
	})
	if err != nil {
		t.Fatalf("Failed to create test user: %v", err)
	}

	// 2. Perform Initial Login
	loginPayload, _ := json.Marshal(map[string]string{"email": email, "password": pass})
	reqLogin, _ := http.NewRequest("POST", "/api/v1/auth/login", bytes.NewBuffer(loginPayload))
	reqLogin.Header.Set("Content-Type", "application/json")
	wLogin := httptest.NewRecorder()
	router.ServeHTTP(wLogin, reqLogin)

	if wLogin.Code != http.StatusOK {
		t.Fatalf("Login failed with status %d: %s", wLogin.Code, wLogin.Body.String())
	}

	var accessCookie, refreshCookie *http.Cookie
	for _, c := range wLogin.Result().Cookies() {
		if c.Name == "access_token" {
			accessCookie = c
		}
		if c.Name == "refresh_token" {
			refreshCookie = c
		}
	}

	if accessCookie == nil || refreshCookie == nil {
		t.Fatalf("Login response missing required access/refresh cookies")
	}

	// 3. Verify /me Endpoint (Page Refresh Simulation #1)
	reqMe1, _ := http.NewRequest("GET", "/api/v1/auth/me", nil)
	reqMe1.AddCookie(accessCookie)
	wMe1 := httptest.NewRecorder()
	router.ServeHTTP(wMe1, reqMe1)

	if wMe1.Code != http.StatusOK {
		t.Fatalf("Initial /me check failed with status %d: %s", wMe1.Code, wMe1.Body.String())
	}

	// 4. Token Refresh Simulation (Simulating session renewal after tab reopen)
	reqRefresh, _ := http.NewRequest("POST", "/api/v1/auth/refresh", nil)
	reqRefresh.AddCookie(refreshCookie)
	wRefresh := httptest.NewRecorder()
	router.ServeHTTP(wRefresh, reqRefresh)

	if wRefresh.Code != http.StatusOK {
		t.Fatalf("Token refresh failed with status %d: %s", wRefresh.Code, wRefresh.Body.String())
	}

	var newAccessCookie *http.Cookie
	for _, c := range wRefresh.Result().Cookies() {
		if c.Name == "access_token" {
			newAccessCookie = c
		}
	}

	if newAccessCookie == nil {
		t.Fatalf("Refresh response missing new access_token cookie")
	}

	// 5. Verify /me Endpoint with Refreshed Access Token
	reqMe2, _ := http.NewRequest("GET", "/api/v1/auth/me", nil)
	reqMe2.AddCookie(newAccessCookie)
	wMe2 := httptest.NewRecorder()
	router.ServeHTTP(wMe2, reqMe2)

	if wMe2.Code != http.StatusOK {
		t.Fatalf("Refreshed session /me check failed with status %d: %s", wMe2.Code, wMe2.Body.String())
	}

	t.Log("Session persistence and token refresh verified 100% PASS.")
}
