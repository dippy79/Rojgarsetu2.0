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

// AshbySource fetches jobs from companies that use Ashby as their ATS.
// It consumes the public Ashby Job Board API:
//
//	https://api.ashbyhq.com/posting-api/job-board/{company}
//
// which returns JSON of the form:
//
//	{"jobs":[{ "title", "locationName", "employmentType", "url", "publishedDate" }]}
type AshbySource struct {
	BaseSource
	client *http.Client
	orgs   []ashbyOrg
}

// ashbyOrg is a single Ashby board org to crawl.
type ashbyOrg struct {
	Org     string // Ashby board slug, e.g. "linear"
	Company string // Display company name
}

// ashbyResponse mirrors the Ashby jobs API response envelope.
type ashbyResponse struct {
	Jobs []ashbyJob `json:"jobs"`
}

// ashbyJob mirrors a single Ashby job object.
type ashbyJob struct {
	Title          string `json:"title"`
	LocationName   string `json:"locationName"`
	EmploymentType string `json:"employmentType"`
	URL            string `json:"url"`
	PublishedDate  string `json:"publishedDate"`
}

// NewAshbySource creates an Ashby source with tech startup company pool.
func NewAshbySource() *AshbySource {
	return &AshbySource{
		BaseSource: BaseSource{NameStr: "ashby", BaseURL: "https://api.ashbyhq.com"},
		client: &http.Client{
			Timeout: 30 * time.Second,
		},
		orgs: []ashbyOrg{
			{Org: "linear", Company: "Linear"},
			{Org: "notion", Company: "Notion"},
			{Org: "figma", Company: "Figma"},
			{Org: "vercel", Company: "Vercel"},
			{Org: "airbnb", Company: "Airbnb"},
		},
	}
}

// Fetch retrieves private jobs from all configured Ashby boards.
func (s *AshbySource) Fetch(ctx context.Context) ([]PrivJobSource, error) {
	log.Info().Msg("Starting crawl for source: Ashby")

	var allJobs []PrivJobSource
	for _, org := range s.orgs {
		jobs, err := s.fetchOrg(ctx, org)
		if err != nil {
			log.Warn().Err(err).Str("org", org.Org).Msg("Ashby fetch failed for org")
			continue
		}
		allJobs = append(allJobs, jobs...)
	}

	log.Info().Int("totalJobs", len(allJobs)).Msg("Ashby fetch completed")
	return allJobs, nil
}

func (s *AshbySource) fetchOrg(ctx context.Context, org ashbyOrg) ([]PrivJobSource, error) {
	url := fmt.Sprintf("https://api.ashbyhq.com/posting-api/job-board/%s", org.Org)
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
		return nil, fmt.Errorf("ashby API returned status: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var ash ashbyResponse
	if err := json.Unmarshal(body, &ash); err != nil {
		return nil, fmt.Errorf("failed to parse ashby JSON: %w", err)
	}

	var jobs []PrivJobSource
	for i := range ash.Jobs {
		j := &ash.Jobs[i]
		job := s.jobToPriv(j, org.Company)
		if isValidPrivJob(&job) {
			jobs = append(jobs, job)
		}
	}

	log.Info().Int("jobs", len(jobs)).Str("org", org.Org).Msg("Ashby org fetch successful")
	return jobs, nil
}

// jobToPriv converts an Ashby job to the shared PrivJobSource model.
func (s *AshbySource) jobToPriv(j *ashbyJob, company string) PrivJobSource {
	job := PrivJobSource{
		Source:      "ashby",
		Company:     company,
		Title:       strings.TrimSpace(j.Title),
		Location:    strings.TrimSpace(j.LocationName),
		URL:         j.URL,
		JobType:     normalizeJobType(j.EmploymentType),
		Description: "",
		CreatedAt:   time.Now(),
	}

	if j.PublishedDate != "" {
		if t, err := time.Parse(time.RFC3339, j.PublishedDate); err == nil {
			job.PostedAt = &t
		}
	}

	return job
}

// Name returns the source name.
func (s *AshbySource) Name() string {
	return s.NameStr
}

// FetchJobs implements the JobSource interface for compatibility with the engine
func (s *AshbySource) FetchJobs() ([]Job, error) {
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
			SourceAttribution: "Source: Ashby ATS API",
			HashChecksum:      "", // Will be set by engine
		}
		jobs = append(jobs, job)
	}

	return jobs, nil
}
