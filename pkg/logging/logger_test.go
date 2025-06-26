package logging

import (
	"bytes"
	"context"
	"encoding/json"
	"log/slog"
	"testing"
)

func TestNewLogger(t *testing.T) {
	tests := []struct {
		name      string
		level     string
		wantLevel slog.Level
	}{
		{"trace level", "trace", LevelTrace},
		{"debug level", "debug", slog.LevelDebug},
		{"info level", "info", slog.LevelInfo},
		{"error level", "error", slog.LevelError},
		{"fatal level", "fatal", LevelFatal},
		{"uppercase level", "INFO", slog.LevelInfo},
		{"mixed case level", "Debug", slog.LevelDebug},
		{"invalid level defaults to fatal", "invalid", LevelFatal},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var buf bytes.Buffer
			logger := NewLoggerWithWriter(tt.level, &buf).(*slogLogger)

			if logger.level != tt.wantLevel {
				t.Errorf("NewLogger() level = %v, want %v", logger.level, tt.wantLevel)
			}
		})
	}
}

func TestLogger_LoggingLevels(t *testing.T) {
	tests := []struct {
		name          string
		loggerLevel   string
		logMethod     string
		expectedLevel string
		shouldLog     bool
	}{
		// Fatal level logger (most restrictive)
		{"fatal logger, fatal message", "fatal", "fatal", "FATAL", true},
		{"fatal logger, error message", "fatal", "error", "ERROR", false},
		{"fatal logger, info message", "fatal", "info", "INFO", false},
		{"fatal logger, debug message", "fatal", "debug", "DEBUG", false},
		{"fatal logger, trace message", "fatal", "trace", "TRACE", false},

		// Error level logger
		{"error logger, fatal message", "error", "fatal", "FATAL", true},
		{"error logger, error message", "error", "error", "ERROR", true},
		{"error logger, info message", "error", "info", "INFO", false},
		{"error logger, debug message", "error", "debug", "DEBUG", false},
		{"error logger, trace message", "error", "trace", "TRACE", false},

		// Info level logger
		{"info logger, fatal message", "info", "fatal", "FATAL", true},
		{"info logger, error message", "info", "error", "ERROR", true},
		{"info logger, info message", "info", "info", "INFO", true},
		{"info logger, debug message", "info", "debug", "DEBUG", false},
		{"info logger, trace message", "info", "trace", "TRACE", false},

		// Debug level logger
		{"debug logger, fatal message", "debug", "fatal", "FATAL", true},
		{"debug logger, error message", "debug", "error", "ERROR", true},
		{"debug logger, info message", "debug", "info", "INFO", true},
		{"debug logger, debug message", "debug", "debug", "DEBUG", true},
		{"debug logger, trace message", "debug", "trace", "TRACE", false},

		// Trace level logger (most verbose)
		{"trace logger, fatal message", "trace", "fatal", "FATAL", true},
		{"trace logger, error message", "trace", "error", "ERROR", true},
		{"trace logger, info message", "trace", "info", "INFO", true},
		{"trace logger, debug message", "trace", "debug", "DEBUG", true},
		{"trace logger, trace message", "trace", "trace", "TRACE", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var buf bytes.Buffer
			logger := NewLoggerWithWriter(tt.loggerLevel, &buf)

			// Skip fatal as it calls os.Exit
			if tt.logMethod == "fatal" {
				return
			}

			// Call the appropriate logging method
			switch tt.logMethod {
			case "error":
				logger.Error("test message", "key", "value")
			case "info":
				logger.Info("test message", "key", "value")
			case "debug":
				logger.Debug("test message", "key", "value")
			case "trace":
				logger.Trace("test message", "key", "value")
			}

			output := buf.String()

			if tt.shouldLog {
				if output == "" {
					t.Errorf("Expected log output but got none")
					return
				}

				// Parse JSON output
				var logEntry map[string]interface{}
				if err := json.Unmarshal([]byte(output), &logEntry); err != nil {
					t.Errorf("Failed to parse log output as JSON: %v", err)
					return
				}

				// Check level
				if level, ok := logEntry["level"].(string); !ok || level != tt.expectedLevel {
					t.Errorf("Expected level %s, got %v", tt.expectedLevel, level)
				}

				// Check message
				if msg, ok := logEntry["msg"].(string); !ok || msg != "test message" {
					t.Errorf("Expected message 'test message', got %v", msg)
				}

				// Check additional fields
				if key, ok := logEntry["key"].(string); !ok || key != "value" {
					t.Errorf("Expected key field 'value', got %v", key)
				}
			} else {
				if output != "" {
					t.Errorf("Expected no log output but got: %s", output)
				}
			}
		})
	}
}

func TestLogger_With(t *testing.T) {
	var buf bytes.Buffer
	logger := NewLoggerWithWriter("info", &buf)

	// Create logger with additional context
	contextLogger := logger.With("component", "test", "version", "1.0")
	contextLogger.Info("test message")

	output := buf.String()
	if output == "" {
		t.Fatal("Expected log output but got none")
	}

	// Parse JSON output
	var logEntry map[string]interface{}
	if err := json.Unmarshal([]byte(output), &logEntry); err != nil {
		t.Fatalf("Failed to parse log output as JSON: %v", err)
	}

	// Check that context fields are present
	if component, ok := logEntry["component"].(string); !ok || component != "test" {
		t.Errorf("Expected component 'test', got %v", component)
	}

	if version, ok := logEntry["version"].(string); !ok || version != "1.0" {
		t.Errorf("Expected version '1.0', got %v", version)
	}
}

