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
	})
	if err != nil {
		return nil, 0, err
	}

	rows, err := p.Queries.GetGovJobs(context.Background(), GetGovJobsParams{
		Column1: f.Department,
		Column2: f.Location,
		Column3: f.Source,
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
	})
	if err != nil {
		return nil, 0, err
	}

	rows, err := p.Queries.GetPrivJobs(context.Background(), GetPrivJobsParams{
		Column1: f.Company,
		Column2: f.Location,
		Column3: f.Source,
		Column4: f.JobType,
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
	})
	if err != nil {
		return nil, 0, err
	}

	rows, err := p.Queries.GetCourses(context.Background(), GetCoursesParams{
		Column1: f.Provider,
		Column2: f.Mode,
		Column3: f.Level,
		Limit:   int32(limit),
		Offset:  int32(offset),
	})
	if err != nil {
		return nil, 0, err
	}

	return rows, int(total), nil
}

type CourseProvider struct {
	Provider string `json:"provider"`
	Count    int    `json:"count"`
}

func (p *PostgresDB) GetCourseProviders() ([]CourseProvider, error) {
	rows, err := p.DB.QueryContext(context.Background(), `
		SELECT provider, COUNT(*)
		FROM courses
		WHERE is_active = true
		GROUP BY provider
		ORDER BY COUNT(*) DESC, provider ASC
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	providers := make([]CourseProvider, 0)
	for rows.Next() {
		var provider CourseProvider
		if err := rows.Scan(&provider.Provider, &provider.Count); err != nil {
			return nil, err
		}
		providers = append(providers, provider)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return providers, nil
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

type VideoChannel struct {
	Channel string `json:"channel"`
	Count   int    `json:"count"`
}

type VideoCategory struct {
	Category string `json:"category"`
	Count    int    `json:"count"`
}

func (p *PostgresDB) GetVideoChannels() ([]VideoChannel, error) {
	rows, err := p.DB.QueryContext(context.Background(), `
		SELECT channel, COUNT(*)
		FROM youtube_videos
		WHERE is_active = true
		GROUP BY channel
		ORDER BY COUNT(*) DESC, channel ASC
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	channels := make([]VideoChannel, 0)
	for rows.Next() {
		var ch VideoChannel
		if err := rows.Scan(&ch.Channel, &ch.Count); err != nil {
			return nil, err
		}
		channels = append(channels, ch)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return channels, nil
}

func (p *PostgresDB) GetVideoCategories() ([]VideoCategory, error) {
	rows, err := p.DB.QueryContext(context.Background(), `
		SELECT category, COUNT(*)
		FROM youtube_videos
		WHERE is_active = true AND category IS NOT NULL AND category <> ''
		GROUP BY category
		ORDER BY COUNT(*) DESC, category ASC
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	categories := make([]VideoCategory, 0)
	for rows.Next() {
		var cat VideoCategory
		if err := rows.Scan(&cat.Category, &cat.Count); err != nil {
			return nil, err
		}
		categories = append(categories, cat)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return categories, nil
}

// GetVideos returns a page of videos. When exclude is non-empty, videos whose
// category equals exclude are omitted from BOTH the result set and the total
// count, so pagination reflects the post-exclusion total (e.g. Tech tab
// excluding "Government").
func (p *PostgresDB) GetVideos(f VideoFilter, exclude string, page, limit int) ([]GetVideosRow, int, error) {
	_, limit, offset := clampPagination(page, limit)

	// With an exclusion, compute the total over the post-exclusion set.
	if exclude != "" {
		var total int64
		err := p.DB.QueryRowContext(context.Background(), `
			SELECT COUNT(*) FROM youtube_videos
			WHERE is_active = true
			  AND ($1::text = '' OR channel ILIKE '%' || $1 || '%')
			  AND ($2::text = '' OR category = $2)
			  AND ($3::text = '' OR category <> $3)
		`, f.Channel, f.Category, exclude).Scan(&total)
		if err != nil {
			return nil, 0, err
		}

		rows, err := p.DB.QueryContext(context.Background(), `
			SELECT id, channel, channel_id, title, url, thumbnail,
			       description, video_id, published_at, duration,
			       view_count, like_count, category, created_at
			FROM youtube_videos
			WHERE is_active = true
			  AND ($1::text = '' OR channel ILIKE '%' || $1 || '%')
			  AND ($2::text = '' OR category = $2)
			  AND ($3::text = '' OR category <> $3)
			ORDER BY published_at DESC
			LIMIT $4 OFFSET $5
		`, f.Channel, f.Category, exclude, limit, offset)
		if err != nil {
			return nil, 0, err
		}
		defer rows.Close()

		items := make([]GetVideosRow, 0, limit)
		for rows.Next() {
			var i GetVideosRow
			if err := rows.Scan(
				&i.ID, &i.Channel, &i.ChannelID, &i.Title, &i.Url, &i.Thumbnail,
				&i.Description, &i.VideoID, &i.PublishedAt, &i.Duration,
				&i.ViewCount, &i.LikeCount, &i.Category, &i.CreatedAt,
			); err != nil {
				return nil, 0, err
			}
			items = append(items, i)
		}
		if err := rows.Err(); err != nil {
			return nil, 0, err
		}

		return items, int(total), nil
	}

	total, err := p.Queries.GetVideosCount(context.Background(), GetVideosCountParams{
		Column1: f.Channel,
		Column2: f.Category,
	})
	if err != nil {
		return nil, 0, err
	}

	rows, err := p.Queries.GetVideos(context.Background(), GetVideosParams{
		Column1: f.Channel,
		Column2: f.Category,
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
