package logger

import (
	"os"
	"strings"

	"github.com/rs/zerolog"
)

// SanitizeEmail masks an email address for GDPR compliance
func SanitizeEmail(email string) string {
	if email == "" {
		return ""
	}
	parts := strings.Split(email, "@")
	if len(parts) != 2 {
		return "***@***"
	}
	name := parts[0]
	if len(name) <= 2 {
		return "*@" + parts[1]
	}
	return name[:2] + "***@" + parts[1]
}

func NewLogger(traceID string) zerolog.Logger {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	return zerolog.New(os.Stdout).With().
		Timestamp().
		Str("service", "backend").
		Str("trace_id", traceID).
		Logger().
		Level(zerolog.InfoLevel)
}

var Logger = NewLogger("")
