package sources

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/rs/zerolog/log"
)

// LinkedInSource scrapes jobs from LinkedIn
type LinkedInSource struct {
	BaseSource
	client *http.Client
	rssURL string
}

// NewLinkedInSource creates a new LinkedIn source
func NewLinkedInSource() *LinkedInSource {
	return &LinkedInSource{
		BaseSource: BaseSource{NameStr: "linkedin", BaseURL: "https://www.linkedin.com"},
		client: &http.Client{
			Timeout: 30 * time.Second,
		},
		rssURL: "https://www.linkedin.com/jobs/jobs-rss-all",
	}
}

// Fetch retrieves private jobs from LinkedIn
func (s *LinkedInSource) Fetch(ctx context.Context) ([]PrivJobSource, error) {
	log.Info().Msg("Starting crawl for source: LinkedIn")

	var jobs []PrivJobSource

	rssJobs, err := s.fetchFromRSS(ctx)
	if err != nil {
		log.Warn().Err(err).Msg("LinkedIn RSS fetch failed")
		rssJobs, err = s.fetchFromAlternative(ctx)
		if err != nil {
			log.Error().Err(err).Msg("All LinkedIn fetch methods failed")
			return nil, fmt.Errorf("failed to fetch from LinkedIn: %w", err)
		}
	}

	jobs = append(jobs, rssJobs...)
	log.Info().Int("totalJobs", len(jobs)).Msg("LinkedIn fetch completed")
	return jobs, nil
}

// fetchFromRSS fetches from LinkedIn RSS feed
func (s *LinkedInSource) fetchFromRSS(ctx context.Context) ([]PrivJobSource, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", s.rssURL, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("User-Agent", "RojgarSetu/2.0")

	resp, err := s.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("LinkedIn RSS returned status: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	doc, err := parseRSSXML(string(body))
	if err != nil {
		return nil, err
	}

	var jobs []PrivJobSource
	for _, item := range doc.Channel.Items {
		job := PrivJobSource{
			Source:    "linkedin",
			Title:     cleanString(item.Title),
			URL:       extractURL(item.Link),
			CreatedAt: time.Now(),
		}

		// Parse description
		desc := item.Description
		job.Company = extractField(desc, "Company:")
		job.Location = extractField(desc, "Location:")
		job.Salary = extractField(desc, "Salary:")
		job.Experience = extractField(desc, "Experience:")
		job.JobType = normalizeJobType(extractField(desc, "Job Type:"))

		if job.Title != "" && isValidPrivJob(&job) {
			jobs = append(jobs, job)
		}
	}

	log.Info().Int("jobsFromRSS", len(jobs)).Msg("LinkedIn RSS fetch successful")
	return jobs, nil
}

// fetchFromAlternative fetches from alternative LinkedIn endpoints
func (s *LinkedInSource) fetchFromAlternative(ctx context.Context) ([]PrivJobSource, error) {
	// LinkedIn doesn't provide public API, try LinkedIn Lite or other sources
	// For now, return empty slice - LinkedIn requires authentication
	log.Warn().Msg("LinkedIn alternative fetch not available - requires authentication")
	return []PrivJobSource{}, nil
}

// Name returns the source name
func (s *LinkedInSource) Name() string {
	return s.NameStr
}
