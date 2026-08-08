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

func (s *VideoService) GetVideos(filter db.VideoFilter, exclude string, page, limit int) ([]db.GetVideosRow, int, error) {
	return s.db.GetVideos(filter, exclude, page, limit)
}

func (s *VideoService) GetVideoByID(id string) (*db.GetVideoByIDRow, error) {
	return s.db.GetVideoByID(id)
}

func (s *VideoService) GetVideoChannels() ([]db.VideoChannel, error) {
	return s.db.GetVideoChannels()
}

func (s *VideoService) GetVideoCategories() ([]db.VideoCategory, error) {
	return s.db.GetVideoCategories()
}
