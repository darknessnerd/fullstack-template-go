package logger

import (
	"os"
	"strings"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

var Log zerolog.Logger

// Initialize sets up the logger based on configuration
func Initialize(debug bool, logLevel string) {
	// Set log level
	level := parseLogLevel(logLevel)
	zerolog.SetGlobalLevel(level)

	// Configure output format
	if debug {
		// Pretty print for development
		Log = log.Output(zerolog.ConsoleWriter{
			Out:        os.Stdout,
			TimeFormat: time.RFC3339,
		})
	} else {
		// JSON format for production
		Log = zerolog.New(os.Stdout).With().Timestamp().Logger()
	}

	Log.Info().
		Str("level", level.String()).
		Bool("debug", debug).
		Msg("Logger initialized")
}

func parseLogLevel(logLevel string) zerolog.Level {
	switch strings.ToLower(logLevel) {
	case "trace":
		return zerolog.TraceLevel
	case "debug":
		return zerolog.DebugLevel
	case "info":
		return zerolog.InfoLevel
	case "warn", "warning":
		return zerolog.WarnLevel
	case "error":
		return zerolog.ErrorLevel
	case "fatal":
		return zerolog.FatalLevel
	case "panic":
		return zerolog.PanicLevel
	default:
		return zerolog.InfoLevel
	}
}
