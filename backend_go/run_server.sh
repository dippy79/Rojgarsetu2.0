#!/bin/bash
export DATABASE_URL="postgres://amitsharma:amitsharma@localhost:5432/rojgarsetu2?sslmode=disable"
export REDIS_URL="redis://localhost:6379"
export JWT_SECRET="test_secret_key_for_development_please_change_in_production_32chars"
export REFRESH_TOKEN_KEY="test_refresh_key_for_development_32chars_minimum"
export CORS_ORIGINS="http://localhost:8080,http://localhost:3000"
export PORT="8083"
export MIGRATIONS_PATH="file://./migrations"
./server
