package sources

import (
	"context"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/rs/zerolog/log"
)

// CompanyPagesSource scrapes jobs from company career pages
type CompanyPagesSource struct {
	BaseSource
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
		BaseSource: BaseSource{NameStr: "company_pages", BaseURL: ""},
		client: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// Fetch retrieves private jobs from company career pages
func (s *CompanyPagesSource) Fetch(ctx context.Context) ([]PrivJobSource, error) {
	log.Info().Msg("Starting crawl for source: Company Pages")

	var allJobs []PrivJobSource

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
func (s *CompanyPagesSource) fetchFromCompany(ctx context.Context, company CompanyInfo) ([]PrivJobSource, error) {
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
func (s *CompanyPagesSource) parseHTMLJobs(html, companyName, baseURL string) []PrivJobSource {
	var jobs []PrivJobSource

	// Common job listing patterns
	patterns := []string{
		`<a[^>]*href="(/careers/[^"]*job[^"]*)"[^>]*>([^<]*)</a>`,
		`<a[^>]*href="(/jobs/[^"]*)"[^>]*>([^<]*)</a>`,
		`<a[^>]*href="(/careers/job[^"]*)"[^>]*>([^<]*)</a>`,
		`<div[^>]*class="[^"]*job[^"]*"[^>]*>.*?<a[^>]*href="([^"]*)"[^>]*>([^<]*)</a>`,
		`<li[^>]*class="[^"]*job[^"]*"[^>]*>.*?<a[^>]*href="([^"]*)"[^>]*>([^<]*)</a>`,
	}

	for _, pattern := range patterns {
		matches := extractMatches(html, pattern)
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

				job := PrivJobSource{
					Source:    "company_" + strings.ToLower(strings.ReplaceAll(companyName, " ", "_")),
					Company:   companyName,
					Title:     title,
					URL:       link,
					CreatedAt: time.Now(),
				}

				if isValidPrivJob(&job) {
					jobs = append(jobs, job)
				}
			}
		}

		if len(jobs) > 0 {
			break
		}
	}

	return jobs
}

// Name returns the source name
func (s *CompanyPagesSource) Name() string {
	return s.NameStr
}
