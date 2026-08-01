package sources

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/chromedp/chromedp"
	"github.com/rojgarsetu/crawler/internal/browser"
	"github.com/rojgarsetu/crawler/internal/proxy"
	"github.com/rs/zerolog/log"
)

// NaukriSource scrapes jobs from Naukri.com
type NaukriSource struct {
	BaseSource
	browserPool *browser.Pool
	proxy       *proxy.Rotator
}

// NewNaukriSource creates a new Naukri source
func NewNaukriSource(browserPool *browser.Pool, proxy *proxy.Rotator) *NaukriSource {
	return &NaukriSource{
		BaseSource:  BaseSource{NameStr: "naukri", BaseURL: "https://www.naukri.com"},
		browserPool: browserPool,
		proxy:       proxy,
	}
}

// Fetch retrieves job listings from Naukri
func (s *NaukriSource) Fetch(ctx context.Context) ([]JobSource, error) {
	// For Naukri, we need browser automation due to dynamic content
	if s.browserPool == nil {
		log.Error().Msg("Browser pool is nil, cannot fetch jobs")
		return nil, fmt.Errorf("browser pool required for Naukri")
	}

	log.Info().Msg("Starting crawl for source: naukri")

	var jobs []JobSource

	// Pass the original context - browser pool will handle timeout
	err := s.browserPool.Run(ctx, func(ctx context.Context) error {
		log.Info().Msg("Navigating to Naukri jobs page")

		// First test with example.com to verify browser works
		var testTitle string
		if err := chromedp.Run(ctx,
			chromedp.Navigate("https://example.com"),
			chromedp.Title(&testTitle),
		); err != nil {
			log.Error().Err(err).Msg("Test navigation to example.com failed")
			return fmt.Errorf("test navigation failed: %w", err)
		}
		log.Info().Str("testTitle", testTitle).Msg("Test navigation successful - browser is working")

		// Now navigate to Naukri
		log.Info().Msg("Navigating to naukri.com/jobs-in-india")
		log.Info().Str("URL", "https://www.naukri.com/jobs-in-india").Msg("URL being crawled")
		if err := chromedp.Run(ctx,
			chromedp.Navigate("https://www.naukri.com/jobs-in-india"),
		); err != nil {
			log.Error().Err(err).Msg("Failed to navigate to Naukri")
			return fmt.Errorf("failed to navigate: %w", err)
		}

		log.Info().Msg("Page loaded successfully")

		// Wait for content to fully load using chromedp proper wait conditions
		// instead of a hardcoded time.Sleep. Wait for the target job card selectors.
		selectors := []string{".jobTuple", ".srp-jobtuple", "[data-job-id]", ".styles_jobTuple", ".job-card", ".job-card-wrapper", ".tuple", "article", ".jobTupleHeader"}
		waitCtx, waitCancel := context.WithTimeout(ctx, 10*time.Second)
		for _, sel := range selectors {
			err := chromedp.Run(waitCtx,
				chromedp.WaitVisible(sel, chromedp.ByQuery),
			)
			if err == nil {
				log.Debug().Str("selector", sel).Msg("Found visible element via WaitVisible")
				break
			}
		}
		waitCancel()

		// Get page content using chromedp.Run
		var html string
		if err := chromedp.Run(ctx,
			chromedp.OuterHTML("html", &html),
		); err != nil {
			log.Error().Err(err).Msg("Failed to get HTML")
			return fmt.Errorf("failed to get HTML: %w", err)
		}

		log.Info().Int("htmlLength", len(html)).Msg("Got HTML content")

		// Parse with goquery
		doc, err := goquery.NewDocumentFromReader(strings.NewReader(html))
		if err != nil {
			log.Error().Err(err).Msg("Failed to parse HTML")
			return fmt.Errorf("failed to parse HTML: %w", err)
		}

		// Extract job listings - Updated selectors for Naukri Next.js (2026)
		// Naukri uses classes: .jobTuple, .comp-name, .locWdth
		jobSelectors := []string{
			".jobTuple",
			".srp-jobtuple",
			"[data-job-id]",
			".styles_jobTuple",
			".job-card",
			".job-card-wrapper",
			".tuple",
			"article",
			".jobTupleHeader",
		}

		for _, selector := range jobSelectors {
			doc.Find(selector).Each(func(i int, sel *goquery.Selection) {
				// Updated selectors for Naukri 2026 Next.js structure
				// Job title: .title or a.title
				title := sel.Find(".title, a.title").First().Text()
				if title == "" {
					title = sel.Find(".job-title, [title], h2, h3").First().Text()
				}
				if title == "" {
					title = sel.AttrOr("title", "")
				}
				if title == "" {
					title = sel.Text()
				}

				if title != "" && len(title) > 2 && len(title) < 200 {
					// Job URL: get href from a.title or first anchor
					url, _ := sel.Find("a.title").First().Attr("href")
					if url == "" {
						url, _ = sel.Find("a").First().Attr("href")
					}
					if url == "" {
						url, _ = sel.Attr("href")
					}

					// Company name: .comp-name (Naukri 2026 Next.js class)
					company := sel.Find(".comp-name").First().Text()
					if company == "" {
						// Fallback to older selectors
						company = sel.Find(".companyInfo, .company, .subTitle, .company-name").First().Text()
					}

					// Location: .locWdth (Naukri 2026 Next.js class)
					location := sel.Find(".locWdth").First().Text()
					if location == "" {
						// Fallback to older selectors
						location = sel.Find(".location, .locationTxt, .location-info, .locationTag").First().Text()
					}

					jobs = append(jobs, JobSource{
						Source:         "naukri",
						Title:          strings.TrimSpace(title),
						Company:        strings.TrimSpace(company),
						Location:       strings.TrimSpace(location),
						ApplicationURL: url,
					})
				}
			})
			if len(jobs) > 0 {
				log.Info().Str("selector", selector).Int("jobsFound", len(jobs)).Msg("Found jobs with selector")
				break
			}
		}

		log.Info().Int("jobsExtracted", len(jobs)).Msg("Job elements found")

		// Print sample job data if jobs were found
		if len(jobs) > 0 {
			log.Info().Msg("Sample Job Found")
			log.Info().Str("Title", jobs[0].Title).Msg("Sample Job Found - Title")
			log.Info().Str("Company", jobs[0].Company).Msg("Sample Job Found - Company")
			log.Info().Str("Location", jobs[0].Location).Msg("Sample Job Found - Location")
		} else {
			log.Warn().Msg("No jobs extracted. Check selectors.")
		}

		// If still no jobs, log the page title for debugging
		if len(jobs) == 0 {
			var pageTitle string
			chromedp.Run(ctx, chromedp.Title(&pageTitle))
			log.Warn().Str("pageTitle", pageTitle).Msg("No jobs found, page title captured for debugging")

			// Log first 500 chars of HTML for debugging
			if len(html) > 500 {
				log.Debug().Str("htmlPreview", html[:500]).Msg("HTML preview")
			}
		}

		return nil
	})

	if err != nil {
		log.Error().Err(err).Msg("Error during fetch")
		return nil, err
	}

	log.Info().Int("totalJobs", len(jobs)).Msg("Fetch completed")
	return jobs, nil
}

