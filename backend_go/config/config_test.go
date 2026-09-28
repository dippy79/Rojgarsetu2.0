package config

import "testing"

func TestLoadRequiresJWTSecrets(t *testing.T) {
	t.Setenv("JWT_SECRET", "short_secret")
	t.Setenv("REFRESH_TOKEN_KEY", "short_key")

	defer func() {
		if r := recover(); r == nil {
			t.Fatal("expected Load to panic when JWT secrets are under 32 chars")
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

func TestLoadRejectsMissingSecretsInProduction(t *testing.T) {
	t.Setenv("ENVIRONMENT", "production")
	t.Setenv("JWT_SECRET", "")
	t.Setenv("JWT_SECRET_FILE", "")
	t.Setenv("REFRESH_TOKEN_KEY", "")
	t.Setenv("REFRESH_TOKEN_KEY_FILE", "")

	defer func() {
		if r := recover(); r == nil {
			t.Fatal("expected Load to reject missing production secrets")
		}
	}()

	Load()
}

func TestLoadRejectsDevelopmentSecretsInProduction(t *testing.T) {
	t.Setenv("ENVIRONMENT", "production")
	t.Setenv("JWT_SECRET", "super-secret-jwt-key-minimum-32-characters-long")
	t.Setenv("JWT_SECRET_FILE", "")
	t.Setenv("REFRESH_TOKEN_KEY", "super-secret-refresh-key-minimum-32-chars")
	t.Setenv("REFRESH_TOKEN_KEY_FILE", "")

	defer func() {
		if r := recover(); r == nil {
			t.Fatal("expected Load to reject known development secrets in production")
		}
	}()

	Load()
}
