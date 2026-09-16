@echo off
set DATABASE_URL=postgres://amitsharma:amitsharma@localhost:5432/rojgarsetu2?sslmode=disable
set REDIS_URL=redis://localhost:6379
set JWT_SECRET=test_secret_key_for_development_please_change_in_production_32chars
set REFRESH_TOKEN_KEY=test_refresh_key_for_development_32chars_minimum
set CORS_ORIGINS=http://localhost:8080,http://localhost:3000
set PORT=8083
set MIGRATIONS_PATH=file://./migrations
./server
