package sources

import (
	"context"
	"github.com/rojgarsetu/crawler/internal/browser"
	"github.com/rs/zerolog/log"
)

type EmploymentNewsSource struct {
	BaseSource
	engine *HeuristicGovEngine
}

func NewEmploymentNewsSource(pool *browser.Pool) *EmploymentNewsSource {
	return &EmploymentNewsSource{
		BaseSource: BaseSource{NameStr: "employment_news"},
		engine:     NewHeuristicGovEngine(pool),
	}
}

func (s *EmploymentNewsSource) Fetch(ctx context.Context) ([]GovJobSource, error) {
	jobs, _ := s.engine.ScrapePortal(ctx, "https://www.employmentnews.gov.in/NewEmp/Weekly_Vacancy.aspx", "Employment News", "CENTRAL", "ALL_INDIA")
	log.Info().Int("jobs", len(jobs)).Msg("Employment News jobs fetch completed")
	return jobs, nil
}

func (s *EmploymentNewsSource) Name() string {
	return s.NameStr
}
