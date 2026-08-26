package gov

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/rojgarsetu/crawler/internal/shared"
	"github.com/rs/zerolog/log"
)

// UPSCSource scrapes jobs from UPSC official website
type UPSCSource struct {
	shared.BaseSource
	client *http.Client
}

// NewUPSCSource creates a new UPSC source
func NewUPSCSource() *UPSCSource {
	return &UPSCSource{
		BaseSource: shared.BaseSource{NameStr: "upsc", BaseURL: "https://www.upsc.gov.in"},
		client: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// Fetch retrieves government jobs from UPSC
func (s *UPSCSource) Fetch(ctx context.Context) ([]shared.GovJobSource, error) {
	log.Info().Msg("Starting crawl for source: UPSC")

	var jobs []shared.GovJobSource

	// Try official UPSC RSS
	rssJobs, err := s.fetchFromRSS(ctx)
	if err != nil {
		log.Warn().Err(err).Msg("UPSC RSS fetch failed, trying website")
		rssJobs, err = s.fetchFromWebsite(ctx)
		if err != nil {
			log.Error().Err(err).Msg("All UPSC fetch methods failed")
			return nil, fmt.Errorf("failed to fetch from UPSC: %w", err)
		}
	}

	jobs = append(jobs, rssJobs...)
	log.Info().Int("totalJobs", len(jobs)).Msg("UPSC fetch completed")
	return jobs, nil
}

// fetchFromRSS fetches from UPSC RSS feed
func (s *UPSCSource) fetchFromRSS(ctx context.Context) ([]shared.GovJobSource, error) {
	rssURL := "https://www.upsc.gov.in/feeds/rss/whatsnew"
	req, err := http.NewRequestWithContext(ctx, "GET", rssURL, nil)
	if err != nil {
		return nil, err
	}

	shared.SetUserAgentAndCheck(req, s.BaseURL)
	if !shared.CheckRobotsTxt(s.BaseURL, "/feeds/rss/whatsnew") {
		return nil, fmt.Errorf("blocked by robots.txt")
	}
	dl := shared.NewDomainLimiter()
	if !dl.Allow("upsc.gov.in") {
		return nil, fmt.Errorf("throttled")
	}

	resp, err := s.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if err := shared.CheckStatusAndPause(resp, "upsc.gov.in"); err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("UPSC RSS returned status: %d", resp.StatusCode)
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
			Source:    "upsc",
			Title:     shared.CleanString(item.Title),
			ApplyURL:  shared.ExtractURL(item.Link),
			CreatedAt: time.Now(),
		}

		desc := item.Description
		job.Department = shared.ExtractField(desc, "Department:")
		job.Location = shared.ExtractField(desc, "Location:")
		job.Salary = shared.ExtractField(desc, "Salary:")
		job.LastDate = shared.ParseDateString(shared.ExtractField(desc, "Last Date:"))

		if job.Title != "" && shared.IsValidJob(&job) {
			jobs = append(jobs, job)
		}
	}

	log.Info().Int("jobsFromRSS", len(jobs)).Msg("UPSC RSS fetch successful")
	return jobs, nil
}

// fetchFromWebsite fetches from UPSC official website
func (s *UPSCSource) fetchFromWebsite(ctx context.Context) ([]shared.GovJobSource, error) {
	urls := []string{
		"https://www.upsc.gov.in/examinations/examinations-notifications",
		"https://www.upsc.gov.in/whats-new",
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
		if !dl.Allow("upsc.gov.in") {
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

		jobs := s.parseHTMLJobs(string(body), url)
		allJobs = append(allJobs, jobs...)
	}

	log.Info().Int("jobsFromWebsite", len(allJobs)).Msg("UPSC website fetch successful")
	return allJobs, nil
}

// parseHTMLJobs parses jobs from UPSC HTML page
func (s *UPSCSource) parseHTMLJobs(html, baseURL string) []shared.GovJobSource {
	var jobs []shared.GovJobSource

	// More specific patterns to match actual examination notices and notifications
	// and exclude general header/sidebar links.
	patterns := []string{
		`<a[^>]*href="(/[^"]*exam-notice/[^"]*)"[^>]*>([^<]*)</a>`,
		`<a[^>]*href="(/[^"]*recruitment-advertisement/[^"]*)"[^>]*>([^<]*)</a>`,
		`<a[^>]*href="(/[^"]*notification[^"]*\.pdf)"[^>]*>([^<]*)</a>`,
	}

	for _, pattern := range patterns {
		matches := shared.ExtractMatches(html, pattern)
		for _, match := range matches {
			if len(match) >= 3 {
				title := strings.TrimSpace(match[2])
				link := match[1]

				if title != "" && !isUPSCHeaderLink(title) {
					job := shared.GovJobSource{
						Source:    "upsc",
						Title:     title,
						ApplyURL:  "https://www.upsc.gov.in" + link,
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

// isUPSCHeaderLink checks if a title looks like a general navigation link
func isUPSCHeaderLink(title string) bool {
	t := strings.ToLower(title)
	headers := []string{"examinations-notifications", "active examinations", "forthcoming examinations", "whats-new", "read more"}
	for _, h := range headers {
		if t == h || strings.Contains(t, h) && len(t) < 30 {
			return true
		}
	}
	return false
}

// Name returns the source name
func (s *UPSCSource) Name() string {
	return s.NameStr
}

// isUPSCRelevant checks if notice is relevant for UPSC
func isUPSCRelevant(title string) bool {
	titleLower := strings.ToLower(title)
	relevant := []string{"examination", "notification", "recruitment", "vacancy", "advertisement"}
	irrelevant := []string{"result", "answer key", "cutoff", "syllabus"}

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
