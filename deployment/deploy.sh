#!/bin/bash
# RojgarSetu Production Deployment Script v2.0 - VPS/Linux

set -e

echo "🚀 RojgarSetu PHASE H Production Deploy..."

cd "$(dirname "$0")"

# 1. Git pull latest
echo "📥 Git pull origin main..."
git pull origin main

# 2. Load/copy .env.production
echo "⚙️  Setup .env..."
if [ ! -f .env.production ]; then
  if [ -f .env.production.example ]; then
    cp .env.production.example .env.production
    echo "ℹ️  Copied .env.production.example → .env.production"
    echo "⚠️  EDIT .env.production NOW with REAL SECRETS before re-running!"
    echo "Required: JWT_SECRET (openssl rand -hex 32), POSTGRES_PASSWORD (20+ chars)"
    exit 1
  else
    echo "❌ .env.production.example missing! Create it first."
    exit 1
  fi
fi
export $(grep -v &#x27;^#&#x27; .env.production | grep -v &#x27;=&#x27; | xargs)

# 3. Docker compose pull latest images
echo "🐳 Docker compose pull..."
docker compose pull

# 4. Graceful stop (preserve volumes)
echo "🛑 Docker compose down..."
docker compose down

# 5. Start all services
echo "🚀 Docker compose up -d..."
docker compose up -d

# 6. Wait &amp; health check
echo "⏳ Waiting 30s for services..."
sleep 30

# 7. Verify status
echo "📊 Services status:"
docker compose ps

if ! docker compose ps | grep -q Up.*healthy; then
  echo "⚠️  Not all services healthy. Checking logs..."
  docker compose logs postgres | tail -10
  docker compose logs migrate | tail -10
  docker compose logs backend | tail -10
fi

# 8. Backend API health
echo "🔬 Backend /health check..."
if curl -f -s http://localhost:8080/health | grep -q &#x27;"status":"healthy&#x27;; then
  echo "✅ Backend healthy ✓"
else
  echo "❌ Backend /health failed!"
  echo "🔍 Logs:"
  docker compose logs backend | tail -20
  exit 1
fi

# 9. Postgres verification
echo "🗄️  Postgres check..."
if docker compose exec postgres pg_isready -U rojgar; then
  echo "✅ Postgres ready ✓"
else
  echo "❌ Postgres not ready"
  exit 1
fi

echo "🎉 PHASE H Production Deployment SUCCESSFUL!"
echo "🌐 Backend API: http://localhost:8080/health"
echo "📋 Full status: docker compose ps"
echo "📄 Logs: docker compose logs -f"
echo "🧪 Smoke test: ./test.sh"
echo "🔒 Nginx SSL: sudo certbot --nginx -d api.yourdomain.com"

