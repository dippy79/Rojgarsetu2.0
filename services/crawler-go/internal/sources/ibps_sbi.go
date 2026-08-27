package sources

import (
	"context"
	"github.com/rojgarsetu/crawler/internal/browser"
	"github.com/rs/zerolog/log"
)

type IbpsSbiSource struct {
	BaseSource
	engine *HeuristicGovEngine
}

func NewIbpsSbiSource(pool *browser.Pool) *IbpsSbiSource {
	return &IbpsSbiSource{
		BaseSource: BaseSource{NameStr: "ibps_sbi"},
		engine:     NewHeuristicGovEngine(pool),
	}
}

var bankingPortals = []string{
	"https://www.ibps.in",
	"https://bank.sbi/careers",
	"https://www.rbi.org.in",
}

func (s *IbpsSbiSource) Fetch(ctx context.Context) ([]GovJobSource, error) {
	var allJobs []GovJobSource
	for _, url := range bankingPortals {
		jobs, _ := s.engine.ScrapePortal(ctx, url, "Banking", "CENTRAL", "ALL_INDIA")
		allJobs = append(allJobs, jobs...)
	}
	log.Info().Int("jobs", len(allJobs)).Msg("Banking jobs fetch completed")
	return allJobs, nil
}

func (s *IbpsSbiSource) FetchJobs() ([]Job, error) {
	govJobs, _ := s.Fetch(context.Background())
	var jobs []Job
	for _, gj := range govJobs {
		jobs = append(jobs, Job{
			Title:             gj.Title,
			CompanyOrDept:     gj.Department,
			ApplyURL:          gj.ApplyURL,
			SourceAttribution: "Source: Banking Portal",
			HashChecksum:      gj.HashChecksum,
			Category:          gj.Category,
			StateName:         gj.StateName,
		})
	}
	return jobs, nil
}

func (s *IbpsSbiSource) Name() string {
	return s.NameStr
}
