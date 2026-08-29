package browser

import (
	"context"
	"fmt"
	"sync"

	"github.com/mxschmitt/playwright-go"
	"github.com/rs/zerolog/log"
)

// Pool manages a pool of browser contexts using Playwright
type Pool struct {
	pw      *playwright.Playwright
	browser playwright.Browser
	mu      sync.Mutex
}

// NewPool creates a new Playwright-based browser pool
func NewPool(size int) (*Pool, error) {
	// pw.Run() starts the playwright driver.
	// Installation is handled during Docker build to ensure zero runtime latency.
	pw, err := playwright.Run()
	if err != nil {
		return nil, fmt.Errorf("could not start playwright: %w", err)
	}

	browser, err := pw.Chromium.Launch(playwright.BrowserTypeLaunchOptions{
		Headless: playwright.Bool(true),
		Args: []string{
			"--no-sandbox",
			"--disable-setuid-sandbox",
			"--disable-dev-shm-usage",
			"--disable-gpu",
		},
	})
	if err != nil {
		pw.Stop()
		return nil, fmt.Errorf("could not launch browser: %w", err)
	}

	log.Info().Msg("Playwright browser pool initialized successfully")
	return &Pool{
		pw:      pw,
		browser: browser,
	}, nil
}

// Run executes a function with a fresh Playwright page
func (p *Pool) Run(ctx context.Context, fn func(playwright.Page) error) error {
	p.mu.Lock()
	// Create a new context for each run to isolate cookies/cache
	browserContext, err := p.browser.NewContext(playwright.BrowserNewContextOptions{
		UserAgent: playwright.String("Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36"),
	})
	p.mu.Unlock()

	if err != nil {
		return fmt.Errorf("could not create browser context: %w", err)
	}
	defer browserContext.Close()

	page, err := browserContext.NewPage()
	if err != nil {
		return fmt.Errorf("could not create page: %w", err)
	}
	defer page.Close()

	// Handle panic within the scraper logic
	var runErr error
	func() {
		defer func() {
			if r := recover(); r != nil {
				log.Error().Interface("panic", r).Msg("Recovered from panic in browser runner")
				runErr = fmt.Errorf("scraper panic: %v", r)
			}
		}()
		runErr = fn(page)
	}()

	return runErr
}

// Close shuts down the playwright instance
func (p *Pool) Close() error {
	log.Info().Msg("Closing browser pool")
	if p.browser != nil {
		p.browser.Close()
	}
	if p.pw != nil {
		return p.pw.Stop()
	}
	return nil
}
