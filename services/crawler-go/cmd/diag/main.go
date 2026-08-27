package main

import (
	"context"
	"fmt"
	"time"

	"github.com/rojgarsetu/crawler/internal/browser"
	"github.com/rojgarsetu/crawler/internal/jobs/gov"
	"github.com/rojgarsetu/crawler/internal/jobs/priv"
	"github.com/rojgarsetu/crawler/internal/shared"
)

func main() {
	fmt.Println("=== RojgarSetu 2.0 Crawler Diagnostic (Dry-Run) ===")
	ctx := context.Background()

	pool, err := browser.NewPool(1)
	if err != nil {
		fmt.Printf("[-] Failed to init browser pool: %v\n", err)
		return
	}
	defer pool.Close()

	govSources := []shared.GovJobFetcher{
		gov.NewUPSCSource(pool),
		gov.NewRRBSource(pool),
		gov.NewSSCSource(pool),
	}

	for _, s := range govSources {
		fmt.Printf("[DIAG] Testing %s...\n", s.Name())
		start := time.Now()
		items, err := s.Fetch(ctx)
		duration := time.Since(start)

		if err != nil {
			fmt.Printf("  [-] FAILED: %v\n", err)
		} else {
			fmt.Printf("  [+] SUCCESS: Found %d items in %s\n", len(items), duration)
			if len(items) > 0 {
				fmt.Printf("  [!] Sample: %s (%s)\n", items[0].Title, items[0].ApplyURL)
			}
		}
	}

	privSources := []shared.PrivJobFetcher{
		priv.NewIndeedSource(),
		priv.NewApnaSource(),
	}

	for _, s := range privSources {
		fmt.Printf("[DIAG] Testing %s...\n", s.Name())
		start := time.Now()
		items, err := s.Fetch(ctx)
		duration := time.Since(start)

		if err != nil {
			fmt.Printf("  [-] FAILED: %v\n", err)
		} else {
			fmt.Printf("  [+] SUCCESS: Found %d items in %s\n", len(items), duration)
		}
	}

	fmt.Println("=== Diagnostic Complete ===")
}
