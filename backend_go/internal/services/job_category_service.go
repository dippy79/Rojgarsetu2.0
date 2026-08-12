package services

import (
	"github.com/rojgarsetu/backend/internal/db"
)

type JobCategoryService struct {
	db *db.PostgresDB
}

func NewJobCategoryService(database *db.PostgresDB) *JobCategoryService {
	return &JobCategoryService{db: database}
}

func (s *JobCategoryService) GetJobCategories() ([]db.JobCategory, error) {
	return s.db.GetJobCategories()
}

func (s *JobCategoryService) GetJobCategoryBySlug(slug string) (*db.JobCategory, error) {
	return s.db.GetJobCategoryBySlug(slug)
}

func (s *JobCategoryService) GetJobCategoryByID(id string) (*db.JobCategory, error) {
	return s.db.GetJobCategoryByID(id)
}
