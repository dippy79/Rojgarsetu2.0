package services

import (
	"context"
	"database/sql"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/rojgarsetu/backend/config"
	"github.com/rojgarsetu/backend/internal/auth"
	"github.com/rojgarsetu/backend/internal/db"
	"github.com/rojgarsetu/backend/internal/middleware"
	"github.com/rs/zerolog/log"
)

type AuthService struct {
	db            *db.PostgresDB
	userSvc       *UserService
	cfg           *config.Config
	cleanupTicker *time.Ticker
	done          chan struct{}
}

func NewAuthService(db *db.PostgresDB, userSvc *UserService, cfg *config.Config) *AuthService {
	s := &AuthService{
		db:            db,
		userSvc:       userSvc,
		cfg:           cfg,
		cleanupTicker: time.NewTicker(1 * time.Hour),
		done:          make(chan struct{}),
	}
	go s.cleanupExpiredTokens()
	return s
}

func (s *AuthService) Stop() {
	s.cleanupTicker.Stop()
	close(s.done)
}

func (s *AuthService) cleanupExpiredTokens() {
	for {
		select {
		case <-s.cleanupTicker.C:
			ctx := context.Background()
			_, err := s.db.CleanupExpiredTokens(ctx)
			if err != nil {
				log.Error().Err(err).Msg("failed to cleanup expired refresh tokens")
			} else {
				log.Debug().Msg("cleanup expired refresh tokens completed")
			}
		case <-s.done:
			return
		}
	}
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
}

type LoginResponse struct {
	AccessToken  string  `json:"access_token"`
	RefreshToken string  `json:"refresh_token"`
	User         db.User `json:"user"`
}

func (s *AuthService) Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := s.userSvc.Login(c, req.Email, req.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
		return
	}

	// Generate tokens
	accessToken, err := auth.GenerateAccessToken(s.cfg, user.ID.String(), user.Email, user.Role)
	if err != nil {
		log.Error().Err(err).Msg("failed to generate access token")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "token generation failed"})
		return
	}

	refreshID := auth.GenerateRefreshTokenID()
	ip := c.ClientIP()
	ua := c.GetHeader("User-Agent")
	sessionHash := auth.HashSession(ip, ua)

	expiresAt := time.Now().Add(s.cfg.JWT.RefreshTokenExpiry)
	rtoken, err := s.db.CreateRefreshToken(c, db.CreateRefreshTokenParams{
		UserID:    user.ID,
		IpHash:    sql.NullString{String: sessionHash, Valid: true},
		UaHash:    sql.NullString{String: ua, Valid: ua != ""},
		ExpiresAt: expiresAt,
	})
	if err != nil {
		log.Error().Err(err).Msg("failed to create refresh token")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create session"})
		return
	}

	c.JSON(http.StatusOK, LoginResponse{
		AccessToken:  accessToken,
		RefreshToken: rtoken.ID.String(),
		User:         *user,
	})
}

type RefreshRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

type RefreshResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

func (s *AuthService) Refresh(c *gin.Context) {
	var req RefreshRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	refreshIDStr := req.RefreshToken
	if refreshIDStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "refresh_token required"})
		return
	}

	refreshIDStr = auth.GenerateRefreshTokenID()
	if refreshIDStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "refresh_token required"})
		return
	}

	// Validate old token
	rtoken, err := s.db.GetRefreshTokenByToken(c, refreshID)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid refresh token"})
		return
	}

	// Check session binding
	ip := c.ClientIP()
	ua := c.GetHeader("User-Agent")
	sessionHash := auth.HashSession(ip, ua)
	if rtoken.IpHash.String != sessionHash {
		log.Warn().Str("expected_hash", sessionHash).Str("stored_hash", rtoken.IpHash.String).Msg("session binding mismatch")
		s.db.RevokeRefreshToken(c, refreshID)
		c.JSON(http.StatusUnauthorized, gin.H{"error": "session changed (IP/UA mismatch)"})
		return
	}

	user, err := s.userSvc.GetUserByID(c, rtoken.UserID.String())
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user not found"})
		return
	}

	// Revoke old token
	err = s.db.RevokeRefreshToken(c, refreshID)
	if err != nil {
		log.Error().Err(err).Msg("failed to revoke old refresh token")
	}

	// Issue new pair
	newAccessToken, err := auth.GenerateAccessToken(s.cfg, user.ID.String(), user.Email, user.Role)
	if err != nil {
		log.Error().Err(err).Msg("failed to generate new access token")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "token generation failed"})
		return
	}

	newRefreshID := auth.GenerateRefreshTokenID()
	expiresAt := time.Now().Add(s.cfg.JWT.RefreshTokenExpiry)
	newRtoken, err := s.db.CreateRefreshToken(c, db.CreateRefreshTokenParams{
		UserID:    user.ID,
		IpHash:    sql.NullString{String: sessionHash, Valid: true},
		UaHash:    sql.NullString{String: ua, Valid: ua != ""},
		ExpiresAt: expiresAt,
	})
	if err != nil {
		log.Error().Err(err).Msg("failed to create new refresh token")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "session creation failed"})
		return
	}

	c.JSON(http.StatusOK, RefreshResponse{
		AccessToken:  newAccessToken,
		RefreshToken: newRtoken.ID.String(),
	})
}

type LogoutRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

func (s *AuthService) Logout(c *gin.Context) {
	var req LogoutRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	refreshID, err := uuid.Parse(req.RefreshToken)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid refresh token"})
		return
	}

	err = s.db.RevokeRefreshToken(c, refreshID)
	if err != nil {
		log.Error().Err(err).Msg("failed to revoke refresh token")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "logout failed"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "logged out successfully"})
}

func (s *AuthService) LogoutAll(c *gin.Context) {
	userIDStr := middleware.GetUserID(c)
	if userIDStr == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user ID"})
		return
	}

	err = s.db.RevokeAllTokensForUser(c, userID)
	if err != nil {
		log.Error().Err(err).Msg("failed to revoke all tokens for user")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "logout all failed"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "all sessions logged out"})
}
