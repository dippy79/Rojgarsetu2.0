package priv

import (
	"context"
	"fmt"
	"strings"

	"github.com/mxschmitt/playwright-go"
	"github.com/rojgarsetu/crawler/internal/browser"
	"github.com/rojgarsetu/crawler/internal/proxy"
	"github.com/rojgarsetu/crawler/internal/shared"
	"github.com/rs/zerolog/log"
)

type NaukriSource struct {
	shared.BaseSource
	browserPool *browser.Pool
	proxy       *proxy.Rotator
}

func NewNaukriSource(browserPool *browser.Pool, proxy *proxy.Rotator) *NaukriSource {
	return &NaukriSource{
		BaseSource:  shared.BaseSource{NameStr: "naukri", BaseURL: "https://www.naukri.com"},
		browserPool: browserPool,
		proxy:       proxy,
	}
}

func (s *NaukriSource) Fetch(ctx context.Context) ([]shared.JobSource, error) {
	if s.browserPool == nil {
		return nil, fmt.Errorf("browser pool required for Naukri")
	}

	log.Info().Msg("Starting Playwright crawl for source: naukri")

	var jobs []shared.JobSource

	err := s.browserPool.Run(ctx, func(page playwright.Page) error {
		log.Info().Msg("Navigating to Naukri jobs page")

		if _, err := page.Goto("https://www.naukri.com/jobs-in-india", playwright.PageGotoOptions{
			WaitUntil: playwright.WaitUntilStateNetworkidle,
			Timeout:   playwright.Float(30000),
		}); err != nil {
			return fmt.Errorf("failed to navigate: %w", err)
		}

		cards, _ := page.QuerySelectorAll(".jobTuple, .srp-jobtuple, article, .job-card")
		for _, card := range cards {
			titleEl, _ := card.QuerySelector(".title, a.title, .job-title")
			companyEl, _ := card.QuerySelector(".comp-name, .companyInfo, .company")
			locationEl, _ := card.QuerySelector(".locWdth, .location, .locationTxt")

			if titleEl != nil {
				title, _ := titleEl.TextContent()
				company, _ := companyEl.TextContent()
				location, _ := locationEl.TextContent()
				url, _ := titleEl.GetAttribute("href")

				if title != "" {
					jobs = append(jobs, shared.JobSource{
						Source:         "naukri",
						Title:          strings.TrimSpace(title),
						Company:        strings.TrimSpace(company),
						Location:       strings.TrimSpace(location),
						ApplicationURL: url,
					})
				}
			}
		}
		return nil
	})

	return jobs, err
}

func (s *NaukriSource) Name() string {
	return s.NameStr
}
