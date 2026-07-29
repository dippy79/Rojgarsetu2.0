package sources

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/rojgarsetu/crawler/internal/logger"
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

// BaseSource common functionality for all sources
type BaseSource struct {
	NameStr string
	BaseURL string
}

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
			"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.124 Safari/537.36 Edg/91.0.864.59",
			"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/17.1 Safari/605.1.15",
			"Mozilla/5.0 (X11; Ubuntu; Linux x86_64; rv:123.0) Gecko/20100101 Firefox/123.0",
			"Mozilla/5.0 (iPhone; CPU iPhone OS 17_3 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/17.3 Mobile/15E148 Safari/604.1",
			"Mozilla/5.0 (Linux; Android 14; SM-S928B) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/121.0.0.0 Mobile Safari/537.36",
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
	mu       chan struct{}
}

func NewDomainLimiter() *DomainLimiter {
	return &DomainLimiter{
		limiters: make(map[string]*rate.Limiter),
		mu:       make(chan struct{}, 100), // concurrency limit
	}
}

func (dl *DomainLimiter) Allow(domain string) bool {
	select {
	case dl.mu <- struct{}{}:
		defer func() { <-dl.mu }()
	default:
		return false
	}

	dl.limitersMu.Lock()
	defer dl.limitersMu.Unlock()

	limiter, exists := dl.limiters[domain]
	if !exists {
		limiter = rate.NewLimiter(rateLimit, burst) // adaptive
		dl.limiters[domain] = limiter
	}
	return limiter.Allow()
}

// CheckRobotsTxt respects robots.txt
func (b *BaseSource) CheckRobotsTxt(ctx context.Context, path string) bool {
	robotsURL := b.BaseURL + "/robots.txt"
	if !strings.HasPrefix(path, "/") {
		path = "/" + path
	}

	req, _ := http.NewRequestWithContext(ctx, "GET", robotsURL, nil)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		logger.Log.Warn().Err(err).Str("robots_url", robotsURL).Msg("Failed to fetch robots.txt")
		return true // allow if can't check
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	robots := string(body)

	// Simple check: if User-agent: * Disallow: /path
	if strings.Contains(robots, "User-agent: *") && strings.Contains(robots, fmt.Sprintf("Disallow:%s", path)) {
		return false
	}
	return true
}

// DoRequest with safeguards
func (b *BaseSource) DoRequest(ctx context.Context, req *http.Request, domain string, dl *DomainLimiter) (*http.Response, error) {
	if !dl.Allow(domain) {
		crawlerBlocksTotal.WithLabelValues("throttle", domain).Inc()
		return nil, fmt.Errorf("throttled")
	}

	uaRotator := NewUserAgentRotator()
	req.Header.Set("User-Agent", uaRotator.Next())

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode == 403 || resp.StatusCode == 429 {
		crawlerBlocksTotal.WithLabelValues("blocked", domain).Inc()
		logger.Log.Warn().Int("status", resp.StatusCode).Str("domain", domain).Msg("Crawler blocked")
		adaptiveThrottleActive.Set(1) // metric active=true
		// Adaptive: halve rate if >5% blocks (simplified check)
		dl.throttleMultiplier = dl.throttleMultiplier * 0.5
		if dl.throttleMultiplier < 0.1 {
			dl.throttleMultiplier = 0.1
		}
		rateLimit := 10 * dl.throttleMultiplier
		burst := int(20 * dl.throttleMultiplier)
		dl.limiters[domain] = rate.NewLimiter(rate.Limit(rateLimit), burst)
		time.Sleep(5 * time.Minute)
		return nil, fmt.Errorf("blocked: %d", resp.StatusCode)
	}

	return resp, nil
}