// FetchDetails fetches detailed job information
func (s *NaukriSource) FetchDetails(ctx context.Context) (*JobSource, error) {
	// Implementation for fetching individual job details
	return nil, nil
}

// Crawl performs browser-based crawling of Naukri jobs
func (s *NaukriSource) Crawl(ctx context.Context) error {
	if s.browserPool == nil {
		return fmt.Errorf("browser pool required for Naukri")
	}

	return s.browserPool.Run(ctx, func(ctx context.Context) error {
		var html string

		// Load page and get HTML
		err := chromedp.Run(ctx,
			chromedp.Navigate("https://www.naukri.com/jobs-in-india"),
			chromedp.WaitVisible(".jobTuple, .srp-jobtuple, [data-job-id], .styles_jobTuple, .job-card, .job-card-wrapper, .tuple, article, .jobTupleHeader", chromedp.ByQueryAll),
			chromedp.OuterHTML("html", &html),
		)

		if err != nil {
			return err
		}

		// Parse HTML
		doc, err := goquery.NewDocumentFromReader(strings.NewReader(html))
		if err != nil {
			return err
		}

		// Extract job listings
		doc.Find(".jobTuple").Each(func(i int, sel *goquery.Selection) {
			title := sel.Find("a.title").Text()
			company := sel.Find(".subTitle").Text()

			log.Info().Str("title", title).Str("company", company).Msg("Found job")
		})

		return nil
	})
}

