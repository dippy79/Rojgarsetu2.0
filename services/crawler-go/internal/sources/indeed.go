package sources

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/rs/zerolog/log"
)

// IndeedSource scrapes jobs from Indeed
type IndeedSource struct {
	BaseSource
	client *http.Client
	rssURL string
}

// NewIndeedSource creates a new Indeed source
func NewIndeedSource() *IndeedSource {
	return &IndeedSource{
		BaseSource: BaseSource{NameStr: "indeed", BaseURL: "https://www.indeed.com"},
		client: &http.Client{
			Timeout: 30 * time.Second,
		},
		rssURL: "https://www.indeed.com/jobs?q=&l=India&rss=1",
	}
}

// Fetch retrieves private jobs from Indeed
func (s *IndeedSource) Fetch(ctx context.Context) ([]PrivJobSource, error) {
	log.Info().Msg("Starting crawl for source: Indeed")

	var jobs []PrivJobSource

	rssJobs, err := s.fetchFromRSS(ctx)
	if err != nil {
		log.Warn().Err(err).Msg("Indeed RSS fetch failed")
		rssJobs, err = s.fetchFromWebsite(ctx)
		if err != nil {
			log.Error().Err(err).Msg("All Indeed fetch methods failed")
			return nil, fmt.Errorf("failed to fetch from Indeed: %w", err)
		}
	}

	jobs = append(jobs, rssJobs...)
	log.Info().Int("totalJobs", len(jobs)).Msg("Indeed fetch completed")
	return jobs, nil
}

// fetchFromRSS fetches from Indeed RSS feed
func (s *IndeedSource) fetchFromRSS(ctx context.Context) ([]PrivJobSource, error) {
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
		return nil, fmt.Errorf("Indeed RSS returned status: %d", resp.StatusCode)
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
			Source:    "indeed",
			Title:     cleanString(item.Title),
			URL:       extractURL(item.Link),
			CreatedAt: time.Now(),
		}

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

	log.Info().Int("jobsFromRSS", len(jobs)).Msg("Indeed RSS fetch successful")
	return jobs, nil
}

// fetchFromWebsite fetches from Indeed website
func (s *IndeedSource) fetchFromWebsite(ctx context.Context) ([]PrivJobSource, error) {
	// Try to scrape Indeed search results
	searchURLs := []string{
		"https://www.indeed.com/jobs?q=software+engineer&l=India",
		"https://www.indeed.com/jobs?q=developer&l=India",
	}

	var allJobs []PrivJobSource

	for _, url := range searchURLs {
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

		jobs := s.parseHTMLJobs(string(body))
		allJobs = append(allJobs, jobs...)
	}

	log.Info().Int("jobsFromWebsite", len(allJobs)).Msg("Indeed website fetch successful")
	return allJobs, nil
}

// parseHTMLJobs parses jobs from Indeed HTML
func (s *IndeedSource) parseHTMLJobs(html string) []PrivJobSource {
	var jobs []PrivJobSource

	// Indeed job card pattern
	pattern := `<a[^>]*href="(/job/[^"]*)"[^>]*>([^<]*)</a>`
	matches := extractMatches(html, pattern)

	for _, match := range matches {
		if len(match) >= 3 {
			title := strings.TrimSpace(match[2])
			link := match[1]

			if len(title) > 5 {
				job := PrivJobSource{
					Source:    "indeed",
					Title:     title,
					URL:       "https://www.indeed.com" + link,
					CreatedAt: time.Now(),
				}
				if isValidPrivJob(&job) {
					jobs = append(jobs, job)
				}
			}
		}
	}

	return jobs
}

// Name returns the source name
func (s *IndeedSource) Name() string {
	return s.NameStr
}
