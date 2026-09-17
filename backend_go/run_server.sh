#!/bin/bash
# Use variables from .env if possible, otherwise use these defaults
export DATABASE_URL="postgres://amitsharma:amitsharma@127.0.0.1:5435/rojgarsetu2?sslmode=disable"
export REDIS_URL="redis://localhost:6379"
export JWT_SECRET="f7e3c9a1b2d3e4f5a6b7c8d9e0f1a2b3c4d5e6f7a8b9c0d1e2f3a4b5c6d7e8f9"
export REFRESH_TOKEN_KEY="a1b2c3d4e5f6a1b2c3d4e5f6a1b2c3d4e5f6a1b2c3d4e5f6a1b2c3d4e5f6a1b2"
export CORS_ORIGINS="http://localhost:8080,http://localhost:3000,http://localhost:3001"
export PORT="8083"
export MIGRATIONS_PATH="file://./migrations"
./server
