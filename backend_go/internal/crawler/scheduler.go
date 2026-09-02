package crawler

import (
    "database/sql"
    "log"
    "time"
)

type Scheduler struct {
    engine   *Engine
    db       *sql.DB
    interval time.Duration
    stopChan chan struct{}
}

func NewScheduler(db *sql.DB, interval time.Duration) *Scheduler {
    return &Scheduler{
        engine:   NewEngine(db),
        db:       db,
        interval: interval,
        stopChan: make(chan struct{}),
    }
}

// Start launches the periodic crawler ticker in a goroutine
func (s *Scheduler) Start() {
    ticker := time.NewTicker(s.interval)
    log.Printf("[Crawler Scheduler] Started with interval: %v", s.interval)

    go func() {
        for {
            select {
            case <-ticker.C:
                log.Println("[Crawler Scheduler] Running scheduled crawl cycle...")
                s.runCycle()
            case <-s.stopChan:
                ticker.Stop()
                log.Println("[Crawler Scheduler] Stopped successfully")
                return
            }
        }
    }()
}

// Stop shuts down the scheduler
func (s *Scheduler) Stop() {
    close(s.stopChan)
}

func (s *Scheduler) runCycle() {
    rows, err := s.db.Query("SELECT id, name, source_type, base_url FROM crawler_sources WHERE is_active = true")
    if err != nil {
        log.Printf("[Crawler Scheduler] Error fetching active sources: %v", err)
        return
    }
    defer rows.Close()

    for rows.Next() {
        var id int
        var name, sType, baseURL string
        if err := rows.Scan(&id, &name, &sType, &baseURL); err != nil {
            continue
        }

        res, err := s.engine.RunCrawlForSource(id, name, sType, baseURL)
        if err != nil {
            log.Printf("[Crawler Scheduler] Source %s (ID: %d) failed: %v", name, id, err)
        } else {
            log.Printf("[Crawler Scheduler] Source %s completed: %d found, %d added, %d dupes",
                name, res.JobsFound, res.JobsAdded, res.Duplicates)
        }
    }
}
