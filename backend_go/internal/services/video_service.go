package services

import (
	"github.com/rojgarsetu/backend/internal/db"
)

type VideoService struct {
	db *db.PostgresDB
}

func NewVideoService(database *db.PostgresDB) *VideoService {
	return &VideoService{db: database}
}

func (s *VideoService) GetVideos(filter db.VideoFilter, page, limit int) ([]db.Video, int, error) {
	return s.db.GetVideos(filter, page, limit)
}

func (s *VideoService) GetVideoByID(id string) (*db.Video, error) {
	return s.db.GetVideoByID(id)
}
