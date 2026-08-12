// Code generated manually. DO NOT EDIT SQLC FILES.
package db

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"time"

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

func nullStringPtr(ns sql.NullString) *string {
	if !ns.Valid {
		return nil
	}
	s := ns.String
	return &s
}

func stringFromNullString(ns sql.NullString) string {
	if ns.Valid {
		return ns.String
	}
	return ""
}

func nullBoolPtr(nb sql.NullBool) *bool {
	if !nb.Valid {
		return nil
	}
	b := nb.Bool
	return &b
}

func boolFromNullBool(nb sql.NullBool) bool {
	return nb.Valid && nb.Bool
}

func nullTimePtr(nt sql.NullTime) *time.Time {
	if !nt.Valid {
		return nil
	}
	t := nt.Time
	return &t
}

func timeFromNullTime(nt sql.NullTime) time.Time {
	if nt.Valid {
		return nt.Time
	}
	return time.Time{}
}

func uuidString(u uuid.UUID) string {
	return u.String()
}

func uuidFromString(s string) (uuid.UUID, error) {
	return uuid.Parse(s)
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

// ── Job Categories ───────────────────────────────────────────────────────────

func (p *PostgresDB) GetJobCategories() ([]JobCategory, error) {
	return p.Queries.GetJobCategories(context.Background())
}

func (p *PostgresDB) GetJobCategoryBySlug(slug string) (*JobCategory, error) {
	row, err := p.Queries.GetJobCategoryBySlug(context.Background(), slug)
	if err != nil {
		return nil, err
	}
	return &row, nil
}

func (p *PostgresDB) GetJobCategoryByID(id string) (*JobCategory, error) {
	uid, err := uuid.Parse(id)
	if err != nil {
		return nil, fmt.Errorf("invalid id: %w", err)
	}
	row, err := p.Queries.GetJobCategoryByID(context.Background(), uid)
	if err != nil {
		return nil, err
	}
	return &row, nil
}

// ── Job Trades ───────────────────────────────────────────────────────────────

func (p *PostgresDB) GetJobTrades(f JobTradeFilter, page, limit int) ([]JobTrade, int, error) {
	_, limit, _ = clampPagination(page, limit)

	var categoryID uuid.UUID
	if f.CategoryID != "" {
		parsed, err := uuid.Parse(f.CategoryID)
		if err != nil {
			return nil, 0, fmt.Errorf("invalid category_id: %w", err)
		}
		categoryID = parsed
	} else {
		categoryID = uuid.Nil
	}

	total, err := p.Queries.GetJobTradesCount(context.Background(), GetJobTradesCountParams{
		Column1: categoryID,
		Column2: f.DemandLevel,
	})
	if err != nil {
		return nil, 0, err
	}

	rows, err := p.Queries.GetJobTrades(context.Background(), GetJobTradesParams{
		Column1: categoryID,
		Column2: f.DemandLevel,
	})
	if err != nil {
		return nil, 0, err
	}

	return rows, int(total), nil
}

func (p *PostgresDB) GetJobTradeByID(id string) (*JobTrade, error) {
	uid, err := uuid.Parse(id)
	if err != nil {
		return nil, fmt.Errorf("invalid id: %w", err)
	}
	row, err := p.Queries.GetJobTradeByID(context.Background(), uid)
	if err != nil {
		return nil, err
	}
	return &row, nil
}

func (p *PostgresDB) GetJobTradeBySlug(slug string, categoryID string) (*JobTrade, error) {
	catUID, err := uuid.Parse(categoryID)
	if err != nil {
		return nil, fmt.Errorf("invalid category_id: %w", err)
	}
	row, err := p.Queries.GetJobTradeBySlug(context.Background(), GetJobTradeBySlugParams{
		Slug:       slug,
		CategoryID: catUID,
	})
	if err != nil {
		return nil, err
	}
	return &row, nil
}

func (p *PostgresDB) GetJobTradesByCategory(categoryID string) ([]JobTrade, error) {
	catUID, err := uuid.Parse(categoryID)
	if err != nil {
		return nil, fmt.Errorf("invalid category_id: %w", err)
	}
	return p.Queries.GetJobTradesByCategory(context.Background(), catUID)
}

// ── User Enrollments ─────────────────────────────────────────────────────────

func (p *PostgresDB) GetUserEnrollments(f UserEnrollmentFilter, userID string, page, limit int) ([]UserEnrollment, int, error) {
	_, limit, offset := clampPagination(page, limit)

	userUID, err := uuid.Parse(userID)
	if err != nil {
		return nil, 0, fmt.Errorf("invalid user_id: %w", err)
	}

	total, err := p.Queries.GetUserEnrollmentsCount(context.Background(), GetUserEnrollmentsCountParams{
		UserID:  userUID,
		Column2: f.Status,
	})
	if err != nil {
		return nil, 0, err
	}

	rows, err := p.Queries.GetUserEnrollments(context.Background(), GetUserEnrollmentsParams{
		UserID:  userUID,
		Column2: f.Status,
		Limit:   int32(limit),
		Offset:  int32(offset),
	})
	if err != nil {
		return nil, 0, err
	}

	return rows, int(total), nil
}

func (p *PostgresDB) GetUserEnrollmentByID(id string) (*UserEnrollment, error) {
	uid, err := uuid.Parse(id)
	if err != nil {
		return nil, fmt.Errorf("invalid id: %w", err)
	}
	row, err := p.Queries.GetUserEnrollmentByID(context.Background(), uid)
	if err != nil {
		return nil, err
	}
	return &row, nil
}

func (p *PostgresDB) GetUserEnrollmentByUserAndTrade(userID string, tradeID string) (*UserEnrollment, error) {
	userUID, err := uuid.Parse(userID)
	if err != nil {
		return nil, fmt.Errorf("invalid user_id: %w", err)
	}
	tradeUID, err := uuid.Parse(tradeID)
	if err != nil {
		return nil, fmt.Errorf("invalid trade_id: %w", err)
	}
	row, err := p.Queries.GetUserEnrollmentByUserAndTrade(context.Background(), GetUserEnrollmentByUserAndTradeParams{
		UserID:  userUID,
		TradeID: tradeUID,
	})
	if err != nil {
		return nil, err
	}
	return &row, nil
}

func (p *PostgresDB) GetExpiringEnrollments() ([]UserEnrollment, error) {
	return p.Queries.GetExpiringEnrollments(context.Background())
}

func (p *PostgresDB) GetExpiringEnrollmentsWithTrade() ([]GetExpiringEnrollmentsWithTradeRow, error) {
	return p.Queries.GetExpiringEnrollmentsWithTrade(context.Background())
}

func (p *PostgresDB) CreateUserEnrollment(userID string, tradeID string, expiresAt time.Time, metadata map[string]interface{}) (*UserEnrollment, error) {
	userUID, err := uuid.Parse(userID)
	if err != nil {
		return nil, fmt.Errorf("invalid user_id: %w", err)
	}
	tradeUID, err := uuid.Parse(tradeID)
	if err != nil {
		return nil, fmt.Errorf("invalid trade_id: %w", err)
	}

	metadataJSON, err := json.Marshal(metadata)
	if err != nil {
		return nil, fmt.Errorf("invalid metadata: %w", err)
	}

	row, err := p.Queries.CreateUserEnrollment(context.Background(), CreateUserEnrollmentParams{
		UserID:    userUID,
		TradeID:   tradeUID,
		ExpiresAt: expiresAt,
		Metadata:  metadataJSON,
	})
	if err != nil {
		return nil, err
	}
	return &row, nil
}

func (p *PostgresDB) UpdateUserEnrollment(id string, status string, expiresAt *time.Time, completedAt *time.Time, progressPct int32, metadata map[string]interface{}) (*UserEnrollment, error) {
	uid, err := uuid.Parse(id)
	if err != nil {
		return nil, fmt.Errorf("invalid id: %w", err)
	}

	metadataJSON, err := json.Marshal(metadata)
	if err != nil {
		return nil, fmt.Errorf("invalid metadata: %w", err)
	}

	var expiresAtVal sql.NullTime
	if expiresAt != nil {
		expiresAtVal = sql.NullTime{Time: *expiresAt, Valid: true}
	}

	var completedAtVal sql.NullTime
	if completedAt != nil {
		completedAtVal = sql.NullTime{Time: *completedAt, Valid: true}
	}

	row, err := p.Queries.UpdateUserEnrollment(context.Background(), UpdateUserEnrollmentParams{
		ID:          uid,
		Status:      status,
		ExpiresAt:   expiresAtVal.Time,
		CompletedAt: completedAtVal,
		ProgressPct: progressPct,
		Metadata:    metadataJSON,
	})
	if err != nil {
		return nil, err
	}
	return &row, nil
}

func (p *PostgresDB) UpdateEnrollmentProgress(id string, progressPct int32) (*UserEnrollment, error) {
	uid, err := uuid.Parse(id)
	if err != nil {
		return nil, fmt.Errorf("invalid id: %w", err)
	}
	row, err := p.Queries.UpdateEnrollmentProgress(context.Background(), UpdateEnrollmentProgressParams{
		ID:          uid,
		ProgressPct: progressPct,
	})
	if err != nil {
		return nil, err
	}
	return &row, nil
}

func (p *PostgresDB) CompleteEnrollment(id string) (*UserEnrollment, error) {
	uid, err := uuid.Parse(id)
	if err != nil {
		return nil, fmt.Errorf("invalid id: %w", err)
	}
	row, err := p.Queries.CompleteEnrollment(context.Background(), uid)
	if err != nil {
		return nil, err
	}
	return &row, nil
}

func (p *PostgresDB) CancelEnrollment(id string) (*UserEnrollment, error) {
	uid, err := uuid.Parse(id)
	if err != nil {
		return nil, fmt.Errorf("invalid id: %w", err)
	}
	row, err := p.Queries.CancelEnrollment(context.Background(), uid)
	if err != nil {
		return nil, err
	}
	return &row, nil
}

// ── User Notification Logs ───────────────────────────────────────────────────

func (p *PostgresDB) GetUserNotificationLogs(f NotificationLogFilter, userID string, page, limit int) ([]UserNotificationLog, int, error) {
	_, limit, offset := clampPagination(page, limit)

	userUID, err := uuid.Parse(userID)
	if err != nil {
		return nil, 0, fmt.Errorf("invalid user_id: %w", err)
	}

	total, err := p.Queries.GetUserNotificationLogsCount(context.Background(), GetUserNotificationLogsCountParams{
		UserID:  userUID,
		Column2: f.NotificationType,
	})
	if err != nil {
		return nil, 0, err
	}

	rows, err := p.Queries.GetUserNotificationLogs(context.Background(), GetUserNotificationLogsParams{
		UserID:  userUID,
		Column2: f.NotificationType,
		Limit:   int32(limit),
		Offset:  int32(offset),
	})
	if err != nil {
		return nil, 0, err
	}

	return rows, int(total), nil
}

func (p *PostgresDB) GetNotificationLogByID(id string) (*UserNotificationLog, error) {
	uid, err := uuid.Parse(id)
	if err != nil {
		return nil, fmt.Errorf("invalid id: %w", err)
	}
	row, err := p.Queries.GetNotificationLogByID(context.Background(), uid)
	if err != nil {
		return nil, err
	}
	return &row, nil
}

func (p *PostgresDB) GetDailyNotificationCount(userID string) (int64, error) {
	userUID, err := uuid.Parse(userID)
	if err != nil {
		return 0, fmt.Errorf("invalid user_id: %w", err)
	}
	return p.Queries.GetDailyNotificationCount(context.Background(), userUID)
}

func (p *PostgresDB) GetDailyNotificationCountByType(userID string, notificationType string) (int64, error) {
	userUID, err := uuid.Parse(userID)
	if err != nil {
		return 0, fmt.Errorf("invalid user_id: %w", err)
	}
	return p.Queries.GetDailyNotificationCountByType(context.Background(), GetDailyNotificationCountByTypeParams{
		UserID:           userUID,
		NotificationType: notificationType,
	})
}

func (p *PostgresDB) CreateNotificationLog(userID string, enrollmentID *string, notificationType string, channel string, title string, message string, payload map[string]interface{}) (*UserNotificationLog, error) {
	userUID, err := uuid.Parse(userID)
	if err != nil {
		return nil, fmt.Errorf("invalid user_id: %w", err)
	}

	var enrollmentUID uuid.NullUUID
	if enrollmentID != nil {
		parsed, err := uuid.Parse(*enrollmentID)
		if err != nil {
			return nil, fmt.Errorf("invalid enrollment_id: %w", err)
		}
		enrollmentUID = uuid.NullUUID{UUID: parsed, Valid: true}
	}

	payloadJSON, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("invalid payload: %w", err)
	}

	row, err := p.Queries.CreateNotificationLog(context.Background(), CreateNotificationLogParams{
		UserID:           userUID,
		EnrollmentID:     enrollmentUID,
		NotificationType: notificationType,
		Channel:          channel,
		Title:            title,
		Message:          message,
		Payload:          payloadJSON,
	})
	if err != nil {
		return nil, err
	}
	return &row, nil
}

func (p *PostgresDB) MarkNotificationRead(id string) (*UserNotificationLog, error) {
	uid, err := uuid.Parse(id)
	if err != nil {
		return nil, fmt.Errorf("invalid id: %w", err)
	}
	row, err := p.Queries.MarkNotificationRead(context.Background(), uid)
	if err != nil {
		return nil, err
	}
	return &row, nil
}

func (p *PostgresDB) MarkNotificationClicked(id string) (*UserNotificationLog, error) {
	uid, err := uuid.Parse(id)
	if err != nil {
		return nil, fmt.Errorf("invalid id: %w", err)
	}
	row, err := p.Queries.MarkNotificationClicked(context.Background(), uid)
	if err != nil {
		return nil, err
	}
	return &row, nil
}
