package services

import (
	"github.com/rojgarsetu/backend/internal/db"
)

type NotificationService struct {
	db *db.PostgresDB
}

func NewNotificationService(database *db.PostgresDB) *NotificationService {
	return &NotificationService{db: database}
}

func (s *NotificationService) GetUserNotificationLogs(filter db.NotificationLogFilter, userID string, page, limit int) ([]db.UserNotificationLog, int, error) {
	return s.db.GetUserNotificationLogs(filter, userID, page, limit)
}

func (s *NotificationService) GetNotificationLogByID(id string) (*db.UserNotificationLog, error) {
	return s.db.GetNotificationLogByID(id)
}

func (s *NotificationService) GetDailyNotificationCount(userID string) (int64, error) {
	return s.db.GetDailyNotificationCount(userID)
}

func (s *NotificationService) GetDailyNotificationCountByType(userID string, notificationType string) (int64, error) {
	return s.db.GetDailyNotificationCountByType(userID, notificationType)
}

func (s *NotificationService) CreateNotificationLog(userID string, enrollmentID *string, notificationType string, channel string, title string, message string, payload map[string]interface{}) (*db.UserNotificationLog, error) {
	return s.db.CreateNotificationLog(userID, enrollmentID, notificationType, channel, title, message, payload)
}

func (s *NotificationService) MarkNotificationRead(id string) (*db.UserNotificationLog, error) {
	return s.db.MarkNotificationRead(id)
}

func (s *NotificationService) MarkNotificationClicked(id string) (*db.UserNotificationLog, error) {
	return s.db.MarkNotificationClicked(id)
}

// CanSendNotification checks if user has not exceeded the daily limit of 2 notifications
func (s *NotificationService) CanSendNotification(userID string, notificationType string) (bool, error) {
	count, err := s.GetDailyNotificationCount(userID)
	if err != nil {
		return false, err
	}
	// Strict max 2 notifications per day per user
	return count < 2, nil
}

// CanSendNotificationByType checks if user has not exceeded the daily limit for a specific notification type
func (s *NotificationService) CanSendNotificationByType(userID string, notificationType string) (bool, error) {
	count, err := s.GetDailyNotificationCountByType(userID, notificationType)
	if err != nil {
		return false, err
	}
	// Strict max 2 notifications per day per user per type
	return count < 2, nil
}
