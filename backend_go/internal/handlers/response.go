package handlers

import (
	"database/sql"
	"time"

	"github.com/rojgarsetu/backend/internal/db"
)

// ---- Null helpers ----
func nullStr(ns sql.NullString) string {
	if ns.Valid {
		return ns.String
	}
	return ""
}

func nullInt32Ptr(ni sql.NullInt32) *int32 {
	if ni.Valid {
		v := ni.Int32
		return &v
	}
	return nil
}

func nullInt64Ptr(ni sql.NullInt64) *int64 {
	if ni.Valid {
		v := ni.Int64
		return &v
	}
	return nil
}

func nullBool(nb sql.NullBool) bool {
	return nb.Valid && nb.Bool
}

func nullTimePtr(nt sql.NullTime) *time.Time {
	if nt.Valid {
		t := nt.Time
		return &t
	}
	return nil
}

// ---- PrivJob responses ----
type PrivJobResponse struct {
	ID          string     `json:"id"`
	Company     string     `json:"company"`
	Title       string     `json:"title"`
	Location    string     `json:"location"`
	URL         string     `json:"url"`
	Salary      string     `json:"salary"`
	Experience  string     `json:"experience"`
	JobType     string     `json:"job_type"`
	Skills      []string   `json:"skills"`
	Description string     `json:"description"`
	Source      string     `json:"source"`
	PostedAt    *time.Time `json:"posted_at"`
	CreatedAt   *time.Time `json:"created_at"`
}

func toPrivJobResponse(r db.GetPrivJobsRow) PrivJobResponse {
	return PrivJobResponse{
		ID:          r.ID.String(),
		Company:     r.Company,
		Title:       r.Title,
		Location:    nullStr(r.Location),
		URL:         nullStr(r.Url),
		Salary:      nullStr(r.Salary),
		Experience:  nullStr(r.Experience),
		JobType:     nullStr(r.JobType),
		Skills:      r.Skills,
		Description: nullStr(r.Description),
		Source:      r.Source,
		PostedAt:    nullTimePtr(r.PostedAt),
		CreatedAt:   nullTimePtr(r.CreatedAt),
	}
}

func toPrivJobByIDResponse(r db.GetPrivJobByIDRow) PrivJobResponse {
	return PrivJobResponse{
		ID:          r.ID.String(),
		Company:     r.Company,
		Title:       r.Title,
		Location:    nullStr(r.Location),
		URL:         nullStr(r.Url),
		Salary:      nullStr(r.Salary),
		Experience:  nullStr(r.Experience),
		JobType:     nullStr(r.JobType),
		Skills:      r.Skills,
		Description: nullStr(r.Description),
		Source:      r.Source,
		PostedAt:    nullTimePtr(r.PostedAt),
		CreatedAt:   nullTimePtr(r.CreatedAt),
	}
}

// ---- GovJob responses ----
type GovJobResponse struct {
	ID           string     `json:"id"`
	Title        string     `json:"title"`
	Department   string     `json:"department"`
	Location     string     `json:"location"`
	ApplyURL     string     `json:"apply_url"`
	LastDate     *time.Time `json:"last_date"`
	Source       string     `json:"source"`
	Eligibility  string     `json:"eligibility"`
	VacancyCount *int32     `json:"vacancy_count"`
	Salary       string     `json:"salary"`
	ExamDate     *time.Time `json:"exam_date"`
	CreatedAt    *time.Time `json:"created_at"`
}

func toGovJobResponse(r db.GetGovJobsRow) GovJobResponse {
	return GovJobResponse{
		ID:           r.ID.String(),
		Title:        r.Title,
		Department:   nullStr(r.Department),
		Location:     nullStr(r.Location),
		ApplyURL:     nullStr(r.ApplyUrl),
		LastDate:     nullTimePtr(r.LastDate),
		Source:       r.Source,
		Eligibility:  nullStr(r.Eligibility),
		VacancyCount: nullInt32Ptr(r.VacancyCount),
		Salary:       nullStr(r.Salary),
		ExamDate:     nullTimePtr(r.ExamDate),
		CreatedAt:    nullTimePtr(r.CreatedAt),
	}
}

func toGovJobByIDResponse(r db.GetGovJobByIDRow) GovJobResponse {
	return GovJobResponse{
		ID:           r.ID.String(),
		Title:        r.Title,
		Department:   nullStr(r.Department),
		Location:     nullStr(r.Location),
		ApplyURL:     nullStr(r.ApplyUrl),
		LastDate:     nullTimePtr(r.LastDate),
		Source:       r.Source,
		Eligibility:  nullStr(r.Eligibility),
		VacancyCount: nullInt32Ptr(r.VacancyCount),
		Salary:       nullStr(r.Salary),
		ExamDate:     nullTimePtr(r.ExamDate),
		CreatedAt:    nullTimePtr(r.CreatedAt),
	}
}

