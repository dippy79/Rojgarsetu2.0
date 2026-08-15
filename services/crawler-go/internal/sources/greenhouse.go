package sources

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
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
	BaseSource
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

// NewGreenhouseSource creates a Greenhouse source with expanded company pool
// including major Indian and global tech companies.
func NewGreenhouseSource() *GreenhouseSource {
	return &GreenhouseSource{
		BaseSource: BaseSource{NameStr: "greenhouse", BaseURL: "https://boards-api.greenhouse.io"},
		client: &http.Client{
			Timeout: 30 * time.Second,
		},
		orgs: []greenhouseOrg{
			{Org: "razorpay", Company: "Razorpay"},
			{Org: "swiggy", Company: "Swiggy"},
			{Org: "zomato", Company: "Zomato"},
			{Org: "phonepe", Company: "PhonePe"},
			{Org: "cred", Company: "Cred"},
			{Org: "meesho", Company: "Meesho"},
			{Org: "groww", Company: "Groww"},
			{Org: "uber", Company: "Uber"},
			{Org: "stripe", Company: "Stripe"},
			{Org: "figma", Company: "Figma"},
		},
	}
}

// Fetch retrieves private jobs from all configured Greenhouse boards.
func (s *GreenhouseSource) Fetch(ctx context.Context) ([]PrivJobSource, error) {
	log.Info().Msg("Starting crawl for source: Greenhouse")

	var allJobs []PrivJobSource
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

func (s *GreenhouseSource) fetchOrg(ctx context.Context, org greenhouseOrg) ([]PrivJobSource, error) {
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

	var jobs []PrivJobSource
	for i := range gh.Jobs {
		j := &gh.Jobs[i]
		job := s.jobToPriv(j, org.Company)
		if isValidPrivJob(&job) {
			jobs = append(jobs, job)
		}
	}

	log.Info().Int("jobs", len(jobs)).Str("org", org.Org).Msg("Greenhouse org fetch successful")
	return jobs, nil
}

// jobToPriv converts a Greenhouse job to the shared PrivJobSource model.
func (s *GreenhouseSource) jobToPriv(j *greenhouseJob, company string) PrivJobSource {
	job := PrivJobSource{
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

// GenerateSHA256Hash creates a SHA-256 hash for deduplication
func (s *GreenhouseSource) GenerateSHA256Hash(title, url string) string {
	hash := sha256.Sum256([]byte(title + url))
	return hex.EncodeToString(hash[:])
}

// Name returns the source name.
func (s *GreenhouseSource) Name() string {
	return s.NameStr
}

// FetchJobs implements the JobSource interface for compatibility with the engine
func (s *GreenhouseSource) FetchJobs() ([]Job, error) {
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
			SourceAttribution: "Source: Greenhouse ATS API",
			HashChecksum:      "", // Will be set by engine
		}
		jobs = append(jobs, job)
	}

	return jobs, nil
}
