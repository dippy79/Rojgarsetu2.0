package services

import (
    "github.com/rojgarsetu/backend/internal/db"
)

type PrivJobService struct {
    db *db.PostgresDB
}

func NewPrivJobService(database *db.PostgresDB) *PrivJobService {
    return &PrivJobService{db: database}
}

func (s *PrivJobService) GetPrivJobs(filter db.PrivJobFilter, page, limit int) ([]db.GetPrivJobsRow, int, error) {
    return s.db.GetPrivJobs(filter, page, limit)
}

func (s *PrivJobService) GetPrivJobByID(id string) (*db.GetPrivJobByIDRow, error) {
    return s.db.GetPrivJobByID(id)
}
