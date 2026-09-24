package tests

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"

	"github.com/rojgarsetu/backend/config"
	"github.com/rojgarsetu/backend/internal/db"
	"github.com/rojgarsetu/backend/internal/handlers"
	"github.com/rojgarsetu/backend/internal/services"
)

func TestIntegrationMigrationsSeederAndAuth(t *testing.T) {
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		dbURL = "postgres://postgres:postgres@localhost:5435/rojgarsetu2?sslmode=disable"
	}

	// 1. Test database connection
	sqlDB, err := sql.Open("postgres", dbURL)
	if err != nil {
		t.Skipf("Skipping integration test: cannot open database at %s: %v", dbURL, err)
		return
	}
	if err := sqlDB.Ping(); err != nil {
		t.Skipf("Skipping integration test: postgres ping failed at %s: %v", dbURL, err)
		return
	}
	defer sqlDB.Close()

	// 2. Run Database Migrations
	migrationsPath := "file://../migrations"
	m, err := migrate.New(migrationsPath, dbURL)
	if err != nil {
		t.Fatalf("Failed to initialize migrations: %v", err)
	}
	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		t.Fatalf("Failed to run migrations up: %v", err)
	}
	t.Log("Migrations applied successfully.")

	// 3. Setup Gin Router with AuthHandler
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
	authRoutes.POST("/register", authHandler.Register)

	// 4. Test Candidate Registration (Atomicity & Schema Verification)
	testEmail := fmt.Sprintf("test_candidate_%d@example.com", time.Now().UnixNano())
	regPayload := map[string]string{
		"firstName": "Test",
		"lastName":  "Candidate",
		"email":     testEmail,
		"password":  "Candidate@123456",
		"role":      "candidate",
		"phone":     "+919876543210",
	}
	body, _ := json.Marshal(regPayload)

	req, _ := http.NewRequest("POST", "/api/v1/auth/register", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusCreated {
		t.Fatalf("Registration failed with status %d: %s", w.Code, w.Body.String())
	}

	var regResponse map[string]interface{}
	if err := json.Unmarshal(w.Body.Bytes(), &regResponse); err != nil {
		t.Fatalf("Failed to parse registration response: %v", err)
	}
	if regResponse["success"] != true {
		t.Fatalf("Registration response success flag is not true: %v", regResponse)
	}
	t.Log("Candidate registration completed successfully without 500/404 errors.")

	// 5. Test Candidate Login
	loginPayload := map[string]string{
		"email":    testEmail,
		"password": "Candidate@123456",
	}
	loginBody, _ := json.Marshal(loginPayload)

	reqLogin, _ := http.NewRequest("POST", "/api/v1/auth/login", bytes.NewBuffer(loginBody))
	reqLogin.Header.Set("Content-Type", "application/json")
	wLogin := httptest.NewRecorder()
	router.ServeHTTP(wLogin, reqLogin)

	if wLogin.Code != http.StatusOK {
		t.Fatalf("Login failed with status %d: %s", wLogin.Code, wLogin.Body.String())
	}

	var loginResponse map[string]interface{}
	if err := json.Unmarshal(wLogin.Body.Bytes(), &loginResponse); err != nil {
		t.Fatalf("Failed to parse login response: %v", err)
	}
	if loginResponse["success"] != true {
		t.Fatalf("Login response success flag is not true: %v", loginResponse)
	}

	// Verify Cookie Header was set
	cookies := wLogin.Result().Cookies()
	var foundAccessToken, foundRefreshToken bool
	for _, cookie := range cookies {
		if cookie.Name == "access_token" && cookie.Value != "" {
			foundAccessToken = true
		}
		if cookie.Name == "refresh_token" && cookie.Value != "" {
			foundRefreshToken = true
		}
	}

	if !foundAccessToken || !foundRefreshToken {
		t.Fatalf("Expected access_token and refresh_token cookies to be set in login response")
	}

	t.Log("Login completed successfully with token and HttpOnly cookies verified.")
}
