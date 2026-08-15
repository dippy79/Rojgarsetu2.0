package internal

import (
	"io"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"
)

// PoliteCrawler implements robots.txt respect and rate limiting
type PoliteCrawler struct {
	UserAgent    string
	DelayMs      int
	RobotsStrict bool
	httpClient   *http.Client
}

// New creates a new PoliteCrawler
func New(userAgent string, delayMs int, strict bool) *PoliteCrawler {
	return &PoliteCrawler{
		UserAgent:    userAgent,
		DelayMs:      delayMs,
		RobotsStrict: strict,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// IsAllowed checks if robots.txt allows crawling the target URL
func (p *PoliteCrawler) IsAllowed(robotsURL, targetURL string) bool {
	if !p.RobotsStrict {
		return true
	}

	if robotsURL == "" {
		return true
	}

	// Fetch robots.txt
	resp, err := p.httpClient.Get(robotsURL)
	if err != nil {
		log.Printf("[WARN] failed to fetch robots.txt from %s: %v", robotsURL, err)
		return true // assume allowed if we can't fetch
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Printf("[WARN] robots.txt returned status %d from %s", resp.StatusCode, robotsURL)
		return true
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("[WARN] failed to read robots.txt: %v", err)
		return true
	}

	// Simple robots.txt parsing (user-agent specific)
	robotsTxt := string(body)
	targetPath := p.getPath(targetURL)

	// Check if our user-agent is disallowed
	lines := strings.Split(robotsTxt, "\n")
	var currentUA string
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(strings.ToLower(line), "user-agent:") {
			currentUA = strings.TrimSpace(line[11:])
		} else if strings.HasPrefix(strings.ToLower(line), "disallow:") && currentUA == "*" || currentUA == p.getUserAgentShort() {
			disallowedPath := strings.TrimSpace(line[9:])
			if strings.HasPrefix(targetPath, disallowedPath) {
				return false
			}
		}
	}

	return true
}

// Fetch makes a polite HTTP request with rate limiting
func (p *PoliteCrawler) Fetch(url string) (*http.Response, error) {
	time.Sleep(time.Duration(p.DelayMs) * time.Millisecond)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("User-Agent", p.UserAgent)
	req.Header.Set("Accept-Language", "en-IN,en;q=0.9")
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8")

	return p.httpClient.Do(req)
}

// getUserAgentShort returns short user-agent name for robots.txt matching
func (p *PoliteCrawler) getUserAgentShort() string {
	if strings.Contains(p.UserAgent, "RojgarSetuBot") {
		return "RojgarSetuBot"
	}
	return "*"
}

// getPath extracts the path from a URL
func (p *PoliteCrawler) getPath(rawURL string) string {
	u, err := url.Parse(rawURL)
	if err != nil {
		return "/"
	}
	return u.Path
}
