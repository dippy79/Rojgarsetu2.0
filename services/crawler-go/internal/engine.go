package internal

import (
	"crypto/sha256"
	"database/sql"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/rojgarsetu/crawler/internal/sources"
)

// JobSource interface for all scrapers
type JobSource interface {
	FetchJobs() ([]sources.Job, error)
}

// Engine coordinates crawling from all sources
type Engine struct {
	db      *sql.DB
	client  *PoliteCrawler
	sources []JobSource
}

// NewEngine creates a new crawler engine
func NewEngine(db *sql.DB, client *PoliteCrawler) *Engine {
	return &Engine{
		db:     db,
		client: client,
	}
}

// AddSource adds a job source to the engine
func (e *Engine) AddSource(source JobSource) {
	e.sources = append(e.sources, source)
}

// CrawlResult represents the result of a crawl run
type CrawlResult struct {
	Found       int
	Added       int
	Duplicates  int
	Errors      int
	SourcesRun  int
	Duration    time.Duration
	SourceStats []SourceStat
}

// SourceStat represents per-source statistics
type SourceStat struct {
	Name       string
	Found      int
	Added      int
	Duplicates int
	Error      string
}

// Run executes the crawl across all sources
func (e *Engine) Run() CrawlResult {
	start := time.Now()
	result := CrawlResult{
		SourceStats: make([]SourceStat, 0),
	}

	log.Println("=== Starting crawl run ===")

	for _, source := range e.sources {
		sourceName := "unknown"
		sourceStat := SourceStat{}

		// Get source name by type assertion
		switch source.(type) {
		case *sources.UPSCScraper:
			sourceName = "UPSC"
		case *sources.SSCScraper:
			sourceName = "SSC"
		case *sources.RailwayScraper:
			sourceName = "Railway"
		case *sources.NCSScraper:
			sourceName = "NCS"
		case *sources.AdzunaScraper:
			sourceName = "Adzuna"
		case *sources.JoobleScraper:
			sourceName = "Jooble"
		default:
			sourceName = "Unknown"
		}

		sourceStat.Name = sourceName
		result.SourcesRun++

		jobs, err := source.FetchJobs()
		if err != nil {
			log.Printf("[ERROR] %s: %v", sourceName, err)
			sourceStat.Error = err.Error()
			result.Errors++
			result.SourceStats = append(result.SourceStats, sourceStat)
			continue
		}

		sourceStat.Found = len(jobs)
		result.Found += len(jobs)

		// Process each job
		for i := range jobs {
			// Generate hash if not set
			if jobs[i].HashChecksum == "" {
				jobs[i].HashChecksum = GenerateHash(jobs[i].Title, jobs[i].ApplyURL)
			}

			if isDuplicate(e.db, jobs[i].HashChecksum) {
				sourceStat.Duplicates++
				result.Duplicates++
				continue
			}

			// Insert into crawled_jobs
			if err := e.insertCrawledJob(jobs[i]); err != nil {
				log.Printf("[ERROR] Failed to insert job from %s: %v", sourceName, err)
				continue
			}

			// Also insert into jobs_government if it's a govt job
			if jobs[i].SourceAttribution == "Source: UPSC Official Portal (upsc.gov.in)" ||
				jobs[i].SourceAttribution == "Source: SSC Official Portal (ssc.gov.in)" ||
				jobs[i].SourceAttribution == "Source: Railway RRB Official Portal (rrbapply.gov.in)" ||
				jobs[i].SourceAttribution == "Source: NCS Portal (ncs.gov.in)" {
				if err := e.insertGovernmentJob(jobs[i]); err != nil {
					log.Printf("[ERROR] Failed to insert govt job from %s: %v", sourceName, err)
				}
			}

			sourceStat.Added++
			result.Added++
		}

		result.SourceStats = append(result.SourceStats, sourceStat)
		log.Printf("[OK] %s: found=%d added=%d duplicates=%d", sourceName, sourceStat.Found, sourceStat.Added, sourceStat.Duplicates)
	}

	// Log overall statistics
	result.Duration = time.Since(start)
	log.Printf("=== Crawl complete: total_found=%d total_added=%d total_duplicates=%d errors=%d duration=%s ===",
		result.Found, result.Added, result.Duplicates, result.Errors, result.Duration)

	// Insert into crawler_logs
	if err := e.insertCrawlerLog(result); err != nil {
		log.Printf("[ERROR] Failed to insert crawler log: %v", err)
	}

	return result
}

