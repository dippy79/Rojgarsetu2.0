package gov

import (
	"context"
	"github.com/rojgarsetu/crawler/internal/browser"
	"github.com/rojgarsetu/crawler/internal/sources"
)

type SSCSource struct {
	sources.BaseSource
	engine *sources.HeuristicGovEngine
}

func NewSSCSource(pool *browser.Pool) *SSCSource {
	return &SSCSource{
		BaseSource: sources.BaseSource{NameStr: "ssc", BaseURL: "https://ssc.gov.in"},
		engine:     sources.NewHeuristicGovEngine(pool),
	}
}

func (s *SSCSource) Fetch(ctx context.Context) ([]sources.GovJobSource, error) {
	jobs, _ := s.engine.ScrapePortal(ctx, "https://ssc.gov.in/Portal/ResultDetails", "SSC", "CENTRAL", "ALL_INDIA")
	return jobs, nil
}

func (s *SSCSource) Name() string {
	return s.NameStr
}
