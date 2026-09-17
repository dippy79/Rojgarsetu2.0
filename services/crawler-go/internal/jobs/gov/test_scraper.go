package gov

import (
	"context"
	"fmt"
	"log"
	"github.com/rojgarsetu/crawler/internal/browser"
)

func TestUPSC() {
	pool, _ := browser.NewPool(1)
	defer pool.Close()
	s := NewUPSCSource(pool)
	jobs, err := s.Fetch(context.Background())
	if err != nil {
		log.Fatalf("UPSC Fetch failed: %v", err)
	}
	fmt.Printf("UPSC Scraped %d jobs\n", len(jobs))
	for _, j := range jobs {
		fmt.Printf("- %s (%s)\n", j.Title, j.ApplyURL)
	}
}
