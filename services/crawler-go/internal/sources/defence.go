package sources

import (
	"context"
	"github.com/rojgarsetu/crawler/internal/browser"
	"github.com/rs/zerolog/log"
)

type DefenceSource struct {
	BaseSource
	engine *HeuristicGovEngine
}

func NewDefenceSource(pool *browser.Pool) *DefenceSource {
	return &DefenceSource{
		BaseSource: BaseSource{NameStr: "defence"},
		engine:     NewHeuristicGovEngine(pool),
	}
}

var defencePortals = []string{
	"https://joinindianarmy.nic.in",
	"https://afcat.cdac.in",
	"https://joinindiannavy.gov.in",
	"https://joinindiancoastguard.cdac.in",
}

func (s *DefenceSource) Fetch(ctx context.Context) ([]GovJobSource, error) {
	var allJobs []GovJobSource
	for _, url := range defencePortals {
		jobs, _ := s.engine.ScrapePortal(ctx, url, "Defence", "CENTRAL", "ALL_INDIA")
		allJobs = append(allJobs, jobs...)
	}
	log.Info().Int("jobs", len(allJobs)).Msg("Defence jobs fetch completed")
	return allJobs, nil
}

func (s *DefenceSource) FetchJobs() ([]Job, error) {
	govJobs, _ := s.Fetch(context.Background())
	var jobs []Job
	for _, gj := range govJobs {
		jobs = append(jobs, Job{
			Title:             gj.Title,
			CompanyOrDept:     gj.Department,
			ApplyURL:          gj.ApplyURL,
			SourceAttribution: "Source: Defence Official Portal",
			HashChecksum:      gj.HashChecksum,
			Category:          gj.Category,
			StateName:         gj.StateName,
		})
	}
	return jobs, nil
}

func (s *DefenceSource) Name() string {
	return s.NameStr
}
