package main

import (
	"os"
)

func init() {
	os.Setenv("DATABASE_URL", "postgres://amitsharma:amitsharma@localhost:5432/rojgarsetu2?sslmode=disable")
	os.Setenv("REDIS_URL", "redis://localhost:6379")
	os.Setenv("JWT_SECRET", "test_secret_key_for_development_please_change_in_production_32chars")
	os.Setenv("REFRESH_TOKEN_KEY", "test_refresh_key_for_development_32chars_minimum")
	os.Setenv("CORS_ORIGINS", "http://localhost:8080,http://localhost:3000")
	os.Setenv("PORT", "8083")
	os.Setenv("MIGRATIONS_PATH", "file://./migrations")
}
