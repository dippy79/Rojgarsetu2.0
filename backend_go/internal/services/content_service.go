package services

import (
	"context"

	"github.com/rojgarsetu/backend/internal/db"
)

type ContentService struct {
	db *db.PostgresDB
}

func NewContentService(d *db.PostgresDB) *ContentService {
	return &ContentService{db: d}
}

func (s *ContentService) GetGovJobs(ctx context.Context, f db.GovJobFilter, page, limit int) ([]db.GetGovJobsRow, int, error) {
	return s.db.GetGovJobs(f, page, limit)
}

func (s *ContentService) GetGovJobByID(ctx context.Context, id string) (*db.GetGovJobByIDRow, error) {
	return s.db.GetGovJobByID(id)
}

func (s *ContentService) GetPrivJobs(ctx context.Context, f db.PrivJobFilter, page, limit int) ([]db.GetPrivJobsRow, int, error) {
	return s.db.GetPrivJobs(f, page, limit)
}

func (s *ContentService) GetPrivJobByID(ctx context.Context, id string) (*db.GetPrivJobByIDRow, error) {
	return s.db.GetPrivJobByID(id)
}

func (s *ContentService) GetCourses(ctx context.Context, f db.CourseFilter, page, limit int) ([]db.GetCoursesRow, int, error) {
	return s.db.GetCourses(f, page, limit)
}

func (s *ContentService) GetCourseByID(ctx context.Context, id string) (*db.GetCourseByIDRow, error) {
	return s.db.GetCourseByID(id)
}

func (s *ContentService) GetVideos(ctx context.Context, f db.VideoFilter, exclude string, page, limit int) ([]db.GetVideosRow, int, error) {
	return s.db.GetVideos(f, exclude, page, limit)
}

func (s *ContentService) GetVideoByID(ctx context.Context, id string) (*db.GetVideoByIDRow, error) {
	return s.db.GetVideoByID(id)
}
