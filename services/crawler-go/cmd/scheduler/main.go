package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strconv"
	"sync"
	"time"

	"github.com/rojgarsetu/crawler/internal/browser"
	"github.com/rojgarsetu/crawler/internal/proxy"
	"github.com/rojgarsetu/crawler/internal/scheduler"
	"github.com/rojgarsetu/crawler/internal/store"
)

var (
	storeInstance *store.PostgresStore
	browserPool   *browser.Pool
	proxyRotator  *proxy.Rotator
	crawlMu       sync.Mutex
)

func main() {
	// ── Database connection ──────────────────────────────────────────────────
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		log.Fatal("DATABASE_URL environment variable is required")
	}
	var err error
	storeInstance, err = store.NewPostgresStore(dsn)
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}
	defer storeInstance.Close()
	log.Println("connected to database")

	// ── Browser pool (optional — Naukri needs it) ────────────────────────────
	// If Chrome is not available, Naukri will log a warning and skip.
	pool, err := browser.NewPool(1)
	if err != nil {
		log.Printf("WARNING: browser pool not available (Naukri will be skipped): %v", err)
	} else {
		browserPool = pool
	}

	// ── Proxy rotator (optional, Naukri uses it) ─────────────────────────────
	proxyRotator = proxy.NewRotator()

	// ── HTTP routes ──────────────────────────────────────────────────────────
	http.HandleFunc("/health", healthHandler)
	http.HandleFunc("/trigger", triggerHandler)
	http.HandleFunc("/crawl", triggerHandler) // alias

	portStr := os.Getenv("PORT")
	if portStr == "" {
		portStr = "8082"
	}

	srv := &http.Server{
		Addr:         ":" + portStr,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 0, // no timeout for long-running crawl
		IdleTimeout:  30 * time.Second,
	}

	// ── Background scheduler ─────────────────────────────────────────────────
	intervalHours := 6
	if v := os.Getenv("CRAWL_INTERVAL_HOURS"); v != "" {
		if n, err := strconv.Atoi(v); err == nil && n > 0 {
			intervalHours = n
		}
	}
	ticker := time.NewTicker(time.Duration(intervalHours) * time.Hour)
	go func() {
		log.Printf("background scheduler will run every %d hours", intervalHours)
		// Run once on startup
		runCrawl()
		for range ticker.C {
			runCrawl()
		}
	}()

	log.Printf("crawler scheduler starting on %s", srv.Addr)
	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("server error: %v", err)
	}
}

// ── Handlers ──────────────────────────────────────────────────────────────────

func healthHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
}

func triggerHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	summary := runCrawl()
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(summary)
}

// ── Crawl runner ──────────────────────────────────────────────────────────────

func runCrawl() scheduler.RunSummary {
	if !crawlMu.TryLock() {
		log.Println("crawl already in progress, skipping")
		return scheduler.RunSummary{}
	}
	defer crawlMu.Unlock()

	log.Println("=== starting crawl run ===")
	ctx := context.Background()
	summary := scheduler.RunAll(ctx, storeInstance, browserPool, proxyRotator)

	log.Printf("=== crawl complete: %d sources run, %d succeeded, %d failed, %d total saved in %s",
		summary.SourcesRun, summary.Succeeded, summary.Failed, summary.TotalSaved, summary.Duration)
	for _, r := range summary.SourceResults {
		if r.Error != "" {
			log.Printf("  [FAIL] %s: %s", r.Name, r.Error)
		} else {
			log.Printf("  [OK]   %s: fetched=%d saved=%d", r.Name, r.Fetched, r.Saved)
		}
	}
	return summary
}
