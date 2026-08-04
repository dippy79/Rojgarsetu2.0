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

// EmploymentNewsSource scrapes jobs from Employment News
type EmploymentNewsSource struct {
	BaseSource
	client *http.Client
	rssURL string
}

// NewEmploymentNewsSource creates a new Employment News source
func NewEmploymentNewsSource() *EmploymentNewsSource {
	return &EmploymentNewsSource{
		BaseSource: BaseSource{NameStr: "employment_news", BaseURL: "https://employmentnews.gov.in"},
		client: &http.Client{
			Timeout: 30 * time.Second,
		},
		rssURL: "https://www.employmentnews.gov.in/rss/eng/weekly-vacancy.xml",
	}
}

// Fetch retrieves government jobs from Employment News
func (s *EmploymentNewsSource) Fetch(ctx context.Context) ([]GovJobSource, error) {
	log.Info().Msg("Starting crawl for source: Employment News")

	// FLAG (Phase A): The Employment News RSS + website endpoints are dead.
	//   - https://www.employmentnews.gov.in/rss/eng/weekly-vacancy.xml -> 404 (verified live)
	//   - https://www.employmentnews.gov.in/ (site root) -> 404 (verified live)
	// The portal's current structure does not expose a plain-HTTP RSS/HTML feed
	// this fetcher can reach. Rather than silently returning 0 jobs, surface a
	// clear diagnostic so the RunSummary flags this source as needing a different
	// approach (e.g. a maintained aggregate feed or browser-driven crawl of the
	// weekly-vacancy PDF listings). No code change here resurrects it.
	err := fmt.Errorf("Employment News feed is dead: /rss/eng/weekly-vacancy.xml=404, site root=404 (verified live). Portal requires a different approach (maintained aggregate feed or browser-driven crawl of weekly-vacancy PDFs).")
	log.Warn().Msg(err.Error())
	return nil, err
}

// fetchFromRSS fetches from Employment News RSS feed
func (s *EmploymentNewsSource) fetchFromRSS(ctx context.Context) ([]GovJobSource, error) {
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
		return nil, fmt.Errorf("Employment News RSS returned status: %d", resp.StatusCode)
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
			Source:    "employment_news",
			Title:     cleanString(item.Title),
			ApplyURL:  extractURL(item.Link),
			CreatedAt: time.Now(),
		}

		desc := item.Description
		job.Department = extractField(desc, "Department:")
		job.Location = extractField(desc, "Location:")
		job.Salary = extractField(desc, "Salary:")
		job.LastDate = parseDateString(extractField(desc, "Last Date:"))

		if vc := extractField(desc, "Vacancies:"); vc != "" {
			if v, err := strconv.Atoi(vc); err == nil {
				job.VacancyCount = &v
			}
		}

		if job.Title != "" && isValidJob(&job) {
			jobs = append(jobs, job)
		}
	}

	log.Info().Int("jobsFromRSS", len(jobs)).Msg("Employment News RSS fetch successful")
	return jobs, nil
}

// fetchFromWebsite fetches from Employment News website
func (s *EmploymentNewsSource) fetchFromWebsite(ctx context.Context) ([]GovJobSource, error) {
	urls := []string{
		"https://www.employmentnews.gov.in/weekly-vacancy.php",
		"https://www.employmentnews.gov.in/result.php",
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

		jobs := s.parseHTMLJobs(string(body))
		allJobs = append(allJobs, jobs...)
	}

	log.Info().Int("jobsFromWebsite", len(allJobs)).Msg("Employment News website fetch successful")
	return allJobs, nil
}

// parseHTMLJobs parses jobs from Employment News HTML
func (s *EmploymentNewsSource) parseHTMLJobs(html string) []GovJobSource {
	var jobs []GovJobSource

	patterns := []string{
		`<a[^>]*href="(/[^"]*vacancy[^"]*)"[^>]*>([^<]*)</a>`,
		`<a[^>]*href="(/[^"]*advertisement[^"]*)"[^>]*>([^<]*)</a>`,
	}

	for _, pattern := range patterns {
		matches := extractMatches(html, pattern)
		for _, match := range matches {
			if len(match) >= 3 {
				title := strings.TrimSpace(match[2])
				link := match[1]

				if isEmploymentNewsRelevant(title) {
					job := GovJobSource{
						Source:    "employment_news",
						Title:     title,
						ApplyURL:  "https://www.employmentnews.gov.in" + link,
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
func (s *EmploymentNewsSource) Name() string {
	return s.NameStr
}

// isEmploymentNewsRelevant checks if notice is relevant
func isEmploymentNewsRelevant(title string) bool {
	titleLower := strings.ToLower(title)
	relevant := []string{"vacancy", "recruitment", "advertisement", "notification"}
	irrelevant := []string{"result", "answer key", "cutoff"}

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
