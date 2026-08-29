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

	_ "github.com/lib/pq"
	"github.com/rojgarsetu/crawler/internal/browser"
	"github.com/rojgarsetu/crawler/internal/proxy"
	"github.com/rojgarsetu/crawler/internal/scheduler"
	"github.com/rojgarsetu/crawler/internal/store"
)

var (
	pgStore      *store.PostgresStore
	browserPool  *browser.Pool
	proxyRotator *proxy.Rotator
	crawlMu      sync.Mutex
	lastRunTime  time.Time
	jobsToday    int
)

func main() {
	// ── Database connection ──────────────────────────────────────────────────
	databaseURL := os.Getenv("DATABASE_URL")
	if databaseURL == "" {
		log.Fatal("DATABASE_URL environment variable is required")
	}

	var err error
	pgStore, err = store.NewPostgresStore(databaseURL)
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}
	defer pgStore.Close()
	log.Println("connected to database via PostgresStore")

	// ── Browser pool setup ────────────────────────────────────────────────────
	browserPool, err = browser.NewPool(2) // Max 2 concurrent browser instances
	if err != nil {
		log.Fatalf("failed to initialize browser pool: %v", err)
	}
	defer browserPool.Close()
	log.Println("browser pool initialized")

	// ── Proxy rotator setup ───────────────────────────────────────────────────
	proxyRotator = proxy.NewRotator()
	log.Println("proxy rotator initialized")


	// ── HTTP routes ──────────────────────────────────────────────────────────
	http.HandleFunc("/health", healthHandler)

	// Protected routes (Admin only)
	http.Handle("/stats", adminAuth(http.HandlerFunc(statsHandler)))
	http.Handle("/crawl", adminAuth(http.HandlerFunc(crawlHandler)))
	http.Handle("/sources", adminAuth(http.HandlerFunc(sourcesHandler)))
	http.Handle("/forms", adminAuth(http.HandlerFunc(formsHandler)))
	http.Handle("/legal/takedown", adminAuth(http.HandlerFunc(takedownHandler)))

	portStr := os.Getenv("PORT")
	if portStr == "" {
		portStr = "8082"
	}

	srv := &http.Server{
		Addr:         ":" + portStr,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  30 * time.Second,
	}

	// ── Background scheduler ─────────────────────────────────────────────────
	intervalHours := 6
	if v := os.Getenv("CRAWLER_INTERVAL_HOURS"); v != "" {
		if n, err := strconv.Atoi(v); err == nil && n > 0 {
			intervalHours = n
		}
	}

	ticker := time.NewTicker(time.Duration(intervalHours) * time.Hour)
	go func() {
		log.Printf("background scheduler will run every %d hours", intervalHours)

		// Run once on startup
		go runCrawl()

		for range ticker.C {
			go runCrawl()
		}
	}()

	// ── Error recovery ──────────────────────────────────────────────────────
	defer func() {
		if r := recover(); r != nil {
			log.Printf("Crawler recovered from panic: %v", r)
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
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status":     "ok",
		"last_run":   lastRunTime.Format(time.RFC3339),
		"jobs_today": jobsToday,
	})
}

func statsHandler(w http.ResponseWriter, r *http.Request) {
	var totalCrawled, newAdded, duplicates, errors int
	db := pgStore.DB()

	// Get stats from crawler_logs
	db.QueryRow("SELECT COALESCE(SUM(jobs_found), 0) FROM crawler_logs WHERE created_at > NOW() - INTERVAL '24 hours'").Scan(&totalCrawled)
	db.QueryRow("SELECT COALESCE(SUM(jobs_added), 0) FROM crawler_logs WHERE created_at > NOW() - INTERVAL '24 hours'").Scan(&newAdded)
	db.QueryRow("SELECT COALESCE(SUM(duplicates_found), 0) FROM crawler_logs WHERE created_at > NOW() - INTERVAL '24 hours'").Scan(&duplicates)

	// Get error count
	db.QueryRow("SELECT COUNT(*) FROM crawler_logs WHERE status != 'COMPLETED' AND created_at > NOW() - INTERVAL '24 hours'").Scan(&errors)

	// Get active sources
	rows, err := db.Query("SELECT name, category FROM crawler_sources WHERE is_active = true ORDER BY id")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var sources []map[string]string
	for rows.Next() {
		var name, category string
		if err := rows.Scan(&name, &category); err != nil {
			continue
		}
		sources = append(sources, map[string]string{
			"name":     name,
			"category": category,
		})
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"total_crawled": totalCrawled,
		"new_added":     newAdded,
		"duplicates":    duplicates,
		"errors":        errors,
		"sources":       sources,
	})
}

func crawlHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Trigger crawl in background
	go runCrawl()

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"status":  "started",
		"message": "crawl running in background",
	})
}

func sourcesHandler(w http.ResponseWriter, r *http.Request) {
	rows, err := pgStore.DB().Query("SELECT id, name, category, source_type, base_url, is_active FROM crawler_sources ORDER BY id")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var sources []map[string]interface{}
	for rows.Next() {
		var id int
		var name, category, sourceType, baseURL string
		var isActive bool
		if err := rows.Scan(&id, &name, &category, &sourceType, &baseURL, &isActive); err != nil {
			continue
		}
		sources = append(sources, map[string]interface{}{
			"id":          id,
			"name":        name,
			"category":    category,
			"source_type": sourceType,
			"base_url":    baseURL,
			"is_active":   isActive,
		})
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(sources)
}

func formsHandler(w http.ResponseWriter, r *http.Request) {
	rows, err := pgStore.DB().Query("SELECT id, title, conducting_body, form_type, official_website, is_taken_down FROM gov_forms_info ORDER BY created_at DESC LIMIT 50")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var forms []map[string]interface{}
	for rows.Next() {
		var id int
		var title, conductingBody, formType, officialWebsite string
		var isTakenDown bool
		if err := rows.Scan(&id, &title, &conductingBody, &formType, &officialWebsite, &isTakenDown); err != nil {
			continue
		}
		forms = append(forms, map[string]interface{}{
			"id":               id,
			"title":            title,
			"conducting_body":  conductingBody,
			"form_type":        formType,
			"official_website": officialWebsite,
			"is_taken_down":    isTakenDown,
		})
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(forms)
}

func takedownHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var request struct {
		JobID     int    `json:"job_id"`
		FormID    int    `json:"form_id"`
		Requester string `json:"requester"`
		Reason    string `json:"reason"`
	}

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	if request.Requester == "" || request.Reason == "" {
		http.Error(w, "requester and reason are required", http.StatusBadRequest)
		return
	}

	query := `
		INSERT INTO takedown_requests (job_id, form_id, requester, reason, status)
		VALUES ($1, $2, $3, $4, 'PENDING')
	`

	_, err := pgStore.DB().Exec(query, request.JobID, request.FormID, request.Requester, request.Reason)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{
		"status":  "success",
		"message": "takedown request submitted",
	})
}

// ── Crawl runner ──────────────────────────────────────────────────────────────

func runCrawl() {
	if !crawlMu.TryLock() {
		log.Println("crawl already in progress, skipping")
		return
	}
	defer crawlMu.Unlock()

	log.Println("=== starting crawl run (expanded sources) ===")
	summary := scheduler.RunAll(context.Background(), pgStore, browserPool, proxyRotator)

	lastRunTime = time.Now()
	jobsToday += summary.TotalSaved

	log.Printf("=== crawl complete: sources_run=%d succeeded=%d failed=%d total_saved=%d duration=%s ===",
		summary.SourcesRun, summary.Succeeded, summary.Failed, summary.TotalSaved, summary.Duration)
}

// ── Middleware ──────────────────────────────────────────────────────────────

func adminAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		adminKey := os.Getenv("ADMIN_API_KEY")
		if adminKey == "" {
			// If not set, deny all for security
			http.Error(w, "API configuration error", http.StatusInternalServerError)
			return
		}

		key := r.Header.Get("X-Admin-Key")
		if key != adminKey {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		next.ServeHTTP(w, r)
	})
}
