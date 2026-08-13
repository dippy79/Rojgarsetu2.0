package services

import (
    "context"
    "database/sql"
    "errors"
)

type Notification struct {
    ID        string `json:"id"`
    UserID    string `json:"user_id"`
    Title     string `json:"title"`
    Body      string `json:"body"`
    IsRead    bool   `json:"is_read"`
    CreatedAt string `json:"created_at"`
}

type NotificationService struct {
    DB *sql.DB
}

func (s *NotificationService) GetUserNotifications(ctx context.Context, userID string) ([]Notification, error) {
    rows, err := s.DB.QueryContext(ctx, "SELECT id, user_id, title, body, is_read, created_at FROM notifications WHERE user_id=$1 ORDER BY is_read ASC, created_at DESC", userID)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var list []Notification
    for rows.Next() {
        var n Notification
        rows.Scan(&n.ID, &n.UserID, &n.Title, &n.Body, &n.IsRead, &n.CreatedAt)
        list = append(list, n)
    }
    return list, nil
}

func (s *NotificationService) MarkAsRead(ctx context.Context, id string) error {
    _, err := s.DB.ExecContext(ctx, "UPDATE notifications SET is_read=true WHERE id=$1", id)
    return err
}

func (s *NotificationService) SendNotification(ctx context.Context, userID, title, body string) error {
    var count int
    err := s.DB.QueryRowContext(ctx, 
        "SELECT COUNT(*) FROM user_notification_logs WHERE user_id=$1 AND created_at > NOW() - INTERVAL '24 hours'", 
        userID).Scan(&count)
    
    if err != nil {
        return err
    }
    if count >= 2 {
        return errors.New("daily_limit_reached")
    }

    _, err = s.DB.ExecContext(ctx, 
        "INSERT INTO notifications (user_id, title, body, is_read, created_at) VALUES ($1, $2, $3, false, NOW())", 
        userID, title, body)
    return err
}
