package sources

import (
	"context"
	"fmt"
	"strings"

	"github.com/mxschmitt/playwright-go"
	"github.com/rojgarsetu/crawler/internal/browser"
)

type NaukriScraper struct {
	BaseSource
	browserPool *browser.Pool
}

func NewNaukriScraper(pool *browser.Pool) *NaukriScraper {
	return &NaukriScraper{
		BaseSource:  BaseSource{NameStr: "naukri", BaseURL: "https://www.naukri.com"},
		browserPool: pool,
	}
}

func (s *NaukriScraper) FetchJobs() ([]Job, error) {
	var jobs []Job

	err := s.browserPool.Run(context.Background(), func(page playwright.Page) error {
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
					jobs = append(jobs, Job{
						Title:             strings.TrimSpace(title),
						CompanyOrDept:     strings.TrimSpace(company),
						Location:          strings.TrimSpace(location),
						ApplyURL:          url,
						SourceAttribution: "Source: Naukri Official Portal",
					})
				}
			}
		}
		return nil
	})

	return jobs, err
}

func (s *NaukriScraper) Name() string {
	return s.NameStr
}
