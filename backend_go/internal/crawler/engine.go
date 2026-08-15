package crawler

import (
    "database/sql"
    "encoding/json"
    "fmt"
    "log"
    "strings"
    "time"
)

type Engine struct {
    db      *sql.DB
    fetcher *Fetcher
    mapper  *MappingEngine
}

func NewEngine(db *sql.DB) *Engine {
    return &Engine{
        db:      db,
        fetcher: NewFetcher(),
        mapper:  NewMappingEngine(db),
    }
}

func (e *Engine) RunCrawlForSource(sourceID int, sourceName, sourceType, baseURL string) (*CrawlResult, error) {
    startTime := time.Now()
    result := &CrawlResult{
        SourceName: sourceName,
        Status:     "RUNNING",
    }

    var scrapedJobs []ScrapedJob
    var fetchErr error

    switch sourceType {
    case "RSS":
        items, err := e.fetcher.FetchRSS(baseURL)
        if err != nil {
            fetchErr = err
            break
        }
        for _, item := range items {
            scrapedJobs = append(scrapedJobs, ScrapedJob{
                ExternalJobID:           item.GUID,
                Title:                   item.Title,
                Organization:            sourceName,
                JobType:                 "GOVT",
                ApplyURL:                item.Link,
                OfficialNotificationURL: item.Link,
            })
        }
    default:
        scrapedJobs = append(scrapedJobs, ScrapedJob{
            ExternalJobID: fmt.Sprintf("%s-%d", strings.ToLower(sourceName), time.Now().Unix()),
            Title:         fmt.Sprintf("Latest Recruitment Notification - %s", sourceName),
            Organization:  sourceName,
            JobType:       "GOVT",
            ApplyURL:      baseURL,
        })
    }

    if fetchErr != nil {
        result.Status = "FAILED"
        result.ErrorMessage = fetchErr.Error()
        result.ExecutionTime = time.Since(startTime).Milliseconds()
        e.logCrawlRun(result)
        return result, fetchErr
    }

    result.JobsFound = len(scrapedJobs)

    for _, job := range scrapedJobs {
        jobHash := GenerateJobHash(sourceName, job.Title, job.Organization, job.ApplyURL)

        // De-duplication check against DB
        var exists bool
        err := e.db.QueryRow("SELECT EXISTS(SELECT 1 FROM crawled_jobs WHERE job_hash = $1)", jobHash).Scan(&exists)
        if err != nil {
            result.ErrorsCount++
            continue
        }

        if exists {
            result.Duplicates++
            continue
        }

        // Auto-map Trade and Category IDs
        categoryID, tradeID := e.mapper.AutoMap(job.Title, job.Organization)

        // Save unique job with mapped tags
        payloadJSON, _ := json.Marshal(job)
        _, err = e.db.Exec(`
            INSERT INTO crawled_jobs (
                source_id, external_job_id, job_hash, title, organization, 
                job_type, category_id, trade_id, qualification_required, total_vacancies, salary_range, 
                job_location, official_notification_url, apply_url, raw_payload
            ) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15)
        `, sourceID, job.ExternalJobID, jobHash, job.Title, job.Organization,
            job.JobType, categoryID, tradeID, job.QualificationRequired, job.TotalVacancies, job.SalaryRange,
            job.JobLocation, job.OfficialNotificationURL, job.ApplyURL, payloadJSON)

        if err != nil {
            log.Printf("[Crawler Engine Error] Insert failed for %s: %v", jobHash, err)
            result.ErrorsCount++
        } else {
            result.JobsAdded++
        }
    }

    result.Status = "SUCCESS"
    result.ExecutionTime = time.Since(startTime).Milliseconds()

    // Update crawler_sources table last_crawled_at
    _, _ = e.db.Exec("UPDATE crawler_sources SET last_crawled_at = NOW() WHERE id = $1", sourceID)
    e.logCrawlRun(result)

    return result, nil
}

func (e *Engine) logCrawlRun(res *CrawlResult) {
    _, _ = e.db.Exec(`
        INSERT INTO crawler_logs (source_name, status, jobs_found, jobs_added, duplicates_found, errors_count, error_message, execution_time_ms)
        VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
    `, res.SourceName, res.Status, res.JobsFound, res.JobsAdded, res.Duplicates, res.ErrorsCount, res.ErrorMessage, res.ExecutionTime)
}
