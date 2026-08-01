package services

import (
	"database/sql"
	"errors"
	"net/http"
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
	Password string `json:"password" binding:"required,min=8"`
}

type RegisterRequest struct {
	Name     string `json:"name" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=8"`
	Role     string `json:"role" binding:"required,oneof=candidate company"`
	Phone    string `json:"phone"`
}

type TokenResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

func (s *AuthService) Register(c *gin.Context) {
	var req RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// EmailExists check
	_, err := s.userSvc.GetUserByEmail(c, req.Email)
	if err == nil {
		c.JSON(http.StatusConflict, gin.H{"error": "Email already exists"})
		return
	}

	userReq := db.RegisterRequest{
		Name:     req.Name,
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
	c.JSON(http.StatusCreated, gin.H{"data": user})
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
	expiresAt := time.Now().Add(24 * time.Hour)
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

	c.JSON(http.StatusOK, TokenResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	})
}

func (s *AuthService) Refresh(c *gin.Context) {
	var req struct {
		RefreshToken string `json:"refresh_token" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	token, err := s.tokenSvc.GetRefreshToken(c, req.RefreshToken)
	if err != nil || time.Now().After(token.ExpiresAt) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid refresh token"})
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

	c.JSON(http.StatusOK, gin.H{"access_token": accessToken})
}

func (s *AuthService) Logout(c *gin.Context) {
	userID := middleware.GetUserID(c)
	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}
	_ = s.tokenSvc.RevokeAllTokensForUser(c, userID)
	c.JSON(http.StatusOK, gin.H{"message": "Logged out all sessions"})
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
