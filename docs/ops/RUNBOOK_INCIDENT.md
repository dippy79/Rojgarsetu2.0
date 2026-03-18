# RUNBOOK: Critical Incidents

## 1. DB Outage (pg_up=0)

**Symptoms:** Prometheus pg_up=0, query timeouts, empty job tables

**Steps:**
1. Triage: `kubectl get pods -n rojgarsetu | grep postgres`
2. Logs: `kubectl logs -n rojgarsetu postgres-0`
3. Pause traffic: `kubectl scale deployment backend --replicas=0`
4. Restore: Run db-dry-run-restore.ps1 → if OK, full pg_restore from S3 WAL-G
5. Verify: `SELECT COUNT(*) FROM gov_jobs UNION SELECT COUNT(*) FROM priv_jobs`
6. Resume: `kubectl scale deployment backend --replicas=2`

**SLA:** RTO <30min. Escalate if WAL gap >4h.

## 2. TLS Expiry

**Symptoms:** TLSHandshakeFailures alert, 525/443 errors in nginx logs

**Steps:**
1. Check: `kubectl get secret tls-secret -o yaml | grep expiry`
2. Renew: `certbot renew` or manual Let's Encrypt
3. Update: `kubectl create secret tls tls-secret --cert=new.crt --key=new.key`
4. Rollout: `kubectl rollout restart ingress/rojgar-ingress`
5. Verify: `curl -k https://rojgarsetu.com/health`

**SLA:** <15min. Auto-renew cron weekly.

## 3. Crawler Block Surge

**Symptoms:** CrawlerBlocksHigh >10/5m, jobs stale >24h

**Steps:**
1. Pause: `kubectl scale deployment crawler --replicas=0`
2. Diagnose: Prometheus rate(crawler_requests_blocked_total[5m])
3. Mitigate: Rotate UA/proxies, check robots.txt per domain
4. Resume: `--replicas=2`, monitor adaptive throttle metric
5. Prevent: Monthly UA update, proxy pool refresh

**SLA:** Crawling resumes <1h. Escalate to legal if IP block.
