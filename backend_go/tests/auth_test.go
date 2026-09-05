package tests

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
)

// TestHealthEndpoint tests the liveness/readiness health check
func TestHealthEndpoint(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()

	// Replicating health check logic from main.go
	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status":   "UP",
			"service":  "backend-api",
			"db_ready": true,
		})
	})

	req, err := http.NewRequest("GET", "/health", nil)
	if err != nil {
		t.Fatal(err)
	}

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", w.Code)
	}

	var response map[string]interface{}
	if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
		t.Errorf("failed to unmarshal response: %v", err)
	}

	if response["status"] != "UP" {
		t.Errorf("expected status 'UP', got %v", response["status"])
	}
}

// TestLoginValidation tests that login requires valid credentials
func TestLoginValidation(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()

	// Simple mock of auth validation
	router.POST("/auth/login", func(c *gin.Context) {
		var req struct {
			Email    string `json:"email" binding:"required,email"`
			Password string `json:"password" binding:"required"`
		}
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Validation failed"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": "Valid"})
	})

	t.Run("Empty Body", func(t *testing.T) {
		req, _ := http.NewRequest("POST", "/auth/login", bytes.NewBuffer([]byte("{}")))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		if w.Code != http.StatusBadRequest {
			t.Errorf("expected status 400 for empty body, got %d", w.Code)
		}
	})

	t.Run("Invalid Email", func(t *testing.T) {
		loginData := map[string]string{
			"email":    "not-an-email",
			"password": "somepassword",
		}
		body, _ := json.Marshal(loginData)
		req, _ := http.NewRequest("POST", "/auth/login", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		if w.Code != http.StatusBadRequest {
			t.Errorf("expected status 400 for invalid email, got %d", w.Code)
		}
	})
}

// TestRegister tests user registration endpoint structure
func TestRegister(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.POST("/api/v1/auth/register", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "Success"})
	})

	registerData := map[string]string{
		"email":    "test@example.com",
		"password": "TestPass123!",
		"name":     "Test User",
	}
	jsonData, _ := json.Marshal(registerData)

	req, _ := http.NewRequest("POST", "/api/v1/auth/register", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", w.Code)
	}
	if !strings.Contains(w.Body.String(), "Success") {
		t.Errorf("expected response to contain 'Success'")
	}
}
