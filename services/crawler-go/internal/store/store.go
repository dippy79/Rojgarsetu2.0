package store

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/lib/pq"
	"github.com/rojgarsetu/crawler/internal/parser"
	"github.com/rojgarsetu/crawler/internal/sources"
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

// SaveGovJob upserts a government job into jobs_government.
// Idempotent on (source, apply_url).
func (s *PostgresStore) SaveGovJob(job *sources.GovJobSource) error {
	if job == nil {
		return fmt.Errorf("job is nil")
	}

	lastDate := parseDateToTime(job.LastDate)
	examDate := parseDateToTime(job.ExamDate)

	query := `
	INSERT INTO jobs_government
	(title, department, location, apply_url, last_date, source, eligibility,
	 vacancy_count, salary, exam_date, notification_pdf_url, is_active)
	VALUES
	($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,true)
	ON CONFLICT (source, apply_url)
	DO UPDATE SET
	    title = EXCLUDED.title,
	    department = EXCLUDED.department,
	    location = EXCLUDED.location,
	    eligibility = EXCLUDED.eligibility,
	    vacancy_count = EXCLUDED.vacancy_count,
	    salary = EXCLUDED.salary,
	    exam_date = EXCLUDED.exam_date,
	    notification_pdf_url = EXCLUDED.notification_pdf_url,
	    is_active = true,
	    updated_at = CURRENT_TIMESTAMP
	`
	_, err := s.db.Exec(
		query,
		job.Title,
		job.Department,
		job.Location,
		job.ApplyURL,
		lastDate,
		job.Source,
		job.Eligibility,
		job.VacancyCount,
		job.Salary,
		examDate,
		job.NotificationURL,
	)
	if err != nil {
		return fmt.Errorf("failed to save gov job '%s': %w", job.Title, err)
	}
	return nil
}

// SavePrivJob upserts a private job into jobs_private.
// Idempotent on (source, url).
func (s *PostgresStore) SavePrivJob(job *sources.PrivJobSource) error {
	if job == nil {
		return fmt.Errorf("job is nil")
	}

	var postedAt any
	if job.PostedAt != nil {
		postedAt = *job.PostedAt
	}

	query := `
	INSERT INTO jobs_private
	(company, title, location, url, salary, experience, job_type, skills,
	 description, source, posted_at, is_active)
	VALUES
	($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,true)
	ON CONFLICT (source, url)
	DO UPDATE SET
	    company = EXCLUDED.company,
	    title = EXCLUDED.title,
	    location = EXCLUDED.location,
	    salary = EXCLUDED.salary,
	    experience = EXCLUDED.experience,
	    job_type = EXCLUDED.job_type,
	    skills = EXCLUDED.skills,
	    description = EXCLUDED.description,
	    posted_at = EXCLUDED.posted_at,
	    is_active = true,
	    updated_at = CURRENT_TIMESTAMP
	`
	_, err := s.db.Exec(
		query,
		job.Company,
		job.Title,
		job.Location,
		job.URL,
		job.Salary,
		job.Experience,
		job.JobType,
		pq.Array(job.Skills),
		job.Description,
		job.Source,
		postedAt,
	)
	if err != nil {
		return fmt.Errorf("failed to save priv job '%s': %w", job.Title, err)
	}
	return nil
}

