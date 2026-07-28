package services

import (
    "github.com/rojgarsetu/backend/internal/db"
)

type CourseService struct {
    db *db.PostgresDB
}

func NewCourseService(database *db.PostgresDB) *CourseService {
    return &CourseService{db: database}
}

func (s *CourseService) GetCourses(filter db.CourseFilter, page, limit int) ([]db.GetCoursesRow, int, error) {
    return s.db.GetCourses(filter, page, limit)
}

func (s *CourseService) GetCourseByID(id string) (*db.GetCourseByIDRow, error) {
    return s.db.GetCourseByID(id)
}
