# RUNBOOK: Crawler Blocked (403/429)

## Symptoms
- crawler_requests_blocked_total >10/5m (Prometheus)
- Logs show "Crawler blocked" domain=naukri.com

## Steps
1. **Immediate Pause**
   ```
   docker-compose scale crawler-service=0
   ```

2. **Diagnose**
   - Check robots.txt: curl domain/robots.txt | grep Disallow
   - Test manual: curl -A "Mozilla/5.0..." domain/jobs

3. **Mitigate**
   - Rotate UA list in base.go (add 5 new).
   - Enable proxies: env CRAWLER_PROXY_ENABLED=true
   - Scale workers: CRAWLER_WORKERS=2
   - docker-compose up -d crawler-service

4. **Prevent**
   - Monitor block rate alert.
   - Weekly proxy rotation cron.

5. **Rollback**
   ```
   docker-compose down crawler-service
   git revert crawler changes
   ```

**SLA: Resolve <1h. Escalate to proxy provider if >3 domains affected.**

