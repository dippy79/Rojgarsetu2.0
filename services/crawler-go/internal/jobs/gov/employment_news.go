package gov

import (
	"context"
	"github.com/rojgarsetu/crawler/internal/browser"
	"github.com/rojgarsetu/crawler/internal/sources"
)

type EmploymentNewsSource struct {
	sources.BaseSource
	engine *sources.HeuristicGovEngine
}

func NewEmploymentNewsSource(pool *browser.Pool) *EmploymentNewsSource {
	return &EmploymentNewsSource{
		BaseSource: sources.BaseSource{NameStr: "employment_news", BaseURL: "https://www.employmentnews.gov.in"},
		engine:     sources.NewHeuristicGovEngine(pool),
	}
}

func (s *EmploymentNewsSource) Fetch(ctx context.Context) ([]sources.GovJobSource, error) {
	jobs, _ := s.engine.ScrapePortal(ctx, "https://www.employmentnews.gov.in/NewEmp/Weekly_Vacancy.aspx", "Employment News", "CENTRAL", "ALL_INDIA")
	return jobs, nil
}

func (s *EmploymentNewsSource) Name() string {
	return s.NameStr
}
