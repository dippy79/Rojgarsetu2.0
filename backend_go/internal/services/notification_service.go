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

// CanSendNotification checks if a notification can be sent (max 2 per day)
func (s *NotificationService) CanSendNotification(userID string, notificationType string) (bool, error) {
	count, err := s.db.GetDailyNotificationCount(userID)
	if err != nil {
		return false, err
	}
	return count < 2, nil
}

// Alias for compatibility
func (s *NotificationService) CanSendNotificationByType(userID string, notificationType string) (bool, error) {
	return s.CanSendNotification(userID, notificationType)
}

func (s *NotificationService) CreateNotificationLog(
	userID string,
	enrollmentID *string,
	notificationType string,
	channel string,
	title string,
	message string,
	payload map[string]interface{},
) (*db.UserNotificationLog, error) {
	return s.db.CreateNotificationLog(userID, enrollmentID, notificationType, channel, title, message, payload)
}

func (s *NotificationService) GetUserNotificationLogs(filter db.NotificationLogFilter, userID string, page, limit int) ([]db.UserNotificationLog, int, error) {
	return s.db.GetUserNotificationLogs(filter, userID, page, limit)
}

func (s *NotificationService) GetNotificationLogByID(id string) (*db.UserNotificationLog, error) {
	return s.db.GetNotificationLogByID(id)
}

func (s *NotificationService) MarkNotificationRead(id string) (*db.UserNotificationLog, error) {
	return s.db.MarkNotificationRead(id)
}

func (s *NotificationService) MarkNotificationClicked(id string) (*db.UserNotificationLog, error) {
	return s.db.MarkNotificationClicked(id)
}
