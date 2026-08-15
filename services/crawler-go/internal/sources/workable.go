package sources

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/rs/zerolog/log"
)

// WorkableSource fetches jobs from companies that use Workable as their ATS.
// It consumes the public Workable API:
//
//	https://apply.workable.com/api/v1/accounts/{company}/jobs
//
// which returns JSON of the form:
//
//	{"jobs":[{ "title", "location", "url", "shortlink", "published_at" }]}
type WorkableSource struct {
	BaseSource
	client *http.Client
	orgs   []workableOrg
}

// workableOrg is a single Workable company to crawl.
type workableOrg struct {
	Org     string // Workable company slug, e.g. "doist"
	Company string // Display company name
}

// workableResponse mirrors the Workable API response envelope.
type workableResponse struct {
	Jobs []workableJob `json:"jobs"`
}

// workableJob mirrors a single Workable job object.
type workableJob struct {
	Title       string `json:"title"`
	Location    string `json:"location"`
	URL         string `json:"url"`
	Shortlink   string `json:"shortlink"`
	PublishedAt string `json:"published_at"`
}

// NewWorkableSource creates a Workable source with corporate ATS pool.
func NewWorkableSource() *WorkableSource {
	return &WorkableSource{
		BaseSource: BaseSource{NameStr: "workable", BaseURL: "https://apply.workable.com"},
		client: &http.Client{
			Timeout: 30 * time.Second,
		},
		orgs: []workableOrg{
			{Org: "doist", Company: "Doist"},
			{Org: "buffer", Company: "Buffer"},
			{Org: "zapier", Company: "Zapier"},
			{Org: "automattic", Company: "Automattic"},
		},
	}
}

// Fetch retrieves private jobs from all configured Workable companies.
func (s *WorkableSource) Fetch(ctx context.Context) ([]PrivJobSource, error) {
	log.Info().Msg("Starting crawl for source: Workable")

	var allJobs []PrivJobSource
	for _, org := range s.orgs {
		jobs, err := s.fetchOrg(ctx, org)
		if err != nil {
			log.Warn().Err(err).Str("org", org.Org).Msg("Workable fetch failed for org")
			continue
		}
		allJobs = append(allJobs, jobs...)
	}

	log.Info().Int("totalJobs", len(allJobs)).Msg("Workable fetch completed")
	return allJobs, nil
}

func (s *WorkableSource) fetchOrg(ctx context.Context, org workableOrg) ([]PrivJobSource, error) {
	url := fmt.Sprintf("https://apply.workable.com/api/v1/accounts/%s/jobs", org.Org)
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
		return nil, fmt.Errorf("workable API returned status: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var wr workableResponse
	if err := json.Unmarshal(body, &wr); err != nil {
		return nil, fmt.Errorf("failed to parse workable JSON: %w", err)
	}

	var jobs []PrivJobSource
	for i := range wr.Jobs {
		j := &wr.Jobs[i]
		job := s.jobToPriv(j, org.Company)
		if isValidPrivJob(&job) {
			jobs = append(jobs, job)
		}
	}

	log.Info().Int("jobs", len(jobs)).Str("org", org.Org).Msg("Workable org fetch successful")
	return jobs, nil
}

// jobToPriv converts a Workable job to the shared PrivJobSource model.
func (s *WorkableSource) jobToPriv(j *workableJob, company string) PrivJobSource {
	job := PrivJobSource{
		Source:      "workable",
		Company:     company,
		Title:       strings.TrimSpace(j.Title),
		Location:    strings.TrimSpace(j.Location),
		URL:         j.URL,
		JobType:     "",
		Description: "",
		CreatedAt:   time.Now(),
	}

	if j.PublishedAt != "" {
		if t, err := time.Parse(time.RFC3339, j.PublishedAt); err == nil {
			job.PostedAt = &t
		}
	}

	return job
}

// Name returns the source name.
func (s *WorkableSource) Name() string {
	return s.NameStr
}

// FetchJobs implements the JobSource interface for compatibility with the engine
func (s *WorkableSource) FetchJobs() ([]Job, error) {
	ctx := context.Background()
	privJobs, err := s.Fetch(ctx)
	if err != nil {
		return nil, err
	}

	var jobs []Job
	for _, privJob := range privJobs {
		job := Job{
			Title:             privJob.Title,
			CompanyOrDept:     privJob.Company,
			Location:          privJob.Location,
			QualificationReq:  privJob.JobType,
			SalaryOrPayScale:  "",
			ApplyURL:          privJob.URL,
			SourceAttribution: "Source: Workable ATS API",
			HashChecksum:      "", // Will be set by engine
		}
		jobs = append(jobs, job)
	}

	return jobs, nil
}
