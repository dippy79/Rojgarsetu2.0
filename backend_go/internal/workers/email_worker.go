package workers

import (
	"context"
	"database/sql"
	"log"
	"net/smtp"
	"os"
	"time"

	"github.com/rojgarsetu/backend/internal/db"
)

type EmailWorker struct {
	db *db.PostgresDB
}

func NewEmailWorker(database *db.PostgresDB) *EmailWorker {
	return &EmailWorker{db: database}
}

func (w *EmailWorker) Start(ctx context.Context) {
	ticker := time.NewTicker(5 * time.Minute)
	defer ticker.Stop()

	log.Println("Email worker started")

	for {
		select {
		case <-ctx.Done():
			log.Println("Email worker stopping...")
			return
		case <-ticker.C:
			w.processQueue(ctx)
		}
	}
}

func (w *EmailWorker) processQueue(ctx context.Context) {
	emails, err := w.db.Queries.GetPendingEmails(ctx, 10)
	if err != nil {
		log.Printf("Error fetching pending emails: %v", err)
		return
	}

	for _, email := range emails {
		if err := w.sendEmail(email); err != nil {
			log.Printf("Failed to send email to %s: %v", email.ToEmail, err)
			status := "pending"
			if email.Attempts >= 2 {
				status = "failed"
			}
			_, _ = w.db.Queries.UpdateEmailStatus(ctx, db.UpdateEmailStatusParams{
				ID:           email.ID,
				Status:       status,
				ErrorMessage: sql.NullString{String: err.Error(), Valid: true},
			})
		} else {
			log.Printf("Successfully sent email to %s", email.ToEmail)
			_, _ = w.db.Queries.UpdateEmailStatus(ctx, db.UpdateEmailStatusParams{
				ID:     email.ID,
				Status: "sent",
			})
		}
	}
}

func (w *EmailWorker) sendEmail(email db.EmailQueue) error {
	host := os.Getenv("SMTP_HOST")
	port := os.Getenv("SMTP_PORT")
	user := os.Getenv("SMTP_USER")
	pass := os.Getenv("SMTP_PASS")

	if host == "" || port == "" || user == "" || pass == "" {
		// Mock send if credentials are missing
		log.Printf("[MOCK EMAIL] To: %s, Subject: %s", email.ToEmail, email.Subject)
		return nil
	}

	auth := smtp.PlainAuth("", user, pass, host)
	msg := []byte("To: " + email.ToEmail + "\r\n" +
		"Subject: " + email.Subject + "\r\n" +
		"\r\n" +
		email.Body + "\r\n")

	err := smtp.SendMail(host+":"+port, auth, user, []string{email.ToEmail}, msg)
	return err
}
