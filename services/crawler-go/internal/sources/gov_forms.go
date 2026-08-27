package sources

import (
	"context"
	"github.com/rojgarsetu/crawler/internal/browser"
	"github.com/rojgarsetu/crawler/internal/shared"
	"github.com/rs/zerolog/log"
)

type GovFormsSource struct {
	BaseSource
	engine *HeuristicGovEngine
}

func NewGovFormsSource(pool *browser.Pool) *GovFormsSource {
	return &GovFormsSource{
		BaseSource: BaseSource{NameStr: "gov_forms"},
		engine:     NewHeuristicGovEngine(pool),
	}
}

var formsPortals = []string{
	"https://www.india.gov.in/my-government/forms",
	"https://www.upsc.gov.in/examinations/written-result",
	"https://ssc.gov.in/Portal/ResultDetails",
	"https://www.ncs.gov.in",
}

func (s *GovFormsSource) Fetch(ctx context.Context) ([]shared.GovFormSource, error) {
	var allForms []shared.GovFormSource

	for _, url := range formsPortals {
		_, forms := s.engine.ScrapePortal(ctx, url, "Gov Portal", "CENTRAL", "ALL_INDIA")
		allForms = append(allForms, forms...)
	}

	log.Info().Int("forms", len(allForms)).Msg("Gov forms fetch completed")
	return allForms, nil
}

func (s *GovFormsSource) Name() string {
	return s.NameStr
}
