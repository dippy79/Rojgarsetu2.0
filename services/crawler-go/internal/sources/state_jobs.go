package sources

import (
	"context"
	"github.com/rojgarsetu/crawler/internal/browser"
	"github.com/rojgarsetu/crawler/internal/shared"
	"github.com/rs/zerolog/log"
)

type StateJobsSource struct {
	BaseSource
	engine *HeuristicGovEngine
}

func NewStateJobsSource(pool *browser.Pool) *StateJobsSource {
	return &StateJobsSource{
		BaseSource: BaseSource{NameStr: "state_jobs"},
		engine:     NewHeuristicGovEngine(pool),
	}
}

type StatePortal struct {
	Name      string
	URL       string
	StateName string
}

var statePortals = []StatePortal{
	{Name: "UPPSC", URL: "https://uppsc.up.nic.in", StateName: "Uttar Pradesh"},
	{Name: "BPSC", URL: "https://www.bpsc.bih.nic.in", StateName: "Bihar"},
	{Name: "MPSC", URL: "https://mpsc.gov.in", StateName: "Maharashtra"},
	{Name: "RPSC", URL: "https://rpsc.rajasthan.gov.in", StateName: "Rajasthan"},
}

func (s *StateJobsSource) Fetch(ctx context.Context) ([]shared.GovJobSource, error) {
	var allJobs []shared.GovJobSource

	for _, p := range statePortals {
		jobs, _ := s.engine.ScrapePortal(ctx, p.URL, p.Name, "STATE", p.StateName)
		allJobs = append(allJobs, jobs...)
	}

	log.Info().Int("jobs", len(allJobs)).Msg("State jobs fetch completed")
	return allJobs, nil
}

func (s *StateJobsSource) Name() string {
	return s.NameStr
}
