package gov

import (
	"context"
	"github.com/rojgarsetu/crawler/internal/browser"
	"github.com/rojgarsetu/crawler/internal/sources"
)

type RRBSource struct {
	sources.BaseSource
	engine *sources.HeuristicGovEngine
}

func NewRRBSource(pool *browser.Pool) *RRBSource {
	return &RRBSource{
		BaseSource: sources.BaseSource{NameStr: "rrb", BaseURL: "https://www.rrbapply.gov.in"},
		engine:     sources.NewHeuristicGovEngine(pool),
	}
}

func (s *RRBSource) Fetch(ctx context.Context) ([]sources.GovJobSource, error) {
	jobs, _ := s.engine.ScrapePortal(ctx, "https://rrbapply.gov.in", "RRB", "CENTRAL", "ALL_INDIA")
	return jobs, nil
}

func (s *RRBSource) Name() string {
	return s.NameStr
}
