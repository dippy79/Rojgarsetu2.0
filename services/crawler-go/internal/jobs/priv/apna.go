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

// ApnaSource scrapes jobs from Apna.co
type ApnaSource struct {
	shared.BaseSource
	client *http.Client
}

// NewApnaSource creates a new Apna source
func NewApnaSource() *ApnaSource {
	return &ApnaSource{
		BaseSource: shared.BaseSource{NameStr: "apna", BaseURL: "https://apna.co"},
		client: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// Fetch retrieves private jobs from Apna
func (s *ApnaSource) Fetch(ctx context.Context) ([]shared.PrivJobSource, error) {
	log.Info().Msg("Starting crawl for source: Apna")

	searchURLs := []string{
		"https://apna.co/jobs-in-bangalore",
		"https://apna.co/jobs-in-mumbai",
		"https://apna.co/jobs-in-delhi",
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

	log.Info().Int("totalJobs", len(allJobs)).Msg("Apna fetch completed")
	return allJobs, nil
}

// parseHTMLJobs parses jobs from Apna HTML
func (s *ApnaSource) parseHTMLJobs(html string) []shared.PrivJobSource {
	var jobs []shared.PrivJobSource

	// Apna job card patterns (generalized)
	// Example: <a href="/job/software-engineer-123">Software Engineer</a>
	pattern := `<a[^>]*href="(/job/[^"]*)"[^>]*>([^<]*)</a>`
	matches := shared.ExtractMatches(html, pattern)

	for _, match := range matches {
		if len(match) >= 3 {
			title := strings.TrimSpace(match[2])
			link := match[1]

			if len(title) > 3 {
				job := shared.PrivJobSource{
					Source:    "apna",
					Title:     title,
					URL:       "https://apna.co" + link,
					CreatedAt: time.Now(),
				}
				// Default company if not found
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
func (s *ApnaSource) Name() string {
	return s.NameStr
}
