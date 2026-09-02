package priv

import (
	"github.com/rojgarsetu/crawler/internal/shared"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"regexp"
	"strings"
	"time"

	"github.com/rs/zerolog/log"
)

// CompanyPagesSource scrapes jobs from company career pages
type CompanyPagesSource struct {
	shared.BaseSource
	client *http.Client
}

// CompanyInfo contains company career page URLs
type CompanyInfo struct {
	Name      string
	CareerURL string
	JobsURL   string
}

// CompanyList contains all company career pages
var CompanyList = []CompanyInfo{
	{Name: "Infosys", CareerURL: "https://www.infosys.com/careers.html", JobsURL: "https://www.infosys.com/careers.html"},
	{Name: "TCS", CareerURL: "https://www.tcs.com/careers", JobsURL: "https://www.tcs.com/careers"},
	{Name: "Wipro", CareerURL: "https://careers.wipro.com/careershome", JobsURL: "https://careers.wipro.com/careershome"},
	{Name: "Amazon", CareerURL: "https://www.amazon.jobs", JobsURL: "https://www.amazon.jobs/en-US/landing?base_query=&loc_query=India"},
	{Name: "Flipkart", CareerURL: "https://www.flipkart.com/careers", JobsURL: "https://www.flipkart.com/careers"},
	{Name: "Accenture", CareerURL: "https://www.accenture.com/in-en/careers", JobsURL: "https://www.accenture.com/in-en/careers"},
	{Name: "Cognizant", CareerURL: "https://www.cognizant.com/us-en/careers", JobsURL: "https://www.cognizant.com/us-en/careers"},
	{Name: "Tech Mahindra", CareerURL: "https://careers.techmahindra.com/", JobsURL: "https://careers.techmahindra.com/"},
	{Name: "HCL", CareerURL: "https://www.hcltech.com/careers", JobsURL: "https://www.hcltech.com/careers"},
	{Name: "IBM", CareerURL: "https://www.ibm.com/careers", JobsURL: "https://www.ibm.com/careers/en-in"},
	{Name: "Microsoft", CareerURL: "https://careers.microsoft.com/", JobsURL: "https://careers.microsoft.com/us/en/search"},
	{Name: "Adobe", CareerURL: "https://www.adobe.com/careers.html", JobsURL: "https://www.adobe.com/careers.html"},
	{Name: "Salesforce", CareerURL: "https://careers.salesforce.com/", JobsURL: "https://careers.salesforce.com/"},
	{Name: "Meta", CareerURL: "https://metacareers.com/", JobsURL: "https://metacareers.com/jobs"},
	{Name: "Oracle", CareerURL: "https://www.oracle.com/careers/", JobsURL: "https://www.oracle.com/careers/"},
}

