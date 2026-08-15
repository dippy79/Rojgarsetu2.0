package sources

import (
	"log"
)

// SSCScraper scrapes SSC official portal
type SSCScraper struct {
	client interface{} // interface{} to avoid import cycle
}

// NewSSCScraper creates a new SSC scraper
func NewSSCScraper(client interface{}) *SSCScraper {
	return &SSCScraper{client: client}
}

// FetchJobs fetches jobs from SSC
func (s *SSCScraper) FetchJobs() ([]Job, error) {
	log.Println("[SSC] Fetching jobs (stub implementation)")

	var jobs []Job
	// Stub job for testing
	job := Job{
		Title:             "Combined Graduate Level Examination 2025",
		CompanyOrDept:     "Staff Selection Commission",
		Location:          "All India",
		QualificationReq:  "Bachelor's Degree",
		SalaryOrPayScale:  "As per 7th Pay Commission",
		ApplyURL:          "https://ssc.gov.in/graduate-level-examination",
		SourceAttribution: "Source: SSC Official Portal (ssc.gov.in)",
		HashChecksum:      "", // Will be set by engine
	}
	jobs = append(jobs, job)

	log.Printf("[SSC] Found %d job listings", len(jobs))
	return jobs, nil
}
