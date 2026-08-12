package services

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/rojgarsetu/backend/internal/db"
)

// NotificationEngine handles the periodic checking of expiring enrollments
// and sending notifications with strict daily limits
type NotificationEngine struct {
	enrollmentService   *UserEnrollmentService
	notificationService *NotificationService
	db                  *db.PostgresDB
}

func NewNotificationEngine(enrollmentService *UserEnrollmentService, notificationService *NotificationService, database *db.PostgresDB) *NotificationEngine {
	return &NotificationEngine{
		enrollmentService:   enrollmentService,
		notificationService: notificationService,
		db:                  database,
	}
}

// CheckAndNotifyExpiringEnrollments checks for enrollments expiring in <= 7 days
// and sends notifications with strict max 2/day per user limit
func (e *NotificationEngine) CheckAndNotifyExpiringEnrollments(ctx context.Context) error {
	log.Println("Starting expiring enrollment check...")

	enrollments, err := e.enrollmentService.GetExpiringEnrollmentsWithTrade()
	if err != nil {
		return fmt.Errorf("failed to get expiring enrollments: %w", err)
	}

	log.Printf("Found %d expiring enrollments", len(enrollments))

	for _, enrollment := range enrollments {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			if err := e.processEnrollment(ctx, enrollment); err != nil {
				log.Printf("Error processing enrollment %s: %v", enrollment.ID, err)
				// Continue with other enrollments
			}
		}
	}

	log.Println("Expiring enrollment check completed")
	return nil
}

func (e *NotificationEngine) processEnrollment(ctx context.Context, enrollment db.GetExpiringEnrollmentsWithTradeRow) error {
	userID := enrollment.UserID.String()
	enrollmentID := enrollment.ID.String()

	// Calculate days until expiry
	daysUntilExpiry := int(time.Until(enrollment.ExpiresAt).Hours() / 24)
	if daysUntilExpiry < 0 {
		daysUntilExpiry = 0
	}

	// Determine notification type based on days remaining
	var notificationType string
	var title string
	var message string

	if daysUntilExpiry <= 1 {
		notificationType = "expiry_final"
		title = "Enrollment Expires Today"
		message = fmt.Sprintf("Your enrollment in %s expires today. Complete your training to avoid losing progress.", enrollment.TradeName)
	} else if daysUntilExpiry <= 3 {
		notificationType = "expiry_warning"
		title = "Enrollment Expiring Soon"
		message = fmt.Sprintf("Your enrollment in %s expires in %d days. Please complete your training.", enrollment.TradeName, daysUntilExpiry)
	} else {
		notificationType = "expiry_warning"
		title = "Enrollment Expiring Soon"
		message = fmt.Sprintf("Your enrollment in %s expires in %d days.", enrollment.TradeName, daysUntilExpiry)
	}

	// Check if we can send notification (strict max 2/day per user)
	canSend, err := e.notificationService.CanSendNotificationByType(userID, notificationType)
	if err != nil {
		return fmt.Errorf("failed to check notification limit: %w", err)
	}

	if !canSend {
		log.Printf("Daily notification limit reached for user %s, type %s", userID, notificationType)
		return nil // Skip this notification but don't error
	}

	// Create notification payload
	payload := map[string]interface{}{
		"enrollment_id":     enrollmentID,
		"trade_id":          enrollment.TradeID.String(),
		"trade_name":        enrollment.TradeName,
		"trade_slug":        enrollment.TradeSlug,
		"category_name":     enrollment.CategoryName,
		"category_slug":     enrollment.CategorySlug,
		"category_color":    enrollment.CategoryColor,
		"category_icon":     enrollment.CategoryIcon,
		"days_until_expiry": daysUntilExpiry,
		"expires_at":        enrollment.ExpiresAt.Format(time.RFC3339),
		"progress_pct":      enrollment.ProgressPct,
	}

	// Send in-app notification
	_, err = e.notificationService.CreateNotificationLog(
		userID,
		&enrollmentID,
		notificationType,
		"in_app",
		title,
		message,
		payload,
	)
	if err != nil {
		return fmt.Errorf("failed to create notification: %w", err)
	}

	log.Printf("Sent %s notification to user %s for enrollment %s", notificationType, userID, enrollmentID)
	return nil
}

// StartPeriodicCheck starts a goroutine that periodically checks for expiring enrollments
// interval: how often to check (e.g., 6 hours)
func (e *NotificationEngine) StartPeriodicCheck(interval time.Duration) {
	go func() {
		ticker := time.NewTicker(interval)
		defer ticker.Stop()

		// Run immediately on start
		ctx := context.Background()
		if err := e.CheckAndNotifyExpiringEnrollments(ctx); err != nil {
			log.Printf("Initial expiring enrollment check failed: %v", err)
		}

		for {
			select {
			case <-ticker.C:
				ctx := context.Background()
				if err := e.CheckAndNotifyExpiringEnrollments(ctx); err != nil {
					log.Printf("Periodic expiring enrollment check failed: %v", err)
				}
			}
		}
	}()
}
