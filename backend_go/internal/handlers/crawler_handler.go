package handlers

import (
    "database/sql"
    "net/http"
    "strconv"

    "github.com/gin-gonic/gin"
    "github.com/rojgarsetu/backend/internal/crawler"
)

type CrawlerHandler struct {
    db     *sql.DB
    engine *crawler.Engine
}

func NewCrawlerHandler(db *sql.DB) *CrawlerHandler {
    return &CrawlerHandler{
        db:     db,
        engine: crawler.NewEngine(db),
    }
}

// POST /api/v1/crawler/crawl
func (h *CrawlerHandler) TriggerCrawl(c *gin.Context) {
    sourceIDStr := c.Query("source_id")

    var query string
    var args []interface{}

    if sourceIDStr != "" {
        sourceID, err := strconv.Atoi(sourceIDStr)
        if err != nil {
            c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid source_id parameter"})
            return
        }
        query = "SELECT id, name, source_type, base_url FROM crawler_sources WHERE id = $1 AND is_active = true"
        args = append(args, sourceID)
    } else {
        query = "SELECT id, name, source_type, base_url FROM crawler_sources WHERE is_active = true"
    }

    rows, err := h.db.Query(query, args...)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error fetching sources: " + err.Error()})
        return
    }
    defer rows.Close()

    var results []*crawler.CrawlResult
    for rows.Next() {
        var id int
        var name, sType, baseURL string
        if err := rows.Scan(&id, &name, &sType, &baseURL); err != nil {
            continue
        }

        res, _ := h.engine.RunCrawlForSource(id, name, sType, baseURL)
        results = append(results, res)
    }

    c.JSON(http.StatusOK, gin.H{
        "success": true,
        "message": "Crawl execution finished",
        "runs":    results,
    })
}

// GET /api/v1/crawler/stats
func (h *CrawlerHandler) GetStats(c *gin.Context) {
    var totalCrawled, totalAdded, totalDuplicates, totalErrors int

    err := h.db.QueryRow(`
        SELECT 
            COALESCE(SUM(jobs_found), 0),
            COALESCE(SUM(jobs_added), 0),
            COALESCE(SUM(duplicates_found), 0),
            COALESCE(SUM(errors_count), 0)
        FROM crawler_logs
    `).Scan(&totalCrawled, &totalAdded, &totalDuplicates, &totalErrors)

    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Error querying stats: " + err.Error()})
        return
    }

    var totalUniqueJobs int
    _ = h.db.QueryRow("SELECT COUNT(*) FROM crawled_jobs").Scan(&totalUniqueJobs)

    c.JSON(http.StatusOK, gin.H{
        "success": true,
        "stats": gin.H{
            "total_crawled":      totalCrawled,
            "total_jobs_added":   totalAdded,
            "total_duplicates":  totalDuplicates,
            "total_errors":      totalErrors,
            "total_unique_jobs": totalUniqueJobs,
        },
    })
}

// GET /api/v1/crawler/health
func (h *CrawlerHandler) GetHealth(c *gin.Context) {
    var activeSources int
    _ = h.db.QueryRow("SELECT COUNT(*) FROM crawler_sources WHERE is_active = true").Scan(&activeSources)

    var recentErrors int
    _ = h.db.QueryRow(`
        SELECT COUNT(*) FROM crawler_logs 
        WHERE status = 'FAILED' AND created_at >= NOW() - INTERVAL '24 hours'
    `).Scan(&recentErrors)

    var lastCrawledAt sql.NullString
    _ = h.db.QueryRow("SELECT MAX(created_at)::text FROM crawler_logs").Scan(&lastCrawledAt)

    status := "HEALTHY"
    if recentErrors > 5 {
        status = "DEGRADED"
    }

    c.JSON(http.StatusOK, gin.H{
        "status":            status,
        "active_sources":    activeSources,
        "recent_errors_24h": recentErrors,
        "last_crawl_at":     lastCrawledAt.String,
    })
}
