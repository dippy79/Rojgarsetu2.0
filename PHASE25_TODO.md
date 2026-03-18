# PHASE25_TODO.md - POST-DEPLOYMENT MONITORING & SCALING

## Approved Plan Steps (Implement Sequentially, Update on Complete)

### 1. Create Monitoring Provisioning Files
- [x] monitoring/grafana/provisioning/datasources/prometheus.yml (auto Prometheus datasource)
- [x] monitoring/grafana/provisioning/dashboards/requests.json (rate, latency, errors)
- [x] monitoring/grafana/provisioning/dashboards/db-pool.json (DB conn pool)
- [x] monitoring/grafana/provisioning/dashboards/tls-health.json (TLS handshake)
- [x] monitoring/grafana/provisioning/dashboards/containers.json (health/status)

### 2. Add Alert Rules
- [x] monitoring/prometheus/rules.yml (errors>5%, DB refused, restarts, TLS fail, crawler blocks>10%)

### 3. Update Prometheus Config
- [x] monitoring/prometheus/prometheus.yml (load rules, postgres_exporter, redis_exporter scrape)

### 4. Security Hardening
- [x] Rotate JWT_SECRET in deployment/docker-compose.yml (64-char random)
- [ ] Add .env mount ro to services
- [ ] Test SQLi fails

### 5. Crawler Anti-Block
- [ ] Add UA rotation, backoff, throttle (10req/s), 403/429 pause, robots.txt in crawler_go
- [ ] Alert on blocks >10%/5m

### 6. Load Test
- [ ] deployment/loadtest/k6.js (1000 VU login/refresh/logout)

### 7. Update docker-compose.yml
- [ ] Add postgres_exporter, redis_exporter services
- [ ] Grafana provisioning vol: - ../monitoring/grafana/provisioning:/etc/grafana/provisioning
- [ ] JWT updates

### 8. Deploy & Verify (STEP 5 commands)
- [ ] Build/up
- [ ] Curl health/live/ready/metrics
- [ ] Logs grep
- [ ] Grafana login, check dashboards/datasource/alerts
- [ ] Run k6 loadtest
- [ ] SQLi tests

### 9. Produce Report → attempt_completion

**Progress: 0/9**
