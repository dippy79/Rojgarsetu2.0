#!/bin/bash
# Use variables from .env if present
if [ -f .env ]; then
  export $(cat .env | xargs)
fi

export DATABASE_URL="${DATABASE_URL:-postgres://postgres:postgres@127.0.0.1:5435/rojgarsetu2?sslmode=disable}"
export REDIS_URL="${REDIS_URL:-redis://localhost:6380}"
export JWT_SECRET="${JWT_SECRET:-super-secret-jwt-key-minimum-32-characters-long}"
export REFRESH_TOKEN_KEY="${REFRESH_TOKEN_KEY:-super-secret-refresh-key-minimum-32-chars}"
export CORS_ORIGINS="${CORS_ORIGINS:-http://localhost:8080,http://localhost:3000,http://localhost:3001}"
export PORT="${PORT:-8083}"
export MIGRATIONS_PATH="${MIGRATIONS_PATH:-file://./migrations}"

./server
