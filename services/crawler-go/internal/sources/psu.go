package sources

import (
	"context"
	"github.com/rojgarsetu/crawler/internal/browser"
	"github.com/rs/zerolog/log"
)

type PSUSource struct {
	BaseSource
	engine *HeuristicGovEngine
}

func NewPSUSource(pool *browser.Pool) *PSUSource {
	return &PSUSource{
		BaseSource: BaseSource{NameStr: "psu"},
		engine:     NewHeuristicGovEngine(pool),
	}
}

var psuPortals = []string{
	"https://www.iocl.com/people-careers",
	"https://www.ntpc.co.in/careers",
	"https://www.ongcindia.com/wps/wcm/connect/en/career/",
}

func (s *PSUSource) Fetch(ctx context.Context) ([]GovJobSource, error) {
	var allJobs []GovJobSource
	for _, url := range psuPortals {
		jobs, _ := s.engine.ScrapePortal(ctx, url, "PSU", "CENTRAL", "ALL_INDIA")
		allJobs = append(allJobs, jobs...)
	}
	log.Info().Int("jobs", len(allJobs)).Msg("PSU jobs fetch completed")
	return allJobs, nil
}

func (s *PSUSource) FetchJobs() ([]Job, error) {
	govJobs, _ := s.Fetch(context.Background())
	var jobs []Job
	for _, gj := range govJobs {
		jobs = append(jobs, Job{
			Title:             gj.Title,
			CompanyOrDept:     gj.Department,
			ApplyURL:          gj.ApplyURL,
			SourceAttribution: "Source: PSU Official Portal",
			HashChecksum:      gj.HashChecksum,
			Category:          gj.Category,
			StateName:         gj.StateName,
		})
	}
	return jobs, nil
}

func (s *PSUSource) Name() string {
	return s.NameStr
}
