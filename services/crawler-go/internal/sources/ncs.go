package sources

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"time"

	"github.com/rs/zerolog/log"
)

// NCSJob represents a job from National Career Service API
type NCSJob struct {
	Title          string `json:"title"`
	CompanyName    string `json:"companyName"`
	Location       string `json:"location"`
	ApplyURL       string `json:"applyUrl"`
	LastDate       string `json:"lastDate"`
	Department     string `json:"department"`
	VacancyCount   string `json:"vacancyCount"`
	Salary         string `json:"salary"`
	JobDescription string `json:"jobDescription"`
}

// NCSSource scrapes jobs from National Career Service API
type NCSSource struct {
	BaseSource
	client *http.Client
	apiURL string
}

// NewNCSSource creates a new NCS source
func NewNCSSource() *NCSSource {
	return &NCSSource{
		BaseSource: BaseSource{NameStr: "ncs", BaseURL: "https://www.ncs.gov.in"},
		client: &http.Client{
			Timeout: 30 * time.Second,
		},
		apiURL: "https://www.ncs.gov.in/feeds/rss/jobs",
	}
}

// Fetch retrieves government jobs from NCS
func (s *NCSSource) Fetch(ctx context.Context) ([]GovJobSource, error) {
	log.Info().Msg("Starting crawl for source: NCS (National Career Service)")

	// FLAG (Phase A): The legacy NCS endpoints are dead.
	//   - https://www.ncs.gov.in/feeds/rss/jobs  -> 404 (verified live)
	//   - https://www.ncs.gov.in/_v/api/JobSearch/ -> 404 (verified live)
	//   - https://www.ncs.gov.in/ (site root)  -> 404
	// The current NCS portal (a JS Single-Page App) does not expose a plain-HTTP
	// RSS/JSON feed that this fetcher can reach. Rather than silently returning 0
	// jobs, we surface a clear diagnostic so the RunSummary flags this source as
	// needing a different approach (e.g. a browser-driven crawl of the NCS portal,
	// or a maintained aggregate feed). No code change here resurrects it.
	err := fmt.Errorf("NCS job feed is dead: /feeds/rss/jobs=404, /_v/api/JobSearch/=404, site root=404 (all verified live). NCS portal is a JS SPA requiring a different approach (browser-driven crawl or maintained feed).")
	log.Warn().Msg(err.Error())
	return nil, err
}

// fetchFromRSS fetches jobs from NCS RSS feed
func (s *NCSSource) fetchFromRSS(ctx context.Context) ([]GovJobSource, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", s.apiURL, nil)
	if err != nil {
		return nil, err
	}

	SetUserAgentAndCheck(req, s.BaseURL)
	req.Header.Set("Accept", "application/rss+xml, application/xml, text/xml")
	if !CheckRobotsTxt(s.BaseURL, "/feeds/rss/jobs") {
		return nil, fmt.Errorf("blocked by robots.txt")
	}
	dl := NewDomainLimiter()
	if !dl.Allow("ncs.gov.in") {
		return nil, fmt.Errorf("throttled")
	}

	resp, err := s.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("NCS RSS returned status: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	// Parse RSS XML
	doc, err := parseRSSXML(string(body))
	if err != nil {
		return nil, err
	}

	var jobs []GovJobSource
	for _, item := range doc.Channel.Items {
		// Extract NCS specific fields from description
		job := GovJobSource{
			Source:    "ncs",
			Title:     cleanString(item.Title),
			ApplyURL:  extractURL(item.Link),
			CreatedAt: time.Now(),
		}

		// Parse description for more details
		desc := item.Description
		job.Department = extractField(desc, "Department:")
		job.Location = extractField(desc, "Location:")
		job.Salary = extractField(desc, "Salary:")
		// Extract vacancy count
		if vcStr := extractField(desc, "Vacancies:"); vcStr != "" {
			if vc, err := strconv.Atoi(vcStr); err == nil {
				job.VacancyCount = &vc
			}
		}
		job.LastDate = parseDateString(extractField(desc, "Last Date:"))

		if job.Title != "" && isValidJob(&job) {
			jobs = append(jobs, job)
		}
	}

	log.Info().Int("jobsFromRSS", len(jobs)).Msg("NCS RSS fetch successful")
	return jobs, nil
}

// fetchFromAlternative fetches from alternative NCS endpoints
func (s *NCSSource) fetchFromAlternative(ctx context.Context) ([]GovJobSource, error) {
	// Alternative: Try JSON API if available
	altURLs := []string{
		"https://www.ncs.gov.in/_v/api/JobSearch/",
	}

	for _, url := range altURLs {
		req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
		if err != nil {
			continue
		}

		req.Header.Set("User-Agent", "RojgarSetu/2.0")
		req.Header.Set("Accept", "application/json")

		resp, err := s.client.Do(req)
		if err != nil {
			continue
		}
		defer resp.Body.Close()

		if resp.StatusCode == http.StatusOK {
			var ncsJobs []NCSJob
			if err := json.NewDecoder(resp.Body).Decode(&ncsJobs); err == nil {
				var jobs []GovJobSource
				for _, job := range ncsJobs {
					govJob := GovJobSource{
						Source:     "ncs",
						Title:      cleanString(job.Title),
						Department: cleanString(job.Department),
						Location:   cleanString(job.Location),
						ApplyURL:   cleanString(job.ApplyURL),
						Salary:     cleanString(job.Salary),
						CreatedAt:  time.Now(),
					}
					govJob.LastDate = parseDateString(job.LastDate)
					if vc, err := strconv.Atoi(job.VacancyCount); err == nil {
						govJob.VacancyCount = &vc
					}
					if govJob.Title != "" && isValidJob(&govJob) {
						jobs = append(jobs, govJob)
					}
				}
				log.Info().Int("jobsFromAlt", len(jobs)).Msg("NCS alternative fetch successful")
				return jobs, nil
			}
		}
	}

	return []GovJobSource{}, nil
}

// Name returns the source name
func (s *NCSSource) Name() string {
	return s.NameStr
}
