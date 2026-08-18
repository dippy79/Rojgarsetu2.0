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

// ShineSource scrapes jobs from Shine.com
type ShineSource struct {
	shared.BaseSource
	client *http.Client
}

// NewShineSource creates a new Shine source
func NewShineSource() *ShineSource {
	return &ShineSource{
		BaseSource: shared.BaseSource{NameStr: "shine", BaseURL: "https://www.shine.com"},
		client: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// Fetch retrieves private jobs from Shine.com
func (s *ShineSource) Fetch(ctx context.Context) ([]shared.PrivJobSource, error) {
	log.Info().Msg("Starting crawl for source: Shine")

	searchURLs := []string{
		"https://www.shine.com/job-search/software-jobs",
		"https://www.shine.com/job-search/it-jobs",
	}

	var allJobs []shared.PrivJobSource

	for _, url := range searchURLs {
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

	log.Info().Int("totalJobs", len(allJobs)).Msg("Shine fetch completed")
	return allJobs, nil
}

// parseHTMLJobs parses jobs from Shine HTML
func (s *ShineSource) parseHTMLJobs(html string) []shared.PrivJobSource {
	var jobs []shared.PrivJobSource

	// Shine job card patterns
	pattern := `<a[^>]*href="(/jobs/[^"]*)"[^>]*>([^<]*)</a>`
	matches := shared.ExtractMatches(html, pattern)

	for _, match := range matches {
		if len(match) >= 3 {
			title := strings.TrimSpace(match[2])
			link := match[1]

			if len(title) > 3 {
				job := shared.PrivJobSource{
					Source:    "shine",
					Title:     title,
					URL:       "https://www.shine.com" + link,
					CreatedAt: time.Now(),
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
func (s *ShineSource) Name() string {
	return s.NameStr
}
