package logger

import (
	"os"

	"github.com/rs/zerolog"
)

func NewLogger(traceID string) zerolog.Logger {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	return zerolog.New(os.Stdout).With().
		Timestamp().
		Str("service", "backend").
		Str("trace_id", traceID).
		Logger().
		Level(zerolog.InfoLevel)
}

var Logger = NewLogger()
