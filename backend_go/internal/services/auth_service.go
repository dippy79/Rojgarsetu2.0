package services

import (
	"database/sql"
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/rojgarsetu/backend/config"
	"github.com/rojgarsetu/backend/internal/db"
	"github.com/rojgarsetu/backend/internal/middleware"
)

type AuthService struct {
	userSvc  *UserService
	tokenSvc *TokenService
	cfg      *config.Config
}

func NewAuthService(userSvc *UserService, tokenSvc *TokenService, cfg *config.Config) *AuthService {
	return &AuthService{
		userSvc:  userSvc,
		tokenSvc: tokenSvc,
		cfg:      cfg,
	}
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=12"`
}

type RegisterRequest struct {
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Name      string `json:"name"` // Fallback for single name field
	Email     string `json:"email" binding:"required,email"`
	Password  string `json:"password" binding:"required,min=12"`
	Role      string `json:"role" binding:"required,oneof=candidate company"`
	Phone     string `json:"phone"`
}

type TokenResponse struct {
	Success bool   `json:"success"`
	Data    struct {
		User  *db.User `json:"user"`
		Token string   `json:"token"`
	} `json:"data"`
}

func (s *AuthService) Register(c *gin.Context) {
	var req RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	fullName := req.Name
	if req.FirstName != "" || req.LastName != "" {
		fullName = strings.TrimSpace(req.FirstName + " " + req.LastName)
	}

	if fullName == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Name is required"})
		return
	}

	// Password strength check (Manual check as backup to binding)
	if !validatePassword(req.Password) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Password must be 12+ chars, include uppercase, and a special character"})
		return
	}

	// EmailExists check
	_, err := s.userSvc.GetUserByEmail(c, req.Email)
	if err == nil {
		c.JSON(http.StatusConflict, gin.H{"error": "Email already exists"})
		return
	}

	userReq := db.RegisterRequest{
		Name:     fullName,
		Email:    req.Email,
		Password: req.Password,
		Role:     req.Role,
		Phone:    sql.NullString{String: req.Phone, Valid: req.Phone != ""},
	}
	user, err := s.userSvc.CreateUser(c, userReq)
	if err != nil {
		if errors.Is(err, ErrCompanyNameExists) {
			c.JSON(http.StatusConflict, gin.H{"error": ErrCompanyNameExists.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"success": true, "data": gin.H{"user": user}})
}

func (s *AuthService) GetUser(c *gin.Context, userID string) (*db.User, error) {
	return s.userSvc.GetUserByID(c, userID)
}

func validatePassword(p string) bool {
	if len(p) < 12 {
		return false
	}
	var hasUpper, hasSpecial bool
	for _, char := range p {
		if char >= 'A' && char <= 'Z' {
			hasUpper = true
		}
		if strings.ContainsAny(string(char), "!@#$%^&*()_+-=[]{}|;':\",./<>?") {
			hasSpecial = true
		}
	}
	return hasUpper && hasSpecial
}

func (s *AuthService) Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := s.userSvc.Login(c, req.Email, req.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	// Generate tokens
	expiresAt := time.Now().Add(30 * 24 * time.Hour) // 30 days for refresh
	refreshToken, err := s.tokenSvc.CreateRefreshToken(c, user.ID.String(), sql.NullString{}, sql.NullString{}, expiresAt)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create tokens"})
		return
	}

	accessToken, err := s.generateJWT(user.ID.String(), user.Email, user.Role, 15*time.Minute)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create tokens"})
		return
	}

	// Set HttpOnly Cookies for security
	cookieSecure := os.Getenv("COOKIE_SECURE") == "true"
	c.SetSameSite(http.SameSiteStrictMode)
	c.SetCookie("access_token", accessToken, 900, "/", "", cookieSecure, true)     // 15 min
	c.SetCookie("refresh_token", refreshToken, 2592000, "/", "", cookieSecure, true) // 30 days

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data": gin.H{
			"user":  user,
			"token": accessToken,
		},
	})
}

func (s *AuthService) Refresh(c *gin.Context) {
	// P1: Read from cookie instead of body
	refreshToken, err := c.Cookie("refresh_token")
	if err != nil {
		// Fallback to body for legacy/testing if absolutely necessary, but audit says MUST read from cookie
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Refresh token missing from cookie"})
		return
	}

	token, err := s.tokenSvc.GetRefreshToken(c, refreshToken)
	if err != nil || time.Now().After(token.ExpiresAt) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired refresh token"})
		return
	}

	user, err := s.userSvc.GetUserByID(c, token.UserID.String())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	accessToken, err := s.generateJWT(user.ID.String(), user.Email, user.Role, 15*time.Minute)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create access token"})
		return
	}

	// Update HttpOnly Cookie
	secure := true
	c.SetSameSite(http.SameSiteStrictMode)
	c.SetCookie("access_token", accessToken, 900, "/", "", secure, true)

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data": gin.H{
			"token": accessToken,
		},
	})
}

func (s *AuthService) Logout(c *gin.Context) {
	userID := middleware.GetUserID(c)
	if userID != "" {
		_ = s.tokenSvc.RevokeAllTokensForUser(c, userID)
	}

	// Clear cookies
	c.SetCookie("access_token", "", -1, "/", "", true, true)
	c.SetCookie("refresh_token", "", -1, "/", "", true, true)

	c.JSON(http.StatusOK, gin.H{"success": true, "message": "Logged out all sessions"})
}

func (s *AuthService) generateJWT(userID, email, role string, expires time.Duration) (string, error) {
	claims := middleware.Claims{
		UserID: userID,
		Email:  email,
		Role:   role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(expires)),
			Issuer:    s.cfg.JWT.Issuer,
			Audience:  []string{s.cfg.JWT.Audience},
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(s.cfg.JWT.Secret))
}