func TestLogger_WithContext(t *testing.T) {
	var buf bytes.Buffer
	logger := NewLoggerWithWriter("info", &buf)

	ctx := context.Background()
	contextLogger := logger.WithContext(ctx)

	// Should be able to log without issues
	contextLogger.Info("test message with context")

	output := buf.String()
	if output == "" {
		t.Fatal("Expected log output but got none")
	}

	// Parse JSON output
	var logEntry map[string]interface{}
	if err := json.Unmarshal([]byte(output), &logEntry); err != nil {
		t.Fatalf("Failed to parse log output as JSON: %v", err)
	}

	if msg, ok := logEntry["msg"].(string); !ok || msg != "test message with context" {
		t.Errorf("Expected message 'test message with context', got %v", msg)
	}
}

func TestParseLogLevel(t *testing.T) {
	tests := []struct {
		name  string
		level string
		want  slog.Level
	}{
		{"trace", "trace", LevelTrace},
		{"debug", "debug", slog.LevelDebug},
		{"info", "info", slog.LevelInfo},
		{"error", "error", slog.LevelError},
		{"fatal", "fatal", LevelFatal},
		{"uppercase", "INFO", slog.LevelInfo},
		{"mixed case", "Debug", slog.LevelDebug},
		{"invalid", "invalid", LevelFatal},
		{"empty", "", LevelFatal},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := parseLogLevel(tt.level); got != tt.want {
				t.Errorf("parseLogLevel(%q) = %v, want %v", tt.level, got, tt.want)
			}
		})
	}
}

func TestSlogLogger_IsLevelEnabled(t *testing.T) {
	tests := []struct {
		name        string
		loggerLevel string
		checkLevel  string
		want        bool
	}{
		{"fatal logger, check fatal", "fatal", "fatal", true},
		{"fatal logger, check error", "fatal", "error", false},
		{"error logger, check error", "error", "error", true},
		{"error logger, check info", "error", "info", false},
		{"info logger, check info", "info", "info", true},
		{"info logger, check debug", "info", "debug", false},
		{"debug logger, check debug", "debug", "debug", true},
		{"debug logger, check trace", "debug", "trace", false},
		{"trace logger, check trace", "trace", "trace", true},
		{"trace logger, check debug", "trace", "debug", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var buf bytes.Buffer
			logger := NewLoggerWithWriter(tt.loggerLevel, &buf).(*slogLogger)

			if got := logger.IsLevelEnabled(tt.checkLevel); got != tt.want {
				t.Errorf("IsLevelEnabled(%q) = %v, want %v", tt.checkLevel, got, tt.want)
			}
		})
	}
}

func TestSlogLogger_GetLevel(t *testing.T) {
	tests := []struct {
		name        string
		loggerLevel string
		want        string
	}{
		{"trace", "trace", "trace"},
		{"debug", "debug", "debug"},
		{"info", "info", "info"},
		{"error", "error", "error"},
		{"fatal", "fatal", "fatal"},
		{"invalid defaults to fatal", "invalid", "fatal"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var buf bytes.Buffer
			logger := NewLoggerWithWriter(tt.loggerLevel, &buf).(*slogLogger)

			if got := logger.GetLevel(); got != tt.want {
				t.Errorf("GetLevel() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestLogLevels(t *testing.T) {
	levels := LogLevels()
	expected := []string{"trace", "debug", "info", "error", "fatal"}

	if len(levels) != len(expected) {
		t.Errorf("LogLevels() returned %d levels, expected %d", len(levels), len(expected))
	}

	for i, level := range expected {
		if i >= len(levels) || levels[i] != level {
			t.Errorf("LogLevels()[%d] = %v, want %v", i, levels[i], level)
		}
	}
}

func TestLogger_OutputFormat(t *testing.T) {
	var buf bytes.Buffer
	logger := NewLoggerWithWriter("info", &buf)

	logger.Info("test message", "request_id", "12345", "user", "test@example.com")

	output := buf.String()
	if output == "" {
		t.Fatal("Expected log output but got none")
	}

	// Should be valid JSON
	var logEntry map[string]interface{}
	if err := json.Unmarshal([]byte(output), &logEntry); err != nil {
		t.Fatalf("Log output is not valid JSON: %v\nOutput: %s", err, output)
	}

	// Should contain expected fields
	expectedFields := map[string]string{
		"level":      "INFO",
		"msg":        "test message",
		"request_id": "12345",
		"user":       "test@example.com",
	}

	for key, expectedValue := range expectedFields {
		if value, ok := logEntry[key].(string); !ok || value != expectedValue {
			t.Errorf("Expected %s = %s, got %v", key, expectedValue, value)
		}
	}

	// Should have timestamp
	if _, ok := logEntry["time"]; !ok {
		t.Error("Expected 'time' field in log output")
	}
}

func TestLogger_OutputGoesToStderr(t *testing.T) {
	// This test verifies that NewLogger (without writer) uses stderr
	// We can't easily test this directly, but we can verify the constructor
	logger := NewLogger("info")
	if logger == nil {
		t.Error("NewLogger() returned nil")
	}

	// Should be able to use the logger without issues
	logger.Info("test message")
}