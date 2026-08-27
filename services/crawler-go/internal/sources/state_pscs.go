package sources

import (
	"context"
	"github.com/rojgarsetu/crawler/internal/browser"
	"github.com/rs/zerolog/log"
)

type StatePSCSource struct {
	BaseSource
	engine *HeuristicGovEngine
}

func NewStatePSCSource(pool *browser.Pool) *StatePSCSource {
	return &StatePSCSource{
		BaseSource: BaseSource{NameStr: "state_pscs"},
		engine:     NewHeuristicGovEngine(pool),
	}
}

var statePSCTargets = []struct {
	Name      string
	URL       string
	StateName string
}{
	{Name: "UPPSC", URL: "https://uppsc.up.nic.in", StateName: "Uttar Pradesh"},
	{Name: "BPSC", URL: "https://www.bpsc.bih.nic.in", StateName: "Bihar"},
	{Name: "MPSC", URL: "https://mpsc.gov.in", StateName: "Maharashtra"},
	{Name: "RPSC", URL: "https://rpsc.rajasthan.gov.in", StateName: "Rajasthan"},
	{Name: "WBPSC", URL: "https://wbpsc.gov.in", StateName: "West Bengal"},
	{Name: "TNPSC", URL: "https://www.tnpsc.gov.in", StateName: "Tamil Nadu"},
	{Name: "JKPSC", URL: "https://jkpsc.nic.in", StateName: "Jammu & Kashmir"},
	{Name: "HPSC", URL: "https://hpsc.gov.in", StateName: "Haryana"},
	{Name: "GPSC", URL: "https://gpsc.gujarat.gov.in", StateName: "Gujarat"},
}

func (s *StatePSCSource) Fetch(ctx context.Context) ([]GovJobSource, error) {
	var allJobs []GovJobSource
	for _, target := range statePSCTargets {
		jobs, _ := s.engine.ScrapePortal(ctx, target.URL, target.Name, "STATE", target.StateName)
		allJobs = append(allJobs, jobs...)
	}
	log.Info().Int("jobs", len(allJobs)).Msg("State PSC jobs fetch completed")
	return allJobs, nil
}

func (s *StatePSCSource) FetchJobs() ([]Job, error) {
	govJobs, _ := s.Fetch(context.Background())
	var jobs []Job
	for _, gj := range govJobs {
		jobs = append(jobs, Job{
			Title:             gj.Title,
			CompanyOrDept:     gj.Department,
			ApplyURL:          gj.ApplyURL,
			SourceAttribution: "Source: State PSC Portal",
			HashChecksum:      gj.HashChecksum,
			Category:          gj.Category,
			StateName:         gj.StateName,
		})
	}
	return jobs, nil
}

func (s *StatePSCSource) Name() string {
	return s.NameStr
}
