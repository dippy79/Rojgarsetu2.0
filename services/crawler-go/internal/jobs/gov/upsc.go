package gov

import (
	"context"
	"github.com/rojgarsetu/crawler/internal/browser"
	"github.com/rojgarsetu/crawler/internal/sources"
)

type UPSCSource struct {
	sources.BaseSource
	engine *sources.HeuristicGovEngine
}

func NewUPSCSource(pool *browser.Pool) *UPSCSource {
	return &UPSCSource{
		BaseSource: sources.BaseSource{NameStr: "upsc", BaseURL: "https://www.upsc.gov.in"},
		engine:     sources.NewHeuristicGovEngine(pool),
	}
}

func (s *UPSCSource) Fetch(ctx context.Context) ([]sources.GovJobSource, error) {
	jobs, _ := s.engine.ScrapePortal(ctx, "https://www.upsc.gov.in/examinations/active-exams", "UPSC", "CENTRAL", "ALL_INDIA")
	return jobs, nil
}

func (s *UPSCSource) Name() string {
	return s.NameStr
}
