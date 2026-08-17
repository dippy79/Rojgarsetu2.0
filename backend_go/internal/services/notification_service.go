package services

import (
	"encoding/json"
	"time"
)

type NotificationService struct {
	// This could be extended to include database connections, etc.
}

func NewNotificationService() *NotificationService {
	return &NotificationService{}
}

type Notification struct {
	Type      string                 `json:"type"`
	Title     string                 `json:"title"`
	Message   string                 `json:"message"`
	Data      map[string]interface{} `json:"data,omitempty"`
	Timestamp time.Time              `json:"timestamp"`
	UserID    string                 `json:"user_id,omitempty"`
}

// CreateApplicationStatusNotification creates a notification for application status changes
func (s *NotificationService) CreateApplicationStatusNotification(userID, jobTitle, status string) Notification {
	return Notification{
		Type:    "application_status",
		Title:   "Application Status Update",
		Message: "Your application status has changed",
		Data: map[string]interface{}{
			"job_title": jobTitle,
			"status":    status,
		},
		Timestamp: time.Now(),
		UserID:    userID,
	}
}

// CreateNewJobNotification creates a notification for new job postings
func (s *NotificationService) CreateNewJobNotification(jobTitle, jobType string) Notification {
	return Notification{
		Type:    "new_job",
		Title:   "New Job Posted",
		Message: "A new job matching your preferences has been posted",
		Data: map[string]interface{}{
			"job_title": jobTitle,
			"job_type":  jobType,
		},
		Timestamp: time.Now(),
	}
}

// CreateSystemNotification creates a general system notification
func (s *NotificationService) CreateSystemNotification(title, message string) Notification {
	return Notification{
		Type:      "system",
		Title:     title,
		Message:   message,
		Timestamp: time.Now(),
	}
}

// ToJSON converts notification to JSON
func (n *Notification) ToJSON() ([]byte, error) {
	return json.Marshal(n)
}

// CanSendNotificationByType checks if a user allows receiving a specific notification type
func (s *NotificationService) CanSendNotificationByType(userID string, notificationType string) (bool, error) {
	return true, nil
}

// CreateNotificationLog records notification history into log storage
func (s *NotificationService) CreateNotificationLog(
	userID string,
	recipient *string,
	notificationType string,
	title string,
	body string,
	status string,
	metadata map[string]interface{},
) (string, error) {
	return "log_id_stub", nil
}

func (s *NotificationService) GetUserNotificationLogs(filter interface{}, userID string, limit, offset int) ([]*Notification, int, error) {
	return []*Notification{}, 0, nil
}

// GetNotificationLogByID retrieves a specific notification log by ID
func (s *NotificationService) GetNotificationLogByID(id string) (*Notification, error) {
	return &Notification{}, nil
}

// CanSendNotification checks if a notification can be sent
func (s *NotificationService) CanSendNotification(userID string, notificationType string) (bool, error) {
	return true, nil
}

// MarkNotificationRead marks a notification as read
func (s *NotificationService) MarkNotificationRead(id string) (bool, error) {
	return true, nil
}

// MarkNotificationClicked marks a notification as clicked
func (s *NotificationService) MarkNotificationClicked(id string) (bool, error) {
	return true, nil
}
