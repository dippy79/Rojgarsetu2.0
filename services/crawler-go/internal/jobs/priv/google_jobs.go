package priv

import (
	"github.com/rojgarsetu/crawler/internal/shared"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/rs/zerolog/log"
)

// GoogleJobsSource scrapes jobs from Google Jobs (JSON-LD structured data)
type GoogleJobsSource struct {
	shared.BaseSource
	client *http.Client
}

// NewGoogleJobsSource creates a new Google Jobs source
func NewGoogleJobsSource() *GoogleJobsSource {
	return &GoogleJobsSource{
		BaseSource: shared.BaseSource{NameStr: "google_jobs", BaseURL: "https://careers.google.com"},
		client: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// Fetch retrieves private jobs from Google Jobs
func (s *GoogleJobsSource) Fetch(ctx context.Context) ([]shared.PrivJobSource, error) {
	log.Info().Msg("Starting crawl for source: Google Jobs")

	var jobs []shared.PrivJobSource

	// Fetch from Google Careers
	googleJobs, err := s.fetchFromGoogleCareers(ctx)
	if err != nil {
		log.Warn().Err(err).Msg("Google Careers fetch failed")
	}

	jobs = append(jobs, googleJobs...)

	log.Info().Int("totalJobs", len(jobs)).Msg("Google Jobs fetch completed")
	return jobs, nil
}

// fetchFromGoogleCareers fetches from Google Careers page
func (s *GoogleJobsSource) fetchFromGoogleCareers(ctx context.Context) ([]shared.PrivJobSource, error) {
	urls := []string{
		"https://careers.google.com/jobs/results/",
	}

	var allJobs []shared.PrivJobSource

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

		// Look for JSON-LD structured data
		jobs := s.parseJSONLD(string(body), "Google")
		allJobs = append(allJobs, jobs...)
	}

	log.Info().Int("jobsFromGoogleCareers", len(allJobs)).Msg("Google Careers fetch successful")
	return allJobs, nil
}

// parseJSONLD parses JSON-LD structured data from HTML
func (s *GoogleJobsSource) parseJSONLD(html, company string) []shared.PrivJobSource {
	var jobs []shared.PrivJobSource

	// Find JSON-LD script tags
	pattern := `<script[^>]*type="application/ld\+json"[^>]*>([^<]+)</script>`
	matches := shared.ExtractMatches(html, pattern)

	for _, match := range matches {
		if len(match) >= 2 {
			jsonData := match[1]

			// Try to parse as JobPosting
			var jobPosting struct {
				Title        string `json:"title"`
				Description  string `json:"description"`
				Company      string `json:"company"`
				CompanyName  string `json:"companyName"`
				JobLocation  string `json:"jobLocation"`
				BaseSalary   string `json:"baseSalary"`
				DatePosted   string `json:"datePosted"`
				ValidThrough string `json:"validThrough"`
			}

			if err := json.Unmarshal([]byte(jsonData), &jobPosting); err != nil {
				continue
			}

			if jobPosting.Title != "" {
				job := shared.PrivJobSource{
					Source:      "google_jobs",
					Company:     company,
					Title:       shared.CleanString(jobPosting.Title),
					Description: shared.CleanString(jobPosting.Description),
					Salary:      jobPosting.BaseSalary,
					CreatedAt:   time.Now(),
				}

				if jobPosting.Company != "" {
					job.Company = jobPosting.Company
				} else if jobPosting.CompanyName != "" {
					job.Company = jobPosting.CompanyName
				}

				if jobPosting.JobLocation != "" {
					job.Location = jobPosting.JobLocation
				}

				if shared.IsValidPrivJob(&job) {
					jobs = append(jobs, job)
				}
			}
		}
	}

	return jobs
}

// parseHTMLJobs parses jobs from HTML
func (s *GoogleJobsSource) parseHTMLJobs(html string) []shared.PrivJobSource {
	var jobs []shared.PrivJobSource

	// Look for job cards
	patterns := []string{
		`<h2[^>]*class="[^"]*job-title[^"]*"[^>]*><a[^>]*href="([^"]*)"[^>]*>([^<]*)</a></h2>`,
		`<a[^>]*href="(/jobs/view/[^"]*)"[^>]*>([^<]*)</a>`,
	}

	for _, pattern := range patterns {
		matches := shared.ExtractMatches(html, pattern)
		for _, match := range matches {
			if len(match) >= 3 {
				link := match[1]
				title := strings.TrimSpace(match[2])

				if len(title) > 5 {
					job := shared.PrivJobSource{
						Source:    "google_jobs",
						Title:     title,
						URL:       "https://careers.google.com" + link,
						Company:   "Google",
						CreatedAt: time.Now(),
					}
					if shared.IsValidPrivJob(&job) {
						jobs = append(jobs, job)
					}
				}
			}
		}
	}

	return jobs
}

// Name returns the source name
func (s *GoogleJobsSource) Name() string {
	return s.NameStr
}

