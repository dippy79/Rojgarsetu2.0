package config

import "testing"

func TestLoadRequiresJWTSecrets(t *testing.T) {
	t.Setenv("JWT_SECRET", "")
	t.Setenv("REFRESH_TOKEN_KEY", "")
	t.Setenv("JWT_ISSUER", "")
	t.Setenv("JWT_AUDIENCE", "")
	t.Setenv("ACCESS_TOKEN_EXPIRY", "")
	t.Setenv("REFRESH_TOKEN_EXPIRY", "")
	t.Setenv("RATE_LIMIT", "")
	t.Setenv("LOGIN_RATE_LIMIT", "")
	t.Setenv("DB_MAX_OPEN_CONNS", "")
	t.Setenv("DB_MAX_IDLE_CONNS", "")
	t.Setenv("DB_CONN_MAX_LIFETIME", "")
	t.Setenv("DB_CONN_MAX_IDLE_TIME", "")

	defer func() {
		if r := recover(); r == nil {
			t.Fatal("expected Load to panic when JWT secrets are missing")
		}
	}()

	Load()
}

func TestLoadUsesProvidedSecrets(t *testing.T) {
	t.Setenv("JWT_SECRET", "abcdefghijklmnopqrstuvwxyz123456")
	t.Setenv("REFRESH_TOKEN_KEY", "zyxwvutsrqponmlkjihgfedcba123456")
	t.Setenv("JWT_ISSUER", "test-issuer")
	t.Setenv("JWT_AUDIENCE", "test-audience")

	cfg := Load()
	if cfg.JWT.Secret != "abcdefghijklmnopqrstuvwxyz123456" {
		t.Fatalf("JWT secret not loaded from env: %q", cfg.JWT.Secret)
	}
	if cfg.JWT.RefreshSessionKey != "zyxwvutsrqponmlkjihgfedcba123456" {
		t.Fatalf("refresh token key not loaded from env: %q", cfg.JWT.RefreshSessionKey)
	}
}
