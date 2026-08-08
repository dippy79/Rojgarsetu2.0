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

// EmploymentNewsSource scrapes from Employment News
type EmploymentNewsSource struct {
	shared.BaseSource
	client *http.Client
}

// NewEmploymentNewsSource creates a new Employment News source
func NewEmploymentNewsSource() *EmploymentNewsSource {
	return &EmploymentNewsSource{
		BaseSource: shared.BaseSource{NameStr: "employment_news", BaseURL: "https://www.employmentnews.gov.in"},
		client: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// Fetch retrieves jobs from Employment News
func (s *EmploymentNewsSource) Fetch(ctx context.Context) ([]shared.GovJobSource, error) {
	log.Info().Msg("Starting crawl for source: Employment News")

	// FLAG (Phase A first-pass): The Employment News RSS feed is dead
	// (404 verified live at /rss/eng/weekly-vacancy.xml), and the site root
	// also returns 404. The current portal does not expose a plain-HTTP
	// RSS/JSON feed that this fetcher can reach.
	//
	// Rather than silently returning 0 jobs, we surface a clear diagnostic so
	// the RunSummary flags this source as needing a different approach
	// (e.g. a maintained aggregate feed or a browser-driven crawl of the
	// weekly-vacancy PDFs at employmentnews.gov.in/NewEmp/Weekly_Vacancy.aspx).
	err := fmt.Errorf("Employment News feed is dead: /rss/eng/weekly-vacancy.xml=404, site root=404 (verified live). Portal requires a different approach (maintained aggregate feed or browser-driven crawl of weekly-vacancy PDFs).")
	log.Warn().Msg(err.Error())
	return nil, err
}

// fetchFromRSS fetches from Employment News RSS
func (s *EmploymentNewsSource) fetchFromRSS(ctx context.Context) ([]shared.GovJobSource, error) {
	rssURL := "https://www.employmentnews.gov.in/rss/eng/weekly-vacancy.xml"
	req, err := http.NewRequestWithContext(ctx, "GET", rssURL, nil)
	if err != nil {
		return nil, err
	}

	shared.SetUserAgentAndCheck(req, s.BaseURL)
	if !shared.CheckRobotsTxt(s.BaseURL, "/rss/eng/weekly-vacancy.xml") {
		return nil, fmt.Errorf("blocked by robots.txt")
	}
	dl := shared.NewDomainLimiter()
	if !dl.Allow("employmentnews.gov.in") {
		return nil, fmt.Errorf("throttled")
	}

	resp, err := s.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Employment News RSS returned status: %d", resp.StatusCode)
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
			Source:    "employment_news",
			Title:     shared.CleanString(item.Title),
			ApplyURL:  shared.ExtractURL(item.Link),
			CreatedAt: time.Now(),
		}

		// Parse description for more details
		desc := item.Description
		job.Department = shared.ExtractField(desc, "Department:")
		job.Location = shared.ExtractField(desc, "Location:")
		job.Salary = shared.ExtractField(desc, "Salary:")
		job.LastDate = shared.ParseDateString(shared.ExtractField(desc, "Last Date:"))
		job.NotificationURL = item.Link

		if job.Title != "" && shared.IsValidJob(&job) {
			jobs = append(jobs, job)
		}
	}

	log.Info().Int("jobsFromRSS", len(jobs)).Msg("Employment News RSS fetch successful")
	return jobs, nil
}

// fetchFromWebsite fetches from the main website
func (s *EmploymentNewsSource) fetchFromWebsite(ctx context.Context) ([]shared.GovJobSource, error) {
	urls := []string{
		"https://www.employmentnews.gov.in/NewEmp/Weekly_Vacancy.aspx",
		"https://www.employmentnews.gov.in/NewEmp/Archives.aspx",
	}

	var allJobs []shared.GovJobSource

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

		pageJobs := s.parseHTMLJobs(string(body), url)
		allJobs = append(allJobs, pageJobs...)
	}

	log.Info().Int("jobsFromWebsite", len(allJobs)).Msg("Employment News website fetch successful")
	return allJobs, nil
}

// parseHTMLJobs parses Employment News HTML
func (s *EmploymentNewsSource) parseHTMLJobs(html, baseURL string) []shared.GovJobSource {
	var jobs []shared.GovJobSource

	patterns := []string{
		`<a[^>]*href="([^"]*weekly-vacancy[^"]*)"[^>]*>([^<]*)</a>`,
		`<a[^>]*href="([^"]*vacancy[^"]*)"[^>]*>([^<]*)</a>`,
		`<a[^>]*href="([^"]*notification[^"]*)"[^>]*>([^<]*)</a>`,
	}

	for _, pattern := range patterns {
		matches := shared.ExtractMatches(html, pattern)
		for _, match := range matches {
			if len(match) >= 3 {
				title := strings.TrimSpace(match[2])
				link := match[1]

				// Construct URL
				if strings.HasPrefix(link, "/") {
					link = "https://www.employmentnews.gov.in" + link
				} else if !strings.HasPrefix(link, "http") {
					link = baseURL + "/" + link
				}

				job := shared.GovJobSource{
					Source:          "employment_news",
					Title:           title,
					ApplyURL:        link,
					NotificationURL: link,
					CreatedAt:       time.Now(),
				}

				if shared.IsValidJob(&job) {
					jobs = append(jobs, job)
				}
			}
		}
	}

	return jobs
}

// Name returns the source name
func (s *EmploymentNewsSource) Name() string {
	return s.NameStr
}
