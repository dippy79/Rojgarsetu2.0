# RojgarSetu PHASE H Production Checklist

| Phase | Status | Notes |
|-------|--------|-------|
| [x] VPS provisioned (min 2GB RAM, 20GB disk) | Run deployment/vps_setup.sh | Ubuntu 22.04 |
| [ ] Domain pointed to VPS IP | A record api.yourdomain.com | |
| [x] .env.production filled with real secrets | cp .env.production.example .env.production | JWT_SECRET, POSTGRES_PASSWORD |
| [ ] SSL certificate installed | sudo certbot --nginx -d api.yourdomain.com | Uses deployment/nginx_vps.conf |
| [ ] Firewall configured | UFW 22/80/443 open | vps_setup.sh |
| [ ] Docker running | docker compose ps | |
| [ ] Migrations run | docker compose up -d migrate | 000001-000008 |
| [ ] Backend healthy | curl -f http://localhost:8080/health | {status:ok, db:connected} |
| [x] Flutter APK built with production URL | API_BASE_URL=https://api.yourdomain.com flutter build apk --release | deployment/flutter_release_instructions.md |
| [ ] Test register + login on production | Postman/cURL | Rate limiting 5/min login |
| [ ] Rate limiting working on production | 10 req/s Nginx + backend 60/min/IP | 429 expected |

## Deploy Commands
```bash
./deploy.sh
docker compose ps
curl localhost:8080/health
```

**All checks passed → PRODUCTION READY! 🚀**

