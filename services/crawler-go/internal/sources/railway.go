package sources

import (
	"context"
	"github.com/rojgarsetu/crawler/internal/browser"
)

type RailwayScraper struct {
	BaseSource
	engine *HeuristicGovEngine
}

func NewRailwayScraper(pool *browser.Pool) *RailwayScraper {
	return &RailwayScraper{
		BaseSource: BaseSource{NameStr: "railway"},
		engine:     NewHeuristicGovEngine(pool),
	}
}

func (s *RailwayScraper) FetchJobs() ([]Job, error) {
	govJobs, _ := s.engine.ScrapePortal(context.Background(), "https://rrbapply.gov.in", "Railway RRB", "CENTRAL", "ALL_INDIA")

	var jobs []Job
	for _, gj := range govJobs {
		jobs = append(jobs, Job{
			Title:             gj.Title,
			CompanyOrDept:     gj.Department,
			ApplyURL:          gj.ApplyURL,
			SourceAttribution: "Source: Railway RRB Official Portal (rrbapply.gov.in)",
			HashChecksum:      gj.HashChecksum,
		})
	}
	return jobs, nil
}

func (s *RailwayScraper) Name() string {
	return s.NameStr
}
