package db

import (
	"context"
	"database/sql"
	"fmt"
	"time"
)

type PostgresDB struct {
	*Queries
	DB *sql.DB
}

func NewPostgresDB(connStr string) (*PostgresDB, error) {
	dbConn, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, fmt.Errorf("failed to open database (check DATABASE_URL env): %w", err)
	}

	dbConn.SetMaxOpenConns(25)
	dbConn.SetMaxIdleConns(5)
	dbConn.SetConnMaxLifetime(5 * time.Minute)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := dbConn.PingContext(ctx); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	q := New(dbConn)

	return &PostgresDB{
		Queries: q,
		DB:      dbConn,
	}, nil
}

func (p *PostgresDB) Close() error {
	return p.DB.Close()
}

func (p *PostgresDB) GetDB() *sql.DB {
	return p.DB
}

func (p *PostgresDB) WithTx(tx *sql.Tx) *PostgresDB {
	return &PostgresDB{
		Queries: p.Queries.WithTx(tx),
		DB:      p.DB,
	}
}

// ExecWithRetry executes query with deadlock retry
func (p *PostgresDB) ExecWithRetry(ctx context.Context, query string, args ...interface{}) (sql.Result, error) {
	var result sql.Result
	var err error
	for attempt := 0; attempt < 3; attempt++ {
		result, err = p.DB.ExecContext(ctx, query, args...)
		if err == nil {
			return result, nil
		}
		if pgErr, ok := err.(interface{ Code() string }); ok && pgErr.Code() == "40P01" { // deadlock
			time.Sleep(time.Duration(attempt+1) * 100 * time.Millisecond)
			continue
		}
		return nil, err
	}
	return nil, err
}

// QueryWithRetry for select queries with deadlock retry
func (p *PostgresDB) QueryWithRetry(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error) {
	var rows *sql.Rows
	var err error
	for attempt := 0; attempt < 3; attempt++ {
		rows, err = p.DB.QueryContext(ctx, query, args...)
		if err == nil {
			return rows, nil
		}
		if pgErr, ok := err.(interface{ Code() string }); ok && pgErr.Code() == "40P01" {
			if rows != nil {
				rows.Close()
			}
			time.Sleep(time.Duration(attempt+1) * 100 * time.Millisecond)
			continue
		}
		return nil, err
	}
	return nil, err
}

// GovJob methods using sqlc
func (p *PostgresDB) GetGovJobs(filter GovJobFilter, page, limit int) ([]GovJob, int, error) {
	ctx := context.Background()
	offset := (page - 1) * limit
	listParams := GetGovJobsParams{
		Department: filter.Department,
		Location:   filter.Location,
		Source:     filter.Source,
		Limit:      int32(limit),
		Offset:     int32(offset),
	}
	jobs, err := p.Queries.GetGovJobs(ctx, listParams)
	if err != nil {
		return nil, 0, err
	}
	countParams := GetGovJobsCountParams{
		Department: filter.Department,
		Location:   filter.Location,
		Source:     filter.Source,
	}
	count, err := p.Queries.GetGovJobsCount(ctx, countParams)
	if err != nil {
		return jobs, 0, err
	}
	return jobs, int(count), nil
}

func (p *PostgresDB) GetGovJobByID(id string) (*GovJob, error) {
	ctx := context.Background()
	return p.Queries.GetGovJobByID(ctx, id)
}

// PrivJob methods using sqlc
func (p *PostgresDB) GetPrivJobs(filter PrivJobFilter, page, limit int) ([]PrivJob, int, error) {
	ctx := context.Background()
	offset := (page - 1) * limit
	listParams := GetPrivJobsParams{
		Company:  filter.Company,
		Location: filter.Location,
		Source:   filter.Source,
		JobType:  filter.JobType,
		Limit:    int32(limit),
		Offset:   int32(offset),
	}
	jobs, err := p.Queries.GetPrivJobs(ctx, listParams)
	if err != nil {
		return nil, 0, err
	}
	countParams := GetPrivJobsCountParams{
		Company:  filter.Company,
		Location: filter.Location,
		Source:   filter.Source,
		JobType:  filter.JobType,
	}
	count, err := p.Queries.GetPrivJobsCount(ctx, countParams)
	if err != nil {
		return jobs, 0, err
	}
	return jobs, int(count), nil
}

func (p *PostgresDB) GetPrivJobByID(id string) (*PrivJob, error) {
	ctx := context.Background()
	return p.Queries.GetPrivJobByID(ctx, id)
}

// Course methods using sqlc
func (p *PostgresDB) GetCourses(filter CourseFilter, page, limit int) ([]Course, int, error) {
	ctx := context.Background()
	offset := (page - 1) * limit
	listParams := GetCoursesParams{
		Provider: filter.Provider,
		Mode:     filter.Mode,
		Level:    filter.Level,
		Limit:    int32(limit),
		Offset:   int32(offset),
	}
	courses, err := p.Queries.GetCourses(ctx, listParams)
	if err != nil {
		return nil, 0, err
	}
	countParams := GetCoursesCountParams{
		Provider: filter.Provider,
		Mode:     filter.Mode,
		Level:    filter.Level,
	}
	count, err := p.Queries.GetCoursesCount(ctx, countParams)
	if err != nil {
		return courses, 0, err
	}
	return courses, int(count), nil
}

func (p *PostgresDB) GetCourseByID(id string) (*Course, error) {
	ctx := context.Background()
	return p.Queries.GetCourseByID(ctx, id)
}

// Video methods using sqlc
func (p *PostgresDB) GetVideos(filter VideoFilter, page, limit int) ([]Video, int, error) {
	ctx := context.Background()
	offset := (page - 1) * limit
	listParams := GetVideosParams{
		Channel:  filter.Channel,
		Category: filter.Category,
		Limit:    int32(limit),
		Offset:   int32(offset),
	}
	videos, err := p.Queries.GetVideos(ctx, listParams)
	if err != nil {
		return nil, 0, err
	}
	countParams := GetVideosCountParams{
		Channel:  filter.Channel,
		Category: filter.Category,
	}
	count, err := p.Queries.GetVideosCount(ctx, countParams)
	if err != nil {
		return videos, 0, err
	}
	return videos, int(count), nil
}

func (p *PostgresDB) GetVideoByID(id string) (*Video, error) {
	ctx := context.Background()
	return p.Queries.GetVideoByID(ctx, id)
}
