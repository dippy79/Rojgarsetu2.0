package auth

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/rojgarsetu/backend/config"
)

type Claims struct {
	UserID string `json:"user_id"`
	Email  string `json:"email"`
	Role   string `json:"role"`
	jwt.RegisteredClaims
}

// GenerateAccessToken creates a new access token
func GenerateAccessToken(cfg *config.Config, userID, email, role string) (string, error) {
	claims := Claims{
		UserID: userID,
		Email:  email,
		Role:   role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(cfg.JWT.AccessTokenExpiry)),
			Issuer:    cfg.JWT.Issuer,
			Audience:  []string{cfg.JWT.Audience},
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(cfg.JWT.Secret))
}

// GenerateRefreshTokenID creates new refresh token ID string
func GenerateRefreshTokenID() string {
	return fmt.Sprintf("rt_%d", time.Now().UnixNano())
}

// HashSession binds IP + UA with HMAC
func HashSession(ip, ua string, cfg *config.Config) string {
	h := hmac.New(sha256.New, []byte(cfg.JWT.RefreshSessionKey))
	h.Write([]byte(ip + ua))
	return hex.EncodeToString(h.Sum(nil))
}

// ParseAndValidateAccessToken validates and parses token
func ParseAccessToken(tokenStr string, cfg *config.Config) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenStr, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method")
		}
		return []byte(cfg.JWT.Secret), nil
	})
	if err != nil {
		return nil, err
	}
	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}
	return nil, fmt.Errorf("invalid token")
}