// SaveCourse upserts a course into courses.
// Idempotent on (source, url).
func (s *PostgresStore) SaveCourse(course *sources.CourseSource) error {
	if course == nil {
		return fmt.Errorf("course is nil")
	}

	startDate := parseDateToTime(course.StartDate)
	endDate := parseDateToTime(course.EndDate)

	rating := ""
	if course.Rating != nil {
		rating = strconv.FormatFloat(*course.Rating, 'f', 2, 64)
	}

	query := `
	INSERT INTO courses
	(provider, title, url, duration, mode, level, skills, description,
	 thumbnail_url, price, is_free, source, start_date, end_date,
	 enrollment_count, rating, is_active)
	VALUES
	($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14,$15,$16,true)
	ON CONFLICT (source, url)
	DO UPDATE SET
	    provider = EXCLUDED.provider,
	    title = EXCLUDED.title,
	    duration = EXCLUDED.duration,
	    mode = EXCLUDED.mode,
	    level = EXCLUDED.level,
	    skills = EXCLUDED.skills,
	    description = EXCLUDED.description,
	    thumbnail_url = EXCLUDED.thumbnail_url,
	    price = EXCLUDED.price,
	    is_free = EXCLUDED.is_free,
	    start_date = EXCLUDED.start_date,
	    end_date = EXCLUDED.end_date,
	    enrollment_count = EXCLUDED.enrollment_count,
	    rating = EXCLUDED.rating,
	    is_active = true,
	    updated_at = CURRENT_TIMESTAMP
	`
	_, err := s.db.Exec(
		query,
		course.Provider,
		course.Title,
		course.URL,
		course.Duration,
		course.Mode,
		course.Level,
		pq.Array(course.Skills),
		course.Description,
		course.ThumbnailURL,
		course.Price,
		course.IsFree,
		course.Source,
		startDate,
		endDate,
		course.EnrollmentCount,
		rating,
	)
	if err != nil {
		return fmt.Errorf("failed to save course '%s': %w", course.Title, err)
	}
	return nil
}

// SaveVideo upserts a YouTube video into youtube_videos.
// Idempotent on video_id.
func (s *PostgresStore) SaveVideo(video *sources.YouTubeVideoSource) error {
	if video == nil {
		return fmt.Errorf("video is nil")
	}

	var publishedAt any
	if video.PublishedAt != nil {
		publishedAt = *video.PublishedAt
	}

	query := `
	INSERT INTO youtube_videos
	(channel, channel_id, title, url, thumbnail, description, video_id,
	 published_at, duration, view_count, like_count, category, is_active)
	VALUES
	($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,true)
	ON CONFLICT (video_id)
	DO UPDATE SET
	    channel = EXCLUDED.channel,
	    channel_id = EXCLUDED.channel_id,
	    title = EXCLUDED.title,
	    url = EXCLUDED.url,
	    thumbnail = EXCLUDED.thumbnail,
	    description = EXCLUDED.description,
	    published_at = EXCLUDED.published_at,
	    duration = EXCLUDED.duration,
	    view_count = EXCLUDED.view_count,
	    like_count = EXCLUDED.like_count,
	    category = EXCLUDED.category,
	    is_active = true,
	    updated_at = CURRENT_TIMESTAMP
	`
	_, err := s.db.Exec(
		query,
		video.Channel,
		video.ChannelID,
		video.Title,
		video.URL,
		video.Thumbnail,
		video.Description,
		video.VideoID,
		publishedAt,
		video.Duration,
		video.ViewCount,
		video.LikeCount,
		video.Category,
	)
	if err != nil {
		return fmt.Errorf("failed to save video '%s': %w", video.Title, err)
	}
	return nil
}

// parseDateToTime converts a scraped date *string into a *time.Time.
// Returns nil for nil/empty/unparseable input so the caller can store NULL
// in a TIMESTAMPTZ column instead of failing with an invalid-input syntax
// error when Postgres tries to cast a non-standard date string.
func parseDateToTime(s *string) *time.Time {
	if s == nil {
		return nil
	}
	trimmed := strings.TrimSpace(*s)
	if trimmed == "" {
		return nil
	}
	layouts := []string{
		time.RFC3339,
		"2006-01-02",
		"2006-01-02T15:04:05",
		"2006-01-02 15:04:05",
		"2006/01/02",
		"02/01/2006",
		"02-01-2006",
		"Jan 02, 2006",
		"02 January 2006",
		"02 Jan 2006",
	}
	for _, layout := range layouts {
		if t, err := time.Parse(layout, trimmed); err == nil {
			return &t
		}
	}
	log.Printf("WARN: unparseable date %q, storing NULL", trimmed)
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
