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

// LeverSource fetches jobs from companies that use Lever as their ATS.
// It consumes the public Lever Postings API:
//
//	https://api.lever.co/v0/postings/{org}?mode=json
//
// which returns a JSON array of posting objects with fields:
//
//	{ "text", "hostedUrl", "categories":{"location","commitment","team"},
//	  "createdAt", "company", ... }
type LeverSource struct {
	BaseSource
	client *http.Client
	orgs   []leverOrg
}

// leverOrg is a single Lever posting org to crawl.
type leverOrg struct {
	Org     string // Lever org slug, e.g. "lever"
	Company string // Display company name
}

// leverPosting mirrors a single Lever posting object.
type leverPosting struct {
	Text       string `json:"text"`
	HostedURL  string `json:"hostedUrl"`
	Company    string `json:"company"`
	CreatedAt  int64  `json:"createdAt"` // Unix milliseconds
	Categories struct {
		Location   string `json:"location"`
		Commitment string `json:"commitment"`
		Team       string `json:"team"`
	} `json:"categories"`
}

// NewLeverSource creates a Lever source. The POC org is "lever" (verified live:
// api.lever.co/v0/postings/lever?mode=json returns HTTP 200, currently an empty
// array — i.e. no open postings for that org, which is honest to report).
func NewLeverSource() *LeverSource {
	return &LeverSource{
		BaseSource: BaseSource{NameStr: "lever", BaseURL: "https://api.lever.co"},
		client: &http.Client{
			Timeout: 30 * time.Second,
		},
		orgs: []leverOrg{
			{Org: "lever", Company: "Lever"},
		},
	}
}

// Fetch retrieves private jobs from all configured Lever posting orgs.
func (s *LeverSource) Fetch(ctx context.Context) ([]PrivJobSource, error) {
	log.Info().Msg("Starting crawl for source: Lever")

	var allJobs []PrivJobSource
	for _, org := range s.orgs {
		jobs, err := s.fetchOrg(ctx, org)
		if err != nil {
			log.Warn().Err(err).Str("org", org.Org).Msg("Lever fetch failed for org")
			continue
		}
		allJobs = append(allJobs, jobs...)
	}

	log.Info().Int("totalJobs", len(allJobs)).Msg("Lever fetch completed")
	return allJobs, nil
}

func (s *LeverSource) fetchOrg(ctx context.Context, org leverOrg) ([]PrivJobSource, error) {
	url := fmt.Sprintf("https://api.lever.co/v0/postings/%s?mode=json", org.Org)
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
		return nil, fmt.Errorf("lever API returned status: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var postings []leverPosting
	if err := json.Unmarshal(body, &postings); err != nil {
		return nil, fmt.Errorf("failed to parse lever JSON: %w", err)
	}

	var jobs []PrivJobSource
	for i := range postings {
		p := &postings[i]
		job := s.postingToPriv(p, org.Company)
		if isValidPrivJob(&job) {
			jobs = append(jobs, job)
		}
	}

	log.Info().Int("jobs", len(jobs)).Str("org", org.Org).Msg("Lever org fetch successful")
	return jobs, nil
}

// postingToPriv converts a Lever posting to the shared PrivJobSource model.
func (s *LeverSource) postingToPriv(p *leverPosting, company string) PrivJobSource {
	job := PrivJobSource{
		Source:      "lever",
		Company:     company,
		Title:       strings.TrimSpace(p.Text),
		Location:    strings.TrimSpace(p.Categories.Location),
		URL:         p.HostedURL,
		JobType:     normalizeJobType(p.Categories.Commitment),
		Description: "",
		CreatedAt:   time.Now(),
	}

	if p.CreatedAt > 0 {
		t := time.UnixMilli(p.CreatedAt)
		job.PostedAt = &t
	}

	return job
}

// Name returns the source name.
func (s *LeverSource) Name() string {
	return s.NameStr
}
