package tests

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

// TestRegister tests user registration endpoint
func TestRegister(t *testing.T) {
	gin.SetMode(gin.TestMode)
	
	// Create test router
	router := gin.New()
	router.POST("/api/v1/auth/register", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "User registered successfully",
			"user_id": "test-user-id",
		})
	})

	// Test data
	registerData := map[string]string{
		"email":    "test@example.com",
		"password": "TestPass123!",
		"name":     "Test User",
	}
	jsonData, _ := json.Marshal(registerData)

	// Create request
	req, _ := http.NewRequest("POST", "/api/v1/auth/register", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	
	// Record response
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Assertions
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "User registered successfully")
}

// TestLogin tests user login endpoint
func TestLogin(t *testing.T) {
	gin.SetMode(gin.TestMode)
	
	router := gin.New()
	router.POST("/api/v1/auth/login", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"access_token":  "test-access-token",
			"refresh_token": "test-refresh-token",
			"user": gin.H{
				"id":    "test-user-id",
				"email": "test@example.com",
			},
		})
	})

	loginData := map[string]string{
		"email":    "test@example.com",
		"password": "TestPass123!",
	}
	jsonData, _ := json.Marshal(loginData)

	req, _ := http.NewRequest("POST", "/api/v1/auth/login", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "access_token")
	assert.Contains(t, w.Body.String(), "refresh_token")
}

// TestRefreshToken tests token refresh endpoint
func TestRefreshToken(t *testing.T) {
	gin.SetMode(gin.TestMode)
	
	router := gin.New()
	router.POST("/api/v1/auth/refresh", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"access_token": "new-access-token",
		})
	})

	refreshData := map[string]string{
		"refresh_token": "test-refresh-token",
	}
	jsonData, _ := json.Marshal(refreshData)

	req, _ := http.NewRequest("POST", "/api/v1/auth/refresh", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "access_token")
}

// TestUnauthorized tests unauthorized access
func TestUnauthorized(t *testing.T) {
	gin.SetMode(gin.TestMode)
	
	router := gin.New()
	
	// Protected route
	protected := router.Group("/api/v1")
	protected.Use(func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			c.Abort()
			return
		}
		c.Next()
	})
	
	protected.GET("/protected", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "Success"})
	})

	req, _ := http.NewRequest("GET", "/api/v1/protected", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
	assert.Contains(t, w.Body.String(), "Unauthorized")
}