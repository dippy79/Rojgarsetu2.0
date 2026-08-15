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

// RemoteOKSource fetches jobs from RemoteOK API.
// It consumes the public RemoteOK API:
//
//	https://remoteok.com/api
//
// which returns a JSON array of job objects with fields:
//
//	{ "position", "company", "location", "url", "date", "tags" }
type RemoteOKSource struct {
	BaseSource
	client *http.Client
}

// remoteOKJob mirrors a single RemoteOK job object.
type remoteOKJob struct {
	Position string   `json:"position"`
	Company  string   `json:"company"`
	Location string   `json:"location"`
	URL      string   `json:"url"`
	Date     string   `json:"date"`
	Tags     []string `json:"tags"`
}

// NewRemoteOKSource creates a RemoteOK source.
func NewRemoteOKSource() *RemoteOKSource {
	return &RemoteOKSource{
		BaseSource: BaseSource{NameStr: "remoteok", BaseURL: "https://remoteok.com"},
		client: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// Fetch retrieves remote jobs from RemoteOK API.
func (s *RemoteOKSource) Fetch(ctx context.Context) ([]PrivJobSource, error) {
	log.Info().Msg("Starting crawl for source: RemoteOK")

	url := "https://remoteok.com/api"
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
		return nil, fmt.Errorf("remoteok API returned status: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var jobs []remoteOKJob
	if err := json.Unmarshal(body, &jobs); err != nil {
		return nil, fmt.Errorf("failed to parse remoteok JSON: %w", err)
	}

	var privJobs []PrivJobSource
	for i := range jobs {
		j := &jobs[i]
		// Skip the first item as it's metadata
		if j.Position == "" {
			continue
		}

		job := s.jobToPriv(j)
		if isValidPrivJob(&job) {
			privJobs = append(privJobs, job)
		}
	}

	log.Info().Int("totalJobs", len(privJobs)).Msg("RemoteOK fetch completed")
	return privJobs, nil
}

// jobToPriv converts a RemoteOK job to the shared PrivJobSource model.
func (s *RemoteOKSource) jobToPriv(j *remoteOKJob) PrivJobSource {
	job := PrivJobSource{
		Source:      "remoteok",
		Company:     strings.TrimSpace(j.Company),
		Title:       strings.TrimSpace(j.Position),
		Location:    strings.TrimSpace(j.Location),
		URL:         j.URL,
		JobType:     "Remote",
		Description: "",
		CreatedAt:   time.Now(),
	}

	// Parse date if available
	if j.Date != "" {
		if t, err := time.Parse(time.RFC3339, j.Date); err == nil {
			job.PostedAt = &t
		}
	}

	// Add location tag for India/Global remote
	if strings.Contains(strings.ToLower(j.Location), "india") {
		job.Location = "Remote (India)"
	} else {
		job.Location = "Remote (Global)"
	}

	return job
}

// Name returns the source name.
func (s *RemoteOKSource) Name() string {
	return s.NameStr
}

// FetchJobs implements the JobSource interface for compatibility with the engine
func (s *RemoteOKSource) FetchJobs() ([]Job, error) {
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
			SourceAttribution: "Source: RemoteOK API",
			HashChecksum:      "", // Will be set by engine
		}
		jobs = append(jobs, job)
	}

	return jobs, nil
}
