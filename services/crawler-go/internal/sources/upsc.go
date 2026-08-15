package sources

import (
	"log"
)

// UPSCScraper scrapes UPSC official portal
type UPSCScraper struct {
	client interface{} // interface{} to avoid import cycle
}

// NewUPSCScraper creates a new UPSC scraper
func NewUPSCScraper(client interface{}) *UPSCScraper {
	return &UPSCScraper{client: client}
}

// FetchJobs fetches jobs from UPSC
func (s *UPSCScraper) FetchJobs() ([]Job, error) {
	// For now, we'll skip actual fetching to avoid the import cycle
	// In production, this would use the polite client to fetch and parse
	log.Println("[UPSC] Fetching jobs (stub implementation)")

	var jobs []Job
	// Stub job for testing
	job := Job{
		Title:             "Civil Services Examination 2025",
		CompanyOrDept:     "Union Public Service Commission",
		Location:          "All India",
		QualificationReq:  "Bachelor's Degree",
		SalaryOrPayScale:  "As per 7th Pay Commission",
		ApplyURL:          "https://www.upsc.gov.in/examinations/civil-services-examination",
		SourceAttribution: "Source: UPSC Official Portal (upsc.gov.in)",
		HashChecksum:      "", // Will be set by engine
	}
	jobs = append(jobs, job)

	log.Printf("[UPSC] Found %d job listings", len(jobs))
	return jobs, nil
}
