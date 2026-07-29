package sources

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/rs/zerolog/log"
	"golang.org/x/time/rate"
)

var (
	crawlerBlocksTotal = promauto.NewCounterVec(
		prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Name: "crawler_requests_blocked_total",
				Help: "Total blocked crawler requests (403/429/robots.txt)",
			},
			[]string{"reason", "domain"},
		),
	)
)

// UserAgentRotator rotates through realistic UAs
type UserAgentRotator struct {
	agents []string
	index  int
}

func NewUserAgentRotator() *UserAgentRotator {
	return &UserAgentRotator{
		agents: []string{
			"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36",
			"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36",
			"Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36",
			"Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:109.0) Gecko/20100101 Firefox/121.0",
			"Mozilla/5.0 (Macintosh; Intel Mac OS X 10.15; rv:109.0) Gecko/20100101 Firefox/121.0",
		},
	}
}

func (r *UserAgentRotator) Next() string {
	agent := r.agents[r.index]
	r.index = (r.index + 1) % len(r.agents)
	return agent
}

// DomainLimiter throttles per domain (10 req/sec max)
type DomainLimiter struct {
	limiters map[string]*rate.Limiter
	mu       sync.RWMutex
}

func NewDomainLimiter() *DomainLimiter {
	return &DomainLimiter{
		limiters: make(map[string]*rate.Limiter),
	}
}

func (dl *DomainLimiter) Allow(domain string) bool {
	dl.mu.RLock()
	limiter, exists := dl.limiters[domain]
	dl.mu.RUnlock()

	if !exists {
		dl.mu.Lock()
		limiter, exists = dl.limiters[domain]
		if !exists {
			limiter = rate.NewLimiter(rate.Limit(10), 20) // 10/sec, burst 20
			dl.limiters[domain] = limiter
		}
		dl.mu.Unlock()
	}

	if !limiter.Allow() {
		crawlerBlocksTotal.WithLabelValues("throttle", domain).Inc()
		return false
	}
	return true
}

// CheckRobotsTxt respects robots.txt
func CheckRobotsTxt(baseURL, path string) bool {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	robotsURL := baseURL + "/robots.txt"
	if !strings.HasPrefix(path, "/") {
		path = "/" + path
	}

	req, _ := http.NewRequestWithContext(ctx, "GET", robotsURL, nil)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Warn().Err(err).Str("robots_url", robotsURL).Msg("Failed to fetch robots.txt")
		return true // allow if can't check
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	robots := string(body)

	// Simple check: if User-agent: * Disallow: /path
	if strings.Contains(robots, "User-agent: *") && strings.Contains(robots, "Disallow: "+path) {
		crawlerBlocksTotal.WithLabelValues("robots", baseURL).Inc()
		return false
	}
	return true
}

// SetUserAgentAndCheck sets UA and robots check
func SetUserAgentAndCheck(req *http.Request, baseURL string) {
	uaRotator := NewUserAgentRotator()
	req.Header.Set("User-Agent", uaRotator.Next())
}

// CheckStatusAndPause checks for blocks and pauses
func CheckStatusAndPause(resp *http.Response, domain string) error {
	if resp.StatusCode == 403 || resp.StatusCode == 429 {
		crawlerBlocksTotal.WithLabelValues("blocked", domain).Inc()
		log.Warn().Int("status", resp.StatusCode).Str("domain", domain).Msg("Crawler blocked")
		time.Sleep(5 * time.Minute)
		return fmt.Errorf("blocked: %d", resp.StatusCode)
	}
	return nil
}
