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

// SmartRecruitersSource fetches jobs from companies that use SmartRecruiters as their ATS.
// It consumes the public SmartRecruiters API:
//
//	https://api.smartrecruiters.com/v1/companies/{company}/postings
//
// which returns JSON of the form:
//
//	{"content":[{ "title", "location", "type", "ref", "createdOn" }]}
type SmartRecruitersSource struct {
	BaseSource
	client *http.Client
	orgs   []smartRecruitersOrg
}

// smartRecruitersOrg is a single SmartRecruiters company to crawl.
type smartRecruitersOrg struct {
	Org     string // SmartRecruiters company slug, e.g. "netflix"
	Company string // Display company name
}

// smartRecruitersResponse mirrors the SmartRecruiters API response envelope.
type smartRecruitersResponse struct {
	Content []smartRecruitersJob `json:"content"`
}

// smartRecruitersJob mirrors a single SmartRecruiters job object.
type smartRecruitersJob struct {
	Title     string `json:"title"`
	Location  string `json:"location"`
	Type      string `json:"type"`
	Ref       string `json:"ref"`
	CreatedOn string `json:"createdOn"`
}

// NewSmartRecruitersSource creates a SmartRecruiters source with corporate ATS pool.
func NewSmartRecruitersSource() *SmartRecruitersSource {
	return &SmartRecruitersSource{
		BaseSource: BaseSource{NameStr: "smartrecruiters", BaseURL: "https://api.smartrecruiters.com"},
		client: &http.Client{
			Timeout: 30 * time.Second,
		},
		orgs: []smartRecruitersOrg{
			{Org: "netflix", Company: "Netflix"},
			{Org: "spotify", Company: "Spotify"},
			{Org: "booking", Company: "Booking.com"},
			{Org: "adobe", Company: "Adobe"},
		},
	}
}

// Fetch retrieves private jobs from all configured SmartRecruiters companies.
func (s *SmartRecruitersSource) Fetch(ctx context.Context) ([]PrivJobSource, error) {
	log.Info().Msg("Starting crawl for source: SmartRecruiters")

	var allJobs []PrivJobSource
	for _, org := range s.orgs {
		jobs, err := s.fetchOrg(ctx, org)
		if err != nil {
			log.Warn().Err(err).Str("org", org.Org).Msg("SmartRecruiters fetch failed for org")
			continue
		}
		allJobs = append(allJobs, jobs...)
	}

	log.Info().Int("totalJobs", len(allJobs)).Msg("SmartRecruiters fetch completed")
	return allJobs, nil
}

func (s *SmartRecruitersSource) fetchOrg(ctx context.Context, org smartRecruitersOrg) ([]PrivJobSource, error) {
	url := fmt.Sprintf("https://api.smartrecruiters.com/v1/companies/%s/postings", org.Org)
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
		return nil, fmt.Errorf("smartrecruiters API returned status: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var sr smartRecruitersResponse
	if err := json.Unmarshal(body, &sr); err != nil {
		return nil, fmt.Errorf("failed to parse smartrecruiters JSON: %w", err)
	}

	var jobs []PrivJobSource
	for i := range sr.Content {
		j := &sr.Content[i]
		job := s.jobToPriv(j, org.Company)
		if isValidPrivJob(&job) {
			jobs = append(jobs, job)
		}
	}

	log.Info().Int("jobs", len(jobs)).Str("org", org.Org).Msg("SmartRecruiters org fetch successful")
	return jobs, nil
}

// jobToPriv converts a SmartRecruiters job to the shared PrivJobSource model.
func (s *SmartRecruitersSource) jobToPriv(j *smartRecruitersJob, company string) PrivJobSource {
	job := PrivJobSource{
		Source:      "smartrecruiters",
		Company:     company,
		Title:       strings.TrimSpace(j.Title),
		Location:    strings.TrimSpace(j.Location),
		URL:         fmt.Sprintf("https://www.smartrecruiters.com/%s/jobs/%s", company, j.Ref),
		JobType:     normalizeJobType(j.Type),
		Description: "",
		CreatedAt:   time.Now(),
	}

	if j.CreatedOn != "" {
		if t, err := time.Parse(time.RFC3339, j.CreatedOn); err == nil {
			job.PostedAt = &t
		}
	}

	return job
}

// Name returns the source name.
func (s *SmartRecruitersSource) Name() string {
	return s.NameStr
}

// FetchJobs implements the JobSource interface for compatibility with the engine
func (s *SmartRecruitersSource) FetchJobs() ([]Job, error) {
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
			SourceAttribution: "Source: SmartRecruiters ATS API",
			HashChecksum:      "", // Will be set by engine
		}
		jobs = append(jobs, job)
	}

	return jobs, nil
}
