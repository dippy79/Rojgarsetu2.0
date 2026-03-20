#!/bin/bash
# RojgarSetu Production Deployment Script for VPS/Linux

set -e

echo "🚀 RojgarSetu Production Deploy..."

# 1. Load .env
if [ -f .env.production ]; then
  export $(grep -v '^#' .env.production | xargs)
else
  echo "❌ .env.production not found. Copy .env.production.example and edit."
  exit 1
fi

# 2. Pull latest images
docker compose pull

# 3. Stop existing containers
docker compose down

# 4. Start services (migrate first, then backend)
docker compose up -d postgres migrate backend

# 5. Wait for services
echo "⏳ Waiting for services..."
sleep 10

# 6. Check health
docker compose ps

if ! docker compose ps | grep -q "healthy"; then
  echo "❌ Some services not healthy"
  docker compose logs
  exit 1
fi

# 7. Final verification
if curl -f http://localhost:8080/robots.txt; then
  echo "✅ Backend healthy - /robots.txt OK"
else
  echo "❌ Backend not responding"
  docker compose logs backend
  exit 1
fi

echo "🎉 Deployment complete! Backend: http://localhost:8080"
echo "Run './test.sh' for full smoke tests."
