package priv

import (
	"github.com/rojgarsetu/crawler/internal/shared"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/rs/zerolog/log"
)

// GreenhouseSource fetches jobs from companies that use Greenhouse as their ATS.
// It consumes the public Greenhouse Job Board API:
//
//	https://boards-api.greenhouse.io/v1/boards/{org}/jobs
//
// which returns JSON of the form:
//
//	{"jobs":[{ "id", "title", "company_name", "location":{"name"},
//	            "absolute_url", "updated_at", "first_published" }]}
type GreenhouseSource struct {
	shared.BaseSource
	client *http.Client
	orgs   []greenhouseOrg
}

// greenhouseOrg is a single Greenhouse board org to crawl.
type greenhouseOrg struct {
	Org     string // Greenhouse board slug, e.g. "stripe"
	Company string // Display company name
}

// greenhouseResponse mirrors the Greenhouse jobs API response envelope.
type greenhouseResponse struct {
	Jobs []greenhouseJob `json:"jobs"`
}

// greenhouseJob mirrors a single Greenhouse job object.
type greenhouseJob struct {
	ID          int64  `json:"id"`
	Title       string `json:"title"`
	CompanyName string `json:"company_name"`
	Location    struct {
		Name string `json:"name"`
	} `json:"location"`
	AbsoluteURL    string `json:"absolute_url"`
	UpdatedAt      string `json:"updated_at"`
	FirstPublished string `json:"first_published"`
}

// NewGreenhouseSource creates a Greenhouse source. The third POC company is
// GitLab (verified live: 189 jobs via boards-api.greenhouse.io/v1/boards/gitlab/jobs).
func NewGreenhouseSource() *GreenhouseSource {
	return &GreenhouseSource{
		BaseSource: shared.BaseSource{NameStr: "greenhouse", BaseURL: "https://boards-api.greenhouse.io"},
		client: &http.Client{
			Timeout: 30 * time.Second,
		},
		orgs: []greenhouseOrg{
			{Org: "stripe", Company: "Stripe"},
			{Org: "gitlab", Company: "GitLab"},
		},
	}
}

// Fetch retrieves private jobs from all configured Greenhouse boards.
func (s *GreenhouseSource) Fetch(ctx context.Context) ([]shared.PrivJobSource, error) {
	log.Info().Msg("Starting crawl for source: Greenhouse")

	var allJobs []shared.PrivJobSource
	for _, org := range s.orgs {
		jobs, err := s.fetchOrg(ctx, org)
		if err != nil {
			log.Warn().Err(err).Str("org", org.Org).Msg("Greenhouse fetch failed for org")
			continue
		}
		allJobs = append(allJobs, jobs...)
	}

	log.Info().Int("totalJobs", len(allJobs)).Msg("Greenhouse fetch completed")
	return allJobs, nil
}

func (s *GreenhouseSource) fetchOrg(ctx context.Context, org greenhouseOrg) ([]shared.PrivJobSource, error) {
	url := fmt.Sprintf("https://boards-api.greenhouse.io/v1/boards/%s/jobs", org.Org)
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
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
		return nil, fmt.Errorf("greenhouse API returned status: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var gh greenhouseResponse
	if err := json.Unmarshal(body, &gh); err != nil {
		return nil, fmt.Errorf("failed to parse greenhouse JSON: %w", err)
	}

	var jobs []shared.PrivJobSource
	for i := range gh.Jobs {
		j := &gh.Jobs[i]
		job := s.jobToPriv(j, org.Company)
		if shared.IsValidPrivJob(&job) {
			jobs = append(jobs, job)
		}
	}

	log.Info().Int("jobs", len(jobs)).Str("org", org.Org).Msg("Greenhouse org fetch successful")
	return jobs, nil
}

// jobToPriv converts a Greenhouse job to the shared shared.PrivJobSource model.
func (s *GreenhouseSource) jobToPriv(j *greenhouseJob, company string) shared.PrivJobSource {
	job := shared.PrivJobSource{
		Source:      "greenhouse",
		Company:     company,
		Title:       strings.TrimSpace(j.Title),
		Location:    strings.TrimSpace(j.Location.Name),
		URL:         j.AbsoluteURL,
		JobType:     "",
		Description: "",
		CreatedAt:   time.Now(),
	}

	if j.FirstPublished != "" {
		if t, err := time.Parse(time.RFC3339, j.FirstPublished); err == nil {
			job.PostedAt = &t
		}
	}

	return job
}

// Name returns the source name.
func (s *GreenhouseSource) Name() string {
	return s.NameStr
}

