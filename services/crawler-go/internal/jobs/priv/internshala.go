package priv

import (
	"context"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/rojgarsetu/crawler/internal/shared"
	"github.com/rs/zerolog/log"
)

// InternshalaSource scrapes jobs from Internshala
type InternshalaSource struct {
	shared.BaseSource
	client *http.Client
}

// NewInternshalaSource creates a new Internshala source
func NewInternshalaSource() *InternshalaSource {
	return &InternshalaSource{
		BaseSource: shared.BaseSource{NameStr: "internshala", BaseURL: "https://internshala.com"},
		client: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// Fetch retrieves private jobs/internships from Internshala
func (s *InternshalaSource) Fetch(ctx context.Context) ([]shared.PrivJobSource, error) {
	log.Info().Msg("Starting crawl for source: Internshala")

	urls := []string{
		"https://internshala.com/jobs/matching-preferences",
		"https://internshala.com/internships",
	}

	var allJobs []shared.PrivJobSource

	for _, url := range urls {
		req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
		if err != nil {
			continue
		}

		shared.SetUserAgentAndCheck(req, s.BaseURL)

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

		jobs := s.parseHTMLJobs(string(body))
		allJobs = append(allJobs, jobs...)
	}

	log.Info().Int("totalJobs", len(allJobs)).Msg("Internshala fetch completed")
	return allJobs, nil
}

// parseHTMLJobs parses jobs from Internshala HTML
func (s *InternshalaSource) parseHTMLJobs(html string) []shared.PrivJobSource {
	var jobs []shared.PrivJobSource

	// Internshala job patterns
	pattern := `<a[^>]*href="(/job/detail/[^"]*)"[^>]*>([^<]*)</a>`
	matches := shared.ExtractMatches(html, pattern)

	for _, match := range matches {
		if len(match) >= 3 {
			title := strings.TrimSpace(match[2])
			link := match[1]

			if len(title) > 3 {
				job := shared.PrivJobSource{
					Source:    "internshala",
					Title:     title,
					URL:       "https://internshala.com" + link,
					CreatedAt: time.Now(),
					JobType:   "internship",
				}
				job.Company = "Unspecified Company"

				if shared.IsValidPrivJob(&job) {
					jobs = append(jobs, job)
				}
			}
		}
	}

	return jobs
}

// Name returns the source name
func (s *InternshalaSource) Name() string {
	return s.NameStr
}
