package sources

import (
	"context"
	"io"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/rs/zerolog/log"
)

// RRBSource scrapes jobs from Railway Recruitment Board
type RRBSource struct {
	BaseSource
	client *http.Client
}

// NewRRBSource creates a new RRB source
func NewRRBSource() *RRBSource {
	return &RRBSource{
		BaseSource: BaseSource{NameStr: "rrb", BaseURL: "https://www.rrbcdg.gov.in"},
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
func (s *RRBSource) Fetch(ctx context.Context) ([]GovJobSource, error) {
	log.Info().Msg("Starting crawl for source: RRB (Railway Recruitment Board)")

	var allJobs []GovJobSource

	for name, baseURL := range RRBURLs {
		jobs, err := s.fetchFromRRB(ctx, name, baseURL)
		if err != nil {
			log.Warn().Err(err).Str("rrb", name).Msg("Failed to fetch from RRB")
			continue
		}
		allJobs = append(allJobs, jobs...)
	}

	log.Info().Int("totalJobs", len(allJobs)).Msg("RRB fetch completed")
	return allJobs, nil
}

// fetchFromRRB fetches from a specific RRB website
func (s *RRBSource) fetchFromRRB(ctx context.Context, name, baseURL string) ([]GovJobSource, error) {
	urls := []string{
		baseURL + "/notice.php",
		baseURL + "/latest-news.php",
		baseURL + "/vacancy.php",
		baseURL + "/career.php",
	}

	var jobs []GovJobSource

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
func (s *RRBSource) parseHTMLJobs(html, baseURL, rrbName string) []GovJobSource {
	var jobs []GovJobSource

	patterns := []string{
		`<a[^>]*href="(/[^"]*notice[^"]*)"[^>]*>([^<]*)</a>`,
		`<a[^>]*href="(/[^"]*vacancy[^"]*)"[^>]*>([^<]*)</a>`,
		`<a[^>]*href="(/[^"]*advertisement[^"]*)"[^>]*>([^<]*)</a>`,
		`<a[^>]*href="(/[^"]*recruitment[^"]*)"[^>]*>([^<]*)</a>`,
	}

	for _, pattern := range patterns {
		matches := extractMatches(html, pattern)
		for _, match := range matches {
			if len(match) >= 3 {
				title := strings.TrimSpace(match[2])
				link := match[1]

				if isRRBRelevant(title) {
					job := GovJobSource{
						Source:    "rrb_" + rrbName,
						Title:     title,
						ApplyURL:  baseURL + link,
						CreatedAt: time.Now(),
					}
					if isValidJob(&job) {
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
