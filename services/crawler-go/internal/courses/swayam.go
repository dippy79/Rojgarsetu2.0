package courses

import (
	"github.com/rojgarsetu/crawler/internal/shared"
	"context"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/rs/zerolog/log"
)

// SWAYAMSource scrapes courses from SWAYAM
type SWAYAMSource struct {
	shared.BaseSource
	client *http.Client
	apiURL string
}

// NewSWAYAMSource creates a new SWAYAM source
func NewSWAYAMSource() *SWAYAMSource {
	return &SWAYAMSource{
		BaseSource: shared.BaseSource{NameStr: "swayam", BaseURL: "https://swayam.gov.in"},
		client: &http.Client{
			Timeout: 30 * time.Second,
		},
		apiURL: "https://swayam.gov.in/",
	}
}

// Fetch retrieves courses from SWAYAM.
//
// Verified live endpoint status (Phase A):
//   - /explore      → 404 (dead)
//   - /api/courses  → 404 (dead)
//   - /             → 200, but a static landing page with no JSON-LD, no course
//     anchors, and no SPA data-bootstrap. Course catalog is loaded client-side
//     (JS-rendered) from an internal/CDN API (assets served from apis.com).
//
// Decision point (per approved plan): fixing SWAYAM requires JS rendering
// (chromedp) or an official data-access agreement. We deliberately do NOT wire
// chromedp into this source because that would add browser-pool contention
// with Naukri. Returning 0 courses with a clear log so RunSummary is honest.
func (s *SWAYAMSource) Fetch(ctx context.Context) ([]shared.CourseSource, error) {
	log.Info().Msg("Starting crawl for source: SWAYAM")

	// Single reachability probe so RunSummary logs reflect a live check.
	probe, err := s.fetchHomeForData(ctx)
	if err != nil {
		log.Warn().Err(err).Msg("SWAYAM home fetch failed")
		return nil, fmt.Errorf("failed to fetch from SWAYAM: %w", err)
	}

	if probe {
		log.Warn().Msg("SWAYAM course catalog is JS-rendered / requires internal API - flagged as decision point, returning 0 courses")
	} else {
		log.Warn().Msg("SWAYAM home page contained no embedded course data - flagged as decision point, returning 0 courses")
	}

	return []shared.CourseSource{}, nil
}

// fetchHomeForData fetches the SWAYAM home page and reports whether any
// directly-parseable course data (JSON-LD blocks or /course/ anchors) exists.
func (s *SWAYAMSource) fetchHomeForData(ctx context.Context) (bool, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", s.apiURL, nil)
	if err != nil {
		return false, err
	}
	req.Header.Set("User-Agent", "RojgarSetu/2.0")

	resp, err := s.client.Do(req)
	if err != nil {
		return false, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return false, fmt.Errorf("SWAYAM home returned status: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return false, err
	}
	html := string(body)

	if len(shared.ExtractMatches(html, `<script[^>]*type="application/ld\+json"[^>]*>`)) > 0 {
		return true, nil
	}
	if len(shared.ExtractMatches(html, `href="[^"]*/course/[^"]*"`)) > 0 {
		return true, nil
	}
	if strings.Contains(html, "__NEXT_DATA__") || strings.Contains(html, "window.__INITIAL") {
		return true, nil
	}

	return false, nil
}

// Name returns the source name
func (s *SWAYAMSource) Name() string {
	return s.NameStr
}

