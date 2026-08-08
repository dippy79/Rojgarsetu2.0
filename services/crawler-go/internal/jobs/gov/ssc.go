package gov

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/rojgarsetu/crawler/internal/shared"
	"github.com/rs/zerolog/log"
)

// SSCSource scrapes jobs from SSC (Staff Selection Commission)
type SSCSource struct {
	shared.BaseSource
	client *http.Client
	rssURL string
}

// NewSSCSource creates a new SSC source
func NewSSCSource() *SSCSource {
	return &SSCSource{
		BaseSource: shared.BaseSource{NameStr: "ssc", BaseURL: "https://ssc.gov.in"},
		client: &http.Client{
			Timeout: 30 * time.Second,
		},
		// Verified live: the SSURL is https://ssc.gov.in/rss%20feeds/Monthly.xml
		// (200). The old hardcoded value had a literal space which broke the
		// robots.txt path check and the request. Encoded as %20.
		rssURL: "https://ssc.gov.in/rss%20feeds/Monthly.xml",
	}
}

// Fetch retrieves government jobs from SSC
func (s *SSCSource) Fetch(ctx context.Context) ([]shared.GovJobSource, error) {
	log.Info().Msg("Starting crawl for source: SSC")

	var jobs []shared.GovJobSource

	// Try official SSC RSS
	rssJobs, err := s.fetchFromRSS(ctx)
	if err != nil {
		log.Warn().Err(err).Msg("SSC RSS fetch failed, trying website")
		rssJobs, err = s.fetchFromWebsite(ctx)
		if err != nil {
			log.Error().Err(err).Msg("All SSC fetch methods failed")
			return nil, fmt.Errorf("failed to fetch from SSC: %w", err)
		}
	}

	jobs = append(jobs, rssJobs...)
	log.Info().Int("totalJobs", len(jobs)).Msg("SSC fetch completed")
	return jobs, nil
}

// fetchFromRSS fetches from SSC RSS feed
func (s *SSCSource) fetchFromRSS(ctx context.Context) ([]shared.GovJobSource, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", s.rssURL, nil)
	if err != nil {
		return nil, err
	}

	shared.SetUserAgentAndCheck(req, s.BaseURL)
	if !shared.CheckRobotsTxt(s.BaseURL, "/rss%20feeds/Monthly.xml") {
		return nil, fmt.Errorf("blocked by robots.txt")
	}
	dl := shared.NewDomainLimiter()
	if !dl.Allow("ssc.gov.in") {
		return nil, fmt.Errorf("throttled")
	}

	// Execute HTTP Request
	resp, err := s.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if err := shared.CheckStatusAndPause(resp, "ssc.gov.in"); err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("SSC RSS returned status: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	doc, err := shared.ParseRSSXML(string(body))
	if err != nil {
		return nil, err
	}

	var jobs []shared.GovJobSource
	for _, item := range doc.Channel.Items {
		job := shared.GovJobSource{
			Source:    "ssc",
			Title:     shared.CleanString(item.Title),
			ApplyURL:  shared.ExtractURL(item.Link),
			CreatedAt: time.Now(),
		}

		// Parse description for details
		desc := item.Description
		job.Department = shared.ExtractField(desc, "Department:")
		job.Location = shared.ExtractField(desc, "Location:")
		job.Salary = shared.ExtractField(desc, "Pay Level:")
		job.LastDate = shared.ParseDateString(shared.ExtractField(desc, "Last Date:"))

		// Extract vacancy count
		if vc := shared.ExtractField(desc, "Vacancies:"); vc != "" {
			if v, err := strconv.Atoi(vc); err == nil {
				job.VacancyCount = &v
			}
		}

		if job.Title != "" && shared.IsValidJob(&job) {
			jobs = append(jobs, job)
		}
	}

	log.Info().Int("jobsFromRSS", len(jobs)).Msg("SSC RSS fetch successful")
	return jobs, nil
}

// fetchFromWebsite fetches from SSC official website
func (s *SSCSource) fetchFromWebsite(ctx context.Context) ([]shared.GovJobSource, error) {
	urls := []string{
		"https://ssc.gov.in/notice",
		"https://ssc.gov.in/status-of-your-application",
	}

	var allJobs []shared.GovJobSource

	for _, url := range urls {
		req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
		if err != nil {
			continue
		}

		shared.SetUserAgentAndCheck(req, s.BaseURL)
		if !shared.CheckRobotsTxt(s.BaseURL, url) {
			continue
		}
		dl := shared.NewDomainLimiter()
		if !dl.Allow("ssc.gov.in") {
			continue
		}

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

		// Parse HTML for notice links
		jobs := s.parseHTMLJobs(string(body), url)
		allJobs = append(allJobs, jobs...)
	}

	log.Info().Int("jobsFromWebsite", len(allJobs)).Msg("SSC website fetch successful")
	return allJobs, nil
}

// parseHTMLJobs parses jobs from SSC HTML page
func (s *SSCSource) parseHTMLJobs(html, baseURL string) []shared.GovJobSource {
	var jobs []shared.GovJobSource

	// Look for notice links
	patterns := []string{
		`<a[^>]*href="(/[^"]*notice[^"]*)"[^>]*>([^<]*)</a>`,
		`<a[^>]*href="(/[^"]*advertisement[^"]*)"[^>]*>([^<]*)</a>`,
	}

	for _, pattern := range patterns {
		matches := shared.ExtractMatches(html, pattern)
		for _, match := range matches {
			if len(match) >= 3 {
				title := strings.TrimSpace(match[2])
				link := match[1]

				// Only include relevant notices
				if isRelevantNotice(title) {
					job := shared.GovJobSource{
						Source:    "ssc",
						Title:     title,
						ApplyURL:  "https://ssc.gov.in" + link,
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
func (s *SSCSource) Name() string {
	return s.NameStr
}

// isRelevantNotice checks if notice is a job notification
func isRelevantNotice(title string) bool {
	titleLower := strings.ToLower(title)
	relevant := []string{"notification", "advertisement", "result", "examination", "recruitment", "vacancy"}
	irrelevant := []string{"answer key", "result", "cutoff", "syllabus"}

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
