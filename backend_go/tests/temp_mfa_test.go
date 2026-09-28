package tests

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/rojgarsetu/backend/config"
	"github.com/rojgarsetu/backend/internal/middleware"
)

func TestAdminMFAMiddlewareRequiresMFA(t *testing.T) {
	gin.SetMode(gin.TestMode)
	cfg := &config.Config{}
	cfg.JWT.Secret = "0123456789abcdef0123456789abcdef"
	cfg.JWT.Issuer = "rojgarsetu-backend"
	cfg.JWT.Audience = "rojgarsetu-api"

	router := gin.New()
	router.Use(middleware.AuthMiddleware(cfg))
	router.Use(middleware.AdminMFAMiddleware(cfg))
	router.GET("/api/admin", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"ok": true})
	})

	for _, test := range []struct {
		name       string
		mfa        bool
		wantStatus int
	}{
		{name: "admin without MFA is denied", mfa: false, wantStatus: http.StatusForbidden},
		{name: "admin with MFA is allowed", mfa: true, wantStatus: http.StatusOK},
	} {
		t.Run(test.name, func(t *testing.T) {
			claims := middleware.Claims{
				UserID:      "admin-1",
				Email:       "admin@example.com",
				Role:        "admin",
				MFAVerified: test.mfa,
				RegisteredClaims: jwt.RegisteredClaims{
					ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour)),
					Issuer:    cfg.JWT.Issuer,
					Audience:  jwt.ClaimStrings{cfg.JWT.Audience},
				},
			}
			token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
			signed, err := token.SignedString([]byte(cfg.JWT.Secret))
			if err != nil {
				t.Fatalf("sign token: %v", err)
			}

			req := httptest.NewRequest(http.MethodGet, "/api/admin", nil)
			req.Header.Set("Authorization", "Bearer "+signed)
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			if w.Code != test.wantStatus {
				t.Fatalf("expected %d, got %d body=%s", test.wantStatus, w.Code, w.Body.String())
			}
		})
	}
}