// NewCompanyPagesSource creates a new Company Pages source
func NewCompanyPagesSource() *CompanyPagesSource {
	return &CompanyPagesSource{
		BaseSource: shared.BaseSource{NameStr: "company_pages", BaseURL: ""},
		client: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// Fetch retrieves private jobs from company career pages
func (s *CompanyPagesSource) Fetch(ctx context.Context) ([]shared.PrivJobSource, error) {
	log.Info().Msg("Starting crawl for source: Company Pages")

	var allJobs []shared.PrivJobSource

	for _, company := range CompanyList {
		jobs, err := s.fetchFromCompany(ctx, company)
		if err != nil {
			log.Warn().Err(err).Str("company", company.Name).Msg("Failed to fetch from company")
			continue
		}
		allJobs = append(allJobs, jobs...)
	}

	log.Info().Int("totalJobs", len(allJobs)).Msg("Company Pages fetch completed")
	return allJobs, nil
}

// fetchFromCompany fetches jobs from a specific company
func (s *CompanyPagesSource) fetchFromCompany(ctx context.Context, company CompanyInfo) ([]shared.PrivJobSource, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", company.JobsURL, nil)
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
		return nil, err
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	jobs := s.parseHTMLJobs(string(body), company.Name, company.JobsURL)
	log.Info().Int("jobs", len(jobs)).Str("company", company.Name).Msg("Company fetch successful")
	return jobs, nil
}

// parseHTMLJobs parses jobs from company HTML
func (s *CompanyPagesSource) parseHTMLJobs(html, companyName, baseURL string) []shared.PrivJobSource {
	var jobs []shared.PrivJobSource

	// Common job listing patterns
	patterns := []string{
		`<a[^>]*href="(/careers/[^"]*job[^"]*)"[^>]*>([^<]*)</a>`,
		`<a[^>]*href="(/jobs/[^"]*)"[^>]*>([^<]*)</a>`,
		`<a[^>]*href="(/careers/job[^"]*)"[^>]*>([^<]*)</a>`,
		`<div[^>]*class="[^"]*job[^"]*"[^>]*>.*?<a[^>]*href="([^"]*)"[^>]*>([^<]*)</a>`,
		`<li[^>]*class="[^"]*job[^"]*"[^>]*>.*?<a[^>]*href="([^"]*)"[^>]*>([^<]*)</a>`,
	}

	for _, pattern := range patterns {
		matches := shared.ExtractMatches(html, pattern)
		for _, match := range matches {
			if len(match) >= 3 {
				link := strings.TrimSpace(match[1])
				title := strings.TrimSpace(match[2])

				// Skip if too short or not a job
				if len(title) < 5 {
					continue
				}

				// Make absolute URL
				if !strings.HasPrefix(link, "http") {
					if strings.HasPrefix(link, "/") {
						link = baseURL + link
					} else {
						link = baseURL + "/" + link
					}
				}

				job := shared.PrivJobSource{
					Source:    "company_" + strings.ToLower(strings.ReplaceAll(companyName, " ", "_")),
					Company:   companyName,
					Title:     title,
					URL:       link,
					CreatedAt: time.Now(),
				}

				if shared.IsValidPrivJob(&job) {
					jobs = append(jobs, job)
				}
			}
		}

		if len(jobs) > 0 {
			break
		}
	}

	// If no jobs found via regex (e.g. server-rendered page with JSON-LD),
	// try parsing structured data embedded in <script type="application/ld+json">.
	if len(jobs) == 0 {
		jobs = s.parseJSONLDJobs(html, companyName, baseURL)
	}

	return jobs
}

// jsonLDJob mirrors the subset of schema.org JobPosting we care about.
type jsonLDJob struct {
	Context        string `json:"@context"`
	Type           string `json:"@type"`
	Title          string `json:"title"`
	Description    string `json:"description"`
	DatePosted     string `json:"datePosted"`
	ValidThrough   string `json:"validThrough"`
	EmploymentType string `json:"employmentType"`
	HiringOrg      struct {
		Name string `json:"name"`
	} `json:"hiringOrganization"`
	JobLocation struct {
		Type    string `json:"@type"`
		Address struct {
			AddressLocality string `json:"addressLocality"`
			AddressRegion   string `json:"addressRegion"`
			AddressCountry  string `json:"addressCountry"`
		} `json:"address"`
	} `json:"jobLocation"`
	BaseSalary struct {
		Currency string `json:"currency"`
		Value    struct {
			MinValue float64 `json:"minValue"`
			MaxValue float64 `json:"maxValue"`
		} `json:"value"`
	} `json:"baseSalary"`
}

// parseJSONLDJobs extracts schema.org JobPosting entries from JSON-LD blocks
// embedded in the server-rendered HTML. This is a defensive fallback for
// companies that expose job data via structured data rather than static
// HTML anchors.
func (s *CompanyPagesSource) parseJSONLDJobs(html, companyName, baseURL string) []shared.PrivJobSource {
	var jobs []shared.PrivJobSource

	// Match each <script type="application/ld+json">...</script> block.
	ldRe := regexp.MustCompile(`(?s)<script[^>]*type=["']application/ld\+json["'][^>]*>(.*?)</script>`)
	blocks := ldRe.FindAllStringSubmatch(html, -1)

	for _, block := range blocks {
		if len(block) < 2 {
			continue
		}
		raw := strings.TrimSpace(block[1])

		// Try to unmarshal a single JobPosting object.
		var single jsonLDJob
		if err := json.Unmarshal([]byte(raw), &single); err == nil && single.isJobPosting() {
			if job := s.jsonLDToPrivJob(single, companyName, baseURL); job != nil {
				jobs = append(jobs, *job)
			}
			continue
		}

		// Try to unmarshal an array of JobPosting objects.
		var arr []jsonLDJob
		if err := json.Unmarshal([]byte(raw), &arr); err == nil {
			for _, item := range arr {
				if !item.isJobPosting() {
					continue
				}
				if job := s.jsonLDToPrivJob(item, companyName, baseURL); job != nil {
					jobs = append(jobs, *job)
				}
			}
		}
	}

	if len(jobs) > 0 {
		log.Info().Int("jsonldJobs", len(jobs)).Str("company", companyName).Msg("Parsed jobs from JSON-LD")
	}
	return jobs
}

func (j *jsonLDJob) isJobPosting() bool {
	return j.Type == "JobPosting" && j.Title != ""
}

func (s *CompanyPagesSource) jsonLDToPrivJob(j jsonLDJob, companyName, baseURL string) *shared.PrivJobSource {
	// Resolve the company name: prefer the JSON-LD hiring organization name,
	// falling back to the configured company name.
	company := companyName
	if j.HiringOrg.Name != "" {
		company = j.HiringOrg.Name
	}

	// Derive location from the JSON-LD address, if present.
	location := ""
	if j.JobLocation.Address.AddressLocality != "" {
		location = j.JobLocation.Address.AddressLocality
		if j.JobLocation.Address.AddressRegion != "" {
			location += ", " + j.JobLocation.Address.AddressRegion
		}
	}

	var postedAt *time.Time
	if t, err := time.Parse(time.RFC3339, j.DatePosted); err == nil {
		postedAt = &t
	}

	salary := ""
	if j.BaseSalary.Value.MinValue > 0 || j.BaseSalary.Value.MaxValue > 0 {
		if j.BaseSalary.Value.MinValue > 0 && j.BaseSalary.Value.MaxValue > 0 {
			salary = fmt.Sprintf("%.0f-%.0f %s", j.BaseSalary.Value.MinValue, j.BaseSalary.Value.MaxValue, j.BaseSalary.Currency)
		} else if j.BaseSalary.Value.MinValue > 0 {
			salary = fmt.Sprintf("%.0f %s", j.BaseSalary.Value.MinValue, j.BaseSalary.Currency)
		}
	}

	desc := shared.SanitizeString(j.Description, 2000)

	job := &shared.PrivJobSource{
		Source:      "company_" + strings.ToLower(strings.ReplaceAll(company, " ", "_")),
		Company:     company,
		Title:       j.Title,
		Location:    location,
		URL:         baseURL, // fallback; callers with a real URL should override
		Salary:      salary,
		Experience:  "",
		JobType:     shared.NormalizeJobType(j.EmploymentType),
		Description: desc,
		PostedAt:    postedAt,
		CreatedAt:   time.Now(),
	}

	if shared.IsValidPrivJob(job) {
		return job
	}
	return nil
}

// Name returns the source name
func (s *CompanyPagesSource) Name() string {
	return s.NameStr
}

