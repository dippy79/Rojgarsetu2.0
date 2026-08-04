package sources

import (
	"context"
	"time"
)

// JobSource represents a job posting from a source
type JobSource struct {
	Source         string
	URL            string
	Title          string
	Company        string
	Location       string
	JobType        string
	SalaryMin      *int
	SalaryMax      *int
	Eligibility    string
	Description    string
	ApplicationURL string
	PostedAt       *time.Time
	Skills         []string
}

// Source interface for job sources
type Source interface {
	Name() string
	Fetch(ctx context.Context) ([]JobSource, error)
	FetchDetails(ctx context.Context) (*JobSource, error)
}

// BaseSource provides common functionality

func (b *BaseSource) Name() string {
	return b.NameStr
}
