package tests

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
	"github.com/rojgarsetu/backend/config"
	"github.com/rojgarsetu/backend/internal/db"
	"github.com/rojgarsetu/backend/internal/handlers"
	"github.com/rojgarsetu/backend/internal/services"
)

func TestInputValidationAndXSSSanitization(t *testing.T) {
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
		t.Skipf("Skipping validation test: db ping failed: %v", err)
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
	authRoutes.POST("/register", authHandler.Register)

	// 1. Test Malformed Email -> EXPECT 400 Bad Request
	badEmailPayload := map[string]string{
		"name":     "Valid Name",
		"email":    "not-a-valid-email",
		"password": "ValidPassword123!",
		"role":     "candidate",
	}
	bodyBadEmail, _ := json.Marshal(badEmailPayload)
	reqBadEmail, _ := http.NewRequest("POST", "/api/v1/auth/register", bytes.NewBuffer(bodyBadEmail))
	reqBadEmail.Header.Set("Content-Type", "application/json")
	wBadEmail := httptest.NewRecorder()
	router.ServeHTTP(wBadEmail, reqBadEmail)

	if wBadEmail.Code != http.StatusBadRequest {
		t.Errorf("Expected 400 Bad Request for invalid email format, got %d", wBadEmail.Code)
	} else {
		t.Log("Invalid email format correctly REJECTED with 400 Bad Request.")
	}

	// 2. Test XSS Payload in Name -> EXPECT Sanitized Entry in DB
	xssName := "<script>alert('xss')</script>"
	validEmail := fmt.Sprintf("xss_test_%d@example.com", time.Now().UnixNano())
	xssPayload := map[string]string{
		"name":     xssName,
		"email":    validEmail,
		"password": "ValidPassword123!",
		"role":     "candidate",
	}
	bodyXSS, _ := json.Marshal(xssPayload)
	reqXSS, _ := http.NewRequest("POST", "/api/v1/auth/register", bytes.NewBuffer(bodyXSS))
	reqXSS.Header.Set("Content-Type", "application/json")
	wXSS := httptest.NewRecorder()
	router.ServeHTTP(wXSS, reqXSS)

	if wXSS.Code != http.StatusCreated {
		t.Fatalf("Registration failed with status %d: %s", wXSS.Code, wXSS.Body.String())
	}

	// Verify DB Record is Sanitized
	var savedName string
	err = sqlDB.QueryRow("SELECT name FROM users WHERE email = $1", validEmail).Scan(&savedName)
	if err != nil {
		t.Fatalf("Failed to query created user from DB: %v", err)
	}

	if strings.Contains(savedName, "<script>") {
		t.Errorf("XSS Vulnerability! Raw script tag stored in DB: %s", savedName)
	} else {
		t.Logf("XSS Sanitization PASS. Stored name: %s", savedName)
	}
}
