package shared

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"strings"
	"sync"
	"time"
	"unicode"

	"github.com/PuerkitoBio/goquery"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/rs/zerolog/log"
	"golang.org/x/time/rate"
)

var (
	crawlerBlocksTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "crawler_requests_blocked_total",
			Help: "Total blocked crawler requests (403/429/robots.txt)",
		},
		[]string{"reason", "domain"},
	)

	adaptiveThrottleActive = promauto.NewGauge(
		prometheus.GaugeOpts{
			Name: "crawler_adaptive_throttle_active",
			Help: "Whether adaptive throttling is currently active (1) or not (0)",
		},
	)

	defaultDomainLimiter = NewDomainLimiter()
)

// BaseSource common functionality for all sources
type BaseSource struct {
	NameStr string
	BaseURL string
}

// UserAgentRotator rotates through realistic UAs
type UserAgentRotator struct {
	mu     sync.Mutex
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
	r.mu.Lock()
	defer r.mu.Unlock()
	agent := r.agents[r.index]
	r.index = (r.index + 1) % len(r.agents)
	return agent
}

// DomainLimiter throttles per domain, with adaptive backoff on repeated blocks
type DomainLimiter struct {
	mu                 sync.RWMutex
	limiters           map[string]*rate.Limiter
	throttleMultiplier float64
	sem                chan struct{}
}

func NewDomainLimiter() *DomainLimiter {
	return &DomainLimiter{
		limiters:           make(map[string]*rate.Limiter),
		throttleMultiplier: 1.0,
		sem:                make(chan struct{}, 100),
	}
}

func (dl *DomainLimiter) Allow(domain string) bool {
	select {
	case dl.sem <- struct{}{}:
		defer func() { <-dl.sem }()
	default:
		return false
	}

	dl.mu.RLock()
	limiter, exists := dl.limiters[domain]
	dl.mu.RUnlock()

	if !exists {
		dl.mu.Lock()
		limiter, exists = dl.limiters[domain]
		if !exists {
			limiter = rate.NewLimiter(rate.Limit(10*dl.throttleMultiplier), int(20*dl.throttleMultiplier))
			dl.limiters[domain] = limiter
		}
		dl.mu.Unlock()
	}

	return limiter.Allow()
}

func (dl *DomainLimiter) applyBackoff(domain string) {
	dl.mu.Lock()
	defer dl.mu.Unlock()

	dl.throttleMultiplier *= 0.5
	if dl.throttleMultiplier < 0.1 {
		dl.throttleMultiplier = 0.1
	}

	newRate := 10 * dl.throttleMultiplier
	newBurst := int(20 * dl.throttleMultiplier)
	dl.limiters[domain] = rate.NewLimiter(rate.Limit(newRate), newBurst)
	adaptiveThrottleActive.Set(1)
}

// ---- method-based API (BaseSource) ----

func (b *BaseSource) CheckRobotsTxt(ctx context.Context, path string) bool {
	robotsURL := b.BaseURL + "/robots.txt"
	if !strings.HasPrefix(path, "/") {
		path = "/" + path
	}

	req, _ := http.NewRequestWithContext(ctx, "GET", robotsURL, nil)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Warn().Err(err).Str("robots_url", robotsURL).Msg("Failed to fetch robots.txt")
		return true
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	robots := string(body)

	if strings.Contains(robots, "User-agent: *") && strings.Contains(robots, fmt.Sprintf("Disallow: %s", path)) {
		return false
	}
	return true
}

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
		log.Warn().Int("status", resp.StatusCode).Str("domain", domain).Msg("Crawler blocked")
		dl.applyBackoff(domain)
		time.Sleep(5 * time.Minute)
		return nil, fmt.Errorf("blocked: %d", resp.StatusCode)
	}

	return resp, nil
}

// ---- free-function API (kept for ncs.go / ssc.go compatibility) ----

// SetUserAgentAndCheck sets a rotated User-Agent header on req.
func SetUserAgentAndCheck(req *http.Request, baseURL string) {
	uaRotator := NewUserAgentRotator()
	req.Header.Set("User-Agent", uaRotator.Next())
}

// CheckRobotsTxt (free function) respects robots.txt for a given base URL + path.
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
		return true
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	robots := string(body)

	if strings.Contains(robots, "User-agent: *") && strings.Contains(robots, "Disallow: "+path) {
		crawlerBlocksTotal.WithLabelValues("robots", baseURL).Inc()
		return false
	}
	return true
}

// CheckStatusAndPause checks a response for block signals (403/429) and pauses if blocked.
func CheckStatusAndPause(resp *http.Response, domain string) error {
	if resp.StatusCode == 403 || resp.StatusCode == 429 {
		crawlerBlocksTotal.WithLabelValues("blocked", domain).Inc()
		log.Warn().Int("status", resp.StatusCode).Str("domain", domain).Msg("Crawler blocked")
		defaultDomainLimiter.applyBackoff(domain)
		time.Sleep(5 * time.Minute)
		return fmt.Errorf("blocked: %d", resp.StatusCode)
	}
	return nil
}

// SanitizeString strips HTML tags using goquery, trims whitespace, normalizes
// whitespace runs, and limits length to maxLen. Uses goquery (already vendored
// for naukri.go) rather than a regex `<[^>]*>` approach, because regex-based
// tag stripping fails on malformed/nested HTML fragments commonly found in
// scraped web data.
func SanitizeString(input string, maxLen int) string {
	if input == "" {
		return ""
	}

	// Use goquery to parse HTML and extract clean text content.
	// goquery.NewDocumentFromReader internally uses html.Parse which handles
	// malformed/nested HTML correctly.
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(input))
	if err != nil {
		// If parsing fails (e.g. binary data), fall back to a simple
		// whitespace-trimmed version of the input.
		return strings.TrimSpace(input)
	}

	// Extract all text nodes, joined with spaces.
	text := doc.Text()

	// Normalize whitespace: collapse runs of whitespace to single space.
	// Also trim leading/trailing whitespace.
	var b strings.Builder
	b.Grow(len(text))
	inSpace := false
	for _, r := range text {
		if unicode.IsSpace(r) {
			if !inSpace {
				b.WriteRune(' ')
				inSpace = true
			}
		} else {
			b.WriteRune(r)
			inSpace = false
		}
	}
	result := strings.TrimSpace(b.String())

	// Limit length to maxLen. If truncated, try to cut at a word boundary.
	if maxLen > 0 && len(result) > maxLen {
		truncated := result[:maxLen]
		// Back up to the last space to avoid cutting mid-word (if possible).
		if lastSpace := strings.LastIndex(truncated, " "); lastSpace > maxLen/2 {
			truncated = truncated[:lastSpace]
		}
		result = truncated + "..."
	}

	return result
}
