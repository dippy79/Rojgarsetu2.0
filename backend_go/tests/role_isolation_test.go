package tests

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/rojgarsetu/backend/config"
	"github.com/rojgarsetu/backend/internal/middleware"
)

func TestRoleIsolationBOLA(t *testing.T) {
	cfg := config.Load()
	if cfg.JWT.Secret == "" {
		cfg.JWT.Secret = "super-secret-jwt-key-minimum-32-characters-long"
	}

	gin.SetMode(gin.TestMode)
	router := gin.New()

	api := router.Group("/api/v1")

	// Protected Employer Group
	companyRoutes := api.Group("/company")
	companyRoutes.Use(middleware.AuthMiddleware(cfg))
	companyRoutes.Use(middleware.CompanyMiddleware())
	companyRoutes.GET("/dashboard", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "company dashboard ok"})
	})

	// Protected Candidate Group
	candidateRoutes := api.Group("/candidate")
	candidateRoutes.Use(middleware.AuthMiddleware(cfg))
	candidateRoutes.Use(middleware.CandidateMiddleware())
	candidateRoutes.GET("/profile", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "candidate profile ok"})
	})

	// Generate Candidate JWT Token
	candToken, err := generateTestToken("cand_123", "cand@test.com", "candidate", cfg)
	if err != nil {
		t.Fatalf("Failed to generate candidate token: %v", err)
	}

	// Generate Company JWT Token
	compToken, err := generateTestToken("comp_456", "comp@test.com", "company", cfg)
	if err != nil {
		t.Fatalf("Failed to generate company token: %v", err)
	}

	// 1. Candidate attempts to access Employer Dashboard -> EXPECT 403 Forbidden
	reqComp, _ := http.NewRequest("GET", "/api/v1/company/dashboard", nil)
	reqComp.Header.Set("Authorization", "Bearer "+candToken)
	wComp := httptest.NewRecorder()
	router.ServeHTTP(wComp, reqComp)

	if wComp.Code != http.StatusForbidden {
		t.Errorf("BOLA Failure! Candidate accessed company dashboard. Expected 403, got %d", wComp.Code)
	} else {
		t.Log("Candidate access to company dashboard correctly BLOCKED with 403 Forbidden.")
	}

	// 2. Company attempts to access Candidate Profile -> EXPECT 403 Forbidden
	reqCand, _ := http.NewRequest("GET", "/api/v1/candidate/profile", nil)
	reqCand.Header.Set("Authorization", "Bearer "+compToken)
	wCand := httptest.NewRecorder()
	router.ServeHTTP(wCand, reqCand)

	if wCand.Code != http.StatusForbidden {
		t.Errorf("BOLA Failure! Company accessed candidate profile. Expected 403, got %d", wCand.Code)
	} else {
		t.Log("Company access to candidate profile correctly BLOCKED with 403 Forbidden.")
	}

	// 3. Unauthenticated access -> EXPECT 401 Unauthorized
	reqAnon, _ := http.NewRequest("GET", "/api/v1/company/dashboard", nil)
	wAnon := httptest.NewRecorder()
	router.ServeHTTP(wAnon, reqAnon)

	if wAnon.Code != http.StatusUnauthorized {
		t.Errorf("Expected 401 for unauthenticated request, got %d", wAnon.Code)
	} else {
		t.Log("Unauthenticated access correctly BLOCKED with 401 Unauthorized.")
	}
}

func generateTestToken(userID, email, role string, cfg *config.Config) (string, error) {
	authSvc := servicesForTest(cfg)
	_ = authSvc
	// Use authSvc or jwt library directly
	return createTokenHelper(userID, email, role, cfg)
}

func createTokenHelper(userID, email, role string, cfg *config.Config) (string, error) {
	claims := middleware.Claims{
		UserID: userID,
		Email:  email,
		Role:   role,
	}
	claims.Issuer = cfg.JWT.Issuer
	claims.Audience = []string{cfg.JWT.Audience}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(cfg.JWT.Secret))
}

func servicesForTest(cfg *config.Config) interface{} {
	return nil
}
