package platform

import (
	"log/slog"
	"os"
	"strings"
)

// Different log levels supported in the application
// by the env variable LOG_LEVEL
const (
	DEBUG = "DEBUG"
	INFO  = "INFO"
	WARN  = "WARN"
	ERROR = "ERROR"
)

// InitLogger initializes the logger with the LOG_LEVEL environment variable.
func initLogger() {
	loggerConfig := AppConfig.Logger
	opts := &slog.HandlerOptions{
		Level: loggerConfig.getSlogLevel(),
	}

	var handler slog.Handler
	handler = slog.NewJSONHandler(os.Stdout, opts)
	switch loggerConfig.Format {
	case "json":
		handler = slog.NewJSONHandler(os.Stdout, opts)
	default:
		handler = slog.NewTextHandler(os.Stdout, opts)
	}
	slog.SetDefault(slog.New(handler))
}

func (loggerCfg LoggerConfig) getSlogLevel() slog.Level {
	var logLevel slog.Level
	switch strings.ToUpper(loggerCfg.Level) {
	case DEBUG:
		logLevel = slog.LevelDebug
	case INFO:
		logLevel = slog.LevelInfo
	case WARN:
		logLevel = slog.LevelWarn
	case ERROR:
		logLevel = slog.LevelError
	default:
		// means if variable is not set or invalid
		// we fallback to INFO
		logLevel = slog.LevelInfo
	}
	return logLevel
}
