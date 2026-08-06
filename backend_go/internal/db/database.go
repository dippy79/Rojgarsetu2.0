// Code generated manually. DO NOT EDIT SQLC FILES.
package db

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/google/uuid"
)

// PostgresDB wraps *Queries and adds helper methods with filters + pagination.
type PostgresDB struct {
	Queries *Queries
	DB      *sql.DB
}

// NewPostgresDB creates a PostgresDB from a *sql.DB connection.
func NewPostgresDB(db *sql.DB) *PostgresDB {
	return &PostgresDB{Queries: New(db), DB: db}
}

// GetDB returns the underlying *sql.DB connection for raw queries
func (p *PostgresDB) GetDB() *sql.DB {
	return p.DB
}

// ── pagination helper ────────────────────────────────────────────────────────

func clampPagination(page, limit int) (int, int, int) {
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 20
	}
	offset := (page - 1) * limit
	return page, limit, offset
}

// ── Government Jobs ──────────────────────────────────────────────────────────

func (p *PostgresDB) GetGovJobs(f GovJobFilter, page, limit int) ([]GetGovJobsRow, int, error) {
	_, limit, offset := clampPagination(page, limit)

	total, err := p.Queries.GetGovJobsCount(context.Background(), GetGovJobsCountParams{
		Column1: f.Department,
		Column2: f.Location,
		Column3: f.Source,
		Column4: f.Language,
	})
	if err != nil {
		return nil, 0, err
	}

	rows, err := p.Queries.GetGovJobs(context.Background(), GetGovJobsParams{
		Column1: f.Department,
		Column2: f.Location,
		Column3: f.Source,
		Column4: f.Language,
		Limit:   int32(limit),
		Offset:  int32(offset),
	})
	if err != nil {
		return nil, 0, err
	}

	return rows, int(total), nil
}

func (p *PostgresDB) GetGovJobByID(id string) (*GetGovJobByIDRow, error) {
	uid, err := uuid.Parse(id)
	if err != nil {
		return nil, fmt.Errorf("invalid id: %w", err)
	}
	row, err := p.Queries.GetGovJobByID(context.Background(), uid)
	if err != nil {
		return nil, err
	}
	return &row, nil
}

// ── Private Jobs ─────────────────────────────────────────────────────────────

func (p *PostgresDB) GetPrivJobs(f PrivJobFilter, page, limit int) ([]GetPrivJobsRow, int, error) {
	_, limit, offset := clampPagination(page, limit)

	total, err := p.Queries.GetPrivJobsCount(context.Background(), GetPrivJobsCountParams{
		Column1: f.Company,
		Column2: f.Location,
		Column3: f.Source,
		Column4: f.JobType,
		Column5: f.Language,
	})
	if err != nil {
		return nil, 0, err
	}

	rows, err := p.Queries.GetPrivJobs(context.Background(), GetPrivJobsParams{
		Column1: f.Company,
		Column2: f.Location,
		Column3: f.Source,
		Column4: f.JobType,
		Column5: f.Language,
		Limit:   int32(limit),
		Offset:  int32(offset),
	})
	if err != nil {
		return nil, 0, err
	}

	return rows, int(total), nil
}

func (p *PostgresDB) GetPrivJobByID(id string) (*GetPrivJobByIDRow, error) {
	uid, err := uuid.Parse(id)
	if err != nil {
		return nil, fmt.Errorf("invalid id: %w", err)
	}
	row, err := p.Queries.GetPrivJobByID(context.Background(), uid)
	if err != nil {
		return nil, err
	}
	return &row, nil
}

// ── Courses ──────────────────────────────────────────────────────────────────

func (p *PostgresDB) GetCourses(f CourseFilter, page, limit int) ([]GetCoursesRow, int, error) {
	_, limit, offset := clampPagination(page, limit)

	total, err := p.Queries.GetCoursesCount(context.Background(), GetCoursesCountParams{
		Column1: f.Provider,
		Column2: f.Mode,
		Column3: f.Level,
		Column4: f.Language,
	})
	if err != nil {
		return nil, 0, err
	}

	rows, err := p.Queries.GetCourses(context.Background(), GetCoursesParams{
		Column1: f.Provider,
		Column2: f.Mode,
		Column3: f.Level,
		Column4: f.Language,
		Limit:   int32(limit),
		Offset:  int32(offset),
	})
	if err != nil {
		return nil, 0, err
	}

	return rows, int(total), nil
}

func (p *PostgresDB) GetCourseByID(id string) (*GetCourseByIDRow, error) {
	uid, err := uuid.Parse(id)
	if err != nil {
		return nil, fmt.Errorf("invalid id: %w", err)
	}
	row, err := p.Queries.GetCourseByID(context.Background(), uid)
	if err != nil {
		return nil, err
	}
	return &row, nil
}

// ── Videos ───────────────────────────────────────────────────────────────────

func (p *PostgresDB) GetVideos(f VideoFilter, page, limit int) ([]GetVideosRow, int, error) {
	_, limit, offset := clampPagination(page, limit)

	total, err := p.Queries.GetVideosCount(context.Background(), GetVideosCountParams{
		Column1: f.Channel,
		Column2: f.Category,
		Column3: f.Language,
	})
	if err != nil {
		return nil, 0, err
	}

	rows, err := p.Queries.GetVideos(context.Background(), GetVideosParams{
		Column1: f.Channel,
		Column2: f.Category,
		Column3: f.Language,
		Limit:   int32(limit),
		Offset:  int32(offset),
	})
	if err != nil {
		return nil, 0, err
	}

	return rows, int(total), nil
}

func (p *PostgresDB) GetVideoByID(id string) (*GetVideoByIDRow, error) {
	uid, err := uuid.Parse(id)
	if err != nil {
		return nil, fmt.Errorf("invalid id: %w", err)
	}
	row, err := p.Queries.GetVideoByID(context.Background(), uid)
	if err != nil {
		return nil, err
	}
	return &row, nil
}
