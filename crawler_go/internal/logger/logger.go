package logger

import (
	"os"

	"github.com/rs/zerolog"
)

func NewLogger() zerolog.Logger {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	return zerolog.New(os.Stdout).With().
		Timestamp().
		Str("service", "crawler").
		Logger().
		Level(zerolog.InfoLevel)
}

var Logger = NewLogger()
