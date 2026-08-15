package sources

import (
	"log"
	"os"
)

// JoobleScraper fetches jobs from Jooble API
type JoobleScraper struct {
	client interface{} // interface{} to avoid import cycle
}

// NewJoobleScraper creates a new Jooble scraper
func NewJoobleScraper(client interface{}) *JoobleScraper {
	return &JoobleScraper{client: client}
}

// FetchJobs fetches jobs from Jooble API (graceful skip if no API key)
func (s *JoobleScraper) FetchJobs() ([]Job, error) {
	apiKey := os.Getenv("JOOBLE_API_KEY")

	if apiKey == "" {
		log.Println("[Jooble] API key not configured, skipping")
		return []Job{}, nil
	}

	log.Println("[Jooble] Fetching jobs (stub implementation)")

	var jobs []Job
	// Stub job for testing
	job := Job{
		Title:             "Data Analyst",
		CompanyOrDept:     "Analytics Firm",
		Location:          "Mumbai",
		QualificationReq:  "MBA/Statistics",
		SalaryOrPayScale:  "8-12 LPA",
		ApplyURL:          "https://example.com/job/456",
		SourceAttribution: "Source: Jooble API",
		HashChecksum:      "", // Will be set by engine
	}
	jobs = append(jobs, job)

	log.Printf("[Jooble] Found %d job listings", len(jobs))
	return jobs, nil
}
