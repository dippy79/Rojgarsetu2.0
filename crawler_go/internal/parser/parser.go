package parser

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/rojgarsetu/crawler/internal/sources"
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

// ParseJob converts JobSource to Job
func ParseJob(source *sources.JobSource) *Job {
	if source == nil {
		return nil
	}

	job := &Job{
		Source:         source.Source,
		Title:          strings.TrimSpace(source.Title),
		Company:        strings.TrimSpace(source.Company),
		Location:       strings.TrimSpace(source.Location),
		JobType:        normalizeJobType(source.JobType),
		Eligibility:    strings.TrimSpace(source.Eligibility),
		Description:    strings.TrimSpace(source.Description),
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
