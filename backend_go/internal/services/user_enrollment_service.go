package services

import (
	"time"

	"github.com/rojgarsetu/backend/internal/db"
)

type UserEnrollmentService struct {
	db *db.PostgresDB
}

func NewUserEnrollmentService(database *db.PostgresDB) *UserEnrollmentService {
	return &UserEnrollmentService{db: database}
}

func (s *UserEnrollmentService) GetUserEnrollments(filter db.UserEnrollmentFilter, userID string, page, limit int) ([]db.UserEnrollment, int, error) {
	return s.db.GetUserEnrollments(filter, userID, page, limit)
}

func (s *UserEnrollmentService) GetUserEnrollmentByID(id string) (*db.UserEnrollment, error) {
	return s.db.GetUserEnrollmentByID(id)
}

func (s *UserEnrollmentService) GetUserEnrollmentByUserAndTrade(userID string, tradeID string) (*db.UserEnrollment, error) {
	return s.db.GetUserEnrollmentByUserAndTrade(userID, tradeID)
}

func (s *UserEnrollmentService) GetExpiringEnrollments() ([]db.UserEnrollment, error) {
	return s.db.GetExpiringEnrollments()
}

func (s *UserEnrollmentService) GetExpiringEnrollmentsWithTrade() ([]db.GetExpiringEnrollmentsWithTradeRow, error) {
	return s.db.GetExpiringEnrollmentsWithTrade()
}

func (s *UserEnrollmentService) CreateUserEnrollment(userID string, tradeID string, expiresAt time.Time, metadata map[string]interface{}) (*db.UserEnrollment, error) {
	return s.db.CreateUserEnrollment(userID, tradeID, expiresAt, metadata)
}

func (s *UserEnrollmentService) UpdateUserEnrollment(id string, status string, expiresAt *time.Time, completedAt *time.Time, progressPct int32, metadata map[string]interface{}) (*db.UserEnrollment, error) {
	return s.db.UpdateUserEnrollment(id, status, expiresAt, completedAt, progressPct, metadata)
}

func (s *UserEnrollmentService) UpdateEnrollmentProgress(id string, progressPct int32) (*db.UserEnrollment, error) {
	return s.db.UpdateEnrollmentProgress(id, progressPct)
}

func (s *UserEnrollmentService) CompleteEnrollment(id string) (*db.UserEnrollment, error) {
	return s.db.CompleteEnrollment(id)
}

func (s *UserEnrollmentService) CancelEnrollment(id string) (*db.UserEnrollment, error) {
	return s.db.CancelEnrollment(id)
}
