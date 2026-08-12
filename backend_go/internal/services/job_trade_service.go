package services

import (
	"github.com/rojgarsetu/backend/internal/db"
)

type JobTradeService struct {
	db *db.PostgresDB
}

func NewJobTradeService(database *db.PostgresDB) *JobTradeService {
	return &JobTradeService{db: database}
}

func (s *JobTradeService) GetJobTrades(filter db.JobTradeFilter, page, limit int) ([]db.JobTrade, int, error) {
	return s.db.GetJobTrades(filter, page, limit)
}

func (s *JobTradeService) GetJobTradeByID(id string) (*db.JobTrade, error) {
	return s.db.GetJobTradeByID(id)
}

func (s *JobTradeService) GetJobTradeBySlug(slug string, categoryID string) (*db.JobTrade, error) {
	return s.db.GetJobTradeBySlug(slug, categoryID)
}

func (s *JobTradeService) GetJobTradesByCategory(categoryID string) ([]db.JobTrade, error) {
	return s.db.GetJobTradesByCategory(categoryID)
}