// ---- Course responses ----
type CourseResponse struct {
	ID           string     `json:"id"`
	Provider     string     `json:"provider"`
	Title        string     `json:"title"`
	URL          string     `json:"url"`
	Duration     string     `json:"duration"`
	Mode         string     `json:"mode"`
	Level        string     `json:"level"`
	Skills       []string   `json:"skills"`
	Description  string     `json:"description"`
	ThumbnailURL string     `json:"thumbnail_url"`
	Price        string     `json:"price"`
	IsFree       bool       `json:"is_free"`
	Source       string     `json:"source"`
	CreatedAt    *time.Time `json:"created_at"`
}

func toCourseResponse(r db.GetCoursesRow) CourseResponse {
	return CourseResponse{
		ID:           r.ID.String(),
		Provider:     r.Provider,
		Title:        r.Title,
		URL:          r.Url,
		Duration:     nullStr(r.Duration),
		Mode:         nullStr(r.Mode),
		Level:        nullStr(r.Level),
		Skills:       r.Skills,
		Description:  nullStr(r.Description),
		ThumbnailURL: nullStr(r.ThumbnailUrl),
		Price:        nullStr(r.Price),
		IsFree:       nullBool(r.IsFree),
		Source:       r.Source,
		CreatedAt:    nullTimePtr(r.CreatedAt),
	}
}

func toCourseByIDResponse(r db.GetCourseByIDRow) CourseResponse {
	return CourseResponse{
		ID:           r.ID.String(),
		Provider:     r.Provider,
		Title:        r.Title,
		URL:          r.Url,
		Duration:     nullStr(r.Duration),
		Mode:         nullStr(r.Mode),
		Level:        nullStr(r.Level),
		Skills:       r.Skills,
		Description:  nullStr(r.Description),
		ThumbnailURL: nullStr(r.ThumbnailUrl),
		Price:        nullStr(r.Price),
		IsFree:       nullBool(r.IsFree),
		Source:       r.Source,
		CreatedAt:    nullTimePtr(r.CreatedAt),
	}
}

// ---- Video responses ----
type VideoResponse struct {
	ID          string     `json:"id"`
	Channel     string     `json:"channel"`
	ChannelID   string     `json:"channel_id"`
	Title       string     `json:"title"`
	URL         string     `json:"url"`
	Thumbnail   string     `json:"thumbnail"`
	Description string     `json:"description"`
	VideoID     string     `json:"video_id"`
	PublishedAt *time.Time `json:"published_at"`
	Duration    string     `json:"duration"`
	ViewCount   *int64     `json:"view_count"`
	LikeCount   *int64     `json:"like_count"`
	Category    string     `json:"category"`
	CreatedAt   *time.Time `json:"created_at"`
}

func toVideoResponse(r db.GetVideosRow) VideoResponse {
	return VideoResponse{
		ID:          r.ID.String(),
		Channel:     r.Channel,
		ChannelID:   nullStr(r.ChannelID),
		Title:       r.Title,
		URL:         r.Url,
		Thumbnail:   nullStr(r.Thumbnail),
		Description: nullStr(r.Description),
		VideoID:     r.VideoID,
		PublishedAt: nullTimePtr(r.PublishedAt),
		Duration:    nullStr(r.Duration),
		ViewCount:   nullInt64Ptr(r.ViewCount),
		LikeCount:   nullInt64Ptr(r.LikeCount),
		Category:    nullStr(r.Category),
		CreatedAt:   nullTimePtr(r.CreatedAt),
	}
}

func toVideoByIDResponse(r db.GetVideoByIDRow) VideoResponse {
	return VideoResponse{
		ID:          r.ID.String(),
		Channel:     r.Channel,
		ChannelID:   nullStr(r.ChannelID),
		Title:       r.Title,
		URL:         r.Url,
		Thumbnail:   nullStr(r.Thumbnail),
		Description: nullStr(r.Description),
		VideoID:     r.VideoID,
		PublishedAt: nullTimePtr(r.PublishedAt),
		Duration:    nullStr(r.Duration),
		ViewCount:   nullInt64Ptr(r.ViewCount),
		LikeCount:   nullInt64Ptr(r.LikeCount),
		Category:    nullStr(r.Category),
		CreatedAt:   nullTimePtr(r.CreatedAt),
	}
}
