package sources

import (
	"context"
	"github.com/rojgarsetu/crawler/internal/browser"
	"github.com/rojgarsetu/crawler/internal/shared"
	"github.com/rs/zerolog/log"
)

type GovUniversalSource struct {
	BaseSource
	engine *HeuristicGovEngine
}

func NewGovUniversalSource(pool *browser.Pool) *GovUniversalSource {
	return &GovUniversalSource{
		BaseSource: BaseSource{NameStr: "gov_universal"},
		engine:     NewHeuristicGovEngine(pool),
	}
}

type GovPortal struct {
	Name      string
	URL       string
	Category  string
	StateName string
}

var portals = []GovPortal{
	{Name: "UPSC", URL: "https://www.upsc.gov.in/examinations/active-exams", Category: "CENTRAL", StateName: "ALL_INDIA"},
	{Name: "SSC", URL: "https://ssc.gov.in/Portal/ResultDetails", Category: "CENTRAL", StateName: "ALL_INDIA"},
	{Name: "NCS", URL: "https://www.ncs.gov.in/pages/all-jobs.aspx", Category: "CENTRAL", StateName: "ALL_INDIA"},
	{Name: "IndiaGov", URL: "https://www.india.gov.in/my-government/jobs", Category: "CENTRAL", StateName: "ALL_INDIA"},
	{Name: "Railway", URL: "https://rrbapply.gov.in", Category: "CENTRAL", StateName: "ALL_INDIA"},
	{Name: "IBPS", URL: "https://www.ibps.in", Category: "CENTRAL", StateName: "ALL_INDIA"},
	{Name: "UPPSC", URL: "https://uppsc.up.nic.in", Category: "STATE", StateName: "Uttar Pradesh"},
	{Name: "BPSC", URL: "https://www.bpsc.bih.nic.in", Category: "STATE", StateName: "Bihar"},
	{Name: "MPSC", URL: "https://mpsc.gov.in", Category: "STATE", StateName: "Maharashtra"},
	{Name: "RPSC", URL: "https://rpsc.rajasthan.gov.in", Category: "STATE", StateName: "Rajasthan"},
}

func (s *GovUniversalSource) Fetch(ctx context.Context) ([]shared.GovJobSource, []shared.GovFormSource, error) {
	var allJobs []shared.GovJobSource
	var allForms []shared.GovFormSource

	for _, p := range portals {
		jobs, forms := s.engine.ScrapePortal(ctx, p.URL, p.Name, p.Category, p.StateName)
		allJobs = append(allJobs, jobs...)
		allForms = append(allForms, forms...)
	}

	log.Info().Int("jobs", len(allJobs)).Int("forms", len(allForms)).Msg("Universal Gov fetch completed")
	return allJobs, allForms, nil
}

func (s *GovUniversalSource) Name() string {
	return s.NameStr
}
