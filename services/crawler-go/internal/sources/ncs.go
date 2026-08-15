package sources

import (
	"log"
)

// NCSScraper scrapes NCS Portal
type NCSScraper struct {
	client interface{} // interface{} to avoid import cycle
}

// NewNCSScraper creates a new NCS scraper
func NewNCSScraper(client interface{}) *NCSScraper {
	return &NCSScraper{client: client}
}

// FetchJobs fetches jobs from NCS Portal
func (s *NCSScraper) FetchJobs() ([]Job, error) {
	log.Println("[NCS] Fetching jobs (stub implementation)")

	var jobs []Job
	// Stub job for testing
	job := Job{
		Title:             "Various Government Jobs",
		CompanyOrDept:     "National Career Service",
		Location:          "All India",
		QualificationReq:  "Varies by post",
		SalaryOrPayScale:  "Varies by post",
		ApplyURL:          "https://www.ncs.gov.in/job-listings",
		SourceAttribution: "Source: NCS Portal (ncs.gov.in)",
		HashChecksum:      "", // Will be set by engine
	}
	jobs = append(jobs, job)

	log.Printf("[NCS] Found %d job listings", len(jobs))
	return jobs, nil
}
