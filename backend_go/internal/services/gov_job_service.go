package services

import (
	"github.com/rojgarsetu/backend/internal/db"
)

type GovJobService struct {
	db *db.PostgresDB
}

func NewGovJobService(database *db.PostgresDB) *GovJobService {
	return &GovJobService{db: database}
}

func (s *GovJobService) GetGovJobs(filter db.GovJobFilter, page, limit int) ([]db.GovJob, int, error) {
	return s.db.GetGovJobs(filter, page, limit)
}

func (s *GovJobService) GetGovJobByID(id string) (*db.GovJob, error) {
	return s.db.GetGovJobByID(id)
}

