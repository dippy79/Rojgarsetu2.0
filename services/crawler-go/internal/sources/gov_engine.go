package sources

import (
	"context"
	"crypto/sha256"
	"fmt"
	"net/url"
	"strings"
	"time"

	"github.com/mxschmitt/playwright-go"
	"github.com/rojgarsetu/crawler/internal/browser"
	"github.com/rojgarsetu/crawler/internal/shared"
	"github.com/rs/zerolog/log"
)

// HeuristicGovEngine implements keyword-based link extraction for govt portals
type HeuristicGovEngine struct {
	browserPool *browser.Pool
}

func NewHeuristicGovEngine(pool *browser.Pool) *HeuristicGovEngine {
	return &HeuristicGovEngine{browserPool: pool}
}

var govJobKeywords = []string{
	"recruitment", "notification", "vacancy", "apply", "exam", "admit", "post", "job", "bharti", "counseling", "direct recruitment", "advertisement",
}

var formKeywords = []string{
	"form", "admit card", "hall ticket", "result", "answer key", "application form", "certificate download", "written result",
}

func (e *HeuristicGovEngine) ScrapePortal(ctx context.Context, targetURL, sourceName, category, stateName string) ([]shared.GovJobSource, []shared.GovFormSource) {
	var jobs []shared.GovJobSource
	var forms []shared.GovFormSource

	err := e.browserPool.Run(ctx, func(page playwright.Page) error {
		defer func() {
			if r := recover(); r != nil {
				log.Error().Interface("panic", r).Str("url", targetURL).Msg("Recovered from panic in HeuristicGovEngine")
			}
		}()

		if _, err := page.Goto(targetURL, playwright.PageGotoOptions{
			WaitUntil: playwright.WaitUntilStateNetworkidle,
			Timeout:   playwright.Float(20000),
		}); err != nil {
			log.Warn().Err(err).Str("url", targetURL).Msg("Heuristic engine failed to load page")
			return nil
		}

		links, err := page.QuerySelectorAll("a")
		if err != nil {
			return err
		}

		for _, link := range links {
			text, _ := link.TextContent()
			href, _ := link.GetAttribute("href")

			text = strings.TrimSpace(text)
			text = strings.Join(strings.Fields(text), " ")

			if text == "" || href == "" || strings.HasPrefix(href, "javascript") || strings.HasPrefix(href, "#") {
				continue
			}

			lowerText := strings.ToLower(text)
			absURL := e.resolveURL(targetURL, href)

			isJob := false
			for _, kw := range govJobKeywords {
				if strings.Contains(lowerText, kw) {
					isJob = true
					break
				}
			}

			if isJob {
				job := shared.GovJobSource{
					Source:       strings.ToLower(sourceName),
					Title:        text,
					Department:   sourceName,
					ApplyURL:     absURL,
					Category:     category,
					StateName:    stateName,
					HashChecksum: e.generateHash(text + absURL),
					CreatedAt:    time.Now(),
				}
				jobs = append(jobs, job)
			}

			isForm := false
			for _, kw := range formKeywords {
				if strings.Contains(lowerText, kw) {
					isForm = true
					break
				}
			}

			if isForm {
				formType := e.categorizeForm(lowerText)
				form := shared.GovFormSource{
					Title:           text,
					ConductingBody:  sourceName,
					FormType:        formType,
					OfficialWebsite: absURL,
					Category:        category,
					HashChecksum:    e.generateHash(text + absURL),
				}
				forms = append(forms, form)
			}
		}
		return nil
	})

	if err != nil {
		log.Error().Err(err).Str("url", targetURL).Msg("Heuristic engine run error")
	}

	return jobs, forms
}

func (e *HeuristicGovEngine) resolveURL(base, ref string) string {
	u, err := url.Parse(ref)
	if err != nil || u.IsAbs() {
		return ref
	}
	baseURL, _ := url.Parse(base)
	return baseURL.ResolveReference(u).String()
}

func (e *HeuristicGovEngine) categorizeForm(text string) string {
	if strings.Contains(text, "admit") || strings.Contains(text, "hall ticket") {
		return "ADMIT_CARD"
	}
	if strings.Contains(text, "result") || strings.Contains(text, "merit") || strings.Contains(text, "cutoff") {
		return "RESULT"
	}
	if strings.Contains(text, "form") || strings.Contains(text, "apply") {
		return "APPLICATION_FORM"
	}
	return "NOTIFICATION"
}

func (e *HeuristicGovEngine) generateHash(input string) string {
	h := sha256.Sum256([]byte(strings.ToLower(strings.TrimSpace(input))))
	return fmt.Sprintf("%x", h)
}
