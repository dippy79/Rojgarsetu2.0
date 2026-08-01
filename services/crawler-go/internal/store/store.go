package store

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"strings"
	"time"

	_ "github.com/lib/pq"
	"github.com/rojgarsetu/crawler/internal/parser"
)

type PostgresStore struct {
	db *sql.DB
}

func NewPostgresStore(connStr string) (*PostgresStore, error) {

	if connStr == "" {
		return nil, fmt.Errorf("DATABASE_URL environment variable required")
	}

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(5 * time.Minute)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := db.PingContext(ctx); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	return &PostgresStore{db: db}, nil
}

func (s *PostgresStore) Close() error {
	return s.db.Close()
}

func (s *PostgresStore) SaveJob(job *parser.Job) error {

	if job == nil {
		return fmt.Errorf("job is nil")
	}

	// Start a database transaction for atomicity
	tx, err := s.db.Begin()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer func() {
		if err != nil {
			_ = tx.Rollback()
			log.Printf("ERROR: transaction rolled back for job '%s': %v", job.Title, err)
		}
	}()

	companyID, err := s.getOrCreateCompanyTx(tx, job.Company)
	if err != nil {
		log.Printf("ERROR: company insert failed for '%s': %v", job.Company, err)
		return fmt.Errorf("failed to get/create company: %w", err)
	}

	query := `
    INSERT INTO jobs 
    (title, company_id, source, location, job_type,
     salary_min, salary_max, eligibility, description,
     application_url, posted_at, is_active)
    VALUES
    ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,true)
    ON CONFLICT (source, application_url)
    DO UPDATE SET
        title = EXCLUDED.title,
        location = EXCLUDED.location,
        job_type = EXCLUDED.job_type,
        salary_min = EXCLUDED.salary_min,
        salary_max = EXCLUDED.salary_max,
        description = EXCLUDED.description,
        updated_at = CURRENT_TIMESTAMP
    RETURNING id
    `

	var jobID string

	err = tx.QueryRow(
		query,
		job.Title,
		companyID,
		job.Source,
		job.Location,
		job.JobType,
		job.SalaryMin,
		job.SalaryMax,
		job.Eligibility,
		job.Description,
		job.ApplicationURL,
		job.PostedAt,
	).Scan(&jobID)

	if err != nil {
		log.Printf("ERROR: job insert failed for '%s': %v", job.Title, err)
		return fmt.Errorf("failed to save job: %w", err)
	}

	// Commit the transaction
	if err := tx.Commit(); err != nil {
		log.Printf("ERROR: transaction commit failed for '%s': %v", job.Title, err)
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

func (s *PostgresStore) getOrCreateCompanyTx(tx *sql.Tx, name string) (string, error) {

	name = strings.TrimSpace(name)
	if name == "" {
		name = "Unknown Company"
	}

	// Normalize name for case-insensitive comparison
	normalizedName := strings.ToLower(name)

	var id string

	// First try to find an existing company with case-insensitive match
	err := tx.QueryRow(`
        SELECT id FROM companies WHERE LOWER(name) = $1
    `, normalizedName).Scan(&id)

	if err == nil {
		// Found existing company, return its ID
		return id, nil
	}

	if err != sql.ErrNoRows {
		return "", fmt.Errorf("failed to lookup company: %w", err)
	}

	// Company not found, insert it
	err = tx.QueryRow(`
        INSERT INTO companies (name)
        VALUES ($1)
        RETURNING id
    `, name).Scan(&id)

	if err != nil {
		return "", fmt.Errorf("failed to create company: %w", err)
	}

	return id, nil
}

func (s *PostgresStore) getOrCreateCompany(name string) (string, error) {

	name = strings.TrimSpace(name)
	if name == "" {
		name = "Unknown Company"
	}

	var id string

	err := s.db.QueryRow(`
        INSERT INTO companies (name)
        VALUES ($1)
        ON CONFLICT ((LOWER(name)))
        DO UPDATE SET name = EXCLUDED.name
        RETURNING id
    `, name).Scan(&id)

	if err != nil {
		return "", fmt.Errorf("failed to get/create company: %w", err)
	}

	return id, nil
}

func (s *PostgresStore) GetJobCount() (int, error) {

	var count int

	err := s.db.QueryRow(`
    SELECT COUNT(*) 
    FROM jobs 
    WHERE is_active = true
    `).Scan(&count)

	if err != nil {
		return 0, err
	}

	return count, nil
}

func (s *PostgresStore) GetJobsBySource(source string, limit int) ([]parser.Job, error) {

	rows, err := s.db.Query(`
    SELECT 
        j.title,
        j.location,
        j.job_type,
        j.source,
        c.name as company
    FROM jobs j
    LEFT JOIN companies c ON j.company_id = c.id
    WHERE j.source = $1 
    AND j.is_active = true
    ORDER BY j.posted_at DESC
    LIMIT $2
    `, source, limit)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var jobs []parser.Job

	for rows.Next() {

		var job parser.Job

		err := rows.Scan(
			&job.Title,
			&job.Location,
			&job.JobType,
			&job.Source,
			&job.Company,
		)

		if err != nil {
			return nil, err
		}

		jobs = append(jobs, job)
	}

	return jobs, nil
}
