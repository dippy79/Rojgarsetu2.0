package main

import (
	"os"
)

func init() {
	if os.Getenv("DATABASE_URL") == "" {
		os.Setenv("DATABASE_URL", "postgres://postgres:postgres@localhost:5435/rojgarsetu2?sslmode=disable")
	}
	if os.Getenv("REDIS_URL") == "" {
		os.Setenv("REDIS_URL", "redis://localhost:6380")
	}
	if os.Getenv("JWT_SECRET") == "" {
		os.Setenv("JWT_SECRET", "super-secret-jwt-key-minimum-32-characters-long")
	}
	if os.Getenv("REFRESH_TOKEN_KEY") == "" {
		os.Setenv("REFRESH_TOKEN_KEY", "super-secret-refresh-key-minimum-32-chars")
	}
	if os.Getenv("CORS_ORIGINS") == "" {
		os.Setenv("CORS_ORIGINS", "http://localhost:8080,http://localhost:3000,http://localhost:3001")
	}
	if os.Getenv("PORT") == "" {
		os.Setenv("PORT", "8083")
	}
	if os.Getenv("MIGRATIONS_PATH") == "" {
		os.Setenv("MIGRATIONS_PATH", "file://./migrations")
	}
}
