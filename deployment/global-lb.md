# Global Load Balancing - Route53 Geo-Routing

## Configuration
- Primary: Latency-based routing to regional ELBs (us-east-1, europe-west1, ap-south-1)
- Failover: HealthCheck /health, reroute if >50% pods unhealthy
- DNS: rojgarsetu.com -> geo records

## Cloudflare Alternative
Use Workers for latency routing.

## Verification
dig us.rojgarsetu.com @8.8.8.8
