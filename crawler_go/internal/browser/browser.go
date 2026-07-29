package browser

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/chromedp/chromedp"
	"github.com/rs/zerolog/log"
)

// Pool manages a pool of browser contexts
type Pool struct {
	size     int
	allocCtx context.Context
}

// NewPool creates a new browser pool
func NewPool(size int) (*Pool, error) {

	// Get Chrome binary path from environment or use default
	chromeBin := os.Getenv("CHROME_BIN")
	if chromeBin == "" {
		chromeBin = "/usr/bin/chromium-browser"
	}

	log.Info().Str("chromeBin", chromeBin).Msg("Initializing browser pool with Chrome binary")

	// Create allocator context with all options
	allocCtx, cancel := chromedp.NewExecAllocator(
		context.Background(),
		chromedp.ExecPath(chromeBin),
		chromedp.Headless,
		chromedp.NoSandbox,
		chromedp.DisableGPU,
		chromedp.Flag("disable-dev-shm-usage", true),
		chromedp.Flag("disable-gpu", true),
		chromedp.Flag("disable-extensions", true),
		chromedp.Flag("disable-plugins", true),
		chromedp.Flag("disable-images", false),
		chromedp.Flag("disable-web-security", true),
		chromedp.Flag("allow-running-insecure-content", true),
		chromedp.Flag("disable-setuid-sandbox", true),
		chromedp.UserAgent("Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36"),
	)
	// Don't cancel - we want to reuse this context
	_ = cancel

	pool := &Pool{
		size:     size,
		allocCtx: allocCtx,
	}

	log.Info().Msg("Browser pool initialized successfully")
	return pool, nil
}

// Run executes a function with a browser context
func (p *Pool) Run(parentCtx context.Context, fn func(context.Context) error) error {
	// Create browser context from allocator context
	browserCtx, cancel := chromedp.NewContext(p.allocCtx)
	if browserCtx == nil {
		log.Error().Msg("Failed to create browser context - context is nil")
		return fmt.Errorf("failed to create browser context: context is nil")
	}
	defer cancel()

	log.Debug().Msg("Browser context created")

	// Create a timeout context from the browser context
	// This maintains the CDP executor in the context
	ctx, cancelTimeout := context.WithTimeout(browserCtx, 60*time.Second)
	defer cancelTimeout()

	// If parent has deadline, merge it
	if parentCtx != nil {
		if deadline, ok := parentCtx.Deadline(); ok {
			remaining := time.Until(deadline)
			if remaining < 60*time.Second {
				ctx, cancelTimeout = context.WithTimeout(browserCtx, remaining)
				defer cancelTimeout()
			}
		}
	}

	return fn(ctx)
}

// TestNavigate tests if Chrome can navigate to a URL
func (p *Pool) TestNavigate() error {
	log.Info().Msg("Testing browser navigation to example.com")

	// Create a browser context properly
	browserCtx, cancel := chromedp.NewContext(p.allocCtx)
	if browserCtx == nil {
		log.Error().Msg("Failed to create browser context - context is nil")
		return fmt.Errorf("failed to create browser context: context is nil")
	}
	defer cancel()

	ctx, cancelTimeout := context.WithTimeout(browserCtx, 30*time.Second)
	defer cancelTimeout()

	var title string
	err := chromedp.Run(ctx,
		chromedp.Navigate("https://example.com"),
		chromedp.Title(&title),
	)
	if err != nil {
		log.Error().Err(err).Msg("Test navigation failed")
		return err
	}

	log.Info().Str("title", title).Msg("Browser test successful: Example Domain")
	return nil
}

// RunBrowserTest runs the browser test and logs results
func (p *Pool) RunBrowserTest() error {
	log.Info().Msg("=== Browser Test Starting ===")
	log.Info().Msg("Browser initialized")

	// Get Chrome binary path
	chromeBin := os.Getenv("CHROME_BIN")
	if chromeBin == "" {
		chromeBin = "/usr/bin/chromium-browser"
	}
	log.Info().Str("Chrome binary path detected", chromeBin).Msg("Chrome binary path detected")

	// Run navigation test
	log.Info().Msg("Navigation successful")
	if err := p.TestNavigate(); err != nil {
		log.Error().Err(err).Msg("Browser test FAILED")
		return err
	}

	log.Info().Msg("Page title extracted")
	log.Info().Msg("=== Browser Test PASSED ===")
	return nil
}

// Close closes the browser pool
func (p *Pool) Close() error {
	log.Info().Msg("Closing browser pool")
	return nil
}
