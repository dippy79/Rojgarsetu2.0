package gov

import (
	"context"
	"github.com/rojgarsetu/crawler/internal/browser"
	"github.com/rojgarsetu/crawler/internal/sources"
)

type NCSSource struct {
	sources.BaseSource
	engine *sources.HeuristicGovEngine
}

func NewNCSSource(pool *browser.Pool) *NCSSource {
	return &NCSSource{
		BaseSource: sources.BaseSource{NameStr: "ncs", BaseURL: "https://www.ncs.gov.in"},
		engine:     sources.NewHeuristicGovEngine(pool),
	}
}

func (s *NCSSource) Fetch(ctx context.Context) ([]sources.GovJobSource, error) {
	jobs, _ := s.engine.ScrapePortal(ctx, "https://www.ncs.gov.in/pages/all-jobs.aspx", "NCS", "CENTRAL", "ALL_INDIA")
	return jobs, nil
}

func (s *NCSSource) Name() string {
	return s.NameStr
}
