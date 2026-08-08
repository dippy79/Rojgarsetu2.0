package gov

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/rojgarsetu/crawler/internal/shared"
	"github.com/rs/zerolog/log"
)

// RRBSource scrapes jobs from Railway Recruitment Board
type RRBSource struct {
	shared.BaseSource
	client *http.Client
}

// NewRRBSource creates a new RRB source
func NewRRBSource() *RRBSource {
	return &RRBSource{
		BaseSource: shared.BaseSource{NameStr: "rrb", BaseURL: "https://www.rrbcdg.gov.in"},
		client: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// RRBURLs contains all RRB websites
var RRBURLs = map[string]string{
	"rrbcdg":          "https://www.rrbcdg.gov.in",
	"rrbmumbai":       "https://www.rrbmumbai.gov.in",
	"rrbsecunderabad": "https://www.rrbsecunderabad.nic.in",
	"rrbahmedabad":    "https://www.rrbahmedabad.gov.in",
	"rrbchandigarh":   "https://www.rrbchandigarh.gov.in",
	"rrbchennai":      "https://www.rrbchennai.gov.in",
	"rrbguwahati":     "https://www.rrbguwahati.gov.in",
	"rrbkolkata":      "https://www.rrbkolkata.gov.in",
	"rrbpatna":        "https://www.rrbpatna.gov.in",
	"rrbranchi":       "https://www.rrbranchi.gov.in",
	"rrbsiliguri":     "https://www.rrbsiliguri.gov.in",
	"rrbmuzaffarpur":  "https://www.rrbmuzaffarpur.gov.in",
	"rrbbhopal":       "https://www.rrbbhopal.gov.in",
	"rrbajmer":        "https://www.rrbajmer.gov.in",
	"rrbbhubaneswar":  "https://www.rrbbhubaneswar.gov.in",
}

// Fetch retrieves government jobs from all RRB sources
func (s *RRBSource) Fetch(ctx context.Context) ([]shared.GovJobSource, error) {
	log.Info().Msg("Starting crawl for source: RRB (Railway Recruitment Board)")

	// FLAG (Phase A first-pass): All 16 RRB regional boards are unreachable
	// from this environment via plain HTTP — the site roots return connection
	// errors (blocked DNS / TLS handshake / WAF), and the legacy paths this
	// fetcher tried (notice.php, vacancy.php, latest-news.php, career.php)
	// return 403/404. Each board runs its own independently hosted site with
	// no shared RSS/JSON feed, so fixing them one-by-one would require
	// per-board URL + selector verification against live sites that are
	// currently not reachable from here.
	//
	// Rather than silently returning 0 jobs, we surface a clear diagnostic so
	// the RunSummary flags this source as needing a different approach (e.g.
	// a maintained aggregate feed of RRB CEN notifications, or a browser-driven
	// crawl of indianrailways.gov.in / the RRB page on sarkari jobs portals).
	// This is the "report status" first-pass; no code change here resurrects it.
	err := fmt.Errorf("RRB regional board sites are unreachable from this environment (all 16 boards: connection errors / 403 / 404 on legacy paths). No shared RSS/JSON feed exists. Requires a different approach (aggregate feed or browser-driven crawl of RRB CEN notifications).")
	log.Warn().Msg(err.Error())
	return nil, err
}

// fetchFromRRB fetches from a specific RRB website
func (s *RRBSource) fetchFromRRB(ctx context.Context, name, baseURL string) ([]shared.GovJobSource, error) {
	urls := []string{
		baseURL + "/notice.php",
		baseURL + "/latest-news.php",
		baseURL + "/vacancy.php",
		baseURL + "/career.php",
	}

	var jobs []shared.GovJobSource

	for _, url := range urls {
		req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
		if err != nil {
			continue
		}

		req.Header.Set("User-Agent", "RojgarSetu/2.0")

		resp, err := s.client.Do(req)
		if err != nil {
			continue
		}

		if resp.StatusCode != http.StatusOK {
			resp.Body.Close()
			continue
		}

		body, err := io.ReadAll(resp.Body)
		resp.Body.Close()
		if err != nil {
			continue
		}

		rrbJobs := s.parseHTMLJobs(string(body), baseURL, name)
		jobs = append(jobs, rrbJobs...)
	}

	log.Info().Int("jobs", len(jobs)).Str("rrb", name).Msg("RRB fetch successful")
	return jobs, nil
}

// parseHTMLJobs parses jobs from RRB HTML page
func (s *RRBSource) parseHTMLJobs(html, baseURL, rrbName string) []shared.GovJobSource {
	var jobs []shared.GovJobSource

	patterns := []string{
		`<a[^>]*href="(/[^"]*notice[^"]*)"[^>]*>([^<]*)</a>`,
		`<a[^>]*href="(/[^"]*vacancy[^"]*)"[^>]*>([^<]*)</a>`,
		`<a[^>]*href="(/[^"]*advertisement[^"]*)"[^>]*>([^<]*)</a>`,
		`<a[^>]*href="(/[^"]*recruitment[^"]*)"[^>]*>([^<]*)</a>`,
	}

	for _, pattern := range patterns {
		matches := shared.ExtractMatches(html, pattern)
		for _, match := range matches {
			if len(match) >= 3 {
				title := strings.TrimSpace(match[2])
				link := match[1]

				if isRRBRelevant(title) {
					job := shared.GovJobSource{
						Source:    "rrb_" + rrbName,
						Title:     title,
						ApplyURL:  baseURL + link,
						CreatedAt: time.Now(),
					}
					if shared.IsValidJob(&job) {
						jobs = append(jobs, job)
					}
				}
			}
		}
	}

	return jobs
}

// Name returns the source name
func (s *RRBSource) Name() string {
	return s.NameStr
}

// isRRBRelevant checks if notice is relevant for RRB
func isRRBRelevant(title string) bool {
	titleLower := strings.ToLower(title)
	relevant := []string{"vacancy", "recruitment", "notification", "advertisement", "rrb", "railway", "CEN"}
	irrelevant := []string{"result", "answer key", "cutoff", " CBT", "exam date"}

	for _, irr := range irrelevant {
		if strings.Contains(titleLower, irr) {
			return false
		}
	}

	for _, rel := range relevant {
		if strings.Contains(titleLower, rel) {
			return true
		}
	}

	return false
}

// extractVacancyCount extracts vacancy count from text
func extractVacancyCount(text string) *int {
	re := regexp.MustCompile(`(\d+)\s*(?:vacancy|post|position)`)
	matches := re.FindStringSubmatch(text)
	if len(matches) > 1 {
		if vc, err := strconv.Atoi(matches[1]); err == nil {
			return &vc
		}
	}
	return nil
}
