package workers

import (
    "context"
    "database/sql"
    "fmt"
    "net/smtp"
    "os"
    "time"
)

type EmailWorker struct {
    DB *sql.DB
}

func (w *EmailWorker) Start(ctx context.Context) {
    ticker := time.NewTicker(5 * time.Minute)
    defer ticker.Stop()

    for {
        select {
        case <-ctx.Done():
            return
        case <-ticker.C:
            w.processQueue(ctx)
        }
    }
}

func (w *EmailWorker) processQueue(ctx context.Context) {
    rows, err := w.DB.QueryContext(ctx, "SELECT id, to_email, subject, body, attempts FROM email_queue WHERE status='pending' LIMIT 10")
    if err != nil {
        return
    }
    defer rows.Close()

    smtpHost := os.Getenv("SMTP_HOST")
    smtpPort := os.Getenv("SMTP_PORT")
    auth := smtp.PlainAuth("", os.Getenv("SMTP_USER"), os.Getenv("SMTP_PASS"), smtpHost)

    for rows.Next() {
        var id, to, subject, body string
        var attempts int
        rows.Scan(&id, &to, &subject, &body, &attempts)

        msg := []byte(fmt.Sprintf("To: %s\r\nSubject: %s\r\n\r\n%s", to, subject, body))
        err := smtp.SendMail(smtpHost+":"+smtpPort, auth, os.Getenv("SMTP_USER"), []string{to}, msg)

        if err == nil {
            w.DB.ExecContext(ctx, "UPDATE email_queue SET status='sent', sent_at=NOW() WHERE id=$1", id)
        } else {
            if attempts+1 >= 3 {
                w.DB.ExecContext(ctx, "UPDATE email_queue SET status='failed', attempts=attempts+1 WHERE id=$1", id)
            } else {
                w.DB.ExecContext(ctx, "UPDATE email_queue SET attempts=attempts+1 WHERE id=$1", id)
            }
        }
    }
}
