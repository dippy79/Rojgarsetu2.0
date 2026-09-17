@echo off
set DATABASE_URL=postgres://amitsharma:amitsharma@127.0.0.1:5435/rojgarsetu2?sslmode=disable
set REDIS_URL=redis://localhost:6379
set JWT_SECRET=f7e3c9a1b2d3e4f5a6b7c8d9e0f1a2b3c4d5e6f7a8b9c0d1e2f3a4b5c6d7e8f9
set REFRESH_TOKEN_KEY=a1b2c3d4e5f6a1b2c3d4e5f6a1b2c3d4e5f6a1b2c3d4e5f6a1b2c3d4e5f6a1b2
set CORS_ORIGINS=http://localhost:8080,http://localhost:3000,http://localhost:3001
set PORT=8084
set MIGRATIONS_PATH=file://./migrations
set JWT_ISSUER=rojgarsetu-backend
set JWT_AUDIENCE=rojgarsetu-api
server.exe
