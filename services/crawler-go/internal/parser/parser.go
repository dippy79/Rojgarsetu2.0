package parser

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/rojgarsetu/crawler/internal/sources"
)

// Max field lengths for sanitized text
const (
	maxTitleLen       = 500
	maxCompanyLen     = 300
	maxLocationLen    = 300
	maxEligibilityLen = 2000
	maxDescriptionLen = 10000
)

// Job represents a parsed job
type Job struct {
	Source         string    `json:"source"`
	Title          string    `json:"title"`
	Company        string    `json:"company"`
	Location       string    `json:"location"`
	JobType        string    `json:"job_type"`
	SalaryMin      *int      `json:"salary_min,omitempty"`
	SalaryMax      *int      `json:"salary_max,omitempty"`
	Eligibility    string    `json:"eligibility"`
	Description    string    `json:"description"`
	ApplicationURL string    `json:"application_url"`
	PostedAt       time.Time `json:"posted_at"`
}

// ParseJob converts JobSource to Job, sanitizing all text fields
func ParseJob(source *sources.JobSource) *Job {
	if source == nil {
		return nil
	}

	job := &Job{
		Source:         strings.TrimSpace(source.Source),
		Title:          sources.SanitizeString(source.Title, maxTitleLen),
		Company:        sources.SanitizeString(source.Company, maxCompanyLen),
		Location:       sources.SanitizeString(source.Location, maxLocationLen),
		JobType:        normalizeJobType(source.JobType),
		Eligibility:    sources.SanitizeString(source.Eligibility, maxEligibilityLen),
		Description:    sources.SanitizeString(source.Description, maxDescriptionLen),
		ApplicationURL: strings.TrimSpace(source.ApplicationURL),
		PostedAt:       time.Now(),
	}

	if source.PostedAt != nil {
		job.PostedAt = *source.PostedAt
	}

	// Parse salary
	if source.SalaryMin != nil {
		job.SalaryMin = source.SalaryMin
	}
	if source.SalaryMax != nil {
		job.SalaryMax = source.SalaryMax
	}

	return job
}

// ValidateJob validates job data
func ValidateJob(job *Job) error {
	if job == nil {
		return fmt.Errorf("job is nil")
	}
	if job.Title == "" {
		return fmt.Errorf("job title is required")
	}
	if job.Source == "" {
		return fmt.Errorf("job source is required")
	}
	return nil
}

// ParseSalary extracts salary from string
func ParseSalary(salaryStr string) (min, max *int) {
	re := regexp.MustCompile(`(\d+(?:,\d+)*)`)
	matches := re.FindAllString(salaryStr, -1)

	if len(matches) == 0 {
		return nil, nil
	}

	// Remove commas and convert to int
	clean := func(s string) int {
		s = strings.ReplaceAll(s, ",", "")
		n, _ := strconv.Atoi(s)
		return n
	}

	minVal := clean(matches[0])
	min = &minVal

	if len(matches) > 1 {
		maxVal := clean(matches[len(matches)-1])
		max = &maxVal
	}

	return min, max
}

// NormalizeJobType normalizes job type
func normalizeJobType(jobType string) string {
	jt := strings.ToLower(strings.TrimSpace(jobType))

	typeMappings := map[string]string{
		"full time":   "full-time",
		"fulltime":    "full-time",
		"part time":   "part-time",
		"parttime":    "part-time",
		"contract":    "contract",
		"contractual": "contract",
		"internship":  "internship",
		"intern":      "internship",
		"temporary":   "temporary",
		"permanent":   "permanent",
	}

	for k, v := range typeMappings {
		if strings.Contains(jt, k) {
			return v
		}
	}

	return "full-time" // default
}
