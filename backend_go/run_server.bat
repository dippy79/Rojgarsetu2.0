@echo off
if "%DATABASE_URL%"=="" set DATABASE_URL=postgres://postgres:postgres@127.0.0.1:5435/rojgarsetu2?sslmode=disable
if "%REDIS_URL%"=="" set REDIS_URL=redis://localhost:6380
if "%JWT_SECRET%"=="" set JWT_SECRET=super-secret-jwt-key-minimum-32-characters-long
if "%REFRESH_TOKEN_KEY%"=="" set REFRESH_TOKEN_KEY=super-secret-refresh-key-minimum-32-chars
if "%CORS_ORIGINS%"=="" set CORS_ORIGINS=http://localhost:8080,http://localhost:3000,http://localhost:3001
if "%PORT%"=="" set PORT=8083
if "%MIGRATIONS_PATH%"=="" set MIGRATIONS_PATH=file://./migrations
go run cmd/server/main.go