// insertCrawledJob inserts a job into crawled_jobs table
func (e *Engine) insertCrawledJob(job sources.Job) error {
	query := `
		INSERT INTO crawled_jobs (source_id, job_type, title, company_or_dept, location, 
			qualification_req, salary_or_pay_scale, apply_url, source_attribution, hash_checksum)
		VALUES (
			(SELECT id FROM crawler_sources WHERE name = $1),
			'GOVT',
			$2, $3, $4, $5, $6, $7, $8, $9
		)
	`

	sourceName := e.extractSourceName(job.SourceAttribution)
	_, err := e.db.Exec(query, sourceName, job.Title, job.CompanyOrDept, job.Location,
		job.QualificationReq, job.SalaryOrPayScale, job.ApplyURL, job.SourceAttribution, job.HashChecksum)
	return err
}

// insertGovernmentJob inserts a job into jobs_government table
func (e *Engine) insertGovernmentJob(job sources.Job) error {
	query := `
		INSERT INTO jobs_government (title, department, location, qualification, salary, 
			apply_link, last_date, source, is_active)
		VALUES ($1, $2, $3, $4, $5, $6, NOW() + INTERVAL '30 days', $7, true)
		ON CONFLICT DO NOTHING
	`

	_, err := e.db.Exec(query, job.Title, job.CompanyOrDept, job.Location,
		job.QualificationReq, job.SalaryOrPayScale, job.ApplyURL, job.SourceAttribution)
	return err
}

// insertCrawlerLog inserts a crawl log entry
func (e *Engine) insertCrawlerLog(result CrawlResult) error {
	query := `
		INSERT INTO crawler_logs (source_id, jobs_found, jobs_added, duplicates_found, status, error_message)
		VALUES (NULL, $1, $2, $3, 'COMPLETED', NULL)
	`

	_, err := e.db.Exec(query, result.Found, result.Added, result.Duplicates)
	if err != nil {
		return err
	}

	return nil
}

// extractSourceName extracts source name from attribution string
func (e *Engine) extractSourceName(attribution string) string {
	if attribution == "Source: UPSC Official Portal (upsc.gov.in)" {
		return "UPSC"
	}
	if attribution == "Source: SSC Official Portal (ssc.gov.in)" {
		return "SSC"
	}
	if attribution == "Source: Railway RRB Official Portal (rrbapply.gov.in)" {
		return "Railway RRB"
	}
	if attribution == "Source: NCS Portal (ncs.gov.in)" {
		return "NCS Portal"
	}
	if attribution == "Source: Adzuna API" {
		return "Adzuna API"
	}
	if attribution == "Source: Jooble API" {
		return "Jooble API"
	}
	return "Unknown"
}

// isDuplicate checks if a job with the given hash already exists
func isDuplicate(db *sql.DB, hash string) bool {
	var count int
	err := db.QueryRow("SELECT COUNT(*) FROM crawled_jobs WHERE hash_checksum = $1", hash).Scan(&count)
	if err != nil {
		return false
	}
	return count > 0
}

// GenerateHash creates a SHA-256 hash from title and apply URL
func GenerateHash(title, applyURL string) string {
	input := strings.ToLower(strings.TrimSpace(title)) + "|" + strings.TrimSpace(applyURL)
	hash := sha256.Sum256([]byte(input))
	return fmt.Sprintf("%x", hash)
}
