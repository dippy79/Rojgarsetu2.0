package sources

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/rs/zerolog/log"
)

// UPSCSource scrapes jobs from UPSC (Union Public Service Commission)
type UPSCSource struct {
	BaseSource
	client *http.Client
	rssURL string
}

// NewUPSCSource creates a new UPSC source
func NewUPSCSource() *UPSCSource {
	return &UPSCSource{
		BaseSource: BaseSource{NameStr: "upsc", BaseURL: "https://upsc.gov.in"},
		client: &http.Client{
			Timeout: 30 * time.Second,
		},
		rssURL: "https://www.upsc.gov.in/rss",
	}
}

// Fetch retrieves government jobs from UPSC
func (s *UPSCSource) Fetch(ctx context.Context) ([]GovJobSource, error) {
	log.Info().Msg("Starting crawl for source: UPSC")

	var jobs []GovJobSource

	// Try RSS feed
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
func (s *UPSCSource) fetchFromRSS(ctx context.Context) ([]GovJobSource, error) {
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
		return nil, fmt.Errorf("UPSC RSS returned status: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	doc, err := parseRSSXML(string(body))
	if err != nil {
		return nil, err
	}

	var jobs []GovJobSource
	for _, item := range doc.Channel.Items {
		job := GovJobSource{
			Source:    "upsc",
			Title:     cleanString(item.Title),
			ApplyURL:  extractURL(item.Link),
			CreatedAt: time.Now(),
		}

		// Parse description
		desc := item.Description
		job.Department = extractField(desc, "Department:")
		job.LastDate = parseDateString(extractField(desc, "Last Date:"))

		// Extract vacancy count
		if vc := extractField(desc, "Posts:"); vc != "" {
			if v, err := strconv.Atoi(vc); err == nil {
				job.VacancyCount = &v
			}
		}

		if job.Title != "" && isValidJob(&job) {
			jobs = append(jobs, job)
		}
	}

	log.Info().Int("jobsFromRSS", len(jobs)).Msg("UPSC RSS fetch successful")
	return jobs, nil
}

// fetchFromWebsite fetches from UPSC official website
func (s *UPSCSource) fetchFromWebsite(ctx context.Context) ([]GovJobSource, error) {
	urls := []string{
		"https://www.upsc.gov.in/recruitment",
		"https://www.upsc.gov.in/whats-new",
	}

	var allJobs []GovJobSource

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

		jobs := s.parseHTMLJobs(string(body), url)
		allJobs = append(allJobs, jobs...)
	}

	log.Info().Int("jobsFromWebsite", len(allJobs)).Msg("UPSC website fetch successful")
	return allJobs, nil
}

// parseHTMLJobs parses jobs from UPSC HTML page
func (s *UPSCSource) parseHTMLJobs(html, baseURL string) []GovJobSource {
	var jobs []GovJobSource

	patterns := []string{
		`<a[^>]*href="(/[^"]*recruitment[^"]*)"[^>]*>([^<]*)</a>`,
		`<a[^>]*href="(/[^"]*notice[^"]*)"[^>]*>([^<]*)</a>`,
	}

	for _, pattern := range patterns {
		matches := extractMatches(html, pattern)
		for _, match := range matches {
			if len(match) >= 3 {
				title := strings.TrimSpace(match[2])
				link := match[1]

				if isUPSCRelevant(title) {
					job := GovJobSource{
						Source:    "upsc",
						Title:     title,
						ApplyURL:  "https://www.upsc.gov.in" + link,
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
func (s *UPSCSource) Name() string {
	return s.NameStr
}

// isUPSCRelevant checks if notice is relevant
func isUPSCRelevant(title string) bool {
	titleLower := strings.ToLower(title)
	relevant := []string{"recruitment", "advertisement", "notification", "vacancy", "examination", "civil services"}
	irrelevant := []string{"result", "answer key", "cutoff", "interview"}

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
