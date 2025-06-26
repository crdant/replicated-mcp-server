// Package logging provides structured logging capabilities for the Replicated MCP Server.
// It uses Go's slog package with JSON output directed to stderr, maintaining separation
// between logging and MCP protocol communication on stdout.
package logging

import (
	"context"
	"io"
	"log/slog"
	"os"
	"strings"
)

// Logger interface for structured logging with multiple levels
type Logger interface {
	Fatal(msg string, args ...any)
	Error(msg string, args ...any)
	Info(msg string, args ...any)
	Debug(msg string, args ...any)
	Trace(msg string, args ...any)
	With(args ...any) Logger
	WithContext(ctx context.Context) Logger
}

// slogLogger implements Logger using Go's slog package
type slogLogger struct {
	logger *slog.Logger
	level  slog.Level
}

// Custom log levels
const (
	LevelTrace = slog.Level(-8) // More verbose than Debug (-4)
	LevelFatal = slog.Level(12) // More severe than Error (8)
)

// NewLogger creates a new structured logger with the specified level
// All logs are directed to stderr to keep stdout available for MCP protocol
func NewLogger(level string) Logger {
	return NewLoggerWithWriter(level, os.Stderr)
}

// NewLoggerWithWriter creates a logger with a custom writer (useful for testing)
func NewLoggerWithWriter(level string, writer io.Writer) Logger {
	slogLevel := parseLogLevel(level)

	// Create custom handler options
	opts := &slog.HandlerOptions{
		Level: slogLevel,
		ReplaceAttr: func(_ []string, a slog.Attr) slog.Attr {
			// Customize level names for our custom levels
			if a.Key == slog.LevelKey {
				switch a.Value.Any().(slog.Level) {
				case LevelTrace:
					a.Value = slog.StringValue("TRACE")
				case LevelFatal:
					a.Value = slog.StringValue("FATAL")
				case slog.LevelDebug:
					a.Value = slog.StringValue("DEBUG")
				case slog.LevelInfo:
					a.Value = slog.StringValue("INFO")
				case slog.LevelWarn:
					a.Value = slog.StringValue("WARN")
				case slog.LevelError:
					a.Value = slog.StringValue("ERROR")
				}
			}
			return a
		},
	}

	// Use JSON handler for structured logging
	handler := slog.NewJSONHandler(writer, opts)
	logger := slog.New(handler)

	return &slogLogger{
		logger: logger,
		level:  slogLevel,
	}
}

// Log level constants
const (
	logLevelTrace = "trace"
	logLevelDebug = "debug"
	logLevelInfo  = "info"
	logLevelError = "error"
	logLevelFatal = "fatal"
)

// parseLogLevel converts string level to slog.Level
func parseLogLevel(level string) slog.Level {
	switch strings.ToLower(level) {
	case logLevelTrace:
		return LevelTrace
	case logLevelDebug:
		return slog.LevelDebug
	case logLevelInfo:
		return slog.LevelInfo
	case logLevelError:
		return slog.LevelError
	case logLevelFatal:
		return LevelFatal
	default:
		return LevelFatal // Default to most restrictive
	}
}

// Fatal logs at fatal level and exits the program
func (l *slogLogger) Fatal(msg string, args ...any) {
	l.logger.Log(context.Background(), LevelFatal, msg, args...)
	os.Exit(1)
}

// Error logs at error level
func (l *slogLogger) Error(msg string, args ...any) {
	l.logger.Log(context.Background(), slog.LevelError, msg, args...)
}

// Info logs at info level
func (l *slogLogger) Info(msg string, args ...any) {
	l.logger.Log(context.Background(), slog.LevelInfo, msg, args...)
}

// Debug logs at debug level
func (l *slogLogger) Debug(msg string, args ...any) {
	l.logger.Log(context.Background(), slog.LevelDebug, msg, args...)
}

// Trace logs at trace level (most verbose)
func (l *slogLogger) Trace(msg string, args ...any) {
	l.logger.Log(context.Background(), LevelTrace, msg, args...)
}

// With returns a new logger with additional context fields
func (l *slogLogger) With(args ...any) Logger {
	return &slogLogger{
		logger: l.logger.With(args...),
		level:  l.level,
	}
}

// WithContext returns a new logger with context
func (l *slogLogger) WithContext(_ context.Context) Logger {
	// For now, return the same logger
	// In the future, we could extract values from context
	return l
}

// IsLevelEnabled checks if the given level is enabled for this logger
func (l *slogLogger) IsLevelEnabled(level string) bool {
	return parseLogLevel(level) >= l.level
}

// GetLevel returns the current log level as a string
func (l *slogLogger) GetLevel() string {
	switch l.level {
	case LevelTrace:
		return logLevelTrace
	case slog.LevelDebug:
		return logLevelDebug
	case slog.LevelInfo:
		return logLevelInfo
	case slog.LevelWarn:
		return "warn"
	case slog.LevelError:
		return logLevelError
	case LevelFatal:
		return logLevelFatal
	default:
		return "unknown"
	}
}

// LogLevels returns all valid log level names
func LogLevels() []string {
	return []string{logLevelTrace, logLevelDebug, logLevelInfo, logLevelError, logLevelFatal}
}
