package services

import (
	"context"
	"crypto/sha256"
	"database/sql"
	"fmt"
	"time"

	"github.com/google/uuid"

	"github.com/rojgarsetu/backend/internal/db"
)

type TokenService struct {
	db *db.PostgresDB
}

func NewTokenService(d *db.PostgresDB) *TokenService {
	return &TokenService{db: d}
}

func (s *TokenService) CreateRefreshToken(ctx context.Context, userID string, ipHash, uaHash sql.NullString, expiresAt time.Time) (string, error) {
	uid, err := uuid.Parse(userID)
	if err != nil {
		return "", fmt.Errorf("invalid user id: %w", err)
	}

	result, err := s.db.Queries.CreateRefreshToken(ctx, db.CreateRefreshTokenParams{
		UserID:    uid,
		IpHash:    ipHash,
		UaHash:    uaHash,
		ExpiresAt: expiresAt,
	})
	if err != nil {
		return "", err
	}

	// Generate deterministic hash from result.ID for client-side verification
	hashBytes := sha256.Sum256([]byte(result.ID.String()))
	hashStr := fmt.Sprintf("%x", hashBytes)
	return hashStr, nil
}

func (s *TokenService) GetRefreshToken(ctx context.Context, token string) (*db.RefreshToken, error) {
	hash := sha256.Sum256([]byte(token))
	hashedTokenStr := fmt.Sprintf("%x", hash)
	tokenUID, err := uuid.Parse(hashedTokenStr)
	if err != nil {
		return nil, fmt.Errorf("invalid hashed token: %w", err)
	}
	result, err := s.db.Queries.GetRefreshTokenByToken(ctx, tokenUID)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

func (s *TokenService) RevokeRefreshToken(ctx context.Context, token string) error {
	hash := sha256.Sum256([]byte(token))
	hashedTokenStr := fmt.Sprintf("%x", hash)
	tokenUID, err := uuid.Parse(hashedTokenStr)
	if err != nil {
		return fmt.Errorf("invalid token: %w", err)
	}
	return s.db.Queries.RevokeRefreshToken(ctx, tokenUID)
}

func (s *TokenService) RevokeAllTokensForUser(ctx context.Context, userID string) error {
	uid, err := uuid.Parse(userID)
	if err != nil {
		return fmt.Errorf("invalid user id: %w", err)
	}
	return s.db.Queries.RevokeAllTokensForUser(ctx, uid)
}

func (s *TokenService) CleanupExpiredTokens(ctx context.Context) (int64, error) {
	_, err := s.db.Queries.CleanupExpiredTokens(ctx)
	if err != nil {
		return 0, err
	}
	return 1, nil // or query ROW_COUNT if needed
}
