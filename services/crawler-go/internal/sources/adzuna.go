package sources

import (
	"log"
	"os"
)

// AdzunaScraper fetches jobs from Adzuna API
type AdzunaScraper struct {
	client interface{} // interface{} to avoid import cycle
}

// NewAdzunaScraper creates a new Adzuna scraper
func NewAdzunaScraper(client interface{}) *AdzunaScraper {
	return &AdzunaScraper{client: client}
}

// FetchJobs fetches jobs from Adzuna API (graceful skip if no API key)
func (s *AdzunaScraper) FetchJobs() ([]Job, error) {
	appID := os.Getenv("ADZUNA_APP_ID")
	appKey := os.Getenv("ADZUNA_APP_KEY")

	if appID == "" || appKey == "" {
		log.Println("[Adzuna] API keys not configured, skipping")
		return []Job{}, nil
	}

	log.Println("[Adzuna] Fetching jobs (stub implementation)")

	var jobs []Job
	// Stub job for testing
	job := Job{
		Title:             "Software Engineer",
		CompanyOrDept:     "Tech Company",
		Location:          "Bangalore",
		QualificationReq:  "B.Tech in CS",
		SalaryOrPayScale:  "10-15 LPA",
		ApplyURL:          "https://example.com/job/123",
		SourceAttribution: "Source: Adzuna API",
		HashChecksum:      "", // Will be set by engine
	}
	jobs = append(jobs, job)

	log.Printf("[Adzuna] Found %d job listings", len(jobs))
	return jobs, nil
}
