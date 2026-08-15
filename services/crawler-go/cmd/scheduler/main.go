package main

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strconv"
	"sync"
	"time"

	_ "github.com/lib/pq"
	"github.com/rojgarsetu/crawler/internal"
	"github.com/rojgarsetu/crawler/internal/sources"
)

var (
	db          *sql.DB
	engine      *internal.Engine
	crawlMu     sync.Mutex
	lastRunTime time.Time
	jobsToday   int
)

func main() {
	// ── Database connection ──────────────────────────────────────────────────
	databaseURL := os.Getenv("DATABASE_URL")
	if databaseURL == "" {
		log.Fatal("DATABASE_URL environment variable is required")
	}

	var err error
	db, err = sql.Open("postgres", databaseURL)
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		log.Fatalf("failed to ping database: %v", err)
	}
	log.Println("connected to database")

	// ── Polite client setup ───────────────────────────────────────────────────
	userAgent := os.Getenv("CRAWLER_USER_AGENT")
	if userAgent == "" {
		userAgent = "RojgarSetuBot/2.0 (+https://rojgarsetu.in/bot-policy; support@rojgarsetu.in)"
	}

	delayMs := 2000
	if v := os.Getenv("CRAWLER_REQUEST_DELAY_MS"); v != "" {
		if n, err := strconv.Atoi(v); err == nil && n > 0 {
			delayMs = n
		}
	}

	robotsStrict := true
	if v := os.Getenv("CRAWLER_ROBOTS_TXT_STRICT"); v != "" {
		robotsStrict = v == "true"
	}

	client := internal.New(userAgent, delayMs, robotsStrict)
	log.Printf("polite client initialized: delay=%dms, robots_strict=%v", delayMs, robotsStrict)

	// ── Engine setup ──────────────────────────────────────────────────────────
	engine = internal.NewEngine(db, client)

	// Add all sources
	engine.AddSource(sources.NewUPSCScraper(client))
	engine.AddSource(sources.NewSSCScraper(client))
	engine.AddSource(sources.NewRailwayScraper(client))
	engine.AddSource(sources.NewNCSScraper(client))
	engine.AddSource(sources.NewAdzunaScraper(client))
	engine.AddSource(sources.NewJoobleScraper(client))

	log.Println("crawler engine initialized with 6 sources")

	// ── HTTP routes ──────────────────────────────────────────────────────────
	http.HandleFunc("/health", healthHandler)
	http.HandleFunc("/stats", statsHandler)
	http.HandleFunc("/crawl", crawlHandler)
	http.HandleFunc("/sources", sourcesHandler)
	http.HandleFunc("/forms", formsHandler)
	http.HandleFunc("/legal/takedown", takedownHandler)

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
	rows, err := db.Query("SELECT id, name, category, source_type, base_url, is_active FROM crawler_sources ORDER BY id")
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
	rows, err := db.Query("SELECT id, title, conducting_body, form_type, official_website, is_taken_down FROM gov_forms_info ORDER BY created_at DESC LIMIT 50")
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

	_, err := db.Exec(query, request.JobID, request.FormID, request.Requester, request.Reason)
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

	log.Println("=== starting crawl run ===")
	result := engine.Run()

	lastRunTime = time.Now()
	jobsToday += result.Added

	log.Printf("=== crawl complete: sources_run=%d found=%d added=%d duplicates=%d errors=%d duration=%s ===",
		result.SourcesRun, result.Found, result.Added, result.Duplicates, result.Errors, result.Duration)
}
